---
name: phase-implement
description: Implement exactly one selected phase from a repository plan. Use when the user identifies a PLAN.md phase or equivalent scoped implementation unit and expects surgical implementation plus verification without later-phase work.
---

# Phase Implementation

Implement exactly one selected plan phase.

## Inputs

Establish:

- the selected phase or scoped unit;
- the authoritative plan or contract document;
- applicable repository instructions.

## Workflow

1. Read repository instructions, the selected phase, and directly linked context.
2. Map the files, tests, commands, and local patterns relevant to that phase.
3. State assumptions only when they affect the implementation.
4. If the plan requires complexity that current evidence does not justify, identify the exact defect and correct it through `plan-update`; do not silently redesign or skip the requirement.
5. Before changing production behavior, use `behavior-implement` for each coherent behavior change. Do not proceed without a demonstrated failing test unless that skill's exception is stated first with alternative verification and residual risk.
6. Implement the selected phase without skipping ahead.
7. Add only prerequisites that are required for the selected phase to work.
8. Keep the diff surgical and consistent with the codebase.
9. Run checks required by repository instructions and the phase contract.
10. Fix failures within phase scope.
11. Confirm the phase is complete end to end before reporting success.

## Scope Rules

- Do not add unrelated features, flags, workflows, release automation, configuration, abstractions, documentation, or cleanup.
- Keep phase scope and completion decisions here; `behavior-implement` owns only the red-green-refactor cycle for behavior inside the phase.
- Update the plan or linked contract only when implementation reveals a concrete defect in it.
- Do not commit or push unless explicitly requested.

Finish only when the selected behavior is complete, required checks ran or are explicitly blocked, and no later-phase work was introduced.
