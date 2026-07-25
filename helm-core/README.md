# Helm Core

Helm Core is the lightweight backend control plane for Helm, a self-hosted homelab control platform. It provides REST APIs for host system monitoring, device tracking, and controlled system management.

## Features

- **Real Linux Metrics**: Reads live host data directly from Linux kernel interfaces:
  - CPU utilization & core count (`/proc/stat`)
  - Memory usage & availability (`/proc/meminfo`)
  - Root disk space & usage (`syscall.Statfs`)
  - System uptime (`/proc/uptime`)
  - Aggregate network throughput (`/proc/net/dev`)
- **Zero Heavy Dependencies**: Built with the Go standard library and `chi` router.
- **Production Ready**: Built-in CORS middleware, structured `slog` logging, and graceful shutdown handling.

## API Reference

All endpoints are versioned under `/api/v1`:

| Endpoint | Method | Description |
|---|---|---|
| `/api/v1/health` | `GET` | Health check endpoint |
| `/api/v1/dashboard` | `GET` | Current system metrics & dashboard snapshot |
| `/api/v1/devices` | `GET` | Registered device list |
| `/api/v1/devices/{id}` | `GET` | Retrieve device details by ID |

## Configuration

Helm Core is configured via environment variables:

| Variable | Default | Description |
|---|---|---|
| `HELM_PORT` | `8080` | HTTP listen port |
| `HELM_NAME` | `Helm` | Service identity name returned by `/health` |

## Building & Running

### Requirements
- Go 1.22+
- Linux (Debian / Ubuntu / Alpine)

### Build & Run Locally
```bash
go build -o helm ./cmd/helm
./helm
```

### Cross-Compile for Linux Target
```bash
GOOS=linux GOARCH=amd64 go build -o helm ./cmd/helm
```

## Architecture

Follows a strict 4-layer unidirectional architecture:

```
API Layer (/internal/api)
   ↓
Service Layer (/internal/services)
   ↓
Integration Layer (/internal/integrations)
   ↓
Linux OS (/proc, syscalls)
```

## License

[MIT](LICENSE)
