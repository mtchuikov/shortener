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
	"github.com/mtchuikov/shortener/internal/middlewares"
	"github.com/mtchuikov/shortener/internal/services"
	"github.com/mtchuikov/shortener/internal/storage"
	"github.com/mtchuikov/shortener/internal/storage/inmemory"
	"github.com/mtchuikov/shortener/internal/storage/postgres"
	"github.com/mtchuikov/shortener/pkg/closer"
	"github.com/mtchuikov/shortener/pkg/filewriter"
	"github.com/mtchuikov/shortener/pkg/logutils"
	"github.com/mtchuikov/shortener/pkg/middlewares/decompress"
	"github.com/mtchuikov/shortener/pkg/middlewares/jwtauth"
	"github.com/mtchuikov/shortener/pkg/middlewares/verbose"
	"github.com/mtchuikov/shortener/pkg/pgxutils"
	"github.com/mtchuikov/shortener/pkg/pinger"
	"github.com/mtchuikov/shortener/pkg/strgen"
)

func init() {
	closer.InitGlobal()
	strgen.InitGlobal()
}

func main() {
	rootCtx := context.Background()
	signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL}

	stopCtx, stop := signal.NotifyContext(rootCtx, signals...)
	defer stop()

	defer func() {
		timeoutCtx, cancel := context.WithTimeout(rootCtx, 5*time.Second)
		closer.Global.Close(timeoutCtx)
		cancel()
	}()

	// ---------------------------------------------------------------- //

	zerolog.LevelFieldName = "lvl"
	zerolog.ErrorFieldName = "err"
	zerolog.MessageFieldName = "msg"
	zerolog.TimeFieldFormat = time.RFC1123

	log.Logger = log.Logger.
		Level(zerolog.InfoLevel).With().
		Timestamp().Str("app", "shortener").
		Logger()

	// ---------------------------------------------------------------- //

	err := config.Init()
	if err != nil {
		log.Fatal().Err(err).
			Msg("unable to init config")
	}

	// ---------------------------------------------------------------- //

	if config.LogToFile() {
		fw, err := filewriter.New(config.LogFileLevel())
		if err != nil {
			log.Fatal().Err(err).
				Msg("unable to configure file logger")
		}

		level, err := logutils.NewLevel(config.LogFileLevel())
		if err != nil {
			log.Fatal().Err(err).
				Msg("unable to configure file logger")
		}

		zlevel, _ := level.ToZerolog()
		zfw := filewriter.NewZerologdWriter(fw, zlevel)

		closeLogger := func(ctx context.Context) { zfw.Close() }
		closer.Global.VeryLast = closeLogger

		output := zerolog.MultiLevelWriter(os.Stdout, zfw)
		log.Logger = log.Output(output)
	}

	// ---------------------------------------------------------------- //

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(decompress.Decompress, middleware.Compress(5))

	if config.Verbose() {
		log.Debug().Msg("enabled debug mode")
		log.Logger = log.Level(zerolog.DebugLevel)
		router.Use(verbose.Verbose(&log.Logger))
	}

	ja := jwtauth.New(config.JWTSecret())
	router.Use(middlewares.AutoAuthenticate(ja))

	// ---------------------------------------------------------------- //

	var shortURLsStorage storage.ShortURLsStorage

	if config.DatabaseDSN() != "" {
		pgxPool, err := pgxutils.ConnectPool(stopCtx, config.DatabaseDSN())
		if err != nil {
			log.Fatal().Err(err).
				Msg("unable to connect to database")
		}

		shortURLsStorage = postgres.NewShortURLs(pgxPool, config.BaseURL())

		pgxPinger := pinger.New(pgxPool)
		pgxPinger.Ping(stopCtx, 2*time.Second)

		closer.Global.Add(false,
			func(ctx context.Context) {
				log.Info().Msg("closing database connection...")
				pgxPinger.Close(ctx)
				pgxPool.Close()
				log.Info().Msg("database connection closed")
			})

		pinger := services.NewPinger(pgxPinger)
		handlers.RegisterPing(router, pinger)

	} else {
		inmemory, err := inmemory.NewShortURLs(config.FileStorage(), config.BaseURL())
		if err != nil {
			log.Fatal().Err(err).
				Msg("unable to init inmemory storage")
		}

		shortURLsStorage = inmemory

		closer.Global.Add(false,
			func(ctx context.Context) {
				log.Info().Msg("cleaning up inmemory storage...")
				inmemory.Close()
				log.Info().Msg("inmemory storage cleaned up")
			})
	}

	shortener := services.NewShortener(&log.Logger, shortURLsStorage)

	handlers.RegisterShorten(router, shortener)
	handlers.RegisterResolve(router, shortener)
	handlers.RegisterListURLs(router, shortener)
	handlers.RegisterMarkURLsAsDeleted(router, shortener)

	// ---------------------------------------------------------------- //

	server := http.Server{
		Addr:         config.ServerAddr(),
		Handler:      router,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	closer.Global.VeryFirst = func(ctx context.Context) {
		log.Info().Msg("shutting down server...")
		server.Shutdown(ctx)
		log.Info().Msg("server shut down")
	}

	log.Info().Str("addr", config.ServerAddr()).
		Msg("listening server...")

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to listen server")
		}
	}()

	// ---------------------------------------------------------------- //

	<-stopCtx.Done()
}
