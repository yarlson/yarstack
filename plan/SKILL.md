---
name: plan
description: >-
  Use when an agent writes or rewrites implementation plans, PLAN.md files,
  roadmap phases, sprint-sized technical plans, or execution plans under an
  extremely strict maintainability bar. Trigger when user wants plans that make
  later implementation avoid strict code-quality review traps:
  structural regressions, spaghetti branching, wrong-layer logic, file sprawl,
  bad abstractions, loose contracts, missing tests, stale docs, or vague phases.
---

# Plan

Use skill write plans guide implementation toward clean code before code exist. Output make later strict code-quality review boring.

## Operating Rule

Plan from evidence, not vibes. Each phase small enough implement one focused pass, big enough deliver real value, explicit enough agent execute without inventing architecture.

No strategy theater. Make build order, boundaries, tests, risks, acceptance clear.

## Workflow

1. Ground plan in current repo.
   - Read existing plan/docs first when present.
   - Inspect code/config/tests for live architecture, canonical owners.
   - Re-check version-sensitive framework/library assumptions from local config or official docs.
   - Identify user journey, API contract, runtime behavior, or operator workflow being changed.

2. Define success.
   - State business/user value concrete.
   - State target behavior.
   - State non-goals.
   - State validation proof: commands, screenshots, API checks, migrations, CI checks, logs, or manual smoke paths.

3. Find structural traps before phasing.
   - Where logic land wrong layer?
   - Which files risk crossing 1000 lines?
   - Which flow risk scattered conditionals?
   - Which contract become `any`, `unknown`, map-shaped, or nullable-mode soup?
   - Which helper already exist, should reuse?
   - Which docs drift from implementation?

4. Split into phases.
   - Number phases from 1.
   - Keep phases PR-sized or sprint-sized as user asked.
   - Make each phase independently reviewable, verifiable.
   - Order by dependencies, not topic aesthetics.
   - Prefer thin vertical slices when behavior must prove end to end.
   - Split broad phases until one agent implement phase one go.
   - Merge tiny tasklets only make sense together.

5. Add acceptance and checks.
   - Every phase need acceptance criteria.
   - Every phase need concrete verification.
   - Browser-visible phases need Playwright/E2E or visual screenshot proof when local project supports.
   - Data/migration phases need rollback or compatibility checks.
   - API/CLI phases need contract and failure-mode checks.

6. Run plan review pass.
   - Check phase sizes.
   - Check dependency order.
   - Check missing tests.
   - Check docs impact.
   - Check structural traps.
   - Fix vague, broad, or unimplementable phases before final.

## Required Plan Shape

Use this shape unless repo has stronger existing template:

```markdown
# <Plan Title>

## Goal

<Concrete outcome and user/business value.>

## Evidence

- <Current code/doc/config fact with path or command.>
- <Current behavior or constraint.>

## Non-Goals

- <Explicitly excluded adjacent work.>

## Quality Bar

- <Architecture boundary rule.>
- <File-size/decomposition rule.>
- <Type/contract rule.>
- <Test/validation rule.>

## Phases

### Phase 1: <Outcome Name>

**Business value:** <Why this matters.>

**Implementation:** <Specific files/layers/behavior, not line-by-line code.>

**Acceptance criteria:**

- <Observable behavior.>
- <Verification command/check.>

**Dependencies:** <Earlier phase, external setup, or "None".>

**Review traps to avoid:**

- <Wrong-layer/spaghetti/file-size/type/test trap.>

## Plan Review

- **Phase sizing:** <Pass/fix notes.>
- **Dependency order:** <Pass/fix notes.>
- **Architecture:** <Pass/fix notes.>
- **Validation:** <Pass/fix notes.>
- **Docs impact:** <Pass/fix notes.>
```

## Phase Sizing Rules

- Too broad: phase spans multiple ownership layers without one testable behavior.
- Too broad: phase says "implement provider", "build CLI", "add dashboard", "refactor API".
- Too broad: phase needs many unrelated files, no clear review boundary.
- Too tiny: phase only renames, moves one helper, or adds plumbing with no independent proof.
- Right size: phase changes one user-visible behavior, one contract slice, one provider capability, one workflow step, or one refactor needed before next behavior.
- Right size: phase reviewable without understanding entire roadmap.
- Right size: phase has concrete acceptance criteria, likely fits one focused implementation pass.

## Quality Bar

### Architecture

- Name canonical owner for each concept.
- Keep business logic out of transport/UI glue.
- Keep feature-specific logic out of shared utilities.
- Put provider/backend/runtime behavior in owning module, not scattered callers.
- Avoid parallel ownership paths.

### Structure

- Flag files likely cross 1000 lines.
- Require decomposition before sprawl when clear.
- Prefer focused helper/module extraction over giant files.
- Avoid over-splitting into tiny files that hide simple flow.

### Complexity

- Plan remove branches when possible.
- Avoid scattered special cases.
- Replace mode booleans/nullables with explicit models where useful.
- Keep orchestration separate from decision logic.

### Types And Contracts

- Prefer explicit data shapes over loose maps.
- Validate at boundaries.
- Avoid casts/fallbacks that mask invariants.
- Call out public API or schema changes.

### Tests And Proof

- Specify behavior tests, failure-mode tests, smoke checks.
- No fake-confidence tests.
- Browser-visible work needs E2E or screenshot artifact expectations when feasible.
- Docs-only phases need link/path consistency checks.

## Writing Rules

- Be concrete about paths, commands, owners, proof.
- Use current repo facts; do not invent architecture.
- Keep implementation detail enough for execution, not full code.
- Include non-goals to prevent scope creep.
- Mention docs updates when behavior, config, setup, CLI, API, or operator flow changes.
- Preserve existing plan conventions when stronger than this template.
- If user says save to file, write file and verify on disk.
- If user asks only chat plan, do not write file.

## Anti-Patterns

- "Phase 1: Refactor everything."
- "Add tests" with no test type, command, or behavior.
- "Improve architecture" with no owner or boundary.
- "Implement UI" with no route, state, form, or E2E proof.
- "Update docs" without naming docs.
- "Handle errors" without naming failure modes.
- "Create abstraction" before proving duplication or complexity.
- Planning around future fantasies instead of requested behavior.

## Pre-Final Review Pass

Before final, answer:

- Does plan start from current evidence?
- Are phases numbered from 1?
- Is each phase implementable in one focused pass?
- Is each phase independently reviewable?
- Is dependency order correct?
- Does each phase include business value, implementation, acceptance criteria, dependencies, traps?
- Are architecture boundaries explicit?
- Are tests/checks concrete?
- Are docs impacts captured?
- Are non-goals clear?

Fix plan if any answer fails.

## Final Response

Report:

- Plan path if saved.
- What evidence grounded plan.
- Validation/check performed.
- Any residual uncertainty.

Keep concise.
