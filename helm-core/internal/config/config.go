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
	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		hostname = "helm"
	}

	cfg := &Config{
		Port: 8080,
		Name: hostname,
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
