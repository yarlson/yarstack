# skills

Reusable prompt-driven capabilities for AI coding agents. Each skill is a self-contained folder that teaches your agent a specific workflow.

Skills work with any agent that supports the [skills](https://github.com/anthropics/skills) format. Install with a single command.

## Install

All current repo skills:

```bash
npx skills add https://github.com/yarlson/skills/tree/main/architecture-sparring
npx skills add https://github.com/yarlson/skills/tree/main/ci
npx skills add https://github.com/yarlson/skills/tree/main/ci-with-pr
npx skills add https://github.com/yarlson/skills/tree/main/code-review
npx skills add https://github.com/yarlson/skills/tree/main/coderabbit
npx skills add https://github.com/yarlson/skills/tree/main/create-todo
npx skills add https://github.com/yarlson/skills/tree/main/critical-journey-docs
npx skills add https://github.com/yarlson/skills/tree/main/infra-code-review
npx skills add https://github.com/yarlson/skills/tree/main/project-context
```

## Skills

| Skill                                             | What it does                                                                                          |
| ------------------------------------------------- | ----------------------------------------------------------------------------------------------------- |
| [`architecture-sparring`](architecture-sparring/) | Architecture sparring — force clear product and system design decisions one question at a time        |
| [`ci`](ci/)                                       | CI — commit all changes, push, wait for checks, fix red builds, and iterate until green               |
| [`ci-with-pr`](ci-with-pr/)                       | CI with PR — branch, commit, push, open a draft PR, wait for checks, and fix red builds               |
| [`code-review`](code-review/)                     | Code review — security, bugs, perf — with ranked findings and fixes                                   |
| [`coderabbit`](coderabbit/)                       | CodeRabbit — triage bot review comments, fix legitimate issues, push, and resolve threads             |
| [`create-todo`](create-todo/)                     | Create todo — capture concrete work as committed Markdown todos grouped by project                    |
| [`critical-journey-docs`](critical-journey-docs/) | Critical journey docs — concise Markdown journey, flow, acceptance, telemetry, and coverage artifacts |
| [`infra-code-review`](infra-code-review/)         | IaC review — network exposure, IAM, destructive changes, cost — before plan/apply/deploy              |
| [`project-context`](project-context/)             | Project context — create or update current-state docs under `docs/context/`                           |

## How skills work

Each skill is a folder with a `SKILL.md` file — YAML frontmatter (name, description) plus a full prompt that defines the agent's workflow. Some skills also bundle scripts or templates for deterministic operations. The prompt tells the agent what to do, what order to do it in, and what to watch out for.

When you install a skill, it gets added to your agent's context. Invoke it by describing what you need — "review my changes", "generate a readme", "review my infra code" — and the agent follows the skill's workflow.

## License

[MIT](LICENSE)
