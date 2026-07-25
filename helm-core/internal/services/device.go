package services

import (
	"github.com/nayiswftw/helm/helm-core/internal/domain"
	"github.com/nayiswftw/helm/helm-core/internal/integrations/system"
)

// DeviceService manages device state and queries.
type DeviceService struct {
	devices []domain.Device
}

// NewDeviceService creates a DeviceService with the local server registered.
func NewDeviceService(sys *system.System) *DeviceService {
	return &DeviceService{
		devices: []domain.Device{
			{
				ID:     "local",
				Name:   sys.Hostname(),
				Type:   "server",
				Status: "online",
			},
		},
	}
}

// List returns all registered devices.
func (s *DeviceService) List() []domain.Device {
	return s.devices
}

// GetByID looks up a device by its ID.
func (s *DeviceService) GetByID(id string) (domain.Device, bool) {
	for _, d := range s.devices {
		if d.ID == id {
			return d, true
		}
	}
	return domain.Device{}, false
}
