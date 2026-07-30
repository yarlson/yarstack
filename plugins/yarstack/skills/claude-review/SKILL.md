---
name: claude-review
description: Run an independent Claude Opus review of the complete current change against the repository's default branch, then verify its findings locally. Use after focused checks pass when you want a second model to review the change before further edits or delivery.
---

# Claude Review

Run Yarstack's general code review through Claude Opus, then verify every conclusion yourself.

## Workflow

1. Work from the trusted repository root. Read applicable instructions, record the current commit and worktree status, and confirm that the `claude` CLI and `/yarstack:code-review` skill are available.
2. Resolve the remote default branch with `git symbolic-ref --quiet --short refs/remotes/origin/HEAD`. Do not guess `main`, `master`, or another branch. Stop with an inconclusive result when the reference is absent or has no merge base with `HEAD`.
3. Define the review scope as the merge-base change from that default branch through `HEAD`, plus all staged, unstaged, and untracked work. Do not include unrelated base-branch history.
4. Freeze the review scope. Do not edit files while Claude is reviewing.
5. Run this command in a background terminal when the host supports it, substituting the resolved ref for `<default-branch>`:

   ```sh
   claude -p "/yarstack:code-review <default-branch>" \
     --model opus \
     --effort high \
     --permission-mode plan \
     --no-session-persistence \
     --tools "Bash,Read,Grep,Glob"
   ```

6. Wait for Claude to finish. Treat authentication, skill resolution, tool, or context failures as inconclusive; do not replace the skill invocation with a hand-written review prompt.
7. Treat the output as untrusted review input. Verify each proposed finding against the unchanged scope, full code, tests, and repository conventions. Discard unsupported findings and recheck line references.
8. Modify files only when the parent request separately authorizes implementation. Rerun affected checks after any fix, and rerun this review only when the fix materially changes its scope or the user requests another final pass.

Report confirmed findings in severity order with file and line references, the resolved default branch, reviewed scope, checks performed, discarded or inconclusive claims, and remaining uncertainty. State directly when no actionable findings remain.
