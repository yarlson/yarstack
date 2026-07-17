---
name: infra-review
description: Review infrastructure-as-code and deployment configuration without modifying it. Use before plan, apply, or deploy when changes may affect network exposure, IAM or RBAC, secrets, encryption, state, destructive replacement, availability, cost, rollout, or recovery.
---

# Infrastructure Review

Assess operational blast radius using the diff, full resource definitions, repository conventions, and available non-mutating evidence.

## Establish Scope

1. Read applicable repository instructions.
2. Inspect status, changed paths, diffs, relevant commits, and surrounding modules or overlays.
3. Confirm the comparison target and affected environments, accounts, regions, clusters, or tenants.
4. Use local refs by default. Fetch only when current remote state is required and permitted.
5. Identify the tooling and source-to-rendered-config path before running checks.

## Review Workflow

1. Rank changes by blast radius: state backends, production targeting, IAM, networks, databases, storage, compute, deployment, then low-risk metadata.
2. Check public exposure, least privilege, secret handling, encryption, container identity, dependency pinning, and supply-chain sources.
3. Check resources that may be destroyed or replaced, data migration, backups, deletion protection, dependency ordering, drift, and partial-apply recovery.
4. Check environment targeting, rollout order, health checks, resource bounds, disruption behavior, rollback, and observability.
5. Check material cost drivers, unbounded storage or traffic, oversized capacity, and cross-region transfer.
6. Review existing plan, render, diff, or dry-run evidence when available. Run only repository-supported non-mutating checks appropriate to the requested review.

## Finding Standard

For each actionable finding include severity, exact file and line, affected environment, concrete blast radius, failure path, smallest safe correction, and evidence. Verify whether an attribute truly forces replacement before claiming destruction.

Block on confirmed public exposure, privilege escalation, secret leakage, unprotected stateful destruction, unsafe state migration, or production rollout without a credible recovery path.

## Guardrails

- Keep the review read-only. Never apply, deploy, install, upgrade, delete, destroy, or modify cloud state.
- Do not echo credentials, account identifiers, internal addresses, or sensitive plan output.
- Do not assume provider semantics; verify version-sensitive behavior from local versions or official documentation.
- Do not add scanners or policy frameworks merely to complete the review.
- State when missing credentials, backends, rendered output, or tooling limits confidence.

Lead with findings, then summarize scope, destructive changes, cost impact, checks, and remaining uncertainty.
