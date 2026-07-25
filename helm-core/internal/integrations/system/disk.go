package system

import "syscall"

func (s *System) Disk() DiskInfo {
	info := DiskInfo{MountPoint: "/"}

	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err != nil {
		s.logger.Warn("failed statfs /", "error", err)
		return info
	}

	info.Total = stat.Blocks * uint64(stat.Bsize)
	info.Free = stat.Bavail * uint64(stat.Bsize)
	info.Used = info.Total - info.Free

	if info.Total > 0 {
		info.UsagePercent = (float64(info.Used) / float64(info.Total)) * 100
	}

	return info
}
