---
name: test-gap-review
description: Review tests and verification for missing behavior coverage or false confidence. Use during validation when changed behavior may lack meaningful normal, boundary, failure, cleanup, or side-effect coverage, or when required checks were skipped or claimed without evidence.
---

# Test Gap Review

Determine whether verification proves the promised behavior rather than implementation trivia.

## Workflow

1. Identify the observable behavior and invariants that must be proven.
2. Inspect existing tests and checks for that behavior.
3. Look for missing normal paths, important boundaries, failure paths, preserved state, cleanup, and external side effects.
4. Reject weak assertions, fixtures that cannot fail meaningfully, implementation-detail tests, unsupported manual claims, and skipped required commands.
5. Add or strengthen tests only where they materially improve confidence.
6. Run the relevant checks.
7. Report any remaining verification gap and residual risk.

## Verification Preference

Prefer, in order:

1. the repository's authoritative test command;
2. validation required by the current contract;
3. focused unit or integration tests;
4. focused manual verification with concrete evidence;
5. an explicit blocker report.

For measurable claims, report:

```text
VERIFIED | NOT VERIFIED | INCONCLUSIVE
Claim: ...
Evidence: ...
Reasoning: ...
```

Use `INCONCLUSIVE` when the baseline is invalid, the signal is noisy, the environment differs materially, or the check failed for an unrelated reason.

Do not chase coverage percentages, add fake-confidence tests, snapshot unstable output without an established reason, or broaden test infrastructure unnecessarily.
