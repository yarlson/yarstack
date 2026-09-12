---
name: behavior-implement
description: Implement one observable behavior through proportionate red-green-refactor testing. Use for product changes, and for tooling or infrastructure only when native checks cannot cover a concrete complexity or failure risk with a focused deterministic test.
---

# Test-Driven Implementation

Implement and prove one coherent observable behavior.

Before coding, stop at the first sufficient option: no change, existing code, the standard library or platform, an installed dependency, or one clear line. Otherwise make the minimum change, deleting before adding. Keep one source of truth without joining unrelated behavior. Never trade away validation, data-loss prevention, security, accessibility, or explicit requirements.

## Workflow

1. Read the behavior contract, repository instructions, closest implementation and test patterns, and callers of code you may change.
2. Confirm focused tests are proportionate. For tooling and infrastructure, leave this workflow for direct implementation with native checks unless they cannot prove material behavior, a concrete risk warrants regression coverage, and a focused deterministic boundary fits.
3. Define the observable result, side effects, failure-preserved state, and supported dependency contracts. Make no guarantee beyond their documented behavior.
4. Choose the smallest reliable mechanism. If meeting the contract requires material scope growth, explain the gap and ask for direction instead of approximating the outcome.
5. Use `test-design` when the test boundary, level, cases, or fixtures are non-trivial; otherwise extend the existing pattern.
6. Write the smallest focused test in repository style and confirm it fails because the behavior is missing or wrong. Keep scenario values visible.
7. Write the minimum production code needed to pass without weakening the assertion; refactor only for needed clarity, rerunning the test after each behavior-preserving change.
8. Use `crap-index-assess` when requested or when the repository's CRAP check covers changed methods; treat it as supporting evidence, not a substitute for behavioral tests.
9. Run relevant surrounding tests and repository-required checks.

If an in-scope test proves impractical, state before editing why, what concrete evidence will replace it, and the residual risk. Do not add a harness that materially expands scope.

Do not apply this workflow to purely non-behavioral documentation, formatting, metadata, or mechanical generated-state changes.

Finish when the behavior is proven from a demonstrated red state, or the bounded exception is evidenced, and relevant regression checks pass.
