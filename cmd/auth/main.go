package main

import (
	"auth/internal/app"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	/* Start and Stop app with gracefully shutdown  */
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM|syscall.SIGINT, os.Interrupt)
	defer cancel()

	go app.Start(ctx)
	<-ctx.Done()
	app.Stop()
}
