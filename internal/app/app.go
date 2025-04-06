package app

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mtchuikov/shortener/internal/cache/inmemory"
	"github.com/mtchuikov/shortener/internal/config"
	"github.com/mtchuikov/shortener/internal/handlers"
	"github.com/mtchuikov/shortener/internal/services"
	"github.com/mtchuikov/shortener/pkg/closer"
	"github.com/mtchuikov/shortener/pkg/logging"
	"github.com/mtchuikov/shortener/pkg/middlewares"
	"github.com/mtchuikov/shortener/pkg/pinger"
	"github.com/mtchuikov/shortener/pkg/postgres"
	"github.com/rs/zerolog"
)

type app struct {
	log      zerolog.Logger
	config   config.Config
	closer   *closer.Closer
	pgPool   *pgxpool.Pool
	pgPinger *pinger.Pinger
	server   *http.Server
}

func (a *app) setupPostgres(ctx context.Context) {
	var err error
	a.pgPool, err = postgres.ConnectPool(ctx, a.config.DatabaseDSN)
	if err != nil {
		a.log.Fatal().
			Err(err).
			Msg("failed to connect to postgres")
	}

	a.pgPinger = pinger.New(a.log, a.pgPool)

	closerFn := func(ctx context.Context) {
		a.pgPinger.Close(ctx)
		a.pgPool.Close()

		a.log.Info().Msg("postgres connection shut down")
	}

	task := closer.Task{Fn: closerFn}
	a.closer.Add(task)
}

func (a *app) newRouter() *chi.Mux {
	router := chi.NewMux()
	router.Use(middleware.Recoverer)

	if a.config.Verbose {
		verbose := middlewares.Verbose(a.log)
		router.Use(verbose)
	}

	compress := middleware.Compress(5)
	router.Use(middlewares.Decompress, compress)

	cache, err := inmemory.New(a.config.FileStorage)
	if err != nil {
		a.log.Fatal().
			Err(err).
			Msg("failed to setup cache")
	}

	shortenerService := services.NewShortener(a.config.BaseURL, cache)
	handlers.RegisterShortener(a.log, router, shortenerService)

	resolverService := services.NewResolver(a.config.BaseURL, cache)
	handlers.RegisterResolver(a.log, router, resolverService)

	handlers.RegisterPing(router, a.pgPinger)

	return router
}

func (a *app) setupServer() {
	a.server = &http.Server{
		Addr:         a.config.ServerAddr,
		Handler:      a.newRouter(),
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	closeFn := func(ctx context.Context) {
		err := a.server.Shutdown(ctx)
		if err != nil {
			a.log.Error().
				Err(err).
				Msg("error when shutting down server")
		}

		a.log.Info().Msg("server shut down")
	}

	task := closer.Task{Fn: closeFn}
	a.closer.AddWithPriority(task, 0)
}

func Setup(ctx context.Context) *app {
	config := config.New()
	log := logging.NewZerolog(config.ServiceName)

	app := &app{
		log:    log,
		config: config,
		closer: closer.New(),
	}

	app.setupPostgres(ctx)
	app.setupServer()

	return app
}

func (a *app) Run(ctx context.Context) {
	a.log.Info().Msg("running app...")

	a.log.Info().Msg("running postgres pinger...")
	go a.pgPinger.Ping(ctx, 2*time.Second)

	a.log.Info().
		Str("addr", a.config.ServerAddr).
		Msg("listening serevr...")

	err := a.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		a.log.Fatal().
			Err(err).
			Msg("failed to start server")
	}
}

func (a *app) Shutdown(ctx context.Context, timeout time.Duration) {
	a.log.Info().Msg("shutting down app...")

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	a.closer.Close(timeoutCtx)
	cancel()

	a.log.Info().Msg("app shut down")
}
