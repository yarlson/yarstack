---
name: behavior-implement
description: Implement behavior through a red-green-refactor cycle when focused automated tests are proportionate. Use for product behavior changes and for engineering tooling and infrastructure only when native checks are insufficient and concrete complexity or failure risk warrants dedicated tests.
---

# Test-Driven Implementation

Own the implementation cycle for one coherent observable behavior.

## Workflow

1. Read the requested behavior, repository instructions, and closest implementation and test patterns.
2. Confirm that focused automated tests are proportionate. For engineering tooling and infrastructure, use this workflow only when native checks cannot credibly prove important behavior, concrete complexity or failure risk warrants regression coverage, and a focused deterministic test boundary fits the task. Otherwise leave this skill and use direct implementation with the applicable native checks.
3. Identify the observable outcome, important side effects, and state preserved on failure.
4. Use `test-design` when selecting the test level, cases, fixtures, or test boundary is non-trivial; otherwise extend the existing test pattern directly.
5. Write the smallest focused test and confirm it fails because the behavior is missing or wrong.
6. Change the minimum production code needed to pass without weakening the assertion.
7. Refactor only as needed for clarity, rerunning the focused test after each behavior-preserving change.
8. Run relevant surrounding tests and repository-required checks.

After selecting this workflow, if an in-scope automated test proves impractical, state before editing why, what concrete evidence will replace it, and the residual risk. Do not add disproportionate infrastructure; stop when the required harness would materially expand scope.

Do not apply this workflow to purely non-behavioral documentation, formatting, metadata, or mechanical generated-state changes.

Finish when the behavior is proven from a demonstrated red state, or the bounded exception is evidenced, and relevant regression checks pass.
