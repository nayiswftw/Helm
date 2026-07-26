//go:build linux

package system

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// MemoryUsage returns the current memory usage as a percentage (0-100).
// It reads /proc/meminfo and computes (Total - Available) / Total * 100.
func (s *System) MemoryUsage() (float64, error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, fmt.Errorf("opening /proc/meminfo: %w", err)
	}
	defer f.Close()

	fields := make(map[string]uint64)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		valueStr := strings.TrimSpace(parts[1])
		valueStr = strings.TrimSuffix(valueStr, " kB")
		valueStr = strings.TrimSpace(valueStr)

		v, err := strconv.ParseUint(valueStr, 10, 64)
		if err != nil {
			continue
		}

		fields[key] = v
	}

	total, ok := fields["MemTotal"]
	if !ok || total == 0 {
		return 0, fmt.Errorf("MemTotal not found in /proc/meminfo")
	}

	available, ok := fields["MemAvailable"]
	if !ok {
		return 0, fmt.Errorf("MemAvailable not found in /proc/meminfo")
	}

	usage := float64(total-available) / float64(total) * 100.0
	return usage, nil
}
