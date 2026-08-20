# Architecture

The control plane uses a zone aggregate backed by an in-memory repository for local development. `snapshot_release` acquires a per-zone publish lock, validates every record, allocates a monotonic serial, computes a deterministic digest, and stores an immutable snapshot. Rollback creates a new snapshot and serial rather than mutating an existing snapshot.

The resolver parses DNS headers and labels with bounded compression-pointer recursion. Malformed packets are rejected before policy evaluation. A bounded cache stores positive and negative responses with TTL expiry. UDP requests use a concurrency limiter; TCP framing and DoH POST use the same resolver path.

Anycast nodes are accepted as provider-neutral health reports. `Plan` fences nodes with expired leases, low availability, high latency, or exhausted capacity. Fencing tokens make reconciliation idempotent and prevent stale advertisements.

## Production adapters

The repository and event interfaces are deliberately independent from PostgreSQL, Redis, and object storage. Deployments can provide adapters without changing domain invariants. TLS and gRPC endpoints should be enabled behind a certificate manager in production.

