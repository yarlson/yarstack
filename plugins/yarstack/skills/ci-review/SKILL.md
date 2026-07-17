---
name: ci-review
description: Audit CI and delivery automation for stale, redundant, unsafe, ineffective, or unjustified machinery. Use when reviewing workflows, build scripts, quality gates, caches, artifacts, matrices, triggers, permissions, release jobs, or deployment jobs for evidence-backed cleanup.
---

# CI Review

Determine which automation protects a real contract and which parts should be simplified, fixed, merged, or removed. Keep the review read-only.

## Workflow

1. Map workflow triggers, jobs, called scripts, task-runner commands, environments, permissions, secrets, caches, artifacts, release steps, deployment steps, and required checks.
2. Identify the consumer and enforced contract for each component. Verify claims against repository configuration, supported platforms, deployment paths, and available run evidence.
3. Check for duplicated local and CI logic, mismatched tool versions, unrelated full-suite runs, ineffective path filters, unused matrices, ignored scanners, retry-masked flakes, dead release paths, and copied configuration that does not fit the repository.
4. Check whether caches and artifacts have measurable consumers, whether schedules and manual gates have owners, and whether jobs receive more secrets or permissions than they need.
5. Classify each actionable item as `keep`, `simplify`, `merge`, `fix`, `remove`, or `owner decision`.
6. For any proposed removal or simplification, state the protection being replaced or why no real protection exists, the change risk, and the smallest verification needed afterward.

Do not recommend broad workflow rewrites, remove a safety gate on speculation, or add scanners and process merely to make CI look mature. Lead with evidenced findings and finish with a minimal cleanup sequence.
