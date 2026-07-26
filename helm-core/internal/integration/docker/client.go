//go:build linux

package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/nayiswftw/helm/helm-core/internal/domain"
)

const defaultSocketPath = "/var/run/docker.sock"

// Client interacts directly with Docker Engine REST API via unix domain socket.
type Client struct {
	socketPath string
	httpClient *http.Client
}

// NewClient creates a new Docker client targeting the specified socket path.
func NewClient(socketPath string) *Client {
	if socketPath == "" {
		socketPath = defaultSocketPath
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			var dialer net.Dialer
			return dialer.DialContext(ctx, "unix", socketPath)
		},
	}

	return &Client{
		socketPath: socketPath,
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   15 * time.Second,
		},
	}
}

// IsAvailable checks whether the Docker socket exists and is accessible.
func (c *Client) IsAvailable() bool {
	info, err := os.Stat(c.socketPath)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSocket != 0
}

type dockerContainerRaw struct {
	ID      string   `json:"Id"`
	Names   []string `json:"Names"`
	Image   string   `json:"Image"`
	State   string   `json:"State"`
	Status  string   `json:"Status"`
	Created int64    `json:"Created"`
}

// ListContainers queries GET /containers/json?all=true from Docker Engine.
func (c *Client) ListContainers(ctx context.Context) ([]domain.Container, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://docker/containers/json?all=true", nil)
	if err != nil {
		return nil, fmt.Errorf("creating list containers request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("requesting docker containers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("docker API returned status %d: %s", resp.StatusCode, string(body))
	}

	var rawList []dockerContainerRaw
	if err := json.NewDecoder(resp.Body).Decode(&rawList); err != nil {
		return nil, fmt.Errorf("decoding docker containers response: %w", err)
	}

	containers := make([]domain.Container, 0, len(rawList))
	for _, item := range rawList {
		name := ""
		if len(item.Names) > 0 {
			name = strings.TrimPrefix(item.Names[0], "/")
		}

		// Truncate long container IDs to 12 chars for readability & ESP32 efficiency
		shortID := item.ID
		if len(shortID) > 12 {
			shortID = shortID[:12]
		}

		containers = append(containers, domain.Container{
			ID:      shortID,
			Name:    name,
			Image:   item.Image,
			State:   item.State,
			Status:  item.Status,
			Created: item.Created,
		})
	}

	return containers, nil
}

// StartContainer issues POST /containers/{id}/start.
func (c *Client) StartContainer(ctx context.Context, id string) error {
	return c.postAction(ctx, fmt.Sprintf("http://docker/containers/%s/start", id))
}

// StopContainer issues POST /containers/{id}/stop.
func (c *Client) StopContainer(ctx context.Context, id string) error {
	return c.postAction(ctx, fmt.Sprintf("http://docker/containers/%s/stop", id))
}

// RestartContainer issues POST /containers/{id}/restart.
func (c *Client) RestartContainer(ctx context.Context, id string) error {
	return c.postAction(ctx, fmt.Sprintf("http://docker/containers/%s/restart", id))
}

func (c *Client) postAction(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("creating container action request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing container action: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotModified {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("docker action returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
