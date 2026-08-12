package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/corvian/argus/services/weather/internal/config"
)

func main() {
	slog.Info("Starting the weather worker")
	config := config.NewConfig()

	// Setup Graceful Shutdown Context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := NewDaemon(ctx, config)
	if err != nil {
		slog.Error("failed to initialize worker", "error", err)
		panic("Failed to initialize the worker")
	}
}
