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
	"github.com/rs/zerolog"
)

func setupHandlers(conf config.Config, lg zerolog.Logger, mux *chi.Mux) {
	cache := inmemory.New()

	shortenerService := services.NewShortener(conf.BaseURL, cache)
	handlers.RegisterShortener(lg, mux, shortenerService)

	resolverService := services.NewResolver(conf.BaseURL, cache)
	handlers.RegisterResolver(lg, mux, resolverService)
}

func newRouter(conf config.Config, lg zerolog.Logger) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	if conf.Verbose {
		verbose := middlewares.Verbose(lg)
		router.Use(verbose)
	}

	compress := middleware.Compress(5)
	router.Use(middlewares.Decompress, compress)

	setupHandlers(conf, lg, router)
	return router
}

func newServer(conf config.Config, lg zerolog.Logger) http.Server {
	return http.Server{
		Addr:         conf.ServerAddr,
		Handler:      newRouter(conf, lg),
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}
}

func main() {
	rootCtx := context.Background()
	stopCtx, stop := signal.NotifyContext(rootCtx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conf := config.New()
	logger := logtools.NewZerolog(conf.ServiceName, os.Stdout)

	server := newServer(conf, logger)
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
	logger.Info().Msg("server shutdown")
}
