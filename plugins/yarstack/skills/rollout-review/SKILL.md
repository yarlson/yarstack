---
name: rollout-review
description: Review operational readiness for changes affecting deployment, infrastructure, runtime configuration, permissions, migrations, observability, background work, dependencies, rollback, recovery, or production failure modes.
---

# Rollout Review

Review only the operational risk introduced by the current change.

## Workflow

1. Identify what changes at runtime.
2. Check required configuration, defaults, and compatibility.
3. Check startup, shutdown, retry, timeout, cancellation, and failure behavior where relevant.
4. Check whether failures are visible through existing logs, metrics, traces, health checks, or user-facing errors.
5. Check rollback, recovery, and partial-state behavior.
6. Check migration and deployment ordering.
7. Update directly affected operational documentation when necessary.
8. Run available preflight and validation commands.
9. Fix operational gaps within the requested scope.

## Required Answers

- What changes at runtime?
- Which configuration or permissions are required?
- How does failure appear?
- How can an operator verify success?
- How can an operator roll back or recover?
- What is the smallest safe deployment order?

Do not add infrastructure, observability systems, release automation, or toy fallback assumptions unless the requested behavior requires them.
