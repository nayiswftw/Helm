# Helm

A self-hosted homelab control platform.

Monitor and control your entire homelab infrastructure from anywhere in the world.

## Overview

Helm is a lightweight backend server that exposes a stable HTTP API for monitoring and managing homelab infrastructure. The primary client is an ESP32-S3 handheld remote with a 1.9-inch integrated display.

### Features

- **Server Health Monitoring** — CPU, memory, disk, network, uptime
- **Lightweight API** — Minimal JSON payloads optimized for embedded devices
- **Structured Logging** — JSON-formatted logs via `log/slog`
- **Graceful Shutdown** — Clean signal handling for production deployments

### Roadmap

- [ ] Docker container management
- [ ] Dokploy integration
- [ ] Predefined administrative actions
- [ ] Device management (multi-device support)
- [ ] API authentication & TLS
- [ ] Notifications

## Tech Stack

| Component | Choice |
|-----------|--------|
| Language | Go 1.26+ |
| Router | [Chi v5](https://github.com/go-chi/chi) |
| Logging | `log/slog` |
| Config | Environment variables |
| External deps | 1 (Chi) |

## Quick Start

**Requirements:** Go 1.26+, Linux (Debian recommended)

```bash
# Build
make build

# Run
make run

# Or manually
go build -o helm ./cmd/helm/
HELM_PORT=:8080 ./helm
```

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `HELM_PORT` | `:8080` | Address to listen on |
| `HELM_LOG_LEVEL` | `info` | Log level: `debug`, `info`, `warn`, `error` |

## API

All responses are flat JSON — no wrapper objects. Optimized for ESP32 clients.

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Liveness check |
| `GET` | `/api/v1/dashboard` | System metrics snapshot |

### Examples

```bash
curl http://localhost:8080/health
# {"status":"ok"}

curl http://localhost:8080/api/v1/dashboard
# {"hostname":"server","cpu":12.3,"memory":45.6,"disk":61.2,"uptime":86400}
```

### Error Format

```json
{
  "error": {
    "code": "metrics_unavailable",
    "message": "Failed to collect system metrics"
  }
}
```

## Architecture

Dependencies flow strictly downward:

```
API Layer        → HTTP routing, request parsing, response formatting
    ↓
Service Layer    → Business logic, orchestration
    ↓
Integration      → Linux interfaces (/proc/*, statfs)
    ↓
Linux            → Kernel interfaces
```

## Project Structure

```
helm-core/
├── cmd/helm/main.go                          # Entrypoint
├── internal/
│   ├── api/                                  # HTTP handlers & routing
│   │   ├── router.go                         # Chi router setup
│   │   ├── health.go                         # GET /health
│   │   ├── dashboard.go                      # GET /api/v1/dashboard
│   │   └── middleware.go                      # Request logging
│   ├── app/                                  # Dependency injection container
│   │   └── application.go
│   ├── config/                               # Environment configuration
│   │   └── config.go
│   ├── domain/                               # Core domain models
│   │   ├── device.go
│   │   └── metric.go
│   ├── integration/
│   │   └── system/                           # Linux system metrics
│   │       ├── system.go                     # System aggregator
│   │       ├── cpu.go                        # /proc/stat
│   │       ├── memory.go                     # /proc/meminfo
│   │       ├── disk.go                       # syscall.Statfs
│   │       ├── uptime.go                     # /proc/uptime
│   │       └── network.go                    # /proc/net/dev
│   ├── server/                               # HTTP server wrapper
│   │   └── server.go
│   └── service/                              # Business logic
│       └── dashboard.go
├── go.mod
└── go.sum
```

## Target Environment

- **OS:** Debian Linux
- **Deployment:** Native systemd service
- **Primary Client:** ESP32-S3 with 1.9" integrated display

## Development

```bash
# Build for current platform (requires Linux)
make build

# Cross-compile for Linux from any OS
make build-linux

# Run static analysis
make vet

# Clean build artifacts
make clean
```

## Security

Helm is designed to be exposed to the public Internet. Security is a first-class requirement.

- No unrestricted shell access
- Only predefined, validated actions
- API authentication *(planned)*
- TLS *(planned)*
- Role-based permissions *(planned)*

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
