# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Phase 3: Docker Integration (Milestone 11)**
  - Zero-dependency `docker.Client` communicating directly with `/var/run/docker.sock` over Unix domain sockets using stdlib `net/http`
  - `Container` domain model (`ID`, `Name`, `Image`, `State`, `Status`, `Created`)
  - `ContainerService` with automatic availability detection
  - Endpoints: `GET /api/v1/containers`, `POST /api/v1/containers/{id}/start`, `POST /api/v1/containers/{id}/stop`, `POST /api/v1/containers/{id}/restart`
- **Phase 2: Device Management (Milestone 9)**
  - `DeviceService` with local device auto-registration and capability probing (`metrics`, `power_control`, `containers`)
  - Expanded `Device` domain model (`Platform`, `Architecture`, `Status`)
  - Endpoints: `GET /api/v1/devices` and `GET /api/v1/devices/{id}`
- **Phase 2: Action Framework (Milestone 10)**
  - `ActionService` with security-first predefined action registry
  - System power management integration (`Reboot()`, `Shutdown()`) using hardcoded `/usr/bin/systemctl` calls
  - `Action` and `ActionResult` domain models (with `dangerous` flag for ESP32 UI prompts)
  - Endpoints: `GET /api/v1/actions` and `POST /api/v1/actions/{id}/execute` (returns `202 Accepted`)
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
