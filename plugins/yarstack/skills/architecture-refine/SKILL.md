---
name: architecture-refine
description: Refine underdefined product or system architecture through focused interactive decisions. Use when the user wants architecture sparring before a formal implementation plan or contract update.
---

# Architecture Sparring

Resolve only the architecture decisions needed for the user's stated next step.

## Interaction Loop

1. Inspect the relevant repository evidence and known constraints.
2. Identify the highest-leverage unresolved architecture decision.
3. Ask one focused question, grouping inseparable constraints only when the decision cannot be evaluated honestly otherwise.
4. Offer two or three concrete options with material tradeoffs.
5. Evaluate the answer directly and state its architectural consequence.
6. Continue only while another unresolved boundary changes the requested next deliverable.

Pressure-test runtime boundaries, ownership, configuration, persistence, deployment, operations, security, migration, rollback, and support burden. Prefer the smallest viable design and reject flexibility that leaves behavior undefined.

Do not create an implementation plan, edit a specification, conduct extended technical research, design a refactor sequence, implement code, or write durable repository context. Use `plan-create`, `spec-update`, `technical-spike`, or `repo-context-document` for those outputs.

Finish when the decisions needed for the stated next step are settled, the user ends sparring, or a named blocker requires another workflow. Report accepted decisions, rejected alternatives, non-goals, open risks, and remaining blockers.
