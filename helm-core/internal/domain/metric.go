package domain

// SystemMetrics represents the current health snapshot of a device.
// Fields are kept lightweight for ESP32 consumption.
type SystemMetrics struct {
	Hostname    string    `json:"hostname"`
	CPU         float64   `json:"cpu"`                   // percentage 0-100
	Memory      float64   `json:"memory"`                // percentage 0-100
	Disk        float64   `json:"disk"`                  // percentage 0-100
	Uptime      int64     `json:"uptime"`                // seconds
	LoadAverage []float64 `json:"load_average,omitempty"` // [1m, 5m, 15m]
	Temperature float64   `json:"temperature,omitempty"`  // degrees Celsius
}

// NetworkMetrics represents aggregate network I/O counters.
type NetworkMetrics struct {
	BytesSent     uint64 `json:"bytes_sent"`
	BytesReceived uint64 `json:"bytes_received"`
}
