---
name: phase-implement
description: Implement exactly one selected phase from an authoritative repository plan. Use when a named phase or equivalent scoped implementation unit must be completed without later-phase work.
---

# Phase Implementation

Own phase scope and orchestrate only the work required to implement it.

Keep the implementation simple and limited to current requirements. Reuse an existing source of truth when the same rule or knowledge would otherwise be maintained in more than one place, but do not combine unrelated behavior merely because the code looks similar. Add abstractions, flexibility, and configuration only when a concrete current need justifies them.

## Workflow

1. Read the selected phase, applicable repository instructions, and the minimum local evidence needed to identify files, patterns, and required checks.
2. Confirm phase boundaries, prerequisites, acceptance criteria, and unchanged behavior.
3. Use `plan-update` for proven plan-mechanics defects and `spec-update` for non-plan semantic contract defects.
4. Use `behavior-implement` for product behavior. For engineering tooling and infrastructure, use it only when native checks cannot credibly prove important behavior, concrete complexity or failure risk warrants regression coverage, and a focused deterministic test boundary fits the task. Otherwise implement through the repository-native path and run its applicable checks without adding test infrastructure.
5. Implement only the phase and prerequisites required for it to work.
6. Run phase-required and focused checks, fixing only scoped implementation failures.
7. Report implemented behavior, changed scope, checks, alternative evidence, and remaining blockers.

Do not perform independent post-implementation validation, structural review, final closeout, commits, pushes, releases, later phases, or unrelated cleanup. Those remain separate gates unless the user explicitly requests them.

Finish at implementation evidence when the selected phase is complete, required checks pass or are explicitly blocked, and no later-phase work was introduced.
