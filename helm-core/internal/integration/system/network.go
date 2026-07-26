//go:build linux

package system

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nayiswftw/helm/helm-core/internal/domain"
)

// NetworkUsage returns aggregate network I/O counters across all
// non-loopback interfaces. It reads /proc/net/dev.
func (s *System) NetworkUsage() (domain.NetworkMetrics, error) {
	f, err := os.Open("/proc/net/dev")
	if err != nil {
		return domain.NetworkMetrics{}, fmt.Errorf("opening /proc/net/dev: %w", err)
	}
	defer f.Close()

	var metrics domain.NetworkMetrics
	scanner := bufio.NewScanner(f)

	// Skip the first two header lines.
	for i := 0; i < 2 && scanner.Scan(); i++ {
	}

	for scanner.Scan() {
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
		// /proc/net/dev fields after the colon:
		// rx_bytes rx_packets ... (8 fields) tx_bytes tx_packets ...
		if len(fields) < 10 {
			continue
		}

		rxBytes, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			return domain.NetworkMetrics{}, fmt.Errorf("parsing rx_bytes for %s: %w", iface, err)
		}

		txBytes, err := strconv.ParseUint(fields[8], 10, 64)
		if err != nil {
			return domain.NetworkMetrics{}, fmt.Errorf("parsing tx_bytes for %s: %w", iface, err)
		}

		metrics.BytesReceived += rxBytes
		metrics.BytesSent += txBytes
	}

	if err := scanner.Err(); err != nil {
		return domain.NetworkMetrics{}, fmt.Errorf("reading /proc/net/dev: %w", err)
	}

	return metrics, nil
}
