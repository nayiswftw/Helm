package domain

// Container represents a Docker container running on a device.
// Payload format is kept lightweight for ESP32 and remote clients.
type Container struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	State   string `json:"state"`   // e.g. "running", "exited", "paused"
	Status  string `json:"status"`  // e.g. "Up 2 hours"
	Created int64  `json:"created"` // Unix timestamp in seconds
}

// ContainerStats represents real-time CPU and Memory usage for a container.
type ContainerStats struct {
	ID            string  `json:"id"`
	CPUPercentage float64 `json:"cpu"`             // percentage e.g. 1.2%
	MemoryMB      float64 `json:"memory_mb"`       // memory usage in MB
	MemoryLimitMB float64 `json:"memory_limit_mb"` // memory limit in MB
}

// ContainerLogs represents recent log lines from a container.
type ContainerLogs struct {
	ID    string   `json:"id"`
	Lines []string `json:"lines"`
}
