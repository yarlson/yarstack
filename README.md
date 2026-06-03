# Yarstack

<img src="assets/yarstack.png" alt="Yarstack" width="100%">

Yarstack is a workflow stack for disciplined coding agents. It helps an agent frame ambiguous work, plan safely, implement with strict maintainability standards, review quality and risk, ship through CI, and leave useful project context behind.

Built for engineers who want agent help without losing control of scope, code quality, or dirty worktrees.

Shipping skills check dirty-worktree scope and stage only confirmed paths by default.

## Workflow Promise

Yarstack turns broad requests into repeatable engineering passes:

1. Frame product and architecture decisions before code.
2. Plan work in reviewable phases.
3. Implement with strict scope control.
4. Review code, infrastructure, and maintainability risks.
5. Ship with CI feedback.
6. Document critical journeys and repo context.

## Install Matrix

| Surface                        | Status        | Install path                                                              |
| ------------------------------ | ------------- | ------------------------------------------------------------------------- |
| Individual skills              | Available now | `npx skills add https://github.com/yarlson/skills/tree/main/<skill-name>` |
| Full Yarstack core plugin      | Available now | `claude plugin install yarstack-core@yarstack`                            |
| Claude plugin marketplace      | Available now | `claude plugin marketplace add yarlson/skills`                            |
| Direct repo clone or vendoring | Available now | Clone this repo and copy the skill folders your agent host supports       |

## Install Claude Plugin Marketplace

```bash
claude plugin marketplace add yarlson/skills
claude plugin install yarstack-core@yarstack
```

## Install Individual Skills

Core:

```bash
npx skills add https://github.com/yarlson/skills/tree/main/frame
npx skills add https://github.com/yarlson/skills/tree/main/plan
npx skills add https://github.com/yarlson/skills/tree/main/implement
npx skills add https://github.com/yarlson/skills/tree/main/review
npx skills add https://github.com/yarlson/skills/tree/main/infra-review
npx skills add https://github.com/yarlson/skills/tree/main/ship
npx skills add https://github.com/yarlson/skills/tree/main/journey-docs
npx skills add https://github.com/yarlson/skills/tree/main/repo-context-docs
```

Integrations:

```bash
npx skills add https://github.com/yarlson/skills/tree/main/coderabbit
```

Personal utilities:

```bash
npx skills add https://github.com/yarlson/skills/tree/main/create-todo
```

Deprecated:

```bash
npx skills add https://github.com/yarlson/skills/tree/main/ci
```

`ci` is deprecated from the public flagship path because it stages all changes and pushes the current branch directly. Use [`ship`](ship/) for PR-based shipping; keep `ci` only for intentional direct-push flows.

## Current Skill Map

### Core

| Skill                                     | Workflow role     | What it does                                                                                                     |
| ----------------------------------------- | ----------------- | ---------------------------------------------------------------------------------------------------------------- |
| [`frame`](frame/)                         | Frame             | Force clear product and system design decisions one question at a time                                           |
| [`plan`](plan/)                           | Plan              | Write implementation plans that prevent structural regressions, wrong-layer logic, file sprawl, and vague phases |
| [`implement`](implement/)                 | Implement         | Implement code under an extremely strict maintainability bar so strict review finds no structural blockers       |
| [`review`](review/)                       | Review            | Review code for security, bugs, performance, regressions, and missing tests                                      |
| [`infra-review`](infra-review/)           | Review infra      | Review IaC for network exposure, IAM, destructive changes, cost, and deploy risk                                 |
| [`ship`](ship/)                           | Ship with PR      | Check scope, branch, stage confirmed paths, commit, push, open a draft PR, wait for checks, and fix red builds   |
| [`journey-docs`](journey-docs/)           | Document journeys | Create concise Markdown journey, flow, acceptance, telemetry, and coverage artifacts                             |
| [`repo-context-docs`](repo-context-docs/) | Document repo     | Create or update current-state docs under `docs/context/`                                                        |

### Integrations

| Skill                       | Workflow role | What it does                                                                        |
| --------------------------- | ------------- | ----------------------------------------------------------------------------------- |
| [`coderabbit`](coderabbit/) | Integration   | Triage CodeRabbit review comments, fix legitimate issues, push, and resolve threads |

### Personal Utilities

| Skill                         | Workflow role    | What it does                                                         |
| ----------------------------- | ---------------- | -------------------------------------------------------------------- |
| [`create-todo`](create-todo/) | Personal utility | Capture concrete work as committed Markdown todos grouped by project |

### Deprecated

| Skill       | Workflow role | What it does                                                                                                                           |
| ----------- | ------------- | -------------------------------------------------------------------------------------------------------------------------------------- |
| [`ci`](ci/) | Ship direct   | Deprecated public flagship path. Stages all changes and pushes the current branch directly. Use [`ship`](ship/) for PR-based shipping. |

## How Skills Work

Each skill is a folder with a `SKILL.md` file: YAML frontmatter plus a prompt that defines the agent's workflow. Some skills also include scripts or templates for deterministic operations.

Install the specific skills you want, then invoke them by describing the work. Examples: "review my changes", "write an implementation plan", "review my infra code", or "document this critical user journey".

## Validate Repo

```bash
scripts/validate-skills.sh
```

Run this before release-oriented skill changes. It checks skill directory shape, `SKILL.md` frontmatter names, per-skill README titles and install links, root README references, and core plugin skill paths.

## License

[MIT](LICENSE)
