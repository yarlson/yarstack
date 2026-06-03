# ship

Ship creates a branch, commits and pushes changes, opens a draft pull request, waits for CI, fixes red builds, and iterates until the checks are green.

It checks dirty-worktree scope first and stages only confirmed paths so unrelated local changes stay untouched by default.

## Known Limits

- Requires working GitHub CLI authentication and a repo with a pushable remote.
- Cannot infer ambiguous shipping scope safely; it should ask before staging unclear files.
- Does not replace human review of the draft PR before merge.

## Install

```bash
npx skills add https://github.com/yarlson/skills/tree/main/ship
```
