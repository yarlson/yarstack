---
name: completion-review
description: Perform final closeout for a completed implementation or plan phase. Use before reporting completion to compare the result with the requested contract, review the final diff, confirm verification, and separate blockers or unrelated follow-ups without expanding scope.
---

# Final Review

Close the current task honestly and without starting new work.

## Workflow

1. Re-read the request or selected phase contract.
2. Review the final diff and every changed file.
3. Confirm each required behavior is implemented.
4. Confirm meaningful tests and validation ran.
5. Confirm the change stayed within scope and followed repository rules.
6. Confirm directly affected documentation is current.
7. Run any focused security, dependency, rollout, test-gap, or documentation review triggered by the change.
8. Separate blockers and unrelated follow-ups from completed work.
9. Produce a concise final report.

## Guardrails

- Make new edits only when required to complete the requested contract or fix a discovered regression.
- Do not hide failed or skipped checks.
- Do not claim validation that was not performed.
- Do not commit, push, release, or modify external systems unless explicitly requested.

Report what changed, checks run, blockers, remaining risks, and intentionally deferred work. Omit empty sections.
