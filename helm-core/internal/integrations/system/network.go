package system

import (
	"bufio"
	"os"
	"strings"
)

func (s *System) Network() NetworkInfo {
	info := NetworkInfo{}

	file, err := os.Open("/proc/net/dev")
	if err != nil {
		s.logger.Debug("network stats unavailable", "error", err)
		return info
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		if lineNum <= 2 {
			continue
		}

		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		iface := strings.TrimSpace(parts[0])
		if iface == "lo" {
			continue
		}

		fields := strings.Fields(parts[1])
		if len(fields) < 10 {
			continue
		}

		info.RxBytes += parseUint(fields[0])
		info.TxBytes += parseUint(fields[8])
	}

	return info
}
