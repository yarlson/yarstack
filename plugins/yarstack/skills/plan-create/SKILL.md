---
name: plan-create
description: Create a final implementation-ready plan from a user brief, settled decisions, and repository evidence. Use when work needs ordered phases, concrete deliverables, dependencies, risks, and validation before coding.
---

# Plan Creation

Synthesize an implementation-ready plan rather than repeating discovery, architecture sparring, or research.

## Workflow

1. Confirm the brief, authoritative repository evidence, applicable instructions, output location, and whether writing a repository plan file is authorized.
2. Identify only decisions whose absence changes the implementation path. Resolve them with the user or record a blocker; keep safe assumptions explicit.
3. Use `architecture-refine` for unresolved architecture decisions and `technical-spike` for blocking unfamiliar or version-sensitive behavior.
4. Lock ownership boundaries, non-goals, compatibility, rollout order, validation gates, risks, and unresolved blockers.
5. Divide work into implementation-sized phases ordered by real dependencies.
6. Give each phase a goal, concrete deliverables, dependencies, steps, acceptance evidence, and structural traps to avoid.
7. Reread the plan for hidden assumptions, missing prerequisites, contradictory acceptance criteria, failure and recovery gaps, and duplicate ownership.

Use the repository's established plan format. Omit empty ceremony. Name concrete subsystems, files, interfaces, migrations, UI surfaces, or test layers only where they remove implementation ambiguity.

Do not implement, validate, commit, push, or release. Use `refactor-plan` instead when the requested deliverable is only a local behavior-preserving structural sequence.

Finish when implementation-changing decisions are settled or explicit blockers, phases are dependency-ordered, and every phase has observable acceptance evidence.
