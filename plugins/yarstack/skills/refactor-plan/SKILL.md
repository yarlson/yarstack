---
name: refactor-plan
description: Design a minimal safe sequence for a refactor, migration, decomposition, rename, cleanup, or other behavior-preserving structural change. Use before editing when preserving behavior and keeping the change reversible require an explicit local plan.
---

# Refactor Plan

Define the smallest safe structural change before editing.

## Workflow

1. Identify the requested structural outcome and behavior that must remain unchanged.
2. Find existing tests and commands that detect regressions.
3. Inspect the current code shape, ownership boundaries, and local patterns.
4. Define a mechanical, reversible edit sequence.
5. Identify any unavoidable behavior change explicitly.
6. Run relevant checks after implementation.

Use this shape:

```md
## Refactor Plan

Goal: ...
Behavior that must not change: ...
Files likely touched: ...
Safe sequence:

1. ...
2. ...

Verification: ...
```

## Guardrails

- Do not combine unrelated cleanup with the refactor.
- Do not add an abstraction unless it removes concrete complexity in the touched area.
- Do not move code across ownership boundaries without a required reason.
- Do not split or rename files merely to make the diff appear cleaner.

Finish when the requested structure is achieved, behavior remains stable or changes intentionally, and focused checks pass.
