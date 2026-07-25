package services

import (
	"github.com/nayiswftw/helm/helm-core/internal/integrations/system"
)

type Dashboard struct {
	Hostname string             `json:"hostname"`
	Uptime   float64            `json:"uptime"`
	CPU      system.CPUInfo     `json:"cpu"`
	Memory   system.MemoryInfo  `json:"memory"`
	Disk     system.DiskInfo    `json:"disk"`
	Network  system.NetworkInfo `json:"network"`
	Devices  int                `json:"devices"`
}

type DashboardService struct {
	devices *DeviceService
	system  *system.System
}

func NewDashboardService(ds *DeviceService, sys *system.System) *DashboardService {
	return &DashboardService{
		devices: ds,
		system:  sys,
	}
}

func (s *DashboardService) Get() Dashboard {
	return Dashboard{
		Hostname: s.system.Hostname(),
		Uptime:   s.system.Uptime(),
		CPU:      s.system.CPU(),
		Memory:   s.system.Memory(),
		Disk:     s.system.Disk(),
		Network:  s.system.Network(),
		Devices:  len(s.devices.List()),
	}
}
