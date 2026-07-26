# Contributing to Helm

Thank you for your interest in contributing to Helm.

## Development Setup

### Prerequisites

- Go 1.26 or later
- Linux environment (or WSL for Windows)
- Git

### Getting Started

```bash
git clone https://github.com/nayiswftw/helm.git
cd helm/helm-core
go mod download
```

### Building

```bash
# Build for Linux
make build

# Cross-compile from macOS/Windows
make build-linux

# Run static analysis
make vet
```

## Code Standards

### Architecture

Helm follows a strict layered architecture. Dependencies always flow downward:

```
API → Service → Integration → Linux/Docker
```

- **API layer** handles HTTP concerns only (routing, parsing, validation, responses)
- **Service layer** contains business logic and orchestration
- **Integration layer** communicates with external systems (Linux, Docker, APIs)

Never violate this dependency direction.

### Go Conventions

- Follow standard Go conventions and idioms
- Use `gofmt` for formatting (enforced by CI)
- Use `go vet` for static analysis (enforced by CI)
- Prefer the standard library over external dependencies
- Keep packages small and focused
- Use explicit dependency injection — no globals or singletons
- Write clear, self-documenting code

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/):

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

**Types:**

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation only |
| `refactor` | Code change that neither fixes a bug nor adds a feature |
| `test` | Adding or correcting tests |
| `chore` | Maintenance tasks, CI, tooling |
| `perf` | Performance improvement |
| `ci` | CI/CD changes |

**Scopes:** `core`, `api`, `service`, `integration`, `config`, `domain`

**Examples:**

```
feat(api): add dashboard endpoint
fix(integration): handle missing /proc/meminfo fields
docs(core): update README with configuration options
chore(ci): add GitHub Actions workflow
```

### Pull Requests

- Keep PRs focused — one logical change per PR
- Include a clear description of what and why
- Ensure CI passes before requesting review
- Update documentation if behavior changes
- Update `CHANGELOG.md` under `[Unreleased]`

## Project Structure

```
helm-core/
├── cmd/helm/          # Application entrypoint
├── internal/
│   ├── api/           # HTTP handlers and routing
│   ├── app/           # Dependency injection container
│   ├── config/        # Configuration loading
│   ├── domain/        # Core domain models
│   ├── integration/   # External system integrations
│   ├── server/        # HTTP server wrapper
│   └── service/       # Business logic
├── Makefile
├── go.mod
└── go.sum
```

## Security

If you discover a security vulnerability, please **do not** open a public issue. See [SECURITY.md](SECURITY.md) for responsible disclosure instructions.

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.
