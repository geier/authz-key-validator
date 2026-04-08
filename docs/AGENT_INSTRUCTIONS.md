# Agent Instructions for AI Gateway Auth

This file contains instructions for AI agents working on the AI Gateway Auth project.

---

## Project Overview

AI Gateway Auth is an Envoy ext_authz gRPC service that provides API key authentication and scope-based authorization for Istio service mesh.

**Tech Stack:**
- Go 1.21+
- gRPC (Envoy ext_authz v3)
- Kubernetes CRDs
- In-memory cache with O(1) lookups

---

## Repository Structure

```
.
├── cmd/server/          # Main entry point
├── internal/
│   auth/      # Authz logic
│   cache/      # Key cache
│   config/     # Configuration
│   crdwatcher/ # Watch AIGatewayKey CRDs
├── chart/              # Helm chart
├.docs/      # Requirements & architecture
└── deploy/    # Kubernetes manifests
```

---

## Build, Test, Run

```bash
# Build
go build -o bin/server ./cmd/server

# Test
go test ./... -race

# Lint (if golangci-lint installed)
golangci-lint run

# Run locally
go run ./cmd/server
```

---

## Code Style

- Go standard formatting (`gofmt`, `goimports`)
- Error messages should be actionable.
- All exported functions need doc comments.
- Use `zap` for structured logging.

---

## Pull Requests

- Branch naming: `feature/issue-N-description` or `fix/issue-N-description`.
- Squash-merge into `main`.
- Include `Closes #N` in PR description.
- Run `go test ./... -race` before pushing.

---

## Key Files

| File/Dir | Description |
|----------|-------------|
| `internal/auth` | Core authorization logic |
| `internal/cache` | In-memory key cache |
| `internal/crdwatcher` | Watches AIGatewayKey CRDs |
| `internal/config` | Config loading from env vars |
| `cmd/server/main.go` | Entry point |

---

## Quick Reference

```bash
# Build
make build

# Run tests
make test

# Run locally
make run
```

For detailed requirements, see `docs/requirements.md`.
