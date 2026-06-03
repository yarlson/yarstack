---
name: implement
description: >-
  Use when an agent implements code and user wants extremely strict
  maintainability bar: avoid structural regressions, spaghetti branching,
  file-size sprawl, bad abstractions, type-boundary leaks, wrong-layer logic,
  duplicate helpers, brittle orchestration, any trap a strict code-quality
  review block. Trigger: "implement this so strict review finds nothing",
  "build brutally clean", "avoid review traps while coding", "ship code under a
  strict quality bar".
---

# Implement

Use skill while implementing, not after. Goal: make diff boring to review - behavior correct, structure intentional, no obvious maintainability blocker left for strict review.

## Operating Rule

Before editing, decide what make implementation fail brutal maintainability review. Design change so those failure modes never enter diff.

Not permission for broad rewrites. Be ambitious simplifying touched area, keep scope tied to requested behavior.

## Workflow

1. Define success criteria.
   - State requested behavior.
   - State concrete verification command or check.
   - State likely review traps for this change.

2. Read local architecture first.
   - Identify canonical owner of concept being changed.
   - Find existing helpers, models, dispatchers, policy objects, state machines, validators, test patterns before adding new.
   - Check file sizes for touched files when change may add substantial code.
   - Revalidate unfamiliar or version-sensitive framework assumptions from local config, installed versions, official docs.

3. Choose simplest clean shape.
   - Prefer design that deletes concepts, branches, special cases.
   - Prefer existing owned boundary over new parallel path.
   - Prefer small focused extraction over pushing cohesive file into sprawl.
   - Avoid new abstraction unless it removes real complexity or matches local pattern.

4. Implement surgically.
   - Keep behavior changes close to owning layer.
   - Keep transport/UI glue thin.
   - Keep business logic out of adapters unless adapters own that policy.
   - Keep helper APIs typed and explicit.
   - Avoid widening public contracts unless request requires.
   - Delete only complexity created or directly touched by this work.

5. Verify behavior and structure.
   - Run targeted tests or concrete checks.
   - Run formatting/linting when local convention expects.
   - Inspect diff as if reviewing someone else's code.
   - Fix structural blockers before final.

## Design Bar

Pass all checks below before finalizing.

### Structural Simplicity

- Can change be reframed so fewer branches, modes, helpers, layers exist?
- Did implementation delete complexity where practical, not just move it?
- Is code direct enough that final shape feels inevitable?
- Does each new function, type, module, file earn its place?

Blocker patterns:

- Complicated implementation where simpler model delete whole branches.
- Refactor that spreads same complexity across more files.
- Wrapper, helper, service that only forwards calls without clarifying ownership.
- Generic mechanism hiding simple data shape.

Preferred fixes:

- Collapse duplicate branches into one flow.
- Replace flag-heavy flows with small explicit model.
- Delete thin wrappers.
- Move logic to existing owner instead of creating sibling owner.

### Spaghetti Control

- Did diff add ad-hoc conditionals in unrelated paths?
- Are special cases centralized behind right abstraction?
- Are repeated conditionals signaling missing model or helper?
- Is edge-case handling local to concept it belongs to?

Blocker patterns:

- Weird `if` statements scattered through busy functions.
- One-off booleans or nullable modes that complicate normal flow.
- Feature-specific checks inside general-purpose code.
- Temporary branching likely to become permanent.

Preferred fixes:

- Extract pure helper when it names real decision.
- Add small dispatcher/policy object only when it removes scattered branching.
- Make default path handle case naturally.
- Move special-case logic behind owning module boundary.

### File Size And Decomposition

- Will touched file cross 1000 lines because of this diff?
- Is large touched file getting harder to scan?
- Is new code cohesive enough to live in current file?
- Would small module/test-helper extraction make future changes safer?

Blocker patterns:

- File pushed from below 1000 lines to above 1000 lines without strong reason.
- Test file becoming dumping ground for unrelated fixtures.
- Component, command, service accumulating multiple responsibilities.

Preferred fixes:

- Split focused helpers, subcomponents, fixtures, policies.
- Keep orchestration in one place, move reusable rules out.
- Avoid tiny files when one cohesive file still clearer.

### Boundary And Layer Cleanliness

- Is logic in canonical package, service, module, component, layer?
- Does API expose implementation details?
- Does UI/transport glue now know business rules?
- Does shared module now know feature-specific concerns?

Blocker patterns:

- Feature logic leaking into shared utilities.
- Adapter details leaking into domain contracts.
- Business policy buried in request handlers, CLI formatting, UI event code.
- Bespoke helper added where canonical helper already exists.

Preferred fixes:

- Move policy to owner of concept.
- Keep adapters as translation and orchestration surfaces.
- Reuse canonical helpers and local conventions.
- Make ownership visible through names and module placement.

### Type And Contract Clarity

- Did change add avoidable `any`, `unknown`, broad maps, casts, nullable fields, optional params?
- Is fallback behavior masking unclear invariant?
- Can explicit type, enum, result shape, validation boundary simplify control flow?
- Are failure states represented clearly enough for callers?

Blocker patterns:

- Cast-heavy code where shape is known.
- Optional fields used as hidden modes.
- Silent fallback that makes impossible states look valid.
- Broad object contracts passed through multiple layers.

Preferred fixes:

- Introduce smallest explicit type that captures real invariant.
- Validate at boundary.
- Return clear error instead of making callers infer failure.
- Keep casts local and justified when unavoidable.

### Orchestration And Atomicity

- Is independent work serialized for no reason?
- Can related updates leave partial state callers must untangle?
- Does async/control flow hide errors or make retries unclear?
- Is implementation simpler if independent operations run together?

Blocker patterns:

- Long sequential setup where steps are independent.
- Multiple state writes that can fail halfway with no coherent recovery.
- Retry/fallback logic scattered across callers.
- Side effects mixed with decision logic.

Preferred fixes:

- Separate decision logic from side effects.
- Group related updates under one clear operation.
- Parallelize independent work when it simplifies flow.
- Keep retry/fallback policy in one owned place.

## Implementation Tactics

- Read tests before changing behavior; mirror existing high-signal test style.
- Add tests for behavior and failure modes, not implementation trivia.
- Prefer small pure functions for dense decision logic.
- Prefer named intermediate values over nested expressions when intent matters.
- Prefer early returns over deep nesting when local style allows.
- Keep public names precise and boring.
- Remove dead paths introduced by change.
- No feature flags, config knobs, extension points unless request requires.
- No unrelated formatting, docs, generated file changes.

## Pre-Final Review Pass

Before final response, inspect diff and answer:

- Did this fully solve requested behavior?
- Is this simplest clean implementation that fully solves it?
- Did every changed line map to request?
- Did any touched file become too large or less cohesive?
- Did diff add scattered branches or special cases?
- Did any logic land in wrong layer?
- Did any abstraction fail to earn its place?
- Did any type boundary get looser than before?
- Did verification concretely prove behavior?

If any answer exposes blocker, fix before finalizing.

## Final Response

Report:

- What changed.
- What verification ran.
- Any residual risk or out-of-scope issue.

Keep response concise.
