package system

import (
	"log/slog"
	"os"
)

type System struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *System {
	return &System{
		logger: logger,
	}
}

func (s *System) Hostname() string {
	host, err := os.Hostname()
	if err != nil || host == "" {
		s.logger.Warn("failed to resolve dynamic hostname", "error", err)
		return "helm-server"
	}
	return host
}
