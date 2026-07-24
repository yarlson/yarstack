---
name: phase-review
description: Review one validated phase for behavior-preserving structural regressions introduced or materially worsened by the change. Use after phase validation and before final delivery.
---

# Phase Structural Review

Own a narrow maintainability pass after correctness has been established.

## Workflow

1. Read the selected contract, applicable repository rules, final diff, and changed units.
2. Identify only complexity, duplication, unclear ownership, loose types, or unnecessary indirection introduced or materially worsened by the selected change.
3. Use `crap-index-assess` when the repository's configured CRAP check covers changed methods or the phase was intended to reduce their CRAP risk. Treat score changes as supporting evidence, not a structural verdict.
4. Fix a finding only when the parent implementation request authorizes scoped edits and the correction is clearly behavior-preserving.
5. Report a worthwhile correction when it requires non-trivial restructuring or broader authority; do not expand scope during the review.
6. Rerun affected checks after every edit.

Correctness belongs to `phase-validate`; general read-only defect review to `code-review`; test sufficiency to `test-gap-review`; documentation accuracy to `docs-drift-review`; and specialized risks to their owning reviewers.

Do not add behavior, pursue unrelated cleanup, use subjective file-size thresholds, or continue until no imaginable simplification remains.

Finish when structural regressions introduced by the change are fixed or reported with their scope and risk.
