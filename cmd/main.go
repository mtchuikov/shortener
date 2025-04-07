package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	app "github.com/mtchuikov/shortener/internal"
	"github.com/mtchuikov/shortener/internal/config"
	"github.com/mtchuikov/shortener/pkg/closer"
	"github.com/rs/zerolog"
)

func init() {
	zerolog.LevelFieldName = "lvl"
	zerolog.ErrorFieldName = "err"
	zerolog.MessageFieldName = "msg"
	zerolog.TimeFieldFormat = time.RFC1123
}

func main() {
	rootCtx := context.Background()

	signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL}
	stopCtx, stop := signal.NotifyContext(rootCtx, signals...)
	defer stop()

	log := zerolog.New(os.Stdout).
		Level(zerolog.InfoLevel).With().
		Timestamp().Str("app", config.ServiceName()).
		Logger()

	err := config.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	if config.Verbose() {
		log = log.Level(zerolog.DebugLevel)
		log.Debug().Msg("enabled debug mode")
	}

	closer.Init()

	app := app.New(stopCtx, log)

	go app.Run(stopCtx)
	<-stopCtx.Done()

	timeoutCtx, cancel := context.WithTimeout(rootCtx, 5*time.Second)
	defer cancel()

	app.Shutdown(timeoutCtx)
}
