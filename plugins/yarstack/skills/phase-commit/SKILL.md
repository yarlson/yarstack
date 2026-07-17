---
name: phase-commit
description: Create one local commit for a completed implementation phase or other clearly scoped unit of work. Use only after implementation, validation, and review are complete when the user explicitly requests a commit without pushing.
---

# Phase Commit

Create one local commit containing only the completed unit of work.

## Workflow

1. Inspect `git status`, unstaged changes, staged changes, and relevant untracked files.
2. Separate the completed work from pre-existing or unrelated user changes.
3. Stop if the intended commit cannot be isolated safely.
4. Stage only files belonging to the completed scope.
5. Review the staged diff.
6. Create one local commit.
7. Report the commit hash and one-line summary.

## Commit Message

Use the repository's commit convention. If none exists, use Conventional Commits:

```text
<type>(<optional-scope>): <imperative summary>
```

Keep the subject concise, imperative, and without a trailing period. Add a body only when the reason, migration impact, security implications, or breaking behavior is not obvious.

Do not include a changed-file inventory, verification boilerplate, AI attribution, or emoji.

## Guardrails

- Do not create an empty commit.
- Do not stage unrelated changes.
- Do not push, create branches or pull requests, manage CI, publish releases, or create tags.

Finish with the commit hash and concise summary only.
