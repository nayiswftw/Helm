package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port int
	Name string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port: 8080,
		Name: "Helm",
	}

	if v := os.Getenv("HELM_PORT"); v != "" {
		port, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("HELM_PORT invalid: %w", err)
		}
		cfg.Port = port
	}

	if v := os.Getenv("HELM_NAME"); v != "" {
		cfg.Name = v
	}

	return cfg, nil
}

func (c *Config) Addr() string {
	return fmt.Sprintf(":%d", c.Port)
}
