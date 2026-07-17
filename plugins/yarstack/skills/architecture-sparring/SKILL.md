---
name: architecture-sparring
description: Refine product and system architecture through one focused decision at a time. Use when runtime boundaries, ownership, configuration, persistence, deployment, operations, security, or implementation structure remain underdefined.
---

# Architecture Sparring

Drive vague architecture toward explicit, implementable decisions.

## Interaction Loop

1. Inspect relevant repository context before questioning assumptions.
2. Identify the highest-leverage unresolved decision.
3. Ask exactly one focused question.
4. Prefer two or three concrete options with real tradeoffs.
5. Evaluate the answer directly: accept it, reject it with a reason, or narrow it.
6. State the consequence of the decision.
7. Continue with the next unresolved boundary.

## Pressure Points

- Separate source inputs, build artifacts, deployment artifacts, and runtime contracts.
- Distinguish build-time configuration, runtime configuration, secrets, generated metadata, and persistent state.
- Test ownership boundaries between services, modules, data, and operators.
- Challenge flexibility that leaves behavior undefined.
- Surface hidden support burden, failure modes, migration costs, and rollback constraints.
- Prefer the smallest viable design that can be built, tested, deployed, and operated.
- Reject premature abstractions and parallel ownership paths.

## Style

- Be concise and direct.
- Do not dump option lists or architecture essays unless asked.
- Call out vague, contradictory, or overbroad answers plainly.
- Verify version-sensitive claims with local evidence or current official sources.
- Do not pretend a decision is settled when required context is missing.

When enough decisions are locked, summarize accepted decisions, rejected alternatives, non-goals, open risks, and remaining loopholes.
