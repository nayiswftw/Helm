//go:build linux

package system

import (
	"fmt"
	"syscall"
)

// DiskUsage returns the root filesystem usage as a percentage (0-100).
// It uses the statfs syscall on "/" to compute (Total - Free) / Total * 100.
func (s *System) DiskUsage() (float64, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err != nil {
		return 0, fmt.Errorf("statfs on /: %w", err)
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)

	if total == 0 {
		return 0, fmt.Errorf("total disk space is zero")
	}

	usage := float64(total-free) / float64(total) * 100.0
	return usage, nil
}
