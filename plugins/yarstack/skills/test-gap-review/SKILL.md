---
name: test-gap-review
description: Review whether existing tests and verification credibly prove a scoped behavior contract. Use during validation when evidence may omit important behavior or provide false confidence.
---

# Test Gap Review

Audit verification sufficiency without silently becoming a test-writing or production implementation workflow.

## Workflow

1. Establish the authoritative behavioral contract, comparison scope, and claimed verification.
2. Map invariants, assumptions, normal behavior, important boundaries, meaningful failures, preserved state, cleanup, and external effects to existing evidence.
3. Identify weak assertions, fixtures that cannot fail meaningfully, implementation-detail tests, unsupported manual claims, and skipped required commands. Flag example-only evidence when the claim spans a broad input space, operation sequence, concurrent ordering, restart, or recovery, and check that generated failures are reproducible.
4. Review changed test code as maintained product code. Flag unclear naming or structure, repeated setup that should use a focused helper, helpers that hide scenario data or mix responsibilities, and table-driven cases that do not share one execution path.
5. Report each material gap with the affected claim, missing evidence, consequence, smallest proving check, and residual risk.
6. If remediation is authorized, use `test-design` for test changes and `behavior-implement` for production behavior; otherwise remain read-only.

For each material claim report:

```text
VERIFIED | NOT VERIFIED | INCONCLUSIVE
Claim: ...
Evidence: ...
Reasoning: ...
```

Do not chase coverage percentages, demand a stronger method without a concrete risk, add fake-confidence tests, broaden infrastructure, or change production behavior directly.

Finish when material claims have supported statuses and remaining verification risk is explicit.
