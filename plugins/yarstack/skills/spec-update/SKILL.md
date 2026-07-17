---
name: spec-update
description: Make a minimal update to a product, architecture, or behavior contract when implementation or validation exposes a real gap. Use only when correctness cannot be determined without clarifying an authoritative plan or linked specification.
---

# Specification Update

Update a contract only when the current implementation depends on a missing or incorrect decision.

## Valid Reasons

- Product behavior required by the current scope is undefined.
- Code and specification expose an architecture contradiction.
- Correctness cannot be validated because a required contract is missing.
- An authoritative linked document is stale in a way that directly affects the change.

## Workflow

1. Identify the exact missing or incorrect contract.
2. Confirm the current work requires a decision.
3. Update the authoritative plan or linked specification minimally.
4. Preserve the document's existing structure and style.
5. Return to the implementation or validation task.

Do not use specifications as implementation logs, add speculative features, rewrite unrelated sections, or silently invent product behavior when evidence cannot determine the correct contract.
