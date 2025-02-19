package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/matthew-hollick/go-snmp-prometheus-getter/internal/config"
	"github.com/matthew-hollick/go-snmp-prometheus-getter/internal/service"
)

func main() {
	// Parse command line flags
	configFile := flag.String("config", "config.toml", "Path to configuration file")
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	flag.Parse()

	// Set up logging
	var logger *slog.Logger
	switch *logLevel {
	case "debug":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "warn":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
	case "error":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	default:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	// Load configuration
	cfg, err := config.LoadBootstrapConfiguration(*configFile)
	if err != nil {
		logger.Error("loading configuration", "error", err)
		os.Exit(1)
	}

	// Create service
	svc, err := service.NewService(cfg, logger)
	if err != nil {
		logger.Error("creating service", "error", err)
		os.Exit(1)
	}

	// Handle interrupts
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Info("received signal", "signal", sig)
		cancel()
	}()

	// Start service
	if err := svc.Start(ctx); err != nil {
		logger.Error("starting service", "error", err)
		os.Exit(1)
	}
}
