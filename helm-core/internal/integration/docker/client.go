//go:build linux

package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
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
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
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
	return c.postAction(ctx, fmt.Sprintf("http://docker/containers/%s/start", url.PathEscape(id)))
}

// StopContainer issues POST /containers/{id}/stop.
func (c *Client) StopContainer(ctx context.Context, id string) error {
	return c.postAction(ctx, fmt.Sprintf("http://docker/containers/%s/stop", url.PathEscape(id)))
}

// RestartContainer issues POST /containers/{id}/restart.
func (c *Client) RestartContainer(ctx context.Context, id string) error {
	return c.postAction(ctx, fmt.Sprintf("http://docker/containers/%s/restart", url.PathEscape(id)))
}

func (c *Client) postAction(ctx context.Context, reqURL string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, nil)
	if err != nil {
		return fmt.Errorf("creating container action request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing container action: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotModified {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("docker action returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

type dockerStatsRaw struct {
	CPUStats struct {
		CPUUsage struct {
			TotalUsage uint64 `json:"total_usage"`
		} `json:"cpu_usage"`
		SystemCPUUsage uint64 `json:"system_cpu_usage"`
		OnlineCPUs     uint64 `json:"online_cpus"`
	} `json:"cpu_stats"`
	PreCPUStats struct {
		CPUUsage struct {
			TotalUsage uint64 `json:"total_usage"`
		} `json:"cpu_usage"`
		SystemCPUUsage uint64 `json:"system_cpu_usage"`
	} `json:"precpu_stats"`
	MemoryStats struct {
		Usage uint64 `json:"usage"`
		Limit uint64 `json:"limit"`
	} `json:"memory_stats"`
}

// GetContainerStats queries GET /containers/{id}/stats?stream=false from Docker Engine.
func (c *Client) GetContainerStats(ctx context.Context, id string) (domain.ContainerStats, error) {
	reqURL := fmt.Sprintf("http://docker/containers/%s/stats?stream=false", url.PathEscape(id))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return domain.ContainerStats{}, fmt.Errorf("creating stats request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return domain.ContainerStats{}, fmt.Errorf("querying container stats: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return domain.ContainerStats{}, fmt.Errorf("container not found")
	}
	if resp.StatusCode != http.StatusOK {
		return domain.ContainerStats{}, fmt.Errorf("docker stats returned status %d", resp.StatusCode)
	}

	var raw dockerStatsRaw
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return domain.ContainerStats{}, fmt.Errorf("decoding docker stats: %w", err)
	}

	cpuDelta := float64(raw.CPUStats.CPUUsage.TotalUsage) - float64(raw.PreCPUStats.CPUUsage.TotalUsage)
	sysDelta := float64(raw.CPUStats.SystemCPUUsage) - float64(raw.PreCPUStats.SystemCPUUsage)
	cpus := float64(raw.CPUStats.OnlineCPUs)
	if cpus == 0 {
		cpus = 1.0
	}

	cpuPct := 0.0
	if sysDelta > 0 && cpuDelta > 0 {
		cpuPct = (cpuDelta / sysDelta) * cpus * 100.0
	}

	memMB := float64(raw.MemoryStats.Usage) / (1024 * 1024)
	limitMB := float64(raw.MemoryStats.Limit) / (1024 * 1024)

	return domain.ContainerStats{
		ID:            id,
		CPUPercentage: cpuPct,
		MemoryMB:      memMB,
		MemoryLimitMB: limitMB,
	}, nil
}

// GetContainerLogs queries GET /containers/{id}/logs from Docker Engine.
func (c *Client) GetContainerLogs(ctx context.Context, id string, tail int) (domain.ContainerLogs, error) {
	if tail <= 0 {
		tail = 100
	}
	reqURL := fmt.Sprintf("http://docker/containers/%s/logs?stdout=true&stderr=true&tail=%d&timestamps=false", url.PathEscape(id), tail)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return domain.ContainerLogs{}, fmt.Errorf("creating logs request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return domain.ContainerLogs{}, fmt.Errorf("querying container logs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return domain.ContainerLogs{}, fmt.Errorf("container not found")
	}
	if resp.StatusCode != http.StatusOK {
		return domain.ContainerLogs{}, fmt.Errorf("docker logs returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.ContainerLogs{}, fmt.Errorf("reading log output: %w", err)
	}

	lines := parseDockerLogs(body)
	return domain.ContainerLogs{
		ID:    id,
		Lines: lines,
	}, nil
}

// parseDockerLogs parses Docker multiplexed stream or raw TTY logs into text lines.
func parseDockerLogs(data []byte) []string {
	var lines []string
	i := 0
	n := len(data)

	for i < n {
		// Check if there's an 8-byte multiplex header
		if i+8 <= n && (data[i] == 1 || data[i] == 2) && data[i+1] == 0 && data[i+2] == 0 && data[i+3] == 0 {
			size := int(data[i+4])<<24 | int(data[i+5])<<16 | int(data[i+6])<<8 | int(data[i+7])
			if i+8+size <= n {
				chunk := string(data[i+8 : i+8+size])
				for _, line := range strings.Split(strings.TrimRight(chunk, "\r\n"), "\n") {
					if l := strings.TrimRight(line, "\r"); l != "" {
						lines = append(lines, l)
					}
				}
				i += 8 + size
				continue
			}
		}

		// Fallback for raw non-multiplexed TTY logs
		for _, line := range strings.Split(strings.TrimRight(string(data[i:]), "\r\n"), "\n") {
			if l := strings.TrimRight(line, "\r"); l != "" {
				lines = append(lines, l)
			}
		}
		break
	}

	if len(lines) == 0 {
		return []string{}
	}
	return lines
}
