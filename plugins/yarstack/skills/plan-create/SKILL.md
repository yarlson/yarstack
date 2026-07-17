---
name: plan-create
description: Create a final repository-local implementation plan from a user brief and codebase evidence. Use for work that needs decisions, non-goals, risks, ordered implementation phases, concrete deliverables, dependencies, and validation before coding.
---

# Plan Creation

Produce an implementation-ready plan, not a brainstorm transcript or roadmap summary.

## Workflow

1. Inspect the repository enough to understand its stack, boundaries, and existing patterns.
2. Read applicable repository instructions and existing plans or context documents.
3. Ask one sharp question at a time until implementation no longer depends on guessing.
4. Prefer concrete choices with meaningful tradeoffs over broad open-ended questions.
5. Keep decisions in working memory; do not create draft plan files during questioning.
6. Resolve discoverable facts from the repository before asking the user.
7. Use current official sources when a decision depends on version-sensitive external behavior.
8. Lock ownership boundaries, runtime and data models, persistence, trust boundaries, rollout order, validation gates, non-goals, risks, and unresolved decisions.
9. Identify review traps before phasing: wrong-layer logic, file sprawl, scattered branching, loose contracts, duplicate ownership, missing tests, and documentation drift.
10. Write the final plan only when it is ready to implement or when remaining blockers are explicit.

## Plan Contract

Use the repository's established format. If none exists, use:

```md
# Plan

## Decisions

- ...

## Non-Goals

- ...

## Open Risks

- ...

## Phase 1 - Name

Goal: ...

Deliverables:

- ...

Dependencies:

- ...

Unresolved decisions:

- ...

Steps:

1. ...

Validation:

- ...
```

## Quality Bar

- Make phases implementation-sized and ordered by real dependencies.
- Ground each phase in current repository evidence and name the canonical owner of changed behavior.
- Name concrete subsystems, files, APIs, models, migrations, UI surfaces, or test layers.
- Separate work with different ownership, rollout risk, or validation methods.
- Name exact checks when they are discoverable.
- Give each phase observable acceptance criteria, dependencies, and specific structural traps to avoid.
- Flag likely file-size, type-boundary, migration, compatibility, and recovery risks before implementation.
- Reject umbrella phases such as “build backend,” “add UI,” “wire everything,” or “add tests.”
- Preserve unresolved decisions explicitly instead of hiding them behind vague verbs.

Do not implement, commit, push, or run release workflows while creating the plan unless explicitly requested.
