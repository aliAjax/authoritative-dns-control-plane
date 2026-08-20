# Operations Runbook

1. Check `/healthz` and `/readyz` before changing traffic.
2. Create or update records in `draft` state and call `validate`.
3. Review the returned digest and warnings, then call `publish`.
4. Inspect snapshots and use `rollback` with an explicit immutable snapshot ID.
5. Reconcile Anycast intents only after node leases and BGP state are current.
6. Keep DNSSEC KSK and ZSK keys in an external secret store; activate a prepublished key only after propagation.

The in-memory profile is intentionally ephemeral and is not suitable for a multi-instance deployment. Use a transactional serial allocator and durable outbox before production rollout.

