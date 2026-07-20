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
5. For behavior-changing work, name the authoritative contract, invariants and required progress, environmental assumptions, failure model, and material correctness claims. Map each claim to evidence; use `test-design` when selecting verification methods is non-trivial and `technical-spike` when a high-consequence design may need unfamiliar model checking or proof.
6. For work affecting long-running or distributed workloads, settle the source of truth and invariants; graceful and abrupt restart behavior; rolling replacement and mixed-version overlap; concurrent-replica ownership and coordination; duplicate, stale, or out-of-order work; partial dependency failure, retry exhaustion, ambiguous outcomes, and recovery. Require observable acceptance evidence for the applicable cases and omit the rest.
7. Divide work into implementation-sized phases ordered by real dependencies.
8. Give each phase a goal, concrete deliverables, dependencies, steps, acceptance evidence, and structural traps to avoid.
9. Reread the plan for hidden assumptions, missing prerequisites, contradictory acceptance criteria, failure and recovery gaps, and duplicate ownership.

Use the repository's established plan format. Omit empty ceremony. Name concrete subsystems, files, interfaces, migrations, UI surfaces, or test layers only where they remove implementation ambiguity.

Do not implement, validate, commit, push, or release. Keep a local behavior-preserving structural plan proportionate to its actual risk.

Finish when implementation-changing decisions are settled or explicit blockers, phases are dependency-ordered, and every phase has observable acceptance evidence.
