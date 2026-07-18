---
name: refactor-plan
description: Design a minimal reversible sequence for a behavior-preserving refactor, migration, decomposition, rename, or cleanup. Use before editing when structural risk justifies an explicit local plan.
---

# Refactor Plan

Produce a structural plan only; stop before implementation.

## Workflow

1. Identify the requested structural outcome and behavior that must remain unchanged.
2. Inspect current ownership, local patterns, tests, and commands that can detect regressions.
3. Define the smallest mechanical and reversible edit sequence.
4. Identify any unavoidable behavior change and route it to `behavior-implement` or broader planning.
5. Specify focused verification commands to run after implementation.

Report the goal, preserved behavior, likely affected units, ordered sequence, rollback points, risks, and verification commands.

Do not implement, validate, commit, push, combine unrelated cleanup, or add an abstraction without removing concrete complexity. Use `plan-create` for broader multi-phase migrations.

Finish when the reversible sequence and verification contract are explicit.
