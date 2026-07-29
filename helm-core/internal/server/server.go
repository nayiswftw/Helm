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
	tlsCert    string
	tlsKey     string
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

// SetTLS configures optional TLS certificate and key file paths.
// If both are non-empty, the server will use ListenAndServeTLS.
func (s *Server) SetTLS(certFile, keyFile string) {
	s.tlsCert = certFile
	s.tlsKey = keyFile
}

// Start begins listening and serving HTTP (or HTTPS) requests.
// It blocks until the server stops.
func (s *Server) Start() error {
	if s.tlsCert != "" && s.tlsKey != "" {
		s.logger.Info("server starting with TLS", "addr", s.httpServer.Addr)
		return s.httpServer.ListenAndServeTLS(s.tlsCert, s.tlsKey)
	}

	s.logger.Info("server starting", "addr", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("server shutting down")
	return s.httpServer.Shutdown(ctx)
}

