//go:build linux

package system

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Temperature scans /sys/class/thermal/thermal_zone*/temp for CPU thermal sensors.
// It returns the highest recorded temperature in Celsius across all valid zones.
// If no sensors are found or readable, it returns 0.0 without error so systems
// lacking thermal drivers don't fail metrics collection.
func (s *System) Temperature() (float64, error) {
	matches, err := filepath.Glob("/sys/class/thermal/thermal_zone*/temp")
	if err != nil || len(matches) == 0 {
		return 0.0, nil
	}

	var maxTemp float64
	var found bool

	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		rawStr := strings.TrimSpace(string(data))
		val, err := strconv.ParseFloat(rawStr, 64)
		if err != nil || val <= 0 {
			continue
		}

		// Values are typically in millidegrees Celsius (e.g., 45000 -> 45.0 C)
		tempC := val
		if val > 1000 {
			tempC = val / 1000.0
		}

		if !found || tempC > maxTemp {
			maxTemp = tempC
			found = true
		}
	}

	if !found {
		return 0.0, nil
	}

	return maxTemp, nil
}
