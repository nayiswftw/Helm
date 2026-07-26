package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

// Server wraps http.Server with sensible defaults and graceful shutdown.
type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

// New creates a new Server with production-appropriate timeouts.
func New(addr string, handler http.Handler, logger *slog.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		logger: logger,
	}
}

// Start begins listening and serving HTTP requests.
// It blocks until the server stops.
func (s *Server) Start() error {
	s.logger.Info("server starting", "addr", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("server shutting down")
	return s.httpServer.Shutdown(ctx)
}
