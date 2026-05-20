---
name: critical-journey-docs
description: Use when documenting critical user journeys, user flows, acceptance scenarios, service blueprint slices, telemetry slices, journey portfolios, story maps, or coverage matrices as concise Markdown docs-as-code artifacts.
---

# Critical Journey Docs

## Rules

- Output Markdown only.
- Keep generated files concise.
- Use tables, bullets, IDs, checklists, and Mermaid.
- No explanatory prose.
- No posterware.
- No unsupported claims.
- One persona, one goal per journey.
- Canonical source: Git-tracked Markdown.
- Visual tools: links only.
- Missing evidence: `TBD(owner)`.

## IDs

- Journey: `JRN.<domain>.<goal>.<channel?>.<variant?>`
- Scenario: `SCN.<domain>.<goal>.<case>`
- Event: `<domain>_<object>_<action>`
- Status: `draft`, `review`, `active`, `deprecated`
- Risk: `low`, `medium`, `high`, `critical`

## Default Files

Single critical journey:

- `journeys/<domain>/<journey-id>/brief.md`
- `journeys/<domain>/<journey-id>/flow.md`
- `journeys/<domain>/<journey-id>/acceptance.md`
- `journeys/<domain>/<journey-id>/telemetry.md`

Add when relevant:

- `journeys/<domain>/<journey-id>/blueprint.md` for cross-team, regulated, support-heavy, revenue-critical, or operationally risky journeys.
- `journeys/portfolio.md` for multi-journey inventory.
- `journeys/coverage.md` for release readiness or traceability.
- `journeys/story-map.md` for backlog sequencing.

Follow existing repo docs paths when present.

## Templates

Load only needed templates:

- `templates/brief.md`
- `templates/flow.md`
- `templates/acceptance.md`
- `templates/telemetry.md`
- `templates/blueprint.md`
- `templates/portfolio.md`
- `templates/coverage.md`
- `templates/story-map.md`

## Workflow

1. Identify actor, goal, domain, channel, risk, and owner.
2. Inspect existing docs/code/config/tests when writing inside a repo.
3. Select the smallest useful file set.
4. Create or update Markdown files from templates.
5. Link requirements, design, contracts, tests, telemetry, dashboards, and evidence.
6. Update portfolio or coverage files when the journey set changes.
7. Run the repo markdown formatter when available.

## Gates

- Discovery: actor, goal, scope, evidence.
- Design: flow, content, accessibility notes.
- Implementation: scenarios, contracts, recovery paths.
- Pre-release: tests, telemetry, dashboards.
- Operate: drift, incidents, support signals, KPI movement.

## Stop Conditions

- No clear journey: ask one question.
- No evidence: write `TBD(owner)`.
- Multiple personas/goals: split journeys.
- Requested prose output: keep prose minimal and structured.
