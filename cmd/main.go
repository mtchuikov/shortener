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
	"github.com/mtchuikov/shortener/internal/cache/inmemory"
	"github.com/mtchuikov/shortener/internal/config"
	"github.com/mtchuikov/shortener/internal/handlers"
	"github.com/mtchuikov/shortener/internal/services"
	"github.com/mtchuikov/shortener/pkg/logtools"
	"github.com/mtchuikov/shortener/pkg/middlewares"
)

func main() {
	rootCtx := context.Background()

	signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL}
	stopCtx, stop := signal.NotifyContext(rootCtx, signals...)
	defer stop()

	conf := config.New()
	logger := logtools.NewZerolog(conf.ServiceName, os.Stdout)

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	if conf.Verbose {
		verbose := middlewares.Verbose(logger)
		router.Use(verbose)

		logger.Info().Msg("enabled debug mode")
	}

	compress := middleware.Compress(5)
	router.Use(middlewares.Decompress, compress)

	cache, err := inmemory.New(conf.FileStorage)
	if err != nil {
		logger.Info().Err(err).Msg("failed to setup cache")
	}

	shortenerService := services.NewShortener(conf.BaseURL, cache)
	handlers.RegisterShortener(logger, router, shortenerService)

	resolverService := services.NewResolver(conf.BaseURL, cache)
	handlers.RegisterResolver(logger, router, resolverService)

	server := http.Server{
		Addr:         conf.ServerAddr,
		Handler:      router,
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	logger.Info().Msgf("server listening on addr %v...", conf.ServerAddr)
	<-stopCtx.Done()

	logger.Info().Msg("shutting down server...")
	shutdownCtx, shutdown := context.WithTimeout(rootCtx, 3*time.Second)
	defer shutdown()

	server.Shutdown(shutdownCtx)
	cache.Close(shutdownCtx)
	logger.Info().Msg("server shutdown")
}
