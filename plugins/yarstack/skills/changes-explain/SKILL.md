---
name: changes-explain
description: Explain a pull request, commit range, branch, or local diff from a system-design perspective. Use for practical understanding of changed behavior, boundaries, flow, ownership, operations, risks, and verification rather than defect review.
---

# Explain Changes

Translate a confirmed code comparison into system effects without inventing intent or exposing sensitive content.

## Workflow

1. Confirm the comparison scope; stop or qualify the explanation when it cannot be isolated.
2. Read the diff, full changed files, relevant surrounding code, tests, configuration, and linked in-scope brief.
3. Group edits by changed behavior or responsibility rather than file mechanics.
4. Explain previous and new behavior, affected owner or boundary, data and control flow, external interactions, and user, developer, runtime, or operational effects.
5. Describe materially changed failure, compatibility, rollout, CI, verification, and residual-risk behavior.
6. Distinguish observed evidence, documented intent, and inference.

Keep the task read-only. Redact secrets and sensitive diff content. Use `code-review` for defect findings and `rollout-readiness-review` for operational-readiness judgment.

Finish when every material change group is explained and evidence limitations are stated. Omit syntax trivia and empty sections.
