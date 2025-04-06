package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mtchuikov/shortener/internal/app"
)

func main() {
	rootCtx := context.Background()

	signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL}
	stopCtx, stop := signal.NotifyContext(rootCtx, signals...)
	defer stop()

	app := app.Setup(stopCtx)

	go app.Run(stopCtx)
	<-stopCtx.Done()

	app.Shutdown(rootCtx, 3*time.Second)
}

// func main() {
// 	rootCtx := context.Background()

// 	signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL}
// 	stopCtx, stop := signal.NotifyContext(rootCtx, signals...)
// 	defer stop()

// 	conf := config.New()
// 	logger := logging.NewZerolog(conf.ServiceName)

// 	router := chi.NewRouter()
// 	router.Use(middleware.Recoverer)

// 	if conf.Verbose {
// 		verbose := middlewares.Verbose(logger)
// 		router.Use(verbose)

// 		logger.Info().Msg("enabled debug mode")
// 	}

// 	compress := middleware.Compress(5)
// 	router.Use(middlewares.Decompress, compress)

// 	cache, err := inmemory.New(conf.FileStorage)
// 	if err != nil {
// 		logger.Info().Err(err).Msg("failed to setup cache")
// 	}

// 	shortenerService := services.NewShortener(conf.BaseURL, cache)
// 	handlers.RegisterShortener(logger, router, shortenerService)

// 	resolverService := services.NewResolver(conf.BaseURL, cache)
// 	handlers.RegisterResolver(logger, router, resolverService)

// 	server := http.Server{
// 		Addr:         conf.ServerAddr,
// 		Handler:      router,
// 		WriteTimeout: 3 * time.Second,
// 		ReadTimeout:  3 * time.Second,
// 	}

// 	go func() {
// 		err := server.ListenAndServe()
// 		if err != nil && err != http.ErrServerClosed {
// 			logger.Fatal().Err(err).Msg("failed to start server")
// 		}
// 	}()

// 	logger.Info().Msgf("server listening on addr %v...", conf.ServerAddr)
// 	<-stopCtx.Done()

// 	logger.Info().Msg("shutting down server...")
// 	shutdownCtx, shutdown := context.WithTimeout(rootCtx, 3*time.Second)
// 	defer shutdown()

// 	server.Shutdown(shutdownCtx)
// 	cache.Close(shutdownCtx)

// 	logger.Info().Msg("server shutdown")
// }
