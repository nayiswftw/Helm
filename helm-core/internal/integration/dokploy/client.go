//go:build linux

package dokploy

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nayiswftw/helm/helm-core/internal/domain"
)

// Client communicates with the Dokploy REST API via HTTP.
// All requests are authenticated using the x-api-key header.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Dokploy API client.
// If baseURL or apiKey is empty, the client reports as unavailable.
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// IsAvailable reports whether the client has a configured base URL and API key.
func (c *Client) IsAvailable() bool {
	return c.baseURL != "" && c.apiKey != ""
}

// --- Raw Dokploy API response types ---

type rawProject struct {
	ProjectID   string `json:"projectId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
}

type rawApplication struct {
	ApplicationID string `json:"applicationId"`
	Name          string `json:"name"`
	AppType       string `json:"applicationType"`
	Status        string `json:"applicationStatus"`
	ProjectID     string `json:"projectId"`
}

type rawDeployment struct {
	DeploymentID  string `json:"deploymentId"`
	ApplicationID string `json:"applicationId"`
	Status        string `json:"status"`
	CreatedAt     string `json:"createdAt"`
}

// ListProjects fetches all projects from Dokploy.
// GET /api/project.all
func (c *Client) ListProjects(ctx context.Context) ([]domain.DokployProject, error) {
	var raw []rawProject
	if err := c.doGet(ctx, "/api/project.all", &raw); err != nil {
		return nil, fmt.Errorf("listing dokploy projects: %w", err)
	}

	projects := make([]domain.DokployProject, 0, len(raw))
	for _, r := range raw {
		projects = append(projects, domain.DokployProject{
			ProjectID:   r.ProjectID,
			Name:        r.Name,
			Description: r.Description,
			CreatedAt:   r.CreatedAt,
		})
	}
	return projects, nil
}

// GetApplication fetches a single application by ID.
// GET /api/application.one?applicationId=...
func (c *Client) GetApplication(ctx context.Context, appID string) (domain.DokployApplication, error) {
	path := "/api/application.one?applicationId=" + url.QueryEscape(appID)

	var raw rawApplication
	if err := c.doGet(ctx, path, &raw); err != nil {
		return domain.DokployApplication{}, fmt.Errorf("fetching dokploy application %s: %w", appID, err)
	}

	return domain.DokployApplication{
		ApplicationID: raw.ApplicationID,
		Name:          raw.Name,
		AppType:       raw.AppType,
		Status:        raw.Status,
		ProjectID:     raw.ProjectID,
	}, nil
}

// DeployApplication triggers a deployment for the specified application.
// POST /api/application.deploy
func (c *Client) DeployApplication(ctx context.Context, appID string) error {
	body := fmt.Sprintf(`{"applicationId":"%s"}`, appID)
	return c.doPost(ctx, "/api/application.deploy", body)
}

// RedeployApplication triggers a redeployment for the specified application.
// POST /api/application.redeploy
func (c *Client) RedeployApplication(ctx context.Context, appID string) error {
	body := fmt.Sprintf(`{"applicationId":"%s"}`, appID)
	return c.doPost(ctx, "/api/application.redeploy", body)
}

// ListDeployments fetches deployment history for an application.
// GET /api/deployment.all?applicationId=...
func (c *Client) ListDeployments(ctx context.Context, appID string) ([]domain.DokployDeployment, error) {
	path := "/api/deployment.all?applicationId=" + url.QueryEscape(appID)

	var raw []rawDeployment
	if err := c.doGet(ctx, path, &raw); err != nil {
		return nil, fmt.Errorf("listing dokploy deployments for %s: %w", appID, err)
	}

	deployments := make([]domain.DokployDeployment, 0, len(raw))
	for _, r := range raw {
		deployments = append(deployments, domain.DokployDeployment{
			DeploymentID:  r.DeploymentID,
			ApplicationID: r.ApplicationID,
			Status:        r.Status,
			CreatedAt:     r.CreatedAt,
		})
	}
	return deployments, nil
}

// --- Internal HTTP helpers ---

func (c *Client) doGet(ctx context.Context, path string, target any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("dokploy API returned status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}
	return nil
}

func (c *Client) doPost(ctx context.Context, path string, jsonBody string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, strings.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("dokploy API returned status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
}
