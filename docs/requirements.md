# AI Gateway Auth - Requirements Specification

## 1. Overview

AI Gateway Auth is an Envoy ext_authz gRPC service providing API key authentication and scope-based authorization for Istio service mesh.

---

## 2. Functional Requirements

### F1 — External Authorization Server

- Implements `envoy.service.auth.v3.Authorization/Check` gRPC endpoint
- Extracts API key from:
  - `x-api-key` header (preferred)
  - `Authorization: Bearer <key>` (fallback)
- **Allow**: Return OK with `x-api-key-id` and `x-api-key-owner` headers
- **Deny**: Return JSON `{"error": "<message>"}` with appropriate HTTP status

### F2 — Key Validation Pipeline

1. **Lookup** — Find key by value → 401 if not found
2. **Enabled check** → 403 if disabled
3. **Expiration check** → 403 if expired
4. **Scope check** → 403 if unauthorized

### F3 — Scope Authorization

| Scope Type | Behavior |
|------------|----------|
| `workspace` | Grants access to all resources in owner's namespace |
| `model` | Grants access to specific model (scope.value matches path model) |
| `path` | Grants access to specific path patterns |

### F4 — CRD: AIGatewayKey

```yaml
apiVersion: prokube.ai/v1alpha1
kind: AIGatewayKey
metadata:
  name: aigk-<uuid>
spec:
  displayName: "my-key"
  description: "Optional description"
  owner: "user@example.com"
  enabled: true
  expiresAt: "2025-12-31T23:59:59Z"  # optional
  scopes:
    - type: workspace
    - type: model
      value: "gpt-4"
  secretRef:
    name: "my-secret"
    key: "token"
```

### F5 — Controller Behavior

- **Startup**: List all AIGatewayKey resources, populate cache
- **Watch**: Subscribe to changes, incrementally update cache
- **Retry**: Exponential backoff on API failures
- **Cache**: keyID → KeyData, keyValue → keyID reverse lookup

### F6 — In-Memory Cache

Thread-safe LRU cache with O(1) lookups:
- `keyID → KeyData`
- `keyValue → keyID` (reverse lookup)

### F7 — Health & Metrics

- `GET /healthz` — liveness
- `GET /readyz` — readiness
- `GET /metrics` — Prometheus metrics
- `grpc.health.v1.Health` — gRPC health check

---

## 3. Non-Functional Requirements

| ID | Requirement |
|----|-------------|
| NF1 | < 1ms p99 latency for key lookup (in-memory O(1)) |
| NF2 | Graceful shutdown within 15s on SIGTERM |
| NF3 | Structured JSON logs (zap) |
| NF4 | Run as non-root, read-only root filesystem |
| NF5 | Scale to 1M keys (~2KB each in memory) |

---

## 4. Data Model

```go
type KeyData struct {
    ID          string
    Value       string
    Owner       string
    Namespace   string
    Enabled     bool
    ExpiresAt   *time.Time
    Scopes      []Scope
    CreatedAt   time.Time
}

type Scope struct {
    Type  string  // "workspace" | "model" | "path"
    Value *string // nil for workspace type
}
```

---

## 5. Metrics

| Metric | Type | Labels |
|--------|------|--------|
| `auth_checks_total` | Counter | result, reason, source |
| `auth_check_duration_seconds` | Histogram | result |
| `cache_keys_total` | Gauge | state |
| `watcher_reloads_total` | Counter | result |

---

## 6. Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `GRPC_PORT` | 9001 | gRPC port |
| `HTTP_PORT` | 8080 | Health/metrics port |
| `LOG_LEVEL` | info | debug, info, warn, error |
| `NAMESPACE` | (none) | Restrict to namespace |
| `METRICS_PATH` | /metrics | Prometheus endpoint |

---

## 7. RBAC Requirements

```yaml
rules:
- apiGroups: ["aigatewaykey.prokube.ai"]
  resources: ["aigatewaykeys"]
  verbs: ["get", "list", "watch"]
- apiGroups: [""]
  resources: ["secrets"]
  verbs: ["get"]
```

---

**Status:** Draft  
**Version:** 0.1.0  
**Updated:** March 2025
