package config

import (
	"os"
)

// Config holds the application configuration.
// All values are sourced from environment variables with sensible defaults.
type Config struct {
	Port     string // HELM_PORT — address to listen on (default ":8080")
	LogLevel string // HELM_LOG_LEVEL — slog level: debug, info, warn, error (default "info")
}

// Load reads configuration from environment variables.
// Missing variables fall back to defaults.
func Load() Config {
	return Config{
		Port:     envOrDefault("HELM_PORT", ":8080"),
		LogLevel: envOrDefault("HELM_LOG_LEVEL", "info"),
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
