# repo-context-docs

Repo-context-docs creates or updates current-state repository context documentation under `docs/context/` based on the codebase, configuration, and tests.

## Known Limits

- Documents current state only; it is not a changelog or implementation plan.
- Should write only under the target repo's `docs/context/` directory.
- Needs fresh repo inspection because context docs drift quickly.

## Install

```bash
npx skills add https://github.com/yarlson/skills/tree/main/repo-context-docs
```
