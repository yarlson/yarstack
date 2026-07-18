# Skill Format

This is the authoring contract for skills under `plugins/yarstack/skills/`.

## Layout

Each skill must use this layout:

```text
plugins/yarstack/skills/<skill-name>/
└── SKILL.md
```

Add `scripts/`, `references/`, or `assets/` only when the skill has a current
consumer for them. Do not add per-skill readmes, changelogs, templates, examples,
or host-specific copies by default.

## Naming

- Use lowercase kebab-case.
- Name skills `<scope>-<action>`, with the subject or domain first and the
  operation last, such as `code-review`, `phase-validate`, or `plan-create`.
- Use the narrowest recognizable scope. It may contain multiple words when
  needed, such as `docs-drift-review` or `marketing-claims-review`.
- End with a concrete action word. Prefer a short base-form verb and preserve an
  established operation term when it is clearer than forced grammar.
- Keep the name short. Do not use abstract labels that hide the operation or
  action-first names that break scope grouping.
- Match the directory name and frontmatter `name` exactly.
- Name a tool or platform when it materially narrows the trigger.
- Do not encode project history, origin, or the Yarstack name in a reusable skill.

## Frontmatter

Use exactly these fields:

```yaml
---
name: example-skill
description: Perform one clear action. Use when the concrete triggering conditions apply.
---
```

The description is the trigger contract. It must:

- state what the skill does;
- state when it should be used;
- distinguish the skill from nearby skills;
- remain concise enough to justify its always-visible context cost.

Do not move trigger conditions into the body. Do not add metadata unless a current
Codex or Claude Code capability requires it.

## Body

Write the body in concise English using direct, imperative instructions.

1. Start with an H1 human-readable title.
2. State the outcome or responsibility in one short paragraph.
3. Describe the smallest reliable workflow.
4. Put evidence, scope, safety, and authorization rules near the step they govern.
5. End with a concrete completion condition or required result when useful.

Use headings such as `Workflow`, `Inputs`, `Finding Standard`, `Guardrails`, or
`Result` only when they clarify the skill's actual contract. Do not force every
skill into the same section template.

## Content Standard

- Give the skill one clear responsibility and trigger.
- Include only knowledge or procedure that improves agent behavior beyond the
  repository guidance and the model's general capability.
- Prefer repository evidence and observable behavior over assumptions.
- Preserve user authorization boundaries for commits, pushes, network access,
  external writes, destructive actions, and subagent use.
- Keep workflows bounded. Do not require repeated passes until subjective
  perfection.
- Delegate specialized concerns to an existing skill instead of duplicating its
  checklist.
- Use examples, scripts, or references only when they improve repeated execution
  or deterministic reliability.

Do not add generic engineering policy, motivational prose, background history,
fixed ceremony, speculative future workflows, exhaustive output templates, or
mandatory parallelism that the task does not require.

## Review

Before adding or changing a skill, confirm:

1. No existing skill already owns the trigger.
2. The skill is reusable beyond one repository or incident.
3. Its description will select it at the right time without colliding with other
   descriptions.
4. Every instruction has a concrete purpose.
5. The body is shorter than the source material it replaces and contains no
   duplicated policy.
6. Optional files have an immediate consumer.

Validate the complete shared plugin with:

```sh
make validate
```

Do not add parallel schema checks or custom validators for rules already enforced
by `plugin-scanner` or `claude plugin validate --strict`.
