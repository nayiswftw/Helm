package system

import (
	"os"
	"strconv"
	"strings"
)

// Uptime reads system uptime in seconds from /proc/uptime.
func (s *System) Uptime() float64 {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		s.logger.Debug("uptime unavailable", "error", err)
		return 0
	}

	fields := strings.Fields(string(data))
	if len(fields) < 1 {
		s.logger.Warn("/proc/uptime: unexpected format")
		return 0
	}

	uptime, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		s.logger.Warn("failed to parse uptime value", "error", err)
		return 0
	}

	return uptime
}
