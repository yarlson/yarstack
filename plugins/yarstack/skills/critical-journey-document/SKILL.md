---
name: critical-journey-document
description: Create concise Markdown documentation for critical user journeys, flows, acceptance scenarios, service blueprints, telemetry, portfolios, story maps, and coverage. Use when product behavior needs traceable docs-as-code tied to evidence and ownership.
---

# Critical Journey Documentation

Document one actor pursuing one goal. Use structured Markdown and only the artifacts needed for the current decision or delivery gate.

## Core Model

- Journey ID: `JRN.<domain>.<goal>.<optional-channel>.<optional-variant>`
- Scenario ID: `SCN.<domain>.<goal>.<case>`
- Event: `<domain>_<object>_<action>`
- Status: `draft`, `review`, `active`, or `deprecated`
- Risk: `low`, `medium`, `high`, or `critical`
- Missing evidence: `TBD(owner)`

## Artifact Selection

Create the smallest useful set under the repository's established documentation path:

- **Brief:** actor, goal, scope, trigger, preconditions, success, failure, owner, risk, and evidence links.
- **Flow:** ordered touchpoints, decisions, state changes, failure paths, recovery, and accessibility behavior.
- **Acceptance:** named Given/When/Then scenarios for success, failure, recovery, permissions, and accessibility.
- **Telemetry:** funnel steps, event contracts, required properties, privacy constraints, metrics, thresholds, and dashboard ownership.
- **Blueprint:** frontstage, backstage, support systems, handoffs, operational risks, and recovery; add only for cross-team or operationally risky journeys.
- **Portfolio, coverage, or story map:** add only when coordinating multiple journeys, release evidence, or sequencing.

## Workflow

1. Identify actor, goal, domain, channel, risk, and owner.
2. Inspect existing requirements, code, contracts, tests, telemetry, and documentation.
3. Split multiple actors or goals into separate journeys.
4. Create or update only the necessary artifacts.
5. Link claims to existing requirements, contracts, tests, dashboards, and evidence.
6. Mark unsupported facts as `TBD(owner)` rather than inventing them.
7. Run the repository's Markdown formatter when available.

## Guardrails

- Keep files concise and Git-tracked.
- Prefer tables, bullets, checklists, IDs, and small Mermaid diagrams over prose.
- Do not create posterware, unsupported claims, duplicate documentation, or an entire artifact suite by default.
- Do not treat visual design tools as canonical; link to them.

Finish when the chosen artifacts provide enough traceability to design, implement, validate, release, or operate the journey.
