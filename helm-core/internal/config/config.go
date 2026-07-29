package config

import (
	"os"
)

// Config holds the application configuration.
// All values are sourced from environment variables with sensible defaults.
type Config struct {
	Port     string // HELM_PORT — address to listen on (default ":8080")
	LogLevel string // HELM_LOG_LEVEL — slog level: debug, info, warn, error (default "info")

	// Dokploy integration (optional — empty values disable the feature)
	DokployURL    string // HELM_DOKPLOY_URL — Dokploy instance base URL
	DokployAPIKey string // HELM_DOKPLOY_API_KEY — Dokploy API key

	// Authentication (optional — empty HELM_API_KEY disables auth / dev mode)
	APIKey  string // HELM_API_KEY — API key for Helm authentication
	TLSCert string // HELM_TLS_CERT — path to TLS certificate file
	TLSKey  string // HELM_TLS_KEY — path to TLS private key file

	// Notifications (optional — empty HELM_NOTIFY_URL disables notifications)
	NotifyURL string // HELM_NOTIFY_URL — webhook URL for notifications
}

// Load reads configuration from environment variables.
// Missing variables fall back to defaults.
func Load() Config {
	return Config{
		Port:          envOrDefault("HELM_PORT", ":8080"),
		LogLevel:      envOrDefault("HELM_LOG_LEVEL", "info"),
		DokployURL:    envOrDefault("HELM_DOKPLOY_URL", ""),
		DokployAPIKey: envOrDefault("HELM_DOKPLOY_API_KEY", ""),
		APIKey:        envOrDefault("HELM_API_KEY", ""),
		TLSCert:       envOrDefault("HELM_TLS_CERT", ""),
		TLSKey:        envOrDefault("HELM_TLS_KEY", ""),
		NotifyURL:     envOrDefault("HELM_NOTIFY_URL", ""),
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
