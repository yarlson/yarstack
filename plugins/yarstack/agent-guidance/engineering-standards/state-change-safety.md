## State Change Safety

For state changes, identify the source of truth, its invariants, and the state that means the operation succeeded.

Make a multi-step change atomic when one transaction can own it. Otherwise, record enough progress to resume, compensate, or reconcile after partial failure. Do not report success before durable state and side effects meet the contract.

Assume requests and messages may arrive more than once. Make repeated execution safe, or detect the same intent with a stable identifier before causing another side effect. Distinguish a duplicate from a new request that happens to contain the same data.

A timeout, lost response, or cancellation after dispatch may leave the outcome unknown. Check authoritative state before retrying a non-idempotent action or reporting failure. Do not claim exactly-once behavior unless every participating boundary provides it.

Retry only failures that may succeed without another change. Bound attempts and elapsed time. When callers can retry together, use delay and jitter where needed to avoid amplifying load. Preserve cancellation and surface exhaustion according to project conventions.

Define behavior for duplicate, stale, and out-of-order work where the boundary permits them. Before finalizing, test success, partial failure, duplicate delivery, ambiguous outcome, cancellation, retry exhaustion, and recovery while confirming that invariants remain true.
