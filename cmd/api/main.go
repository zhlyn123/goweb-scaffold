package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"goweb-scaffold/internal/app/bootstrap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app := bootstrap.NewApplication()
	if err := app.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "application stopped with error: %v\n", err)
		os.Exit(1)
	}
}
