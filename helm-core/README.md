<div align="center">

# ☸️ Helm Core

**Backend Control Plane for Homelab Infrastructure**

[![CI & Snapshot Release](https://github.com/nayiswftw/Helm/actions/workflows/ci.yml/badge.svg)](https://github.com/nayiswftw/Helm/actions/workflows/ci.yml)
[![Python Version](https://img.shields.io/badge/Python-3.10%2B-blue?style=flat&logo=python)](https://python.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Debian%20%7C%20Linux-orange.svg)](https://www.debian.org)

Helm Core is a self-hosted control plane that monitors and manages your homelab server.  
Designed to run with **zero heavy dependencies**, minimal memory footprint, and high reliability.

[Quick Start](#-quick-start-on-debian) • [API Documentation](#-api-documentation) • [Systemd Setup](#%EF%B8%8F-run-as-a-systemd-service) • [Architecture](#%EF%B8%8F-architecture)

</div>

---

## ✨ Features

- 📊 **Real Kernel Metrics**: Direct `/proc` and `statfs` parsing without external metric agents.
  - **CPU**: Utilization % & core counts (`/proc/stat`).
  - **Memory**: Total, used, available memory & usage % (`/proc/meminfo`).
  - **Disk**: Total, free, used bytes & disk usage % (`syscall.Statfs`).
  - **Uptime**: System uptime in seconds (`/proc/uptime`).
  - **Network**: Aggregate Rx/Tx byte throughput (`/proc/net/dev`).
- ⚡ **Lightweight & Fast**: Written in Python using FastAPI for high performance.
- 🌐 **Web Ready**: Built-in CORS headers for `helm-web` frontend integration.
- 🛡️ **Robust**: Structured `slog` logging and graceful `SIGINT`/`SIGTERM` server shutdown.

---

## 🚀 Quick Start on Debian

### 1. Setup Environment

```bash
# Install Python & Git
sudo apt update && sudo apt install -y python3 python3-pip python3-venv git

# Clone & Setup
git clone https://github.com/nayiswftw/Helm.git
cd Helm/helm-core

# Create virtual environment
python3 -m venv venv
source venv/bin/activate

# Install dependencies
pip install -r requirements.txt

# Run
python main.py
```

## ⚙️ Configuration

Helm Core is configured via environment variables:

| Variable | Default | Description |
|---|---|---|
| `HELM_PORT` | `8080` | HTTP listen port |
| `HELM_NAME` | *(Dynamic OS Hostname)* | Service identity name returned by `/health` |

---

## 📡 API Documentation

All endpoints are versioned under `/api/v1`:

### `GET /api/v1/health`
Checks server health status.

**Example Request:**
```bash
curl -s http://localhost:8080/api/v1/health
```

**Response (`200 OK`):**
```json
{
  "name": "debian-server",
  "status": "ok"
}
```

---

### `GET /api/v1/dashboard`
Returns live host system metrics and registered device summary.

**Example Request:**
```bash
curl -s http://localhost:8080/api/v1/dashboard
```

**Response (`200 OK`):**
```json
{
  "hostname": "debian-server",
  "uptime": 86400.5,
  "cpu": {
    "num_cores": 8,
    "usage_percent": 12.4
  },
  "memory": {
    "total": 17179869184,
    "used": 4294967296,
    "available": 12884901888,
    "usage_percent": 25.0
  },
  "disk": {
    "mount_point": "/",
    "total": 512000000000,
    "used": 102400000000,
    "free": 409600000000,
    "usage_percent": 20.0
  },
  "network": {
    "rx_bytes": 104857600,
    "tx_bytes": 52428800
  },
  "devices": 1
}
```

---

### `GET /api/v1/devices`
Lists all known devices.

**Response (`200 OK`):**
```json
[
  {
    "id": "local",
    "name": "debian-server",
    "type": "server",
    "status": "online"
  }
]
```

---

## 🛠️ Run as a Systemd Service

To run Helm Core reliably in the background on Debian:

1. **Move binary to `/usr/local/bin`**:
   ```bash
   sudo mv helm-core /usr/local/bin/
   ```

2. **Create Service Unit File**:
   ```bash
   sudo nano /etc/systemd/system/helm-core.service
   ```

   Paste the following configuration:
   ```ini
   [Unit]
   Description=Helm Core Control Plane
   After=network.target

   [Service]
   Type=simple
   User=root
   WorkingDirectory=/path/to/Helm/helm-core
   ExecStart=/path/to/Helm/helm-core/venv/bin/python main.py
   Restart=always
   RestartSec=5s
   Environment=HELM_PORT=8080

   [Install]
   WantedBy=multi-user.target
   ```

3. **Enable and Start Service**:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable --now helm-core
   ```

4. **Manage Service**:
   ```bash
   # Check service status
   sudo systemctl status helm-core

   # View live logs
   sudo journalctl -u helm-core -f
   ```

---

## 🏗️ Architecture

Helm Core strictly follows a unidirectional layered architecture:

```
                  ┌────────────────────────┐
                  │       API Layer        │  HTTP Routing & CORS
                  └───────────┬────────────┘
                              │
                              ▼
                  ┌────────────────────────┐
                  │     Service Layer      │  Business Logic & Aggregation
                  └───────────┬────────────┘
                              │
                              ▼
                  ┌────────────────────────┐
                  │   Integration Layer    │  Kernel / OS Interface
                  └───────────┬────────────┘
                              │
                              ▼
                  ┌────────────────────────┐
                  │  Debian / Linux Kernel │  /proc/stat, /proc/meminfo, statfs
                  └────────────────────────┘
```

---

## 📄 License

This project is open-source software licensed under the [MIT License](LICENSE).
