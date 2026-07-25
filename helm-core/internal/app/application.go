package app

import (
	"log/slog"
	"os"

	"github.com/nayiswftw/helm/helm-core/internal/config"
	"github.com/nayiswftw/helm/helm-core/internal/integrations/system"
	"github.com/nayiswftw/helm/helm-core/internal/services"
)

// Application is the central dependency container for Helm Core.
type Application struct {
	Config    *config.Config
	Logger    *slog.Logger
	Devices   *services.DeviceService
	Dashboard *services.DashboardService
	System    *system.System
}

// New initializes and wires all application components.
func New(cfg *config.Config) *Application {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	sys := system.New(logger)
	devices := services.NewDeviceService(sys)

	return &Application{
		Config:    cfg,
		Logger:    logger,
		System:    sys,
		Devices:   devices,
		Dashboard: services.NewDashboardService(devices, sys),
	}
}
