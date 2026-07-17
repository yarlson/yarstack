---
name: phase-review
description: Perform a strict maintainability review of one completed plan phase or scoped implementation. Use after implementation and validation to find and fix safe behavior-preserving structural problems without broad redesign, shipping work, or unrelated cleanup.
---

# Phase Review

Review the selected implementation aggressively for maintainability while preserving behavior and scope.

## Workflow

1. Read the applicable repository rules and selected contract when context is missing or stale.
2. Inspect the diff, changed files, and surrounding modules.
3. Look for simpler framings that delete branches, helpers, modes, wrappers, or layers.
4. Review ownership, boundaries, types, lifecycle, errors, tests, and documentation.
5. Fix every safe finding within the selected scope.
6. Use `$refactor-plan` before a non-trivial behavior-preserving restructure.
7. Re-run affected checks after fixes.
8. Report real findings that cannot be fixed safely without expanding scope.

## Blocking Findings

Treat these as blockers unless fixed or explicitly out of scope:

- incidental complexity preserved when a clear simpler design exists;
- ad hoc branches, flags, nullable modes, or scattered special cases;
- logic placed outside the module or layer that owns the concept;
- duplicate helpers or near-duplicate abstractions;
- thin wrappers, identity abstractions, generic machinery, or fallbacks that obscure a direct flow;
- unnecessary casts, optionality, loose data shapes, or unclear invariants;
- large-file growth that should be decomposed within the current scope;
- partial-update orchestration with avoidable inconsistent states;
- tests that check implementation trivia instead of behavior;
- unrelated churn or explanatory comments standing in for simpler code.

## Review Questions

- Can the same behavior use fewer concepts, branches, modes, or layers?
- Did the change improve or worsen the local architecture?
- Does each addition have a current consumer and purpose?
- Is logic in its canonical owner and using existing helpers?
- Are failure paths, cleanup, and state transitions explicit?
- Would a behavior-preserving refactor leave the tests valid?

## Guardrails

- Do not add product behavior.
- Do not rewrite a subsystem merely because a cleaner design is imaginable.
- Do not expand into later phases, unrelated cleanup, branches, pull requests, CI, or releases.
- Prefer a few high-confidence structural fixes over cosmetic findings.

Finish when no clear scoped simplification or structural regression remains, checks pass, and unresolved findings are explicit.
