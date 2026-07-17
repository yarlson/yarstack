---
name: plan-update
description: Make the smallest evidence-backed correction to an existing implementation plan. Use when implementation or validation proves a phase has an incorrect order, missing prerequisite or validation gate, impossible instruction, contradiction, or incomplete contract.
---

# Plan Update

Correct a plan only when the current work proves it is wrong or incomplete. Do not use the plan as a progress log.

## Valid Reasons

- The selected phase cannot be implemented safely as written.
- Validation proves that a required step or gate is missing.
- Phase order or dependencies are wrong.
- A required validation gate is insufficient or impossible.
- The plan contradicts an authoritative repository contract.
- A concrete prerequisite must become part of the selected scope.

## Workflow

1. Identify the exact plan defect and affected phase.
2. Confirm it with repository or validation evidence.
3. Make the smallest correction.
4. Preserve the existing plan structure and style.
5. Continue the selected implementation or validation work after the correction.

Do not rewrite unrelated phases, add speculative roadmap work, duplicate implementation details already clear from code, or disguise scope expansion as plan maintenance.
