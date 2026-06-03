# implement

Implement writes code under an extremely strict maintainability bar: keep changes surgical, avoid structural regressions, prevent spaghetti branching and file sprawl, and verify behavior plus structure before final review.

## Known Limits

- Does not replace a dedicated review pass for broad security, infra, or production-risk changes.
- Needs concrete success criteria and local verification commands to avoid fake confidence.
- Should not broaden scope into unrelated refactors just to improve nearby code.

## Install

```bash
npx skills add https://github.com/yarlson/skills/tree/main/implement
```
