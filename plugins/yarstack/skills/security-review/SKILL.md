---
name: security-review
description: Review a focused trust-boundary or exploit-path change without modifying it. Use when work introduces or materially changes untrusted input, privilege decisions, secret flow, dangerous sinks, or agent and tool capabilities.
---

# Security Review

Own concrete trust-boundary and exploit-path analysis, not generic changed-code review.

## Workflow

1. Confirm the diff or audit scope, assets, actors, trust boundaries, and relevant privileges.
2. Trace each material untrusted source through parsing, validation, authorization, and dangerous sinks or persistent side effects.
3. Check changed secret handling, command execution, paths and archives, network destinations, authentication, authorization, tenant isolation, logs, plugins, MCP tools, sandboxes, and agent capabilities only where they participate in a concrete path.
4. Separate confirmed exploit paths from defense-in-depth observations and unresolved evidence gaps.
5. Delegate dependency provenance to `dependency-review`, IaC blast radius to `infra-review`, platform lifecycle to its specialist, and operational recovery to `rollout-readiness-review`.

For each finding report severity, affected location and behavior, attacker prerequisites, exploit path, impact, smallest safe correction, verification, and uncertainty. State directly when no confirmed finding exists.

Keep the review read-only. Do not rotate credentials, modify live services, run unsafe payloads, add security frameworks, or broaden scope. Remediation requires separate implementation authority.

Finish when material trust paths in scope are traced and findings or confidence limits are explicit.
