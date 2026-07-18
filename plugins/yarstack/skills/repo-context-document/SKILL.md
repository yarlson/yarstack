---
name: repo-context-document
description: Create, update, or consolidate adopted current-state documentation under docs/context/. Use for durable implemented architecture, ownership, runtime, domain, flow, terminology, or engineering-practice context.
---

# Repository Context Documentation

Own the information architecture and current-state content under an adopted `docs/context/` system.

## Workflow

1. Confirm requested topics, audience, repository adoption of `docs/context/`, dirty-worktree preservation, and authority for moves or removals.
2. Read all existing context files and the minimum repository evidence needed to verify their claims.
3. Classify units by entrypoints, runtime roles, deployment, consumers, state, and contracts rather than directory names.
4. Create or update only justified topic files before updating the terse context index.
5. Remove stale claims and duplicate facts; merge or remove files only when authorized.
6. Verify every edited claim against current implementation.
7. Use `english-text-review` on edited prose and apply material findings without weakening architectural precision or established terminology.

Write only under `docs/context/`. Use `docs-drift-review` for incidental change-caused drift, `docs-review` for general documentation, and `critical-journey-document` for actor-goal journeys. Link to canonical plans, specifications, runbooks, guides, and API references instead of duplicating them.

Record implemented responsibilities, ownership, runtime roles, contracts, state, invariants, dependency direction, and concise system flows. Exclude history, progress, raw file inventories, speculative design, TODO lists, duplicated instructions, and sensitive values.

Do not impose `docs/context/` on a repository that has not adopted it without an explicit request.

Finish with the context files created, updated, consolidated, or removed and the evidence used.
