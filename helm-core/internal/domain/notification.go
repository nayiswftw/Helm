package domain

import "time"

// Notification event type constants.
const (
	EventActionExecuted        = "action.executed"
	EventContainerStateChanged = "container.state_changed"
	EventDokployDeployed       = "dokploy.deployed"
	EventDokployFailed         = "dokploy.failed"
)

// Notification represents an event notification dispatched to external systems.
type Notification struct {
	Event     string            `json:"event"`
	Message   string            `json:"message"`
	Timestamp time.Time         `json:"timestamp"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}
