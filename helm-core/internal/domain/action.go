package domain

// Action represents a predefined, safe administrative operation.
type Action struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DeviceID    string `json:"device_id"`
	Dangerous   bool   `json:"dangerous"` // Flag for client UI confirmation
}

// ActionResult represents the result of executing an action.
type ActionResult struct {
	ActionID string `json:"action_id"`
	Status   string `json:"status"` // e.g. "accepted", "completed", "failed"
	Message  string `json:"message"`
}
