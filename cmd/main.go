package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mtchuikov/shortener/internal/config"
	"github.com/mtchuikov/shortener/internal/handlers"
	"github.com/mtchuikov/shortener/internal/repo"
	"github.com/mtchuikov/shortener/internal/repo/inmemory"
	"github.com/mtchuikov/shortener/internal/repo/postgres"
	"github.com/mtchuikov/shortener/internal/services"
	"github.com/mtchuikov/shortener/pkg/closer"
	"github.com/mtchuikov/shortener/pkg/filewriter"
	"github.com/mtchuikov/shortener/pkg/logging"
	"github.com/mtchuikov/shortener/pkg/middlewares"
	"github.com/mtchuikov/shortener/pkg/pgxutils"
	"github.com/mtchuikov/shortener/pkg/pinger"
	"github.com/mtchuikov/shortener/pkg/strgen"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	rootCtx       = context.Background()
	signals       = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL}
	stopCtx, stop = signal.NotifyContext(rootCtx, signals...)
)

func init() {
	zerolog.LevelFieldName = "lvl"
	zerolog.ErrorFieldName = "err"
	zerolog.MessageFieldName = "msg"
	zerolog.TimeFieldFormat = time.RFC1123

	closer.InitGlobal()
	strgen.InitGlobal()
}

func main() {
	defer func() {
		stop()
		timeoutCtx, cancel := context.WithTimeout(rootCtx, 5*time.Second)

		closer.Global.Close(timeoutCtx)
		cancel()
	}()

	log.Logger = log.Logger.
		Level(zerolog.InfoLevel).With().
		Timestamp().Str("app", "shortener").
		Logger()

	err := config.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init config")
	}

	if config.LogToFile() {
		fw, err := filewriter.New(config.LogFile())
		if err != nil {
			log.Fatal().Err(err).Msg("failed to configure logger")
		}

		lvl, err := logging.NewLevel(config.LogFileLevel())
		if err != nil {
			log.Fatal().Err(err).Msg("failed to configure logger")
		}

		zlvl, _ := lvl.ToZerolog()
		zfw := filewriter.NewZerologdWriter(fw, zlvl)

		output := zerolog.MultiLevelWriter(os.Stdout, zfw)
		log.Logger = log.Output(output)

		closer.Global.VeryLast = func(context.Context) { zfw.Close() }
	}

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middlewares.Decompress, middleware.Compress(5))

	if config.Verbose() {
		log.Debug().Msg("enabled debug mode")
		log.Logger = log.Level(zerolog.DebugLevel)
		router.Use(middlewares.Verbose(&log.Logger))
	}

	var shortenURLsRepo repo.IShortenURLs
	if config.DatabaseDSN() == "" {
		inmemory, err := inmemory.New(config.FileStorage())
		if err != nil {
			log.Fatal().Err(err).Msg("failed to setup inmemory storage")
		}

		closer.Global.Add(closer.Task{
			Fn: func(ctx context.Context) {
				inmemory.Close()
			}})

		shortenURLsRepo = inmemory

	} else {
		pgConn, err := pgxutils.Connect(stopCtx, config.DatabaseDSN())
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect to postgres")
		}

		pgPinger := pinger.New(&log.Logger, pgConn)
		go pgPinger.Ping(stopCtx, 2*time.Second)

		closer.Global.Add(closer.Task{
			Fn: func(ctx context.Context) {
				pgPinger.Close(ctx)
				pgConn.Close(ctx)
			},
		})

		shortenURLsRepo = postgres.NewShortenURLs(pgConn.Conn)

		pinger := services.NewPinger(pgPinger)
		handlers.RegisterPing(&log.Logger, router, pinger)
	}

	resolver := services.NewResolver(&log.Logger, shortenURLsRepo)
	handlers.RegisterResolver(&log.Logger, router, resolver)

	shortener := services.NewShortener(&log.Logger, config.BaseURL(), shortenURLsRepo)
	handlers.RegisterShorten(&log.Logger, router, shortener)

	server := http.Server{
		Addr:         config.ServerAddr(),
		Handler:      router,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Info().Msgf("listening server on %s...", config.ServerAddr())
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to listen server")
		}
	}()

	<-stopCtx.Done()
}
