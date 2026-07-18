---
name: code-review
description: Review a local diff or branch without modifying it. Use to find concrete cross-cutting correctness, performance, compatibility, lifecycle, and maintainability defects before or after a pull request exists.
---

# Code Review

Own the general evidence-backed review of a confirmed comparison scope.

## Workflow

1. Read applicable repository instructions and confirm the comparison target from the request, status, changed paths, and local refs.
2. Stop or qualify the review when the scope cannot be isolated. Do not include unrelated recent commits.
3. Read full changed files and relevant callers, callees, tests, schemas, and configuration.
4. Check normal paths, boundaries, failures, partial state, cleanup, cancellation, ordering, compatibility, resource ownership, hot-path cost, and unnecessary complexity.
5. Use specialist review skills when material domain depth is required; they supplement this review and inherit its read-only scope.
6. Run only trusted read-only checks that the user requested or that cheaply verify a suspected defect. Never execute arbitrary code introduced by an untrusted change.

## Finding Standard

Report only findings with a concrete failure path and evidence. Include severity, verified file and line, affected behavior, failure path, smallest safe correction, and uncertainty. Rank by impact and group repeated instances. State directly when no actionable findings exist.

Do not expose sensitive values, invent findings, report style preferences as defects, or modify files. Leave GitHub thread handling to `coderabbit-triage`.

Finish with findings, confirmed scope, checks performed, and confidence limits.
