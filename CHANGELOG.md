# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Repository Overhaul & Production Hardening**
  - **OpenAPI 3.1 Specification**: Built-in JSON specification served at `GET /api/v1/openapi.json` documenting all API routes, headers, and schemas
  - **Live Docker Container Observability**:
    - `GET /api/v1/containers/{id}/stats` — Returns real-time container CPU percentage and memory usage in MB
    - `GET /api/v1/containers/{id}/logs?tail=100` — Streams recent stdout/stderr log lines (supports multiplexed and raw TTY logs)
  - **Host Observability Enhancements**:
    - `load_average` `[1m, 5m, 15m]` added to `SystemMetrics` (`GET /api/v1/dashboard`) via `/proc/loadavg`
    - `temperature` (CPU thermal zone temperature in Celsius) added to `SystemMetrics` via `/sys/class/thermal/thermal_zone*/temp`
  - **Performance & Resilience**:
    - **1-Second TTL Coalescing Cache** in `DashboardService` preventing `/proc/stat` blocking under high concurrent request load
    - **Dynamic Command Resolution** via `exec.LookPath` with 10-second timeout contexts for reboot/shutdown actions
    - **Docker Host Mount Compatibility** via `HELM_PROC_PATH` environment variable
  - **Production Containerization**:
    - Multi-stage production `Dockerfile` (`golang:1.26-alpine` builder → minimal alpine image)
    - Production `docker-compose.yml` with host network mode, `/var/run/docker.sock`, `/proc`, and `/sys` mounts
- **Phase 6: Notifications**
  - Webhook-based notification system dispatching JSON payloads to any URL (Discord, Slack, ntfy.sh, Gotify, etc.)
  - `Notification` domain model with event type constants (`action.executed`, `container.state_changed`, `dokploy.deployed`, `dokploy.failed`)
  - `WebhookNotifier` integration with fire-and-forget async dispatch
  - `NotificationService` injected into `ActionService` and `DokployService` for automatic event notifications
  - Endpoint: `POST /api/v1/notifications/test` for debugging webhook configuration
  - Configuration: `HELM_NOTIFY_URL` environment variable
- **Phase 5: Authentication & TLS**
  - API key authentication middleware protecting all `/api/v1/*` routes
  - Supports `Authorization: Bearer <token>` and `X-API-Key: <token>` headers
  - Constant-time key comparison via `crypto/subtle` for timing-attack resistance
  - Dev mode: authentication disabled when `HELM_API_KEY` is not set
  - Optional TLS termination with auto-detection (`HELM_TLS_CERT` + `HELM_TLS_KEY`)
  - `/health` remains unauthenticated for load balancer probes
- **Phase 4: Dokploy Integration**
  - Zero-dependency Dokploy REST API client communicating via `x-api-key` header
  - `DokployProject`, `DokployApplication`, `DokployDeployment` domain models
  - `DokployService` with availability guard (graceful 503 when unconfigured)
  - Endpoints: `GET /api/v1/dokploy/projects`, `GET /api/v1/dokploy/applications/{id}`, `POST /api/v1/dokploy/applications/{id}/deploy`, `POST /api/v1/dokploy/applications/{id}/redeploy`, `GET /api/v1/dokploy/applications/{id}/deployments`
  - Configuration: `HELM_DOKPLOY_URL` and `HELM_DOKPLOY_API_KEY` environment variables
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
