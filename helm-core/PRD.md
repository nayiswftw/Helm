# Helm Core — Product Requirements Document

## 1. Summary

Helm Core is the backend control plane for Helm, a self-hosted homelab
control platform. It exposes an HTTP (and later WebSocket) API that lets
clients — a web UI, a CLI, and an ESP32-based hardware remote — monitor and
operate a homelab: view system metrics, manage devices, and execute
controlled actions (restart, shutdown, service control, deploys).

This PRD covers **helm-core only**. Other components (`helm-remote`,
`helm-agent`, `helm-web`, `helm-cli`, `helm-sdk`) are referenced only where
they constrain core's design.

## 2. Problem

Homelab operators currently juggle SSH sessions, multiple dashboards
(Grafana, Portainer, Dokploy, router UI, NAS UI), and ad-hoc scripts to
answer basic questions ("is the server up?", "how much disk is left?",
"restart that container") or take basic actions. There's no single,
lightweight, self-hosted control surface that works from a phone, a CLI,
or a dedicated piece of hardware on the desk.

## 3. Goals

- One backend that answers "what's the state of my homelab?" and executes
  a small, permission-controlled set of actions against it.
- Small enough to run comfortably on the same Debian box it monitors.
- API payloads small enough for an ESP32 to parse and render.
- Extensible without rearchitecting: new integrations (Docker, Dokploy,
  future agents) plug into the same layered structure.

## 4. Non-Goals (for now)

- Not a general-purpose infra-as-code / deployment tool (Dokploy already
  does this — Helm integrates with it, doesn't replace it).
- Not a multi-tenant SaaS. Single homelab owner, single deployment.
- Not a full metrics time-series store (defer to Prometheus/Grafana if the
  user wants history; Helm Core surfaces current-state snapshots).
- No arbitrary shell execution exposed via API, ever.
- No mobile app in this phase (web UI + CLI + ESP32 remote only).

## 5. Users

Single-user / self-hosted. The "user" is the homelab owner operating their
own infrastructure through their own devices. No concept of external/public
users in this phase.

## 6. Current Scope (Phase 1–3)

### 6.1 Phase 1 — Core Foundation (done)
- HTTP server, routing, Application container, service layer, domain
  models in place.

### 6.2 Phase 2 — Real System Data (current focus)
**Requirement:** Replace placeholder metrics with real Debian/Linux data.

| Metric   | Source            |
|----------|--------------------|
| CPU      | `/proc/stat`       |
| Memory   | `/proc/meminfo`    |
| Disk     | `statfs()`         |
| Uptime   | `/proc/uptime`     |
| Network  | `/proc/net/dev`    |

**Acceptance criteria:**
- `GET /api/v1/dashboard` returns live values, not stubs.
- Values refresh on each request (no stale caching in Phase 2 — caching is
  a later optimization, not a correctness requirement now).
- Reading fails gracefully (partial data + logged error, not a 500 for the
  whole dashboard) if one metric source is unavailable.

### 6.3 Phase 3 — Server Actions (next)
**Requirement:** Execute a small, explicit whitelist of actions against the
local system and services.

In scope:
- Restart / shutdown the host.
- Start / stop / restart a named systemd service (from an allowed list,
  not arbitrary unit names).

**Acceptance criteria:**
- Every action is permission-scoped (see §8).
- Every action produces an `ActionRun` record with status
  (`pending` → `running` → `success`/`failed`) queryable via API.
- No action handler ever shells out with user-supplied strings
  unsanitized — actions map to a fixed, enumerated set of operations.

## 7. API Requirements

- REST, versioned under `/api/v1`.
- Success responses return the resource directly (no `{status, data}`
  envelope).
- Errors return `{"error": {"code": "...", "message": "..."}}`.
- Payloads stay small — this is a hard constraint because of the ESP32
  client, not a style preference.
- Endpoints required by end of Phase 3:
  - `GET /api/v1/health`
  - `GET /api/v1/dashboard`
  - `GET /api/v1/devices`
  - `GET /api/v1/devices/{id}`
  - `POST /api/v1/devices/{id}/actions/{action}`
  - `GET /api/v1/actions/{runId}` (status of an in-flight/past action)

## 8. Security Requirements (must exist in shape, even before fully enforced)

- Actions are gated by a permission scope check at the API layer, even
  before a full permission system exists — the check point must exist in
  the code path now so it isn't a retrofit later.
- No plaintext secrets in config or logs.
- No remote device (future `helm-remote`/`helm-agent`) is trusted by
  default; auth is required once those components come online, but the
  API surface should already assume a caller identity exists on every
  action-executing request.

## 9. Success Criteria for Current Milestone

- Dashboard reflects real host state at all times.
- At least one destructive-ish action (service restart) works end-to-end:
  request → permission check → execution → status tracked → status
  queryable.
- Adding a second action type (e.g., shutdown) requires no changes outside
  the services/integrations layers — API and domain stay stable.

## 10. Open Questions

- Wire protocol for `helm-remote`/`helm-agent` (WebSocket vs MQTT) —
  deferred, but should be decided before Phase 5 to avoid reshaping the API
  layer late. Leaning WebSocket for simplicity and to avoid a broker
  dependency.
- Where does auth/token state live — flat file, SQLite? Leaning SQLite via
  a pure-Go driver (no cgo) to keep deployment to a single binary.
- Persistence for `ActionRun` history — likely the same SQLite store,
  introduced in Phase 3 rather than deferred further.