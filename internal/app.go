package app

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mtchuikov/shortener/internal/config"
	"github.com/mtchuikov/shortener/internal/handlers"
	"github.com/mtchuikov/shortener/internal/repository"
	"github.com/mtchuikov/shortener/internal/repository/inmemory"
	"github.com/mtchuikov/shortener/internal/repository/postgres"
	"github.com/mtchuikov/shortener/internal/services"
	"github.com/mtchuikov/shortener/pkg/closer"
	"github.com/mtchuikov/shortener/pkg/middlewares"
	"github.com/mtchuikov/shortener/pkg/pinger"
	"github.com/rs/zerolog"
)

type app struct {
	log           zerolog.Logger
	pgxPool       *pgxpool.Pool
	pgxConn       *pgxpool.Conn
	pgxPoolPinger *pinger.Pinger
	server        *http.Server
}

func (a *app) initPgxPool(ctx context.Context) {
	var err error
	a.pgxPool, err = pgxpool.New(ctx, config.DatabaseDSN())
	if err != nil {
		a.log.Fatal().Err(err).
			Msg("failed to connect to postgres")
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	err = a.pgxPool.Ping(timeoutCtx)
	if err != nil {
		a.log.Fatal().Err(err).
			Msg("failed to ping postgres")
	}

	timeoutCtx, cancel = context.WithTimeout(ctx, time.Second)
	defer cancel()

	a.pgxConn, err = a.pgxPool.Acquire(timeoutCtx)
	if err != nil {
		a.log.Fatal().Err(err).
			Msg("failed to acquire postgres connection")
	}

	a.pgxPoolPinger = pinger.New(a.log, a.pgxPool)

	shutdownPgxPoolFn := func(ctx context.Context) {
		a.log.Info().Msg("shutting down postgres connection...")

		a.pgxConn.Release()
		a.pgxPool.Close()
		a.pgxPoolPinger.Close(ctx)

		a.log.Info().Msg("postgres connection shut down")
	}

	closer.Add(closer.Task{
		Fn: shutdownPgxPoolFn,
	})
}

func (a *app) newHandler(ctx context.Context) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	if config.Verbose() {
		verbose := middlewares.Verbose(a.log)
		router.Use(verbose)
	}

	compress := middleware.Compress(5)
	router.Use(middlewares.Decompress, compress)

	var (
		repo repository.Repo
		err  error
	)

	inmemory, err := inmemory.New(config.FileStorage())
	if err != nil {
		a.log.Fatal().Err(err).
			Msg("failed to setup inmemory repo")
	}

	if config.DatabaseDSN() != "" {
		postgres := postgres.New(a.pgxConn)
		err = postgres.CreateTable(ctx)
		if err != nil {
			a.log.Fatal().Err(err).
				Msg("failed to create postgres tables")
		}
	} else {
		repo = inmemory
	}

	shortener := services.NewShortener(config.BaseURL(), repo)
	handlers.RegisterShortener(a.log, router, shortener)

	resolver := services.NewResolver(repo)
	handlers.RegisterResolver(router, resolver)

	handlers.RegisterPinger(router, a.pgxPoolPinger)

	return router
}

func (a *app) initServer(ctx context.Context) {
	a.server = &http.Server{
		Addr:         config.ServerAddr(),
		Handler:      a.newHandler(ctx),
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	shutdownServerFn := func(ctx context.Context) {
		a.log.Info().Msg("shutting down server...")

		err := a.server.Shutdown(ctx)
		if err != nil {
			a.log.Error().Err(err).
				Msg("server shut down with error")
		}

		a.log.Info().Msg("server shut down")
	}

	closer.AddWithPriority(0, closer.Task{
		Sync: true,
		Fn:   shutdownServerFn,
	})
}

func New(ctx context.Context, log zerolog.Logger) *app {
	app := &app{log: log}

	if config.DatabaseDSN() != "" {
		app.initPgxPool(ctx)
	}

	app.initServer(ctx)

	return app
}

func (a *app) Run(ctx context.Context) {
	a.log.Info().Msg("running app...")

	if config.DatabaseDSN() != "" {
		a.log.Info().Msg("pinging postgres...")
		go a.pgxPoolPinger.Ping(ctx, 2*time.Second)
	}

	a.log.Info().Str("addr", config.ServerAddr()).
		Msg("listening server...")

	err := a.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		a.log.Fatal().Err(err).
			Msg("failed to start server")
	}
}

func (a *app) Shutdown(ctx context.Context) {
	err := closer.Close(ctx)
	if err != nil {
		a.log.Error().Err(err).
			Msg("failed to gracefully shutdown services")
	}

	a.log.Info().Msg("services gracefully shut down")
}
