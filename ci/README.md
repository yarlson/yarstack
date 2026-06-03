# ci

CI commits all changes, pushes the current branch, waits for CI, fixes red builds, and iterates until the checks are green.

Deprecated from the public flagship path because it stages all changes and pushes the current branch directly. Use [`ship`](../ship/) for PR-based shipping; keep `ci` only for intentional direct-push flows.

## Install

```bash
npx skills add https://github.com/yarlson/skills/tree/main/ci
```
