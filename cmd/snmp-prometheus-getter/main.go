package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/sixworks/go-snmp-prometheus-getter/internal/config"
	"github.com/sixworks/go-snmp-prometheus-getter/internal/service"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config.toml", "Path to configuration file")
	flag.Parse()

	// Set up logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Load configuration
	cfg, err := config.LoadBootstrapConfiguration(*configPath)
	if err != nil {
		logger.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Update log level from configuration
	logLevel := &slog.LevelVar{}
	if err := logLevel.UnmarshalText([]byte(cfg.Instance.LogLevel)); err != nil {
		logger.Error("invalid log level", "level", cfg.Instance.LogLevel)
		os.Exit(1)
	}
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	// Create service
	svc, err := service.NewService(cfg, logger)
	if err != nil {
		logger.Error("failed to create service", "error", err)
		os.Exit(1)
	}

	// Set up signal handling
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start service
	if err := svc.Start(ctx); err != nil {
		logger.Error("service error", "error", err)
		os.Exit(1)
	}
}
