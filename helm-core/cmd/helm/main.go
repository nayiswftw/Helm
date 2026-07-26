//go:build linux

package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nayiswftw/helm/helm-core/internal/api"
	"github.com/nayiswftw/helm/helm-core/internal/app"
	"github.com/nayiswftw/helm/helm-core/internal/config"
	"github.com/nayiswftw/helm/helm-core/internal/server"
)

func main() {
	// Load configuration from environment.
	cfg := config.Load()

	// Initialize structured logger.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(cfg.LogLevel),
	}))

	logger.Info("helm starting",
		"port", cfg.Port,
		"log_level", cfg.LogLevel,
	)

	// Build the application container.
	application := app.New(cfg, logger)

	// Set up HTTP router.
	router := api.NewRouter(application)

	// Create and start the HTTP server.
	srv := server.New(cfg.Port, router, logger)

	// Run server in a goroutine so we can listen for shutdown signals.
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logger.Info("shutdown signal received", "signal", sig.String())

	// Give outstanding requests 10 seconds to complete.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("helm stopped")
}

// parseLogLevel converts a string log level to slog.Level.
func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
