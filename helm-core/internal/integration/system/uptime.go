//go:build linux

package system

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Uptime returns the system uptime in seconds.
// It reads /proc/uptime which contains two values:
// the uptime of the system (seconds) and the idle time (seconds).
func (s *System) Uptime() (int64, error) {
	uptimePath := filepath.Join(s.ProcPath, "uptime")
	data, err := os.ReadFile(uptimePath)
	if err != nil {
		return 0, fmt.Errorf("reading %s: %w", uptimePath, err)
	}

	fields := strings.Fields(string(data))
	if len(fields) < 1 {
		return 0, fmt.Errorf("unexpected /proc/uptime format")
	}

	uptimeFloat, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, fmt.Errorf("parsing uptime value %q: %w", fields[0], err)
	}

	return int64(uptimeFloat), nil
}
