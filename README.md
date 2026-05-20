# skills

Reusable prompt-driven capabilities for AI coding agents. Each skill is a self-contained prompt that teaches your agent a specific workflow.

Skills work with any agent that supports the [skills](https://github.com/anthropics/skills) format. Install with a single command, no dependencies.

## Install

Pick the skills you need:

```bash
npx skills add https://github.com/yarlson/skills/tree/main/code-review
npx skills add https://github.com/yarlson/skills/tree/main/infra-code-review
```

## Skills

| Skill                                     | What it does                                                                             |
| ----------------------------------------- | ---------------------------------------------------------------------------------------- |
| [`code-review`](code-review/)             | Code review — security, bugs, perf — with ranked findings and fixes                      |
| [`infra-code-review`](infra-code-review/) | IaC review — network exposure, IAM, destructive changes, cost — before plan/apply/deploy |

## How skills work

Each skill is a `SKILL.md` file — YAML frontmatter (name, description) plus a full prompt that defines the agent's workflow. No code, no runtime dependencies. The prompt tells the agent what to do, what order to do it in, and what to watch out for.

When you install a skill, it gets added to your agent's context. Invoke it by describing what you need — "review my changes", "generate a readme", "review my infra code" — and the agent follows the skill's workflow.

## License

[MIT](LICENSE)
