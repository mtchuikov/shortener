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
	"github.com/mtchuikov/shortener/internal/config"
	"github.com/mtchuikov/shortener/internal/handler"
	"github.com/mtchuikov/shortener/internal/service"
	memstorage "github.com/mtchuikov/shortener/internal/storage/mem"
	chimiddlewares "github.com/mtchuikov/shortener/pkg/middlewares/chi"
	"github.com/rs/zerolog"
)

func newLogger() zerolog.Logger {
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "msg"
	zerolog.TimeFieldFormat = time.RFC1123

	return zerolog.New(os.Stdout).With().
		Timestamp().Str("app", "shorter").
		Logger()
}

func newHandler(config config.Config, logger zerolog.Logger) *handler.Handler {
	storage := memstorage.New()
	service := service.New(logger, config.BaseURL, storage)
	return handler.New(service)
}

func newRouter(config config.Config, logger zerolog.Logger) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)

	if config.Verbose {
		mux.Use(middleware.RequestID)
		mux.Use(chimiddlewares.Verbose(logger))
	}

	handler := newHandler(config, logger)

	mux.Get("/{short_id}", handler.ResolveShortID)
	mux.Post("/", handler.ShortURL)

	return mux
}

func newServer(config config.Config, logger zerolog.Logger) *http.Server {
	mux := newRouter(config, logger)
	return &http.Server{
		Addr:         config.Addr,
		Handler:      mux,
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}
}

func main() {
	rootCtx := context.Background()
	stopCtx, stop := signal.NotifyContext(rootCtx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := newLogger()
	config := config.New()

	server := newServer(config, logger)
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal().Msgf("failed to start server %v", err)
		}
	}()

	logger.Info().Msgf("server listening on %s...", config.Addr)
	<-stopCtx.Done()

	logger.Info().Msg("shutting down server...")
	shutdownCtx, shutdown := context.WithTimeout(rootCtx, 3*time.Second)
	defer shutdown()

	server.Shutdown(shutdownCtx)
	logger.Info().Msg("server shutdown")
}
