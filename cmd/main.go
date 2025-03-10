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
	"github.com/mtchuikov/shortener/internal/repo/inmemory"
	"github.com/mtchuikov/shortener/internal/service"
	"github.com/mtchuikov/shortener/pkg/middlewares"
	"github.com/rs/zerolog"
)

func newLogger(config config.Config) zerolog.Logger {
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "msg"
	zerolog.ErrorFieldName = "err"
	zerolog.TimeFieldFormat = time.RFC1123

	return zerolog.New(os.Stdout).With().
		Timestamp().Str("app", config.ServiceName).
		Logger()
}

func newHandler(config config.Config, logger zerolog.Logger) *handler.Handler {
	repo := inmemory.New()
	service := service.New(config.BaseURL, repo)
	return handler.New(service)
}

func newRouter(config config.Config, logger zerolog.Logger) http.Handler {
	handler := newHandler(config, logger)
	mux := http.NewServeMux()

	mux.HandleFunc("/", middlewares.OnlyMethod(http.MethodPost,
		handler.CreateShortURL))
	mux.HandleFunc("/{id}", middlewares.OnlyMethod(http.MethodGet,
		handler.ResolveShortURL))

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

	config := config.New()
	logger := newLogger(config)

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
