---
name: go-review
description: Review Go code for language-specific correctness and runtime lifecycle risks. Use for Go changes or full-codebase audits involving goroutines, channels, contexts, timers, synchronization, HTTP or database resources, retries, shutdown, or race behavior.
---

# Go Review

Find concrete Go failures that generic code review can miss. Keep the review read-only.

## Workflow

1. Confirm whether the scope is a diff, package, execution path, or full codebase. Read `go.mod`, repository instructions, relevant entry points, callers, tests, and shutdown wiring.
2. Trace every relevant goroutine from creation through cancellation, error propagation, joining, and shutdown. Check for blocked sends or receives, work that outlives its owner, and unbounded fan-out.
3. Check channel ownership, close responsibility, lock ordering, copied locks, mixed atomic and non-atomic access, and concurrent map or slice mutation.
4. Check context propagation, timeout ownership, retry lifetime, ticker and timer cleanup, and version-sensitive timer semantics against the repository's Go version.
5. Check files, response bodies, database rows, transactions, subscriptions, and temporary resources for cleanup on every return and goroutine path.
6. Check typed nil interfaces, unsafe assertions, zero-value assumptions, slice backing-array aliasing, loop capture for the active Go version, panic boundaries, and recover scope.
7. Check HTTP and database timeouts, retry idempotency, transaction commit or rollback ordering, `rows.Err`, pool exhaustion, backoff, and cancellation.
8. Inspect tests for race, cancellation, shutdown, cleanup, retry, and failure coverage. Prefer deterministic coordination or fake time over sleeps. Run focused tests and `go test -race` when appropriate and authorized by the review scope.

## Finding Standard

Report only issues with a concrete failure path. Include the affected location, runtime scenario, impact, smallest safe correction, and verification. State when repository context or tooling leaves a conclusion uncertain.

Use `code-review`, `security-review`, `test-gap-review`, or `refactor-plan` for broader concerns instead of duplicating their checklists.
