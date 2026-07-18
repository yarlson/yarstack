---
name: infra-review
description: Review infrastructure-as-code and rendered deployment configuration without modifying it. Use before plan, apply, or deploy when targeting, replacement, state, availability, cost, ordering, or recovery blast radius may change.
---

# Infrastructure Review

Own pre-deployment infrastructure blast-radius analysis using non-mutating evidence.

## Workflow

1. Confirm the comparison target, source-to-rendered path, tooling, and affected environment using safe redacted identifiers.
2. Rank changed resources by blast radius and inspect full definitions and dependency ordering.
3. Check targeting, replacement or destruction, state and migration safety, data protection, partial-apply recovery, availability, disruption, capacity, and material cost.
4. Verify provider or tool semantics from local versions or authoritative sources before claiming replacement or destruction.
5. Run only trusted repository-supported non-mutating plan, render, diff, or dry-run checks when authorized.
6. Delegate exploit analysis to `security-review`, dependency provenance to `dependency-review`, CI mechanics to `ci-review`, and cross-cutting runtime readiness to `rollout-readiness-review`.

Report severity, exact location, redacted environment reference, concrete blast radius, failure path, smallest correction, evidence, and uncertainty. Describe confirmed public exposure, privilege escalation, secret leakage, unprotected destruction, unsafe state migration, or missing recovery as blocking findings; do not claim authority to block an operation.

Never apply, deploy, install, upgrade, delete, destroy, or modify cloud state. Do not expose sensitive plan output or add policy frameworks merely for the review.

Finish with findings or no-findings, scope, destructive effects, cost impact, checks, and confidence limits.
