//go:build linux

package app

import (
	"log/slog"

	"github.com/nayiswftw/helm/helm-core/internal/config"
	"github.com/nayiswftw/helm/helm-core/internal/integration/docker"
	"github.com/nayiswftw/helm/helm-core/internal/integration/dokploy"
	"github.com/nayiswftw/helm/helm-core/internal/integration/notify"
	"github.com/nayiswftw/helm/helm-core/internal/integration/system"
	"github.com/nayiswftw/helm/helm-core/internal/service"
)

// Application is the dependency injection container.
// It wires together all components and owns their lifecycle.
type Application struct {
	Config        config.Config
	Logger        *slog.Logger
	System        *system.System
	Docker        *docker.Client
	Dashboard     *service.DashboardService
	Devices       *service.DeviceService
	Actions       *service.ActionService
	Containers    *service.ContainerService
	Dokploy       *service.DokployService
	Notifications *service.NotificationService
}

// New creates a fully wired Application.
func New(cfg config.Config, logger *slog.Logger) *Application {
	// Integration layer
	sys := system.New()
	dockerClient := docker.NewClient("")
	dokployClient := dokploy.NewClient(cfg.DokployURL, cfg.DokployAPIKey)
	webhookNotifier := notify.NewWebhookNotifier(cfg.NotifyURL)

	// Service layer — notifications first (other services depend on it)
	notifications := service.NewNotificationService(webhookNotifier, logger)

	dashboard := service.NewDashboardService(sys)
	devices := service.NewDeviceService(sys, dockerClient)
	containers := service.NewContainerService(dockerClient)
	dokploySvc := service.NewDokployService(dokployClient, notifications)

	// Local device ID matches local registered device ID
	localDevices := devices.GetAll()
	localID := "local-server"
	if len(localDevices) > 0 {
		localID = localDevices[0].ID
	}

	actions := service.NewActionService(sys, localID, logger, notifications)

	return &Application{
		Config:        cfg,
		Logger:        logger,
		System:        sys,
		Docker:        dockerClient,
		Dashboard:     dashboard,
		Devices:       devices,
		Actions:       actions,
		Containers:    containers,
		Dokploy:       dokploySvc,
		Notifications: notifications,
	}
}

