package system

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type cpuTimes struct {
	user, nice, system, idle, iowait, irq, softirq, steal uint64
}

func (c cpuTimes) total() uint64 {
	return c.user + c.nice + c.system + c.idle + c.iowait + c.irq + c.softirq + c.steal
}

func (c cpuTimes) active() uint64 {
	return c.user + c.nice + c.system + c.irq + c.softirq + c.steal
}

func (s *System) CPU() CPUInfo {
	info := CPUInfo{NumCores: runtime.NumCPU()}

	t1, err := readCPUTimes()
	if err != nil {
		s.logger.Debug("cpu utilization unavailable, returning core count only", "error", err)
		return info
	}

	time.Sleep(100 * time.Millisecond)

	t2, err := readCPUTimes()
	if err != nil {
		s.logger.Warn("failed to read /proc/stat (second sample)", "error", err)
		return info
	}

	totalDelta := float64(t2.total() - t1.total())
	if totalDelta > 0 {
		activeDelta := float64(t2.active() - t1.active())
		info.UsagePercent = (activeDelta / totalDelta) * 100
	}

	return info
}

func readCPUTimes() (cpuTimes, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return cpuTimes{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)
			if len(fields) < 9 {
				return cpuTimes{}, fmt.Errorf("/proc/stat format error")
			}

			return cpuTimes{
				user:    parseUint(fields[1]),
				nice:    parseUint(fields[2]),
				system:  parseUint(fields[3]),
				idle:    parseUint(fields[4]),
				iowait:  parseUint(fields[5]),
				irq:     parseUint(fields[6]),
				softirq: parseUint(fields[7]),
				steal:   parseUint(fields[8]),
			}, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return cpuTimes{}, err
	}

	return cpuTimes{}, fmt.Errorf("/proc/stat cpu line not found")
}

func parseUint(s string) uint64 {
	val, _ := strconv.ParseUint(s, 10, 64)
	return val
}
