---
name: phase-validate
description: Validate exactly one implemented phase or scoped change against its authoritative contract. Use after implementation to re-derive criterion-level conclusions from current evidence and fix only required implementation gaps.
---

# Phase Validation

Own contract validation between implementation and structural review.

## Workflow

1. Read the authoritative contract and applicable repository rules unconditionally.
2. Enumerate each observable acceptance criterion and map it to implementation and evidence.
3. Inspect changed files, tests, failure paths, cleanup, side effects, and required checks without trusting implementation claims.
4. Invoke focused review skills only when their trigger materially applies; consume read-only findings and keep fixes under this phase's authority.
5. Use `plan-update` or `spec-update` when the contract is incorrect or incomplete. Do not invent behavior.
6. Fix only implementation gaps required by the selected contract.
7. Run focused checks while correcting, then run the authoritative required command once against final content.

For every criterion report `VERIFIED`, `NOT VERIFIED`, or `INCONCLUSIVE`, with evidence and reasoning. Independence means fresh evidence-based conclusions, not mandatory separate-agent execution.

Do not perform broad maintainability refactoring, final closeout, shipping, commits, pushes, external mutations, or unrelated cleanup. Use `phase-review` after validation for structural regressions.

Finish when every criterion has a supported status, scoped defects are fixed, required checks reflect final content, and blockers are explicit.
