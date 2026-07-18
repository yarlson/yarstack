---
name: pr-draft
description: Publish confirmed local changes as a monitored GitHub draft pull request. Use only when the user authorizes branch creation, scoped staging, commit, push, draft PR creation, and CI follow-up.
---

# Draft PR

Own end-to-end draft-PR publication after implementation, validation, and review are complete.

## Workflow

1. Confirm shipping authority, repository rules, base branch, remotes, current branch, exact scope, and completed validation and review evidence.
2. Preserve unrelated staged, unstaged, and untracked work. Stop only when safe isolation is impossible or a scoped path is sensitive or unexpectedly generated.
3. Create a short branch when needed, stage explicit paths with `git add -- <paths>`, and review the staged path list and diff.
4. Run required checks when existing evidence does not cover the exact staged content.
5. Commit using repository conventions and push with an upstream when needed.
6. Follow repository PR templates, title conventions, and required fields. Otherwise write concise evidence-based `Problem`, `Change`, `Tests`, and only necessary `Notes` sections.
7. Open a draft PR and monitor checks for the pushed commit.
8. Diagnose failures from logs. Route behavior changes through `behavior-implement` and `test-design`, and specialist failures through their owning workflow; verify, commit, and push only scoped corrections.
9. Stop when checks pass, the user permits cancellation, or a concrete external blocker prevents meaningful progress.

Do not rewrite history, force-push, merge, mark ready, release, weaken tests, disable validation, expose secrets, or override repository PR conventions.

Finish with the draft PR URL, commits, checks, final CI state, and blockers.
