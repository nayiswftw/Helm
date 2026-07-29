//go:build linux

package system

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// resolveCommand finds the path for systemctl or falls back to the specific binary name.
func resolveCommand(subcmd string, fallback string) (string, []string, error) {
	if path, err := exec.LookPath("systemctl"); err == nil {
		return path, []string{subcmd}, nil
	}
	if path, err := exec.LookPath(fallback); err == nil {
		return path, nil, nil
	}
	return "", nil, fmt.Errorf("neither systemctl nor %s found in PATH", fallback)
}

// Reboot executes reboot safely using dynamic command resolution.
func (s *System) Reboot() error {
	bin, args, err := resolveCommand("reboot", "reboot")
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin, args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reboot system: %w", err)
	}
	return nil
}

// Shutdown executes poweroff safely using dynamic command resolution.
func (s *System) Shutdown() error {
	bin, args, err := resolveCommand("poweroff", "poweroff")
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin, args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to power off system: %w", err)
	}
	return nil
}
