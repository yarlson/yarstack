---
name: changes-explain
description: Explain a pull request, commit range, branch, or local diff from a system-design perspective. Use when readers need practical understanding of changed behavior, boundaries, data flow, ownership, operations, risks, and verification without a line-by-line code walkthrough.
---

# Explain Changes

Translate a code change into its system effects without inventing intent.

## Workflow

1. Confirm the comparison scope and read the diff, full changed files, relevant surrounding code, tests, configuration, and any linked brief or issue available in scope.
2. Group related edits by changed behavior or responsibility rather than enumerating files mechanically.
3. For each meaningful group, explain the previous and new behavior, affected boundary or owner, data and control flow, external interactions, and user, developer, runtime, or operational effect.
4. Describe failure modes, compatibility implications, rollout or CI effects, verification, and unresolved risks where they materially changed.
5. Distinguish documented intent from inference. Do not claim why a change was made unless the request, issue, tests, or repository evidence supports it.

Lead with a concise system-level summary, then explain meaningful change groups in dependency or execution order. Omit syntax trivia, unchanged context, and empty report sections.
