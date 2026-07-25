package system

import (
	"bufio"
	"os"
	"runtime"
	"strings"
)

func (s *System) Memory() MemoryInfo {
	info, err := readProcMeminfo()
	if err != nil {
		s.logger.Debug("system memory unavailable", "error", err)
		return memoryFallback()
	}
	return info
}

func readProcMeminfo() (MemoryInfo, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return MemoryInfo{}, err
	}
	defer file.Close()

	var info MemoryInfo
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "MemTotal:":
			info.Total = parseUint(fields[1]) * 1024
		case "MemAvailable:":
			info.Available = parseUint(fields[1]) * 1024
		}
	}

	if info.Total > 0 {
		if info.Available <= info.Total {
			info.Used = info.Total - info.Available
		}
		info.UsagePercent = (float64(info.Used) / float64(info.Total)) * 100
	}

	return info, nil
}

func memoryFallback() MemoryInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return MemoryInfo{
		Total:     m.Sys,
		Used:      m.Alloc,
		Available: m.Sys - m.Alloc,
	}
}
