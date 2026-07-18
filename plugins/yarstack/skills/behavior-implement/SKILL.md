---
name: behavior-implement
description: Implement production behavior through a mandatory red-green-refactor cycle. Use before writing or changing testable production code, including bug fixes, features, and behavior-changing refactors, whether or not the work belongs to a formal plan phase.
---

# Test-Driven Implementation

Prove each behavior through a focused test before implementing it.

## Define the Contract

1. Read the requested behavior, repository instructions, and closest existing tests and implementation patterns.
2. Identify the observable outcome, important side effects, and state that must remain unchanged on failure.
3. Select the smallest test level that can prove the behavior through a stable public or system boundary.
4. Cover the normal path and only the boundaries, failures, cleanup, or side effects that materially define the contract.
5. Use `test-design` to design and write the test before changing production behavior.

## Red, Green, Refactor

1. Write the smallest test that expresses the missing behavior.
2. Run it and confirm it fails because the behavior is absent or incorrect, not because setup, syntax, fixtures, or the environment are broken.
3. Change the minimum production code needed to pass the test without weakening the assertion.
4. Run the focused test until it passes.
5. Refactor only where needed to keep the changed code clear, rerunning the focused test after each behavior-preserving change.
6. Run the relevant surrounding tests and repository-required checks.

Do not write production behavior before the failing test. Do not broaden the implementation beyond the behavior established by the test and current task.

## Exceptions

Skip the failing-test-first requirement only when it is genuinely impossible or disproportionate, such as an unavailable external environment, a non-deterministic platform boundary with no practical substitute, or a purely non-behavioral change.

Before changing production code, state:

- why a failing automated test cannot be written first;
- what concrete verification will replace it;
- what residual risk remains.

Lack of an existing test harness, anticipated setup effort, or confidence that a change is simple is not sufficient by itself. Add the smallest harness justified by the current behavior; if that would materially broaden scope, stop and surface the constraint.

Finish when the focused test passed from a demonstrated red state, relevant regression checks passed, and any exception is explicit and evidenced.
