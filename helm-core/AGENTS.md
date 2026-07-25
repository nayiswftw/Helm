# Helm Core — Agent Instructions

These are instructions for an AI coding agent working in this repo. This
file governs **how to work**, not **what to build** — for scope and
requirements see `PRD.md`.

## Project Identity

You are working on **helm-core**, the backend control plane of Helm, a
self-hosted homelab control platform. Other clients (ESP32 remote, web,
CLI) consume its API but live in separate components — do not build UI or
firmware logic here.

Helm is NOT a Dokploy remote. Dokploy is one possible integration among
several. The goal is control of the whole machine and the services on it.

## Core Philosophy

Helm should be lightweight, self-hosted, secure, fast, reliable, easy to
deploy, easy to extend.

- Prefer simple, boring, maintainable solutions.
- Do not introduce a framework or dependency unless there's a clear,
  stated benefit. Ask first: *"can this be done cleanly with the Go
  standard library?"*
- No global variables, singletons, or hidden dependencies. All shared
  state goes through the `Application` container.

## Backend Language

Go. Do not migrate to another language. Reasons: low memory footprint,
single-binary deploys, strong concurrency primitives, good fit for a
long-running system-facing service on modest hardware.

## Architecture — Never Violate This Direction

```
            API Layer
                |
                v
          Service Layer
                |
                v
      Integration Layer
                |
                v
         System / External
```

### API Layer
- HTTP routing, request validation, response formatting.
- Must NOT touch the filesystem, execute shell commands, or contain
  business logic.
- Every action-executing endpoint must pass through a permission-scope
  check, even before the full permission system exists (see PRD §8). Do
  not skip this to "add it later."

### Service Layer
- Business logic, orchestration, combining data from integrations.
- Examples: `DashboardService`, `DeviceService`, `ActionService`.
- Services depend on integration interfaces, never on concrete
  integration structs directly, so integrations stay swappable.

### Integration Layer
- Talks to the outside world: Linux system, Docker, Dokploy, network,
  future agents.
- Structure:
  ```
  integrations/
      system/
      dokploy/
      docker/
      network/
  ```
- Anything that reads `/proc`, shells out, or calls an external API lives
  here — never in services or API handlers.

## Current Repo Structure

Maintain this structure unless there's a strong, stated reason to change
it:

```
helm-core/
  cmd/
    helm/
      main.go
  internal/
    app/
      application.go
    api/
      router.go
      handler.go
      response.go
      health.go
      dashboard.go
    domain/
      device.go
      action.go       # ActionRun status tracking — add in Phase 3
    services/
      dashboard.go
      device.go
      action.go       # add in Phase 3
    integrations/
      system/
        system.go
        cpu.go
        memory.go
        disk.go
        uptime.go
    server/
      server.go
```

## Dependency Management

Use an `Application` container for all shared dependencies:

```go
type Application struct {
    Devices   *services.DeviceService
    Dashboard *services.DashboardService
    Actions   *services.ActionService
    System    *system.System
}
```

## Domain Concepts

- **Device** — a controllable/monitorable entity (homelab server, NAS,
  router, desktop).
- **Capability** — something a device can do (power control, metrics,
  containers, processes, notifications).
- **Action** — something Helm can execute (shutdown, restart, deploy,
  wake, start/stop service). Actions are enumerated, never
  arbitrary/user-constructed commands.
- **ActionRun** — a single execution of an action, with status
  (`pending`/`running`/`success`/`failed`), timestamps, and error detail.
  This is the audit trail; every action handler must create and update one.
- **Metric** — CPU, memory, disk, temperature, network, uptime.
- **Integration** — a connection to another system (Linux, Docker,
  Dokploy, future agents).

## API Design

- REST-style, versioned under `/api/v1`.
- Success responses return the resource directly — do not wrap in
  `{"status": "ok", "data": {...}}`.
- Errors: `{"error": {"code": "...", "message": "..."}}`.
- Keep payloads small — ESP32 clients parse these.

## Security Rules

Never:
- expose unrestricted shell execution
- store plaintext secrets
- trust a remote device/caller implicitly

Every action-executing code path must have a permission-check seam, even
while the full permission system is still a stub — this is a structural
requirement now, not a "future" item.

## Coding Rules

Prefer standard library, small packages, explicit code, clear naming.
Avoid overengineering, large abstractions, unnecessary interfaces,
premature optimization.

## AI Agent Behaviour

1. Read `PRD.md` and this file before making changes; understand existing
   architecture first.
2. Do not redesign the project or rename folders without approval.
3. Do not introduce new patterns unless the current ones genuinely can't
   support the task.
4. Keep changes incremental — one phase/feature at a time.
5. Briefly explain any architectural decision in the PR/commit description.
6. Prefer a working feature over a theoretically cleaner one.
7. When a change would cross a layer boundary (e.g., a service reaching
   into `/proc` directly), stop and route it through the integration layer
   instead — do not rationalize an exception.

## Current Priority

Implement real Debian system metrics (Phase 2 — see `PRD.md` §6.2).

After that, in order: server actions with `ActionRun` tracking (Phase 3),
graceful shutdown, configuration, structured logging, then integrations
(Phase 4). Do not jump ahead of this order without explicit approval.