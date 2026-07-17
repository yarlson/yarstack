## Engineering Quality Gate

Solve the requested task while preserving or improving maintainability, correctness, simplicity, testability, and long-term engineering health.

### Core rules

- **Keep the change small.** Modify only what the task requires. Do not refactor unrelated code, change public interfaces, data models, configuration, or operational behavior, or add speculative abstractions.
- **Preserve correctness.** Handle normal, boundary, and failure paths. Watch for partial state, ordering and concurrency bugs, nil or null values, stale reads, unsafe assumptions, and off-by-one errors. Make invalid states hard to represent where practical.
- **Keep responsibilities clear.** Avoid cleverness, hidden coupling, unnecessary indirection, and mixed concerns. Keep code boring and easy to modify.
- **Fit the codebase.** Follow existing conventions for structure, naming, errors, logging, testing, configuration, dependency wiring, and API shape. Do not introduce a second pattern without a clear reason.
- **Control complexity.** Prefer short units, shallow control flow, guard clauses, and explicit data flow. Add concurrency or complex runtime patterns only when a simpler design is insufficient.
- **Test the contract.** Add or update deterministic tests for changed behavior, important boundaries, failure paths, and cleanup. Do not weaken existing tests without a clear justification.
- **Surface failures.** Return, wrap, log, or expose errors according to project conventions. Add logs, metrics, traces, or health signals only when useful. Keep failures diagnosable without leaking secrets or sensitive data.
- **Protect trust boundaries.** Validate untrusted input and preserve authentication, authorization, tenant isolation, and permissions. Avoid injection, unsafe paths or deserialization, SSRF, secret exposure, insecure defaults, and unnecessary privilege or access.
- **Own resources.** Avoid inefficient work on hot paths and unbounded memory growth. Close, cancel, clean up, or release files, connections, timers, tasks, threads, and goroutines in the correct order. Bound concurrent or background work.
- **Avoid needless dependencies.** Prefer the standard library and existing dependencies. Add a dependency only when it materially reduces complexity or risk, and consider its maintenance, license, size, security, and transitive cost.

### Before finalizing

Review every changed file and confirm:

1. Every change is required for the task.
2. Complexity did not grow unnecessarily.
3. Changed units remain small, readable, and focused.
4. Errors, edge cases, cleanup, and cancellation are handled.
5. Tests cover the meaningful behavior.
6. The implementation follows existing conventions.

Fix any failed gate before claiming the work is complete.

### Report the result

Include:

- what was implemented
- why this is the smallest safe approach
- quality risks considered
- tests and checks run
- intentionally deferred cleanup
- remaining risks or follow-up work
