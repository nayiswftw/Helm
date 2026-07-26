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

		r.Get("/devices", handleListDevices(application.Devices))
		r.Get("/devices/{id}", handleGetDevice(application.Devices))

		r.Get("/actions", handleListActions(application.Actions))
		r.Post("/actions/{id}/execute", handleExecuteAction(application.Actions))

		r.Get("/containers", handleListContainers(application.Containers))
		r.Post("/containers/{id}/start", handleStartContainer(application.Containers))
		r.Post("/containers/{id}/stop", handleStopContainer(application.Containers))
		r.Post("/containers/{id}/restart", handleRestartContainer(application.Containers))
	})

	return r
}


