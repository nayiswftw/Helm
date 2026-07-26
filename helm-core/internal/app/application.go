//go:build linux

package app

import (
	"log/slog"

	"github.com/nayiswftw/helm/helm-core/internal/config"
	"github.com/nayiswftw/helm/helm-core/internal/integration/docker"
	"github.com/nayiswftw/helm/helm-core/internal/integration/system"
	"github.com/nayiswftw/helm/helm-core/internal/service"
)

// Application is the dependency injection container.
// It wires together all components and owns their lifecycle.
type Application struct {
	Config     config.Config
	Logger     *slog.Logger
	System     *system.System
	Docker     *docker.Client
	Dashboard  *service.DashboardService
	Devices    *service.DeviceService
	Actions    *service.ActionService
	Containers *service.ContainerService
}

// New creates a fully wired Application.
func New(cfg config.Config, logger *slog.Logger) *Application {
	sys := system.New()
	dockerClient := docker.NewClient("")

	dashboard := service.NewDashboardService(sys)
	devices := service.NewDeviceService(sys, dockerClient)
	containers := service.NewContainerService(dockerClient)

	// Local device ID matches local registered device ID
	localDevices := devices.GetAll()
	localID := "local-server"
	if len(localDevices) > 0 {
		localID = localDevices[0].ID
	}

	actions := service.NewActionService(sys, localID, logger)

	return &Application{
		Config:     cfg,
		Logger:     logger,
		System:     sys,
		Docker:     dockerClient,
		Dashboard:  dashboard,
		Devices:    devices,
		Actions:    actions,
		Containers: containers,
	}
}
