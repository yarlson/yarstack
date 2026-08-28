---
name: changes-report
description: Report all repository changes that entered a target branch during one required date or period, grouped by product change and written as plain English for a broad audience. Use for daily, weekly, date-range, or relative-period change summaries rather than explaining one known diff.
---

# Report Changes for a Period

Produce one evidence-based account of what changed and why it matters, independent of commit and pull request boundaries.

## Required Input

Require a date or period. Accept a single `YYYY-MM-DD` date, an explicit start and end date, or a clear relative period such as `today`, `this week`, or `last week`. If the user did not specify one, ask which date or period to report and stop. If it cannot be resolved to exact bounds, ask for a precise date or period and stop.

Use the repository's local timezone unless the user specifies another timezone. Treat a single date as the half-open period from 00:00 on that date to 00:00 on the next date. Treat an explicit end date as inclusive by ending at 00:00 on the following date. Resolve calendar weeks from Monday at 00:00. End a relative period that includes the present, such as `today` or `this week`, at the current time; use full calendar bounds for completed periods such as `yesterday` or `last week`.

Before inspecting changes, record the resolved start, exclusive end, timezone, and human-readable period label. Use those exact bounds throughout the report.

## Workflow

1. Confirm the repository and target branch. Use the current repository and branch when the user does not name them.
2. Find each merge, squash, and direct commit added to the target branch during the period, using first-parent history and committer time to define inclusion. Inspect the net change from the branch state immediately before the period through its final in-period state, then inspect individual diffs, full changed files, tests, configuration, and relevant surrounding code as needed. Account for changes inside merge commits without counting their constituent commits twice.
3. Use the authenticated `gh` CLI when available to read GitHub commit metadata, linked pull requests, pull request descriptions, reviews, and issue context. Prefer `gh` over browser scraping or ad hoc GitHub API calls. Fall back to local evidence or another authorized read-only tool when `gh` is unavailable or the repository is not hosted on GitHub. Treat this material as evidence of intent, not proof of behavior, and do not require network access when repository evidence is enough.
4. Apply `changes-explain` to identify changed behavior, ownership, boundaries, external effects, compatibility, operations, and risk. Separate observed behavior, documented intent, and inference.
5. Reconcile overlapping evidence into product changes. Combine commits and pull requests that serve the same user, operator, or business outcome. Split one pull request when it contains distinct outcomes. Do not use commits, pull requests, files, or authors as the report's organizing structure.
6. Explain why each change was made when evidence supports it and what it means for users, operators, developers, or the business. Include implementation detail only when a reader needs it to understand the effect, constraint, risk, or evidence limit. Never invent a rationale.
7. Use `text-improve` on the complete draft. Remove jargon, vague abstractions, repeated meaning, needless qualifiers, and low-level narration while preserving technical precision.

Keep the work read-only. Redact secrets and sensitive content. Report missing branches, unavailable pull request context, shallow history, or other coverage gaps as limits rather than silently claiming completeness.

## Required Format

Use exactly this order:

```text
# Changes for PERIOD

## Summary

One short paragraph that states the main product outcomes, target branch, resolved start and exclusive end, and timezone.

## Product changes

### Product change title

One or more connected paragraphs that explain the reason, changed behavior, practical effect, and material risk or limitation.

Repeat one level-three section for each distinct product change.

## Evidence limits

One short paragraph that states any limits on completeness or confidence. Say that no material evidence limits were found when none remain.
```

Write prose only. Do not use bullets, numbered lists, tables, checklists, commit inventories, pull request inventories, or file-by-file summaries. Keep the summary brief, cover every material product change once, and order product changes by reader impact rather than repository chronology.

Finish only when every in-period change maps to a product-change section or an explicit evidence limit.
