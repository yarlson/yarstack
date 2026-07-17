---
name: repo-context-docs
description: Create, update, or consolidate current-state project documentation under docs/context/. Use when durable architecture, ownership, runtime, domain, flow, terminology, or engineering-practice context must be synchronized with the implemented codebase.
---

# Repository Context Documentation

Maintain concise current-state documentation under `docs/context/`. Document architecture and behavior, not history or a mirrored file tree.

## Truth And Scope

- Read and write only `docs/context/` unless the user explicitly expands scope.
- Treat source code and entrypoints as primary truth, followed by active runtime configuration, tests, existing context docs, then other prose.
- Treat context claims as working hypotheses; verify non-obvious claims against the code before relying on them.
- Classify units by entrypoints, build targets, runtime roles, deployment, consumers, and contracts—not directory names alone.

## Context Model

Use only files justified by the repository:

- `summary.md`: what the system is, architecture, core flow, state, capabilities, and stack.
- `context-map.md`: terse relative-link index with no topic content.
- `terminology.md`: stable active terms.
- `practices.md`: implemented or enforced conventions and invariants.
- `apps/`, `services/`, `packages/`, `platform/`, `domains/`, and `flows/`: on-demand context for meaningful ownership and behavior units.

Keep `summary.md` and `context-map.md` scannable in seconds. Keep one topic per on-demand file and split only when that improves retrieval.

## Workflow

1. Detect whether the repository is a single project or a multi-unit repository.
2. Read all existing context files before deciding what to change.
3. For a clean repository, inspect manifests, entrypoints, wiring, public contracts, runtime configuration, deployment, and meaningful tests.
4. For a dirty repository, map changed files to affected ownership, domain, platform, and flow concepts before updating only those topics.
5. Create or update topic files before updating `context-map.md`.
6. Remove stale claims and duplicate facts.
7. During consolidation, merge overlapping topics, prune derivable trivia, and rebuild the index from surviving files.
8. Read every edited file and verify its claims against current code.

## Record

- responsibilities and ownership boundaries;
- runtime and deployment roles;
- meaningful contracts, state, invariants, and dependency direction;
- end-to-end flows and cross-cutting domains;
- proven practices that are implemented or enforced.

## Do Not Record

- dates, commits, history, progress, recent completions, next steps, or status reports;
- raw file trees, function signatures, test inventories, line counts, or facts found by a trivial search;
- speculative design, aspirational standards, or TODO lists not enforced by current behavior;
- duplicated source documentation or repository instructions acting as a second system prompt;
- secrets, internal credentials, or sensitive runtime values.

Finish with a short factual summary of context files created, updated, consolidated, or removed.
