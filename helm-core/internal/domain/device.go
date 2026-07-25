package domain

// Device represents a controllable or monitorable entity in the homelab.
type Device struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
}
