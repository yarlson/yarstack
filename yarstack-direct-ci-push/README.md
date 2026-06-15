# yarstack-direct-ci-push

Direct CI push commits all changes, pushes the current branch, waits for CI, fixes red builds, and iterates until the checks are green.

Deprecated from the public flagship path because it stages all changes and pushes the current branch directly. Use [`yarstack-draft-pr-shipping`](../yarstack-draft-pr-shipping/) for PR-based shipping; keep `yarstack-direct-ci-push` only for intentional direct-push flows.

## Install

```bash
npx skills add https://github.com/yarlson/skills/tree/main/yarstack-direct-ci-push
```
