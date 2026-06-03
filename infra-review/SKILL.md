---
name: infra-review
description: Use when reviewing infra or IaC changes before plan, apply, or deploy: Terraform, deploy config, cloud config, network exposure, IAM/RBAC, destructive changes, secrets, and cost.
context: fork
---

# Reviewing Infrastructure-as-Code Changes

This skill reviews IaC changes in the local working tree — before they hit `plan`/`apply`/`deploy`. One bad infra change can take down production, expose data to the internet, or run up a six-figure cloud bill. The blast radius is higher than application code.

**Review philosophy:** Infra changes are high-leverage. A one-line security group change can expose every instance in a VPC. Focus on: blast radius, network exposure, IAM/RBAC scope, data destruction, cost impact, and secret leakage. Explain risk in operational terms — "this opens port 22 to 0.0.0.0/0 on all production instances" not "this security group is too open."

## Core Command Pattern

All reviews start by identifying what changed. The universal approach is a **merge-base diff**: everything you changed relative to the integration target.

```bash
# 1. Update remote refs
git fetch origin

# 2. Get the merge-base
MERGE_BASE=$(git merge-base HEAD origin/main)

# 3. Diff: merge-base → working tree
git diff $MERGE_BASE

# 4. Changed file list
git diff $MERGE_BASE --name-only

# 5. Diff stats
git diff $MERGE_BASE --stat
```

**Why merge-base, not `origin/main` directly?**
`git diff origin/main` includes upstream changes that aren't yours. `git merge-base` finds the fork point — so the diff shows only your work, even if main moved ahead. No rebase required.

**Detecting the integration target:**

```bash
# Auto-detect: fall back to origin/main or origin/master
TARGET=$(git rev-parse --verify origin/main 2>/dev/null && echo origin/main) \
  || TARGET=origin/master
```

If the repo uses a non-standard default branch, the user should specify it. Ask rather than guess.

Infra changes can touch any file type — config files, scripts, templates, generated manifests, CI pipelines, application code with deployment logic. Don't filter by extension. Review everything in the diff.

## Review Workflow

Execute ALL phases in order. Never skip phases. In quick mode, skip phases marked _(full mode only)_.

### Phase 1: Gather context

1. Fetch origin and compute the merge-base.
2. Run `git diff $MERGE_BASE --stat` to understand the scope.
3. Get the changed file list.
4. Read the full content of every changed file — you need the full resource definition, not just the diff hunk.
5. Read the diff itself for line-level analysis.
6. Check commit messages with `git log $MERGE_BASE..HEAD --oneline` for intent context.
7. Identify what kind of change this is: new resource, resource modification, resource deletion, config change, dependency update, refactor.
8. Identify the blast radius: what environments are affected? Is this prod, staging, or dev? Check file paths, workspace names, directory structure.

**Output:** 2-3 sentence summary of what changed, what environments are affected, and the estimated blast radius.

### Phase 2: Classify risk

Rank changed files by blast radius:

- **Critical**: IAM policies/roles, network rules (security groups, firewall rules, NACLs), database configs, encryption settings, secrets/vault configs, production environment targeting, state backend changes
- **High**: compute resources (instances, containers, pods), load balancers, DNS records, storage (object stores, block volumes), CI/CD pipeline configs
- **Medium**: monitoring/alerting rules, logging configs, tags/labels, variable/output definitions, non-prod environment changes
- **Low**: comments, formatting, documentation, dev-only configs, .gitignore

Review critical and high files first. If time-constrained, give low-risk files a summary-only pass.

### Phase 3: Security scan

Check every changed line against these categories. Flag critical issues immediately.

**Network exposure:**

- Security groups / firewall rules open to `0.0.0.0/0` or `::/0`?
- SSH (22), RDP (3389), database ports (3306, 5432, 6379, 27017) exposed to the internet?
- Internal services accidentally made public (public subnets, public IPs, public load balancers)?
- VPC peering or transit gateway changes that widen network access?

**IAM / RBAC:**

- Overly broad IAM policies (`"Action": "*"`, `"Resource": "*"`)?
- Admin/root privileges where least-privilege would suffice?
- IAM policies attached to users instead of roles?
- Service accounts with excessive permissions?
- RBAC roles in Kubernetes granting cluster-admin unnecessarily?
- Missing conditions/boundaries on IAM policies?

**Secrets & credentials:**

- Hardcoded secrets, API keys, passwords, tokens in any file?
- Secrets in config files, variable definitions, or environment variables in plain text?
- Secrets should come from a secrets manager (vault, cloud-native secret store, sealed/external secrets, etc.) — never inline.
- State files or generated outputs committed? These often contain secrets in plaintext.

**Encryption:**

- Storage at rest encrypted?
- TLS/SSL enforced for data in transit?
- Key management appropriate (rotation, customer-managed vs provider-managed)?
- Database/cache connections requiring encryption?

**Container security:**

- Running as root (or no explicit non-root user)?
- Using `latest` or unpinned tags instead of specific version/digest?
- Secrets baked into image layers?
- Unnecessary packages installed (attack surface)?
- Build context copying sensitive files (`.git`, `.env`, credentials)?

**Supply chain:**

- Providers, modules, charts, base images pinned to specific versions?
- Third-party dependencies from untrusted or unverified sources?
- Lockfiles committed and consistent?

### Phase 4: Resource & state analysis

**Destructive changes:**

- Resources being destroyed and recreated (name/identifier changes, immutable attribute modifications)? This causes downtime.
- Deletion protection / prevent-destroy guards on stateful resources (databases, persistent volumes, object storage with data)?
- Stateful resources (databases, queues, storage) being modified in ways that force replacement?
- Data backups or snapshots in place before destructive operations?

**State & dependencies:**

- State backend or remote state changes? (Can lock out the team or lose state.)
- Existing resources being adopted into IaC — does the import/migration path exist?
- Resource dependency ordering correct? Explicit dependencies where implicit ordering isn't enough.
- Circular dependency risks between modules/stacks?

**Drift & idempotency:**

- Changes that will cause drift on next apply (codifying manual changes)?
- Ignore-changes / skip-drift rules — justified, or hiding problems?
- Config overrides that will be reset by the next tool run (chart upgrade, stack update)?

**Cost impact:**

- Instance/resource type changes (especially upward)?
- New resources in expensive categories (NAT gateways, load balancers, GPU instances, multi-AZ databases, cross-region replication)?
- Missing auto-scaling (fixed large instance counts)?
- Storage without lifecycle/expiration rules (unbounded growth)?
- Cross-region or cross-cloud data transfer implications?

### Phase 5: Configuration quality

**Resource definitions:**

- Variables/parameters have descriptions and types?
- Sensitive values marked as sensitive / excluded from output and logs?
- Outputs don't expose secrets or internal identifiers unnecessarily?
- Resources have meaningful tags/labels (name, environment, team, managed-by)?
- Reusable modules/components used instead of copy-pasted resource blocks?
- Provider/dependency versions constrained (pinned or bounded)?
- State/backend configured with locking to prevent concurrent corruption?

**Workload definitions:**

- Values parameterized (not hardcoded in templates)?
- Resource requests and limits set on all compute workloads?
- Health checks defined (liveness, readiness, startup)?
- Disruption budgets or availability constraints for critical services?
- Namespace/project/environment isolation correct?
- Template rendering / overlay patches targeting the right resources?

**Deployment & sync:**

- Deployment targeting the correct cluster, region, account, and environment?
- Auto-sync / auto-apply on production — is it intentional and safe?
- Orphan/prune cleanup enabled where intended?
- Source references (repo, branch, path) correct?
- Rollback and failure handling defined?

**Compose / local orchestration:**

- Services have restart policies?
- Volumes mounted for persistent data?
- Network isolation between services?
- Resource limits set?
- Secrets not inline in environment blocks?
- Health checks defined?

### Phase 6: Run checks _(full mode only)_

Discover what tooling exists in the project before running anything. Look at the repo structure, config files, CI pipelines, and Makefiles/Taskfiles/scripts for clues. Run only what's available.

**What to look for and run:**

1. **Formatting / validation**: most IaC tools have a built-in `fmt` or `validate` command. Run it.
2. **Plan / dry-run / template render**: the most important check. Shows exactly what will be created, changed, or destroyed. Always review the output if available.
3. **Linting**: tool-specific linters catch misconfigurations that validation misses.
4. **Security / policy scanning**: static analysis tools that check for misconfigurations, overly permissive rules, missing encryption, etc.
5. **Secret scanning**: `gitleaks`, `trufflehog`, or whatever the project uses.
6. **Image / dependency scanning**: vulnerability scanners for container images and dependencies.

**The plan/dry-run output is the single most important check.** It shows what will actually happen — not what the code looks like. Static analysis catches patterns; plan output catches reality.

### Phase 7: Blast radius assessment _(full mode only)_

- How many environments are affected by this change?
- If this change fails mid-apply, what's the rollback path?
- Are there dependent resources in other stacks/modules that will break?
- Is there a maintenance window needed for this change?
- Does this change require coordination with other teams?
- Can this be applied incrementally, or is it all-or-nothing?

## Producing the Review

### Finding format

For every finding, provide:

- **File + line**: exact path and line number. MUST exist in the diff — verify before reporting.
- **Severity**: `CRITICAL` / `HIGH` / `MEDIUM` / `LOW`
- **Category**: `security` / `network` / `iam` / `data-loss` / `cost` / `reliability` / `config`
- **Title**: one-line summary
- **Blast radius**: what breaks if this goes wrong — specific environment, service, or scope
- **Current code vs fix**: show both. Before/after, not just "consider changing this."
- **Confidence**: if uncertain, say so. Prefix with `[Uncertain]` and explain what would confirm or refute.

Finding template:

```
CRITICAL security: [Title]
File: modules/vpc/main.tf:42
Blast radius: [What's affected — e.g., "all production instances in us-east-1"]

Current:
  [problematic code]

Fix:
  [corrected code]
```

### Output structure

```
## Infra Review Summary

**Recommendation:** BLOCK | APPROVE WITH COMMENTS | APPROVE
**Environments affected:** [prod / staging / dev / all]
**Blast radius:** [high / medium / low]

[2-3 sentence summary of what changed and overall assessment]

## Blocking Issues ([count])

[CRITICAL and HIGH findings with file, line, blast radius, and before/after fix]

## Non-Blocking Suggestions ([count])

[MEDIUM and LOW findings — cost, reliability, config quality]

## Destructive Changes

[Resources being destroyed or replaced — list each with the impact]

## Cost Impact

[New resources, instance type changes, estimated cost delta if possible]

## Metrics

- Files changed: [count]
- Lines added/removed: [+N / -N]
- Resources created/modified/destroyed: [+N / ~N / -N]
- Critical issues: [count]
- Non-blocking suggestions: [count]
```

### Recommendation logic

**BLOCK** — must fix before applying:

- Any network exposure to the public internet that wasn't there before
- Overly broad IAM (`*:*`), missing auth/encryption on data stores
- Hardcoded secrets or secrets in state
- Destructive changes to stateful resources without protection (no deletion guards, no backup, no snapshot)
- State backend changes without team coordination
- Production changes without a rollback path

**APPROVE WITH COMMENTS** — non-blocking, track as follow-up:

- Missing tags/labels
- Cost optimization opportunities
- Non-critical config improvements
- Missing non-prod monitoring
- Architecture suggestions

**APPROVE** — when all of:

- Zero network exposure issues
- IAM follows least-privilege
- No unprotected destructive changes
- Secrets managed properly
- Cost impact reasonable

## Review Modes

The user may request a specific mode:

- **"quick review" / "fast review"**: phases 1-5 only (no tool execution). Focus on security and destructive changes. Default mode.
- **"full review"**: all 7 phases — run plan/dry-run, lint, security scan, blast radius assessment.
- **"security review"**: phases 1-3 only. Deep security focus — network, IAM, secrets, encryption. Skip config quality/cost.

Default to **quick review** unless the user asks otherwise.

## Handling Large Changesets

Changesets over 1000 lines or 50+ resources require chunking:

1. Start with the changed file list and the diff stats.
2. Classify all files by blast radius (Phase 2).
3. Review all critical/high files in full depth — these are non-negotiable.
4. For remaining files, provide a one-line summary per file.
5. If a module refactor touches many files but is structurally the same change, note "same pattern in N files" and review one representative.
6. Tell the user you've prioritized — don't silently skip files.

Over 5000 lines: warn that thorough review at this scale is unreliable. Suggest splitting by module/environment/resource-type.

## Decision Policy

- **ALWAYS** run `git fetch origin` before computing the diff. Stale refs mean wrong diffs.
- **ALWAYS** use the merge-base, not a direct diff against `origin/main`.
- **ALWAYS** read the full resource definition, not just the diff hunk. A security group rule only makes sense in the context of the full security group.
- **ALWAYS** verify file paths and line numbers exist before citing them. Never hallucinate a reference.
- **ALWAYS** provide concrete code evidence for every finding. No vague "this could be problematic".
- **ALWAYS** state the blast radius — what environment, what services, what users are affected.
- **ALWAYS** show before/after code for fixes.
- **ALWAYS** check if a resource change forces replacement (destroy + create). This is the #1 source of IaC outages.
- **Discover before assuming**: examine the repo structure, config files, CI pipelines, and scripts to understand what tools and patterns are in use before running any commands.
- **Limit output**: cap at 15 findings per review. If more issues exist, summarize the rest and note the count.
- **Rank by blast radius**: network exposure and data loss first. Cost and config quality last.
- **One finding per issue**: don't repeat the same missing-tag finding across 20 resources. Flag it once, note "same pattern in N other resources".
- **Verify claims**: if you flag a public security group, verify the CIDR is actually `0.0.0.0/0`. If you flag a destructive change, verify the attribute actually forces replacement.
- **When in doubt, flag it**: infra mistakes are expensive and hard to reverse. Better to over-flag than miss.

## Safety Rules

- **No secrets in output**: never include API keys, tokens, passwords, account IDs, or connection strings in findings — even if found in the diff. Say "hardcoded secret detected at file:line" without echoing the value.
- **Never run mutating commands** (`apply`, `install`, `upgrade`, `delete`, `destroy`, or anything that changes real infrastructure). Plan, template, dry-run, lint, and validate are safe.
- **Read-only**: the review does not modify any code. Suggested patches are advisory.
- **Never approve blindly**: if you can't confidently assess a change (e.g., complex IAM policy, unfamiliar cloud provider, tool you don't know well), say so instead of guessing.
- **Treat plan/dry-run output as sensitive**: it may contain secrets, account IDs, internal hostnames, or resource ARNs/URIs. Summarize findings without echoing raw output.

## Error Handling

- **No remote configured**: fall back to diffing against the first commit or ask the user what to diff against.
- **Detached HEAD**: use `git diff origin/main` directly as fallback.
- **Default branch isn't `main`**: check with `git remote show origin | grep 'HEAD branch'`. Don't assume `main`.
- **Multiple tools in one repo**: identify what each file belongs to from context. Don't run the wrong tool's commands against the wrong files.
- **Tool not installed**: if a lint/plan/scan tool isn't available, skip that check and note it. Don't fail the review.
- **No backend / no credentials configured**: plan/dry-run won't work without backend access. Note this and rely on static analysis only.
- **Empty diff**: if the merge-base diff is empty, tell the user. Check `git status` for uncommitted changes.
- **Generated / lock files in diff**: note them but don't review line-by-line — just verify they're consistent with the source changes.
