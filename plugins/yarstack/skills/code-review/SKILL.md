---
name: code-review
description: Review a local diff or branch without modifying it. Use for general cross-cutting correctness, performance, compatibility, lifecycle, and maintainability findings; combine it with go-review for meaningful Go code and macos-swift-review for meaningful Swift or macOS app code.
---

# Code Review

Own the general evidence-backed review of a confirmed comparison scope.

## Workflow

1. Read applicable repository instructions, confirm the comparison target from the request, status, changed paths, and local refs, and identify the implementation languages in scope.
2. Stop or qualify the review when the scope cannot be isolated. Do not include unrelated recent commits.
3. Read full changed files and relevant callers, callees, tests, schemas, and configuration.
4. Check normal paths, boundaries, failures, partial state, cleanup, cancellation, ordering, compatibility, resource ownership, hot-path cost, and unnecessary complexity.
5. Use `go-review` whenever the scope contains meaningful Go production or test code, a Go package, or a full Go codebase. Use `macos-swift-review` whenever it contains meaningful Swift production or test code, a Swift package, an Xcode project, or a macOS Swift app. Use both for mixed scopes; skip them only when their language code is absent or the change is limited to documentation, metadata, or generated output. Use other specialists when material domain depth is required.
6. Run only trusted read-only checks that the user requested or that cheaply verify a suspected defect. Never execute arbitrary code introduced by an untrusted change.
7. Use `crap-index-assess` within those execution bounds when the review targets CRAP or the repository's configured CRAP check covers changed methods. Treat the result as a maintainability signal and apply the finding standard below.

## Finding Standard

Report only findings with a concrete failure path and evidence. Include severity, verified file and line, affected behavior, failure path, smallest safe correction, and uncertainty. Rank by impact and group repeated instances. State directly when no actionable findings exist.

Do not expose sensitive values, invent findings, report style preferences as defects, or modify files. Leave GitHub thread handling to `coderabbit-triage`.

Finish with findings, confirmed scope, checks performed, and confidence limits.
