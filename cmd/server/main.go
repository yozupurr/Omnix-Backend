package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"omnix/internal/config"
	"omnix/internal/server"
)

func main() {
	cfg := config.Load()
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	srv, err := server.New(cfg, log)
	if err != nil {
		log.Error("init failed", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := srv.Run(ctx); err != nil {
		log.Error("server failed", "err", err)
		os.Exit(1)
	}
}
