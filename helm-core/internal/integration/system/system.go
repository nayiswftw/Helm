//go:build linux

package system

// System provides access to Linux system metrics.
// Each metric is collected via Linux-native interfaces (/proc/*, statfs).
type System struct{}

// New creates a new System integration.
func New() *System {
	return &System{}
}
