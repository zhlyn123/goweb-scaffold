package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"goweb-scaffold/internal/app/bootstrap"
	"goweb-scaffold/internal/platform/config"
)

func main() {
	configPath := flag.String("config", "configs/config.local.yaml", "配置文件路径")
	flag.Parse()

	config, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loading config failed: %v\n", err)
		os.Exit(1)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := bootstrap.NewApplication(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create application failed: %v\n", err)
		os.Exit(1)
	}

	if err := app.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "application stopped with error: %v\n", err)
		os.Exit(1)
	}
}
