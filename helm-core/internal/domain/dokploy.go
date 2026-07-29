package domain

// DokployProject represents a project within the Dokploy deployment platform.
type DokployProject struct {
	ProjectID   string `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

// DokployApplication represents an application managed by Dokploy.
type DokployApplication struct {
	ApplicationID string `json:"application_id"`
	Name          string `json:"name"`
	AppType       string `json:"app_type"` // e.g. "application", "docker-compose"
	Status        string `json:"status"`   // e.g. "idle", "running", "done", "error"
	ProjectID     string `json:"project_id"`
}

// DokployDeployment represents a single deployment record from Dokploy.
type DokployDeployment struct {
	DeploymentID  string `json:"deployment_id"`
	ApplicationID string `json:"application_id"`
	Status        string `json:"status"` // "idle", "running", "done", "error"
	CreatedAt     string `json:"created_at"`
}
