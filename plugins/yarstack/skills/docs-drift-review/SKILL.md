---
name: docs-drift-review
description: Review documentation for drift caused by the current code or configuration change. Use to identify concrete inaccuracies and, when authorized, repair only directly affected documentation.
---

# Documentation Drift Review

Own change-scoped documentation accuracy through separate inspection and repair stages.

## Workflow

1. Confirm the changed behavior and whether the task authorizes documentation edits or only a read-only drift report.
2. Identify canonical, generated, mirrored, and external documentation describing the changed surface.
3. Compare implementation with README files, instructions, APIs, schemas, operational docs, examples, templates, and command help in scope.
4. Report concrete drift with the affected claim, implementation evidence, canonical owner, and required correction.
5. When edits are authorized, update only directly affected canonical documents in their existing style.
6. Change generated documentation or help at its source and regenerate it through the repository's established command.
7. For authorized prose edits, use `text-improve` on the final changed text without expanding the drift scope.

Use `docs-review` for broader documentation quality, `repo-context-document` for deliberate `docs/context/` work, `plan-update` for plan mechanics, and `spec-update` for semantic contract decisions.

Do not restructure documentation, hand-edit generated content, change product behavior, or modify external documentation without authority.

Finish with changed documents, an evidenced no-change result, or unresolved owner decisions.
