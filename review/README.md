# review

Review checks local changes for security issues, bugs, and performance problems with ranked findings and concrete fixes before code becomes a PR.

## Known Limits

- Reviews local diffs; it does not guarantee CI, runtime, or production behavior without running the relevant checks.
- Needs the correct integration target for merge-base diffs; non-standard branch setups may require user input.
- Prioritizes material risks over style nits.

## Install

```bash
npx skills add https://github.com/yarlson/skills/tree/main/review
```
