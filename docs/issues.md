# Implementation Issues

## Issue #1: Implement Key Cache
**Priority:** High
**Files:** `internal/cache/cache.go`

Implement a thread-safe, in-memory cache for API keys with O(1) lookups.

**Interface:**
```go
GetByID(id string) *KeyData
GetByValue(value string) *KeyData
Upsert(KeyData)
Delete(id string)
```

**Requirements:**
- Thread-safe (sync.RWMutex)
- Support 1M+ entries
- Benchmark: <1ms for 1M lookups

---

## Issue #3: Implement Authz Check Endpoint
**Priority:** High
**Files:** `internal/auth/handler.go`, `internal/auth/handler_test.go`

Implement `Authorization/Check` gRPC endpoint:
- Extract API key from `x-api-key` or `Authorization: Bearer <key>`
- Lookup key in cache, validate expiry, scopes
- Return 200 OK with headers or 401/403 with JSON error

---

## Issue #4: CRD Watcher
**Priority:** High
**Files:** `internal/crdwatcher/watcher.go`

Watch AIGatewayKey CRDs across all namespaces and sync to cache.

---

## Issue #5: Scope Authorization
Implement scope validation for `workspace`, `model`, `path` scope types.

## Issue #6: Metrics & Health
Expose /healthz, /readyz, /metrics.

## Issue #7: Helm Chart
Package the service for Kubernetes deployment.
