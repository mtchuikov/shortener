package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mtchuikov/shortener/internal/config"
	"github.com/mtchuikov/shortener/internal/handler"
	"github.com/mtchuikov/shortener/internal/service"
	storage "github.com/mtchuikov/shortener/internal/storage/inmemory"
	"github.com/mtchuikov/shortener/pkg/middlewares"
	"github.com/rs/zerolog"
)

func main() {
	rootCtx := context.Background()
	stopCtx, stop := signal.NotifyContext(rootCtx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	zerolog.LevelFieldName = "lvl"
	zerolog.MessageFieldName = "msg"
	zerolog.TimeFieldFormat = time.RFC1123

	logger := zerolog.New(os.Stdout).With().
		Timestamp().Str("app", "shorter").
		Logger()

	config := config.New()

	storage := storage.New()
	service := service.New(logger, storage)
	handler := handler.New(service)

	mux := http.NewServeMux()
	mux.HandleFunc("/create", middlewares.OnlyMethod(http.MethodPost,
		handler.CreateShortID))
	mux.HandleFunc("/{short_id}", middlewares.OnlyMethod(http.MethodGet,
		handler.ResolveShortID))

	httpServer := http.Server{
		Addr:         config.Addr,
		Handler:      mux,
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal().Msgf("failed to start http server %v", err)
		}
	}()

	logger.Info().Msgf("http server listening on %s...", config.Addr)
	<-stopCtx.Done()

	logger.Info().Msg("gracefully shutting down server...")
	shutdownCtx, shutdown := context.WithTimeout(rootCtx, 5*time.Second)
	defer shutdown()

	httpServer.Shutdown(shutdownCtx)
	logger.Info().Msg("server shutdown gracefully")
}
