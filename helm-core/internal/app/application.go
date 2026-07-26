//go:build linux

package app

import (
	"log/slog"

	"github.com/nayiswftw/helm/helm-core/internal/config"
	"github.com/nayiswftw/helm/helm-core/internal/integration/system"
	"github.com/nayiswftw/helm/helm-core/internal/service"
)

// Application is the dependency injection container.
// It wires together all components and owns their lifecycle.
type Application struct {
	Config    config.Config
	Logger    *slog.Logger
	System    *system.System
	Dashboard *service.DashboardService
}

// New creates a fully wired Application.
func New(cfg config.Config, logger *slog.Logger) *Application {
	sys := system.New()
	dashboard := service.NewDashboardService(sys)

	return &Application{
		Config:    cfg,
		Logger:    logger,
		System:    sys,
		Dashboard: dashboard,
	}
}
