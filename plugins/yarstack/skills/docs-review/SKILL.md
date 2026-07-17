---
name: docs-review
description: Audit, restructure, or improve repository documentation using implemented behavior as the source of truth. Use for deliberate README or documentation-quality work beyond drift caused by a current code change.
---

# Documentation Review

Make the requested documentation accurate, useful, discoverable, and proportionate to the project.

## Workflow

1. Confirm the audience, requested document set, and whether the task is an audit, rewrite, consolidation, or removal pass.
2. Inspect the implementation, commands, configuration, examples, tests, package metadata, and existing canonical documents that support the affected claims.
3. Identify duplicated, stale, unsupported, misplaced, missing, or needlessly detailed content. Classify changes as `keep`, `update`, `merge`, `split`, `move`, `archive`, or `remove`.
4. Edit only the requested or clearly affected canonical documents. Preserve the repository's structure and voice unless that structure is the problem being fixed.
5. For a README, make the opening quickly establish what the project is, who it serves, why it exists, and the shortest verified path to a useful result. Add sections only when the project has real content for them.
6. Verify every changed command, flag, configuration key, file path, package name, example, compatibility statement, and capability claim against repository evidence.

Do not invent features, installation paths, architecture, benchmarks, screenshots, support promises, or roadmap items. Do not impose a standard README template or expand a focused request into a repository-wide documentation program.

Finish with the changed documents, corrected or removed claims, verification performed, and any fact that still needs an owner.
