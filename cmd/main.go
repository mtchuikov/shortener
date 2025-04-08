package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/mtchuikov/shortener/internal/config"
	"github.com/mtchuikov/shortener/internal/handlers"
	"github.com/mtchuikov/shortener/internal/repo/inmemory"
	"github.com/mtchuikov/shortener/internal/repo/postgres"
	"github.com/mtchuikov/shortener/internal/services"
	"github.com/mtchuikov/shortener/pkg/closer"
	"github.com/mtchuikov/shortener/pkg/middlewares"
	"github.com/mtchuikov/shortener/pkg/pgxutils"
	"github.com/mtchuikov/shortener/pkg/pinger"
)

type repo interface {
	CreateShortURL(ctx context.Context, originalURL, shortID string) error
	GetOriginalURL(ctx context.Context, shortID string) (string, error)
	GetShortID(ctx context.Context, originalURL string) (string, error)
}

func main() {
	rootCtx := context.Background()

	signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL}
	stopCtx, stop := signal.NotifyContext(rootCtx, signals...)
	defer stop()

	zerolog.LevelFieldName = "lvl"
	zerolog.ErrorFieldName = "err"
	zerolog.MessageFieldName = "msg"
	zerolog.TimeFieldFormat = time.RFC1123

	log.Logger = log.Logger.
		Level(zerolog.InfoLevel).With().
		Timestamp().Str("app", config.ServiceName()).
		Logger()

	log.Info().Msg("loading config...")
	err := config.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	closer.InitGlobal()
	defer func() {
		timeoutCtx, cancel := context.WithTimeout(rootCtx, 3*time.Second)
		defer cancel()

		log.Info().Msg("shutting down services...")

		err = closer.Global.Close(timeoutCtx)
		if err != nil {
			log.Error().Err(err).
				Msg("took too long to stop all services ")
			return
		}

		log.Info().Msg("all services gracefully shut down")
	}()

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	if config.Verbose() {
		log.Logger = log.Level(zerolog.DebugLevel)
		log.Debug().Msg("enabled debug mode")

		verbose := middlewares.Verbose(log.Logger)
		router.Use(verbose)
	}

	compress := middleware.Compress(5)
	router.Use(middlewares.Decompress, compress)

	var repo repo
	if config.DatabaseDSN() == "" {
		log.Info().Msg("setting up inmemory repo...")
		inmemory, err := inmemory.New(config.FileStorage())
		if err != nil {
			log.Fatal().Err(err).Msg("failed to setup inmemory repo")
		}

		repo = inmemory

		closer.Global.Add(closer.Task{
			Sync: false,
			Fn: func(ctx context.Context) {
				log.Info().Msg("shutting down inmemory repo...")
				inmemory.Close()
				log.Info().Msg("inmemory repo shut down")
			},
		})

	} else {
		log.Info().Msg("setting up postgres conn...")
		pgxPool, err := pgxutils.ConnectPool(stopCtx, config.DatabaseDSN())
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect to postgres")
		}

		conn, err := pgxutils.AcquireConn(stopCtx, pgxPool)
		if err != nil {
			log.Fatal().Err(err).
				Msg("failed to acquire postgres conn")
		}

		postgres := postgres.New(conn)
		err = postgres.CreateTable(stopCtx)
		if err != nil {
			log.Fatal().Err(err).
				Msg("failed to create postgres tables")
		}

		repo = postgres

		pgxPinger := pinger.New(log.Logger, pgxPool)
		closer.Global.Add(closer.Task{
			Sync: false,
			Fn: func(ctx context.Context) {
				log.Info().Msg("shutting down postgres conn...")
				defer log.Info().Msg("postgres conn shut down")

				pgxPinger.Close(ctx)
				conn.Release()
				pgxPool.Close()
			},
		})

		go pgxPinger.Ping(stopCtx, 2*time.Second)

		pinger := services.NewPinger(pgxPinger)
		handlers.RegisterPing(router, pinger)
	}

	shortener := services.NewShortener(config.BaseURL(), repo)
	handlers.RegisterShortener(router, shortener)

	resolver := services.NewResolver(repo)
	handlers.RegisterResolve(router, resolver)

	server := http.Server{
		Addr:         config.ServerAddr(),
		Handler:      router,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	closer.Global.AddWithPriority(0, closer.Task{
		Sync: true,
		Fn: func(ctx context.Context) {
			log.Info().Msg("shutting down http server...")
			server.Shutdown(ctx)
			log.Info().Msg("http server shut down")
		},
	})

	log.Info().Str("addr", config.ServerAddr()).
		Msg("listening http server...")

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to listen server")
		}
	}()

	<-stopCtx.Done()
}
