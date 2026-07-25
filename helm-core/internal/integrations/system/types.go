package system

// CPUInfo represents CPU metrics.
type CPUInfo struct {
	NumCores     int     `json:"num_cores"`
	UsagePercent float64 `json:"usage_percent"`
}

// MemoryInfo represents memory metrics in bytes.
type MemoryInfo struct {
	Total        uint64  `json:"total"`
	Used         uint64  `json:"used"`
	Available    uint64  `json:"available"`
	UsagePercent float64 `json:"usage_percent"`
}

// DiskInfo represents disk space metrics in bytes.
type DiskInfo struct {
	MountPoint   string  `json:"mount_point"`
	Total        uint64  `json:"total"`
	Used         uint64  `json:"used"`
	Free         uint64  `json:"free"`
	UsagePercent float64 `json:"usage_percent"`
}

// NetworkInfo represents network traffic metrics in bytes.
type NetworkInfo struct {
	RxBytes uint64 `json:"rx_bytes"`
	TxBytes uint64 `json:"tx_bytes"`
}
