---
name: behavior-implement
description: Implement behavior through a red-green-refactor cycle when focused automated tests are proportionate. Use for product behavior changes and for engineering tooling and infrastructure only when native checks are insufficient and concrete complexity or failure risk warrants dedicated tests.
---

# Test-Driven Implementation

Own the implementation cycle for one coherent observable behavior.

Keep the implementation simple and limited to current requirements. Reuse an existing source of truth when the same rule or knowledge would otherwise be maintained in more than one place, but do not combine unrelated behavior merely because the code looks similar. Add abstractions, flexibility, and configuration only when a concrete current need justifies them.

## Workflow

1. Read the requested behavior, repository instructions, and closest implementation and test patterns.
2. Confirm that focused automated tests are proportionate. For engineering tooling and infrastructure, use this workflow only when native checks cannot credibly prove important behavior, concrete complexity or failure risk warrants regression coverage, and a focused deterministic test boundary fits the task. Otherwise leave this skill and use direct implementation with the applicable native checks.
3. Identify the observable outcome, important side effects, and state preserved on failure. At each dependency boundary, name the supported contract that makes the outcome reliable. Do not base a guarantee on undocumented or incidental behavior, or claim more than the contract supports.
4. Choose the smallest reliable mechanism with complexity and maintenance cost proportionate to the behavior. If no supported mechanism can meet the requirement without material scope growth, stop before implementation, explain the gap, and ask for direction instead of silently approximating the outcome.
5. Use `test-design` when selecting the test level, cases, fixtures, or test boundary is non-trivial; otherwise extend the existing test pattern directly.
6. Write the smallest focused test and confirm it fails because the behavior is missing or wrong. Treat the test as maintained product code. Reuse existing focused helpers, extract cohesive repeated setup without hiding scenario values, and prefer named table-driven cases only when they share one execution path.
7. Change the minimum production code needed to pass without weakening the assertion.
8. Refactor only as needed for clarity, rerunning the focused test after each behavior-preserving change.
9. Use `crap-index-assess` when the request targets CRAP or the repository's configured CRAP check covers changed methods. Treat the result as maintainability evidence, not a replacement for behavioral tests.
10. Run relevant surrounding tests and repository-required checks.

After selecting this workflow, if an in-scope automated test proves impractical, state before editing why, what concrete evidence will replace it, and the residual risk. Do not add disproportionate infrastructure; stop when the required harness would materially expand scope.

Do not apply this workflow to purely non-behavioral documentation, formatting, metadata, or mechanical generated-state changes.

Finish when the behavior is proven from a demonstrated red state, or the bounded exception is evidenced, and relevant regression checks pass.
