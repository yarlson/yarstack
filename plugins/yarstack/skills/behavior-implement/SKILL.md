---
name: behavior-implement
description: Implement testable production behavior through a red-green-refactor cycle. Use for bug fixes, features, and behavior-changing refactors, whether or not they belong to a formal plan phase.
---

# Test-Driven Implementation

Own the implementation cycle for one coherent observable behavior.

## Workflow

1. Read the requested behavior, repository instructions, and closest implementation and test patterns.
2. Identify the observable outcome, important side effects, and state preserved on failure.
3. Use `test-design` when selecting the test level, cases, fixtures, or test boundary is non-trivial; otherwise extend the existing test pattern directly.
4. Write the smallest focused test and confirm it fails because the behavior is missing or wrong.
5. Change the minimum production code needed to pass without weakening the assertion.
6. Refactor only as needed for clarity, rerunning the focused test after each behavior-preserving change.
7. Run relevant surrounding tests and repository-required checks.

If a maintainable in-scope automated test is genuinely impractical, state before editing why, what concrete evidence will replace it, and the residual risk. Do not add disproportionate infrastructure; stop when the required harness would materially expand scope.

Do not apply this workflow to purely non-behavioral documentation, formatting, metadata, or mechanical generated-state changes.

Finish when the behavior is proven from a demonstrated red state, or the bounded exception is evidenced, and relevant regression checks pass.
