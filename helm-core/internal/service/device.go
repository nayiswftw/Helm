//go:build linux

package service

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/nayiswftw/helm/helm-core/internal/domain"
	"github.com/nayiswftw/helm/helm-core/internal/integration/docker"
	"github.com/nayiswftw/helm/helm-core/internal/integration/system"
)

var ErrDeviceNotFound = errors.New("device not found")

// DeviceService manages device registry and capability detection.
type DeviceService struct {
	mu      sync.RWMutex
	devices map[string]domain.Device
}

// NewDeviceService initializes the service and auto-registers the local server.
func NewDeviceService(sys *system.System, dockerClient *docker.Client) *DeviceService {
	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		hostname = "local-server"
	}

	capabilities := []domain.Capability{
		domain.CapabilityMetrics,
	}

	// Probe if systemctl is available for power control
	if _, err := exec.LookPath("systemctl"); err == nil {
		capabilities = append(capabilities, domain.CapabilityPowerControl)
	} else if _, err := os.Stat("/usr/bin/systemctl"); err == nil {
		capabilities = append(capabilities, domain.CapabilityPowerControl)
	}

	// Probe if Docker is available
	if dockerClient != nil && dockerClient.IsAvailable() {
		capabilities = append(capabilities, domain.CapabilityContainers)
	}

	localDevice := domain.Device{
		ID:           hostname,
		Hostname:     hostname,
		Platform:     runtime.GOOS,
		Architecture: runtime.GOARCH,
		Capabilities: capabilities,
		Status:       "online",
	}

	ds := &DeviceService{
		devices: make(map[string]domain.Device),
	}
	ds.devices[localDevice.ID] = localDevice

	return ds
}

// GetAll returns a slice of all registered devices.
func (s *DeviceService) GetAll() []domain.Device {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]domain.Device, 0, len(s.devices))
	for _, dev := range s.devices {
		result = append(result, dev)
	}
	return result
}

// GetByID retrieves a single device by ID.
func (s *DeviceService) GetByID(id string) (domain.Device, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	dev, ok := s.devices[id]
	if !ok {
		return domain.Device{}, ErrDeviceNotFound
	}
	return dev, nil
}
