package services

import (
	"github.com/nayiswftw/helm/helm-core/internal/domain"
	"github.com/nayiswftw/helm/helm-core/internal/integrations/system"
)

type DeviceService struct {
	system *system.System
}

func NewDeviceService(sys *system.System) *DeviceService {
	return &DeviceService{
		system: sys,
	}
}

func (s *DeviceService) List() []domain.Device {
	return []domain.Device{
		{
			ID:     "local",
			Name:   s.system.Hostname(),
			Type:   "server",
			Status: "online",
		},
	}
}

func (s *DeviceService) GetByID(id string) (domain.Device, bool) {
	for _, dev := range s.List() {
		if dev.ID == id {
			return dev, true
		}
	}
	return domain.Device{}, false
}
