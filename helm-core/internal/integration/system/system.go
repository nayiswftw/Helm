//go:build linux

package system

import "os"

// System provides access to Linux system metrics.
// Each metric is collected via Linux-native interfaces (/proc/*, statfs).
type System struct {
	ProcPath string
}

// New creates a new System integration.
func New() *System {
	procPath := os.Getenv("HELM_PROC_PATH")
	if procPath == "" {
		procPath = "/proc"
	}
	return &System{
		ProcPath: procPath,
	}
}
