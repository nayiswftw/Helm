package system

import (
	"log/slog"
	"os"
)

// System interacts with host system resources and operating system interfaces.
type System struct {
	hostname string
	logger   *slog.Logger
}

// New creates a System integration instance.
func New(logger *slog.Logger) *System {
	host, err := os.Hostname()
	if err != nil {
		logger.Warn("failed to resolve hostname, using fallback", "error", err)
		host = "unknown"
	}

	return &System{
		hostname: host,
		logger:   logger,
	}
}

// Hostname returns the system hostname.
func (s *System) Hostname() string {
	return s.hostname
}
