# AI Gateway Auth

An Envoy ext_authz gRPC service for API key authentication and scope-based authorization in Istio service mesh.

## Overview

AI Gateway Auth validates API keys and enforces scope-based authorization for services behind an Istio gateway.

## Quick Start

```bash
# Build
go build -o bin/server ./cmd/server

# Run tests
go test ./... -race

# Run locally
./bin/server
```

## Architecture

```
Client → Istio Gateway → ext_authz → [Auth Service] → Upstream
                                 ↓
                          Key Cache (in-memory)
                                 ↑
                         AIGatewayKey CRD (watched)
```

## Project Structure

```
cmd/server/        - Application entrypoint
internal/auth/     - Authz Check logic
internal/cache/    - Key cache
internal/config/   - Config loading
internal/crdwatcher/ - AIGatewayKey CRD watcher
```

## Configuration

| Env Var | Default | Description |
|---------|---------|-------------|
| `GRPC_PORT` | 9001 | gRPC port |
| `HTTP_PORT` | 8080 | Health/metrics port |
| `NAMESPACE` | "" | Restrict to namespace |

## Development

```bash
make test
make build
```

See `docs/requirements.md` for full requirements.
