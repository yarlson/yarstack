---
name: go-review
description: Supplement code review with Go-specific correctness and lifecycle analysis. Use alongside code-review whenever the requested scope contains meaningful Go production or test code, a Go package, or a full Go codebase; skip only when no meaningful Go code is in scope.
---

# Go Review

Run the required Go-specific pass alongside general review and find runtime failures it is unlikely to catch. Keep the review read-only.

## Workflow

1. Confirm the scope and Go version, then read relevant entry points, callers, tests, and shutdown wiring.
2. Trace goroutine ownership, cancellation, error propagation, joining, blocked operations, and fan-out bounds.
3. Check channel close ownership, lock ordering, atomics, concurrent collection access, contexts, timeouts, retries, tickers, and timers.
4. Check cleanup of files, bodies, rows, transactions, subscriptions, and temporary resources on every relevant path.
5. Check Go-specific nil, assertion, zero-value, slice-aliasing, loop-capture, panic, HTTP, database, and version-sensitive semantics only where the changed flow uses them.
6. Run trusted focused tests or `go test -race` only when authorized and useful for a concrete risk.

Use `code-review` for general diff ownership, `security-review` for trust boundaries, `test-gap-review` for comprehensive coverage, and `rollout-readiness-review` for deployment readiness. Their invocation does not authorize edits or unsafe execution.

Report affected location, runtime scenario, impact, smallest safe correction, verification, and uncertainty. State directly when no Go-specific finding exists.

Finish when the selected Go lifecycle surface is exhausted and evidence limits are explicit.
