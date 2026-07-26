# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial project scaffolding for `helm-core`
- HTTP server with Chi router and graceful shutdown
- Environment-based configuration (`HELM_PORT`, `HELM_LOG_LEVEL`)
- Application dependency injection container
- Domain models: `SystemMetrics`, `NetworkMetrics`, `Device`, `Capability`
- System integration layer reading from Linux-native interfaces
  - CPU usage via `/proc/stat` (two-sample measurement)
  - Memory usage via `/proc/meminfo`
  - Disk usage via `syscall.Statfs`
  - System uptime via `/proc/uptime`
  - Network I/O via `/proc/net/dev`
- Dashboard service orchestrating metric collection
- API endpoints:
  - `GET /health` — liveness check
  - `GET /api/v1/dashboard` — system metrics snapshot
- Structured JSON request logging middleware
- Production HTTP server timeouts (read: 15s, write: 15s, idle: 60s)
- GitHub Actions CI pipeline (build + vet)
- Makefile with build, cross-compile, run, vet, and clean targets
