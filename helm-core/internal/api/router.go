//go:build linux

package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/nayiswftw/helm/helm-core/internal/app"
)

// NewRouter creates and configures the Chi router with all routes.
func NewRouter(application *app.Application) chi.Router {
	r := chi.NewRouter()

	// Middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(requestLogger(application.Logger))

	// Health check — outside versioned API
	r.Get("/health", handleHealth())

	// Versioned API
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/dashboard", handleDashboard(application.Dashboard, application.Logger))
	})

	return r
}
