---
name: spec-update
description: Make a minimal evidence-backed update to a non-plan product, architecture, API, schema, or behavior contract. Use when implementation or validation cannot determine correctness from the current authoritative contract.
---

# Specification Update

Correct an authoritative semantic contract only when current work depends on it.

## Workflow

1. Identify the missing or incorrect contract, its canonical document and owner, and the current task it blocks.
2. Use `plan-update` for implementation-plan mechanics and `docs-drift-review` for ordinary prose drift.
3. Distinguish an evidence-backed clarification from a new product or architecture decision.
4. Use `architecture-refine` for an unresolved architecture decision. Stop for the user or named owner when available evidence cannot determine intended behavior.
5. Confirm authority before editing linked documents outside the repository or requested scope.
6. Apply only the minimum supported correction in the document's existing structure and style.
7. State the downstream acceptance or validation effect, then return to the invoking task.

Do not create speculative features, use specifications as implementation logs, rewrite unrelated sections, or silently invent behavior.

Finish with the changed contract, supporting evidence, concrete decision, downstream impact, and any remaining blocker.
