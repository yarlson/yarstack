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
4. Implement the selected phase without skipping ahead.
5. Add only prerequisites that are required for the selected phase to work.
6. Keep the diff surgical and consistent with the codebase.
7. Run checks required by repository instructions and the phase contract.
8. Fix failures within phase scope.
9. Confirm the phase is complete end to end before reporting success.

## Scope Rules

- Do not add unrelated features, flags, workflows, release automation, configuration, abstractions, documentation, or cleanup.
- Update the plan or linked contract only when implementation reveals a concrete defect in it.
- Do not commit or push unless explicitly requested.

Finish only when the selected behavior is complete, required checks ran or are explicitly blocked, and no later-phase work was introduced.
