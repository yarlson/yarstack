---
name: ci-review
description: Audit CI workflow architecture and effectiveness without modifying it. Use when triggers, jobs, commands, quality gates, caches, artifacts, matrices, permissions, release paths, or deployment automation may be stale, redundant, unsafe, or unjustified.
---

# CI Review

Own automation topology and evidence that each component protects a real contract.

## Workflow

1. Confirm repository scope, comparison target, supported platforms, required-check ownership, and whether remote run evidence is authorized.
2. Map triggers, jobs, called commands, environments, permissions, secrets, caches, artifacts, release paths, deployment paths, and required checks.
3. Identify each component's consumer and enforced contract using repository and authorized run evidence.
4. Find duplicated logic, mismatched versions, ineffective filters, unused matrices, ignored results, retry-masked flakes, dead paths, unjustified permissions, and caches or artifacts without consumers.
5. Classify actionable items as `keep`, `simplify`, `merge`, `fix`, `remove`, or `owner decision`.
6. Delegate tool provenance to `dependency-review`, exploit paths to `security-review`, infrastructure semantics to `infra-review`, runtime readiness to `rollout-readiness-review`, and behavioral test sufficiency to `test-gap-review`.

For each finding report location, affected contract, evidence, consequence, recommended disposition, change risk, and verification. Recommendations do not authorize implementation.

Do not fetch remote evidence without authority, remove a safety gate on speculation, add scanners for appearance, or rewrite workflows broadly.

Finish with prioritized findings or no-findings, scope, evidence limits, and the smallest advisory cleanup sequence.
