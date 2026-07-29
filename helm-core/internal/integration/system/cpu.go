//go:build linux

package system

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// CPUUsage returns the current CPU usage as a percentage (0-100).
// It takes two samples from /proc/stat separated by a short interval
// to compute instantaneous usage.
func (s *System) CPUUsage() (float64, error) {
	idle1, total1, err := s.readCPUStat()
	if err != nil {
		return 0, fmt.Errorf("reading initial cpu stat: %w", err)
	}

	time.Sleep(200 * time.Millisecond)

	idle2, total2, err := s.readCPUStat()
	if err != nil {
		return 0, fmt.Errorf("reading second cpu stat: %w", err)
	}

	idleDelta := float64(idle2 - idle1)
	totalDelta := float64(total2 - total1)

	if totalDelta == 0 {
		return 0, nil
	}

	usage := (1.0 - idleDelta/totalDelta) * 100.0
	return usage, nil
}

// readCPUStat reads the aggregate CPU line from /proc/stat (or HELM_PROC_PATH/stat).
// Returns idle ticks and total ticks.
func (s *System) readCPUStat() (idle, total uint64, err error) {
	statPath := filepath.Join(s.ProcPath, "stat")
	f, err := os.Open(statPath)
	if err != nil {
		return 0, 0, fmt.Errorf("opening %s: %w", statPath, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "cpu ") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			return 0, 0, fmt.Errorf("unexpected /proc/stat format: %s", line)
		}

		// Fields: cpu user nice system idle iowait irq softirq steal guest guest_nice
		var values []uint64
		for _, field := range fields[1:] {
			v, err := strconv.ParseUint(field, 10, 64)
			if err != nil {
				return 0, 0, fmt.Errorf("parsing cpu field %q: %w", field, err)
			}
			values = append(values, v)
		}

		for _, v := range values {
			total += v
		}

		// idle is the 4th value (index 3)
		if len(values) > 3 {
			idle = values[3]
		}
		// iowait is the 5th value (index 4) — also considered idle
		if len(values) > 4 {
			idle += values[4]
		}

		return idle, total, nil
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, fmt.Errorf("reading /proc/stat: %w", err)
	}

	return 0, 0, fmt.Errorf("cpu line not found in /proc/stat")
}

