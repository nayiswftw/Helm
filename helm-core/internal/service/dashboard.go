//go:build linux

package service

import (
	"fmt"
	"math"
	"os"
	"sync"
	"time"

	"github.com/nayiswftw/helm/helm-core/internal/domain"
	"github.com/nayiswftw/helm/helm-core/internal/integration/system"
)

// DashboardService orchestrates system metric collection
// for the dashboard API endpoint.
type DashboardService struct {
	system *system.System

	mu           sync.Mutex
	lastSnapshot domain.SystemMetrics
	lastFetched  time.Time
}

// NewDashboardService creates a new DashboardService.
func NewDashboardService(sys *system.System) *DashboardService {
	return &DashboardService{
		system: sys,
	}
}

// GetMetrics collects current system metrics and returns them
// as a single SystemMetrics snapshot.
// It caches results for 1 second to prevent blocking on /proc/stat
// two-sample CPU usage collection under concurrent request load.
func (s *DashboardService) GetMetrics() (domain.SystemMetrics, error) {
	s.mu.Lock()
	if time.Since(s.lastFetched) < time.Second && s.lastSnapshot.Hostname != "" {
		res := s.lastSnapshot
		s.mu.Unlock()
		return res, nil
	}
	s.mu.Unlock()

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

	loadavg, _ := s.system.LoadAverage()
	var roundedLoad []float64
	if len(loadavg) == 3 {
		roundedLoad = []float64{round2(loadavg[0]), round2(loadavg[1]), round2(loadavg[2])}
	}

	temp, _ := s.system.Temperature()

	metrics := domain.SystemMetrics{
		Hostname:    hostname,
		CPU:         round1(cpu),
		Memory:      round1(memory),
		Disk:        round1(disk),
		Uptime:      uptime,
		LoadAverage: roundedLoad,
		Temperature: round1(temp),
	}

	s.mu.Lock()
	s.lastSnapshot = metrics
	s.lastFetched = time.Now()
	s.mu.Unlock()

	return metrics, nil
}

// round1 rounds a float to 1 decimal place.
func round1(v float64) float64 {
	return math.Round(v*10) / 10
}

// round2 rounds a float to 2 decimal places.
func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

