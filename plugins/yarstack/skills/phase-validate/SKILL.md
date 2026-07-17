---
name: phase-validate
description: Validate exactly one implemented plan phase or scoped change against its stated contract. Use after implementation to independently inspect behavior, tests, edge cases, documentation, security-sensitive surfaces, and required checks, fixing only scoped gaps.
---

# Phase Validation

Judge the implementation against the selected contract rather than trusting that code presence means completion.

## Workflow

1. Read applicable repository rules and the selected contract when context is missing or stale.
2. Inspect changed files and relevant surrounding modules.
3. Compare observable behavior with the contract.
4. Check normal behavior, important boundaries, meaningful failure paths, preserved state, cleanup, and external side effects.
5. Inspect test quality, documentation drift, security-sensitive surfaces, dependency changes, and skipped required checks.
6. Fix only gaps required to complete the selected scope.
7. Strengthen tests when existing checks do not prove the promised behavior.
8. Run required and targeted checks after fixes.

## Failed Checks

When compilation, type checking, or validation fails:

1. Identify the failing command.
2. Group errors by file and category.
3. Fix the highest-confidence scoped failures first.
4. Re-run the same command after each focused correction.
5. Stop when the remaining failure needs unrelated work or missing external state.

If a tool, service, credential, fixture, or environment is unavailable, report exactly what remains unverified and the resulting risk.

Finish only when the contract is verified by concrete evidence, scoped defects are fixed, and blockers are explicit.
