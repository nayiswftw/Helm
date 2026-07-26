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
