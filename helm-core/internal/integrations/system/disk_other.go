//go:build !linux

package system

func (s *System) Disk() DiskInfo {
	s.logger.Debug("disk stats unavailable on this platform")
	return DiskInfo{}
}
