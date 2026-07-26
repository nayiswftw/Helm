//go:build linux

package system

import (
	"fmt"
	"os/exec"
)

// Reboot executes systemctl reboot.
func (s *System) Reboot() error {
	cmd := exec.Command("/usr/bin/systemctl", "reboot")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reboot system: %w", err)
	}
	return nil
}

// Shutdown executes systemctl poweroff.
func (s *System) Shutdown() error {
	cmd := exec.Command("/usr/bin/systemctl", "poweroff")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to power off system: %w", err)
	}
	return nil
}
