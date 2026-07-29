<div align="center">

# ⎈ Helm

**Your homelab, from anywhere.**

A self-hosted control platform for monitoring and managing your entire infrastructure — from a handheld remote, web dashboard, or CLI.

[![CI](https://github.com/nayiswftw/helm/actions/workflows/ci.yml/badge.svg)](https://github.com/nayiswftw/helm/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/go-1.26+-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-blue?style=flat)](LICENSE)
[![Release](https://img.shields.io/badge/release-snapshot-orange?style=flat)](https://github.com/nayiswftw/helm/releases/tag/snapshot)

---

<br>

```
  GET /api/v1/dashboard
```
```json
{"hostname":"atlas","cpu":12.3,"memory":45.6,"disk":61.2,"uptime":864000}
```

<br>

</div>

## Why Helm?

Most homelab tools let you manage containers. Helm lets you manage **everything** — servers, containers, deployments, power state — through a single, lightweight API designed for embedded devices and low-bandwidth clients.

- 🖥️ **Real system metrics** from Linux-native interfaces — no agents, no overhead
- 📡 **ESP32-first API design** — flat JSON, tiny payloads, minimal bandwidth
- 🏗️ **Layered architecture** — clean separation between API, business logic, and system integration
- 🔒 **Security-first** — designed for public Internet exposure from day one
- ⚡ **Single static binary** — zero runtime dependencies, deploy anywhere

<br>

## Quick Start

> **Requirements:** Go 1.26+ · Linux (Debian recommended)

```bash
# Clone
git clone https://github.com/nayiswftw/helm.git
cd helm/helm-core

# Build
go build -trimpath -o helm ./cmd/helm/

# Run
HELM_PORT=:8080 ./helm
```

```bash
# Health check
curl localhost:8080/health

# System metrics
curl localhost:8080/api/v1/dashboard
```

Or grab a pre-built binary from [**Releases**](https://github.com/nayiswftw/helm/releases/tag/snapshot).

<br>

## API Reference

All responses are flat JSON — no wrappers. Designed for embedded clients.

```
GET  /health                      → Liveness probe
GET  /api/v1/openapi.json         → OpenAPI 3.1 specification
GET  /api/v1/dashboard            → System metrics snapshot
GET  /api/v1/devices              → List all managed devices
GET  /api/v1/devices/{id}         → Get specific device details
GET  /api/v1/actions              → List available administrative actions
POST /api/v1/actions/{id}/execute → Execute a predefined action
GET  /api/v1/containers           → List Docker containers
GET  /api/v1/containers/{id}/stats → Get live container CPU & memory usage
GET  /api/v1/containers/{id}/logs  → Get recent container log lines
POST /api/v1/containers/{id}/start → Start a container
POST /api/v1/containers/{id}/stop → Stop a container
POST /api/v1/containers/{id}/restart → Restart a container
GET  /api/v1/dokploy/projects     → List Dokploy projects
GET  /api/v1/dokploy/applications/{id} → Get Dokploy app details
POST /api/v1/dokploy/applications/{id}/deploy → Deploy a Dokploy app
POST /api/v1/dokploy/applications/{id}/redeploy → Redeploy a Dokploy app
GET  /api/v1/dokploy/applications/{id}/deployments → List deployments
POST /api/v1/notifications/test   → Send test webhook notification
```

<details>
<summary><b>Response Examples</b></summary>

<br>

**`GET /health`**
```json
{"status":"ok"}
```

**`GET /api/v1/dashboard`**
```json
{
  "hostname": "atlas",
  "cpu": 12.3,
  "memory": 45.6,
  "disk": 61.2,
  "uptime": 864000,
  "load_average": [0.15, 0.22, 0.18],
  "temperature": 43.5
}
```

**Error**
```json
{
  "error": {
    "code": "metrics_unavailable",
    "message": "Failed to collect system metrics"
  }
}
```

</details>

<br>

## Configuration

| Variable | Default | Description |
|:---------|:--------|:------------|
| `HELM_PORT` | `:8080` | Listen address |
| `HELM_LOG_LEVEL` | `info` | `debug` · `info` · `warn` · `error` |
| `HELM_PROC_PATH` | `/proc` | Path to host `/proc` filesystem (use `/host/proc` inside Docker) |
| `HELM_API_KEY` | *(empty)* | API key authentication token (dev mode if empty) |
| `HELM_TLS_CERT` | *(empty)* | Path to TLS certificate file |
| `HELM_TLS_KEY` | *(empty)* | Path to TLS private key file |
| `HELM_DOKPLOY_URL` | *(empty)* | Dokploy instance base URL |
| `HELM_DOKPLOY_API_KEY` | *(empty)* | Dokploy API key |
| `HELM_NOTIFY_URL` | *(empty)* | Webhook URL for notifications |

<br>

## Architecture

Dependencies always flow downward. Never upward. Never sideways.

```
┌─────────────────────────────────────────────┐
│  API Layer                                  │
│  HTTP routing · validation · responses      │
├─────────────────────────────────────────────┤
│  Service Layer                              │
│  Business logic · orchestration             │
├─────────────────────────────────────────────┤
│  Integration Layer                          │
│  /proc/stat · /proc/meminfo · statfs        │
├─────────────────────────────────────────────┤
│  Linux Kernel                               │
└─────────────────────────────────────────────┘
```

<details>
<summary><b>Project Structure</b></summary>

<br>

```
helm-core/
├── cmd/helm/
│   └── main.go                    # Entrypoint & graceful shutdown
├── internal/
│   ├── api/                       # HTTP handlers & routing
│   │   ├── router.go              # Chi router + middleware stack
│   │   ├── health.go              # GET /health
│   │   ├── dashboard.go           # GET /api/v1/dashboard
│   │   └── middleware.go          # Structured request logging
│   ├── app/
│   │   └── application.go        # Dependency injection container
│   ├── config/
│   │   └── config.go             # Environment variable loading
│   ├── domain/
│   │   ├── device.go             # Device & Capability models
│   │   └── metric.go             # SystemMetrics & NetworkMetrics
│   ├── integration/
│   │   └── system/               # Linux-native metric collection
│   │       ├── cpu.go            # /proc/stat (two-sample)
│   │       ├── memory.go         # /proc/meminfo
│   │       ├── disk.go           # syscall.Statfs
│   │       ├── uptime.go         # /proc/uptime
│   │       └── network.go        # /proc/net/dev
│   ├── server/
│   │   └── server.go             # HTTP server with timeouts
│   └── service/
│       └── dashboard.go          # Metric collection orchestrator
├── Makefile
├── go.mod
└── go.sum
```

</details>

<br>

## System Metrics

Helm reads directly from Linux kernel interfaces — no external tools, no polling agents.

| Metric | Source | Method |
|:-------|:-------|:-------|
| CPU | `/proc/stat` | Two-sample delta (200ms) |
| Memory | `/proc/meminfo` | `MemAvailable` / `MemTotal` |
| Disk | `syscall.Statfs` | Root filesystem usage |
| Uptime | `/proc/uptime` | Seconds since boot |
| Network | `/proc/net/dev` | Aggregate non-loopback I/O |

<br>

## Development

```bash
# Build for current platform
make build

# Cross-compile for Linux amd64
make build-linux

# Build all targets (amd64 + arm64)
make build-all

# Run static analysis
make vet

# Clean artifacts
make clean
```

> **Windows users:** Use `go build` directly or install Make via [Chocolatey](https://chocolatey.org/) (`choco install make`).

<br>

## Roadmap

| Feature | Status |
|:--------|:-------|
| Server health monitoring | ✅ Done |
| Dashboard API | ✅ Done |
| CI/CD & snapshot releases | ✅ Done |
| Device management | ✅ Done |
| Action framework | ✅ Done |
| Docker integration | ✅ Done |
| Dokploy integration | ✅ Done |
| Authentication & TLS | ✅ Done |
| Notifications | ✅ Done |
| ESP32 remote client | 🔜 Next |
| Web dashboard | 📋 Future |

<br>

## Tech Stack

<table>
<tr><td><b>Language</b></td><td>Go 1.26+</td></tr>
<tr><td><b>Router</b></td><td><a href="https://github.com/go-chi/chi">Chi v5</a></td></tr>
<tr><td><b>Logging</b></td><td><code>log/slog</code> (structured JSON)</td></tr>
<tr><td><b>Config</b></td><td>Environment variables</td></tr>
<tr><td><b>External deps</b></td><td>1</td></tr>
</table>

<br>

## Security

Helm is designed to be exposed to the public Internet.

- No unrestricted shell access
- No arbitrary command execution
- Only predefined, validated actions
- API authentication via bearer token or API key header
- Optional TLS termination
- Role-based access control *(planned)*

See [**SECURITY.md**](SECURITY.md) for vulnerability reporting.

<br>

## Contributing

See [**CONTRIBUTING.md**](CONTRIBUTING.md) for development setup, coding standards, and commit conventions.

<br>

## License

MIT — see [LICENSE](LICENSE).

<div align="center">
<br>
<sub>Built for homelabs that never sleep.</sub>
</div>
