# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in Helm, please report it responsibly.

**Do NOT open a public GitHub issue for security vulnerabilities.**

### How to Report

Send an email to the project maintainer with:

1. A description of the vulnerability
2. Steps to reproduce the issue
3. Potential impact assessment
4. Any suggested fixes (optional)

### What to Expect

- **Acknowledgment** within 48 hours
- **Assessment** within 7 days
- **Fix timeline** communicated after assessment
- **Credit** given in the release notes (unless you prefer anonymity)

## Security Design Principles

Helm is designed to be exposed to the public Internet. The following principles guide all development:

- **No unrestricted shell access** — Only predefined, validated actions are permitted
- **No trust by default** — All client inputs are validated
- **No plaintext secrets** — Sensitive configuration is handled securely
- **No arbitrary command execution** — Actions are whitelisted, never dynamic

## Supported Versions

| Version | Supported |
|---------|-----------|
| `main` (development) | ✅ |

## Planned Security Features

- [ ] API key authentication
- [ ] Device authentication
- [ ] TLS termination
- [ ] Role-based access control
- [ ] Audit logging
- [ ] Rate limiting
