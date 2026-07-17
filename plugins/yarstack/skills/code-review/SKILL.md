---
name: code-review
description: Review local code changes without modifying them. Use before commit or pull request to find concrete correctness, security, performance, compatibility, lifecycle, and test risks with evidence from the diff and surrounding code.
---

# Code Review

Report defects that can materially affect users, data, security, operations, or maintainability. Do not nitpick.

## Establish Scope

1. Read applicable repository instructions.
2. Inspect `git status`, staged and unstaged diffs, recent commits, and changed paths.
3. Determine the intended comparison from the request and repository state.
4. For branch work, use the merge base with the confirmed integration target.
5. Use local refs by default. Fetch only when current remote state is required and permitted.
6. Stop or qualify the review when the intended scope cannot be isolated.

## Review Workflow

1. Read the full changed files and relevant callers, callees, tests, schemas, and configuration.
2. Review highest-risk surfaces first: authorization, secrets, data writes, migrations, concurrency, external calls, file operations, and deployment behavior.
3. Check normal paths, boundaries, failures, partial state, cleanup, cancellation, ordering, and compatibility.
4. Check trust boundaries for injection, traversal, SSRF, privilege escalation, tenant leakage, and sensitive-data exposure.
5. Check hot paths for unbounded work, repeated I/O, N+1 operations, resource leaks, and avoidable serialization.
6. Check ownership, duplication, unnecessary abstraction, loose contracts, and tests that provide false confidence.
7. Run existing read-only checks only when the user requests a full review or when a cheap check is necessary to verify a suspected defect.

## Finding Standard

Report only findings with a concrete failure path and evidence. For each finding include:

- severity;
- exact file and line in the reviewed diff;
- affected behavior;
- failure or exploit path;
- smallest safe correction;
- uncertainty and what would resolve it, when applicable.

Rank findings by impact. Group repeated instances of the same problem. If no actionable findings exist, say so directly.

## Guardrails

- Keep the review read-only.
- Do not echo secrets or sensitive values.
- Do not invent findings to fill a template.
- Do not treat style preferences as defects.
- Do not cite paths or lines without verifying them.
- Do not run arbitrary code introduced by the change.
- State when binaries, generated output, missing context, or unavailable tooling limit confidence.

Lead with findings, then give a short scope and verification summary. Omit empty sections.
