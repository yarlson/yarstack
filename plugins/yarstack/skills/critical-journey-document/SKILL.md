---
name: critical-journey-document
description: Create one concise evidence-linked product journey record for one actor pursuing one goal. Use when product behavior needs traceable ownership, acceptance, telemetry, or service-flow documentation for a current decision or delivery gate.
---

# Critical Journey Documentation

Create one canonical journey record in the repository's established product-documentation location.

## Workflow

1. Confirm what makes the journey critical, its actor, goal, trigger, scope, owner, risk, and canonical document location.
2. Inspect existing requirements, code, contracts, tests, telemetry, and documentation; reuse repository identifiers and terminology.
3. Split multiple actors or goals into separate records only when independently useful.
4. Document the ordered journey, decisions, state changes, success, meaningful failures, recovery, ownership, and evidence links.
5. Add acceptance, telemetry, blueprint, portfolio, coverage, or story-map sections only when a named consumer or delivery gate requires them.
6. Record unresolved facts only when the task permits them and names the decision owner and resolution condition.

Acceptance sections describe product acceptance; `test-design` owns executable tests. Telemetry sections record agreed contracts rather than inventing events. `repo-context-document` owns implemented-system flows under `docs/context/`; `rollout-readiness-review` owns deployment readiness; `spec-update` owns semantic decisions.

Do not create a parallel canonical document, mandatory identifier system, speculative telemetry, or an artifact suite by default.

Finish when the record has evidence links, named ownership, explicit unknowns, and only sections with current consumers.
