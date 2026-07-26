//go:build linux

package service

import (
	"fmt"
	"math"
	"os"

	"github.com/nayiswftw/helm/helm-core/internal/domain"
	"github.com/nayiswftw/helm/helm-core/internal/integration/system"
)

// DashboardService orchestrates system metric collection
// for the dashboard API endpoint.
type DashboardService struct {
	system *system.System
}

// NewDashboardService creates a new DashboardService.
func NewDashboardService(sys *system.System) *DashboardService {
	return &DashboardService{
		system: sys,
	}
}

// GetMetrics collects current system metrics and returns them
// as a single SystemMetrics snapshot.
func (s *DashboardService) GetMetrics() (domain.SystemMetrics, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	cpu, err := s.system.CPUUsage()
	if err != nil {
		return domain.SystemMetrics{}, fmt.Errorf("collecting cpu: %w", err)
	}

	memory, err := s.system.MemoryUsage()
	if err != nil {
		return domain.SystemMetrics{}, fmt.Errorf("collecting memory: %w", err)
	}

	disk, err := s.system.DiskUsage()
	if err != nil {
		return domain.SystemMetrics{}, fmt.Errorf("collecting disk: %w", err)
	}

	uptime, err := s.system.Uptime()
	if err != nil {
		return domain.SystemMetrics{}, fmt.Errorf("collecting uptime: %w", err)
	}

	return domain.SystemMetrics{
		Hostname: hostname,
		CPU:      round1(cpu),
		Memory:   round1(memory),
		Disk:     round1(disk),
		Uptime:   uptime,
	}, nil
}

// round1 rounds a float to 1 decimal place.
func round1(v float64) float64 {
	return math.Round(v*10) / 10
}

