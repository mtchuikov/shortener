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
