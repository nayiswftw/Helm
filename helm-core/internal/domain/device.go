package domain

// Capability represents a feature that a device supports.
type Capability string

const (
	CapabilityMetrics      Capability = "metrics"
	CapabilityContainers   Capability = "containers"
	CapabilityServices     Capability = "services"
	CapabilityPowerControl Capability = "power_control"
	CapabilityNotifications Capability = "notifications"
)

// Device represents a manageable resource in the homelab.
// Examples: a server, NAS, router, or desktop.
type Device struct {
	ID           string       `json:"id"`
	Hostname     string       `json:"hostname"`
	Capabilities []Capability `json:"capabilities"`
}
