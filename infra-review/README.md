# infra-review

Infra-review checks local infrastructure-as-code changes for network exposure, IAM policies, destructive changes, secrets, and cost impact before `plan`/`apply`/`deploy`. Works with any IaC tool, cloud, or format.

## Known Limits

- Reviews changed files; it does not replace provider-native plans, policy checks, or environment-specific approvals.
- Needs environment and workspace context to judge blast radius accurately.
- Should not apply, deploy, or mutate infrastructure during review.

## Install

```bash
npx skills add https://github.com/yarlson/skills/tree/main/infra-review
```
