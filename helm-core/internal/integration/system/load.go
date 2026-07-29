//go:build linux

package system

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// LoadAverage returns the 1m, 5m, and 15m load averages as a slice of floats.
// It reads from /proc/loadavg (or configured HELM_PROC_PATH/loadavg).
func (s *System) LoadAverage() ([]float64, error) {
	loadavgPath := filepath.Join(s.ProcPath, "loadavg")
	data, err := os.ReadFile(loadavgPath)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", loadavgPath, err)
	}

	fields := strings.Fields(string(data))
	if len(fields) < 3 {
		return nil, fmt.Errorf("unexpected loadavg format in %s: %q", loadavgPath, string(data))
	}

	load1, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return nil, fmt.Errorf("parsing 1m load %q: %w", fields[0], err)
	}

	load5, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return nil, fmt.Errorf("parsing 5m load %q: %w", fields[1], err)
	}

	load15, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, fmt.Errorf("parsing 15m load %q: %w", fields[2], err)
	}

	return []float64{load1, load5, load15}, nil
}
