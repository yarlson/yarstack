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

Individual skill names use the `yarstack-*` namespace to avoid collisions with generic installed skills from other catalogs.

Core:

```bash
npx skills add https://github.com/yarlson/skills/tree/main/yarstack-architecture-sparring
npx skills add https://github.com/yarlson/skills/tree/main/yarstack-implementation-planning
npx skills add https://github.com/yarlson/skills/tree/main/yarstack-maintainable-implementation
npx skills add https://github.com/yarlson/skills/tree/main/yarstack-code-review
npx skills add https://github.com/yarlson/skills/tree/main/yarstack-infra-review
npx skills add https://github.com/yarlson/skills/tree/main/yarstack-draft-pr-shipping
npx skills add https://github.com/yarlson/skills/tree/main/yarstack-critical-journey-docs
npx skills add https://github.com/yarlson/skills/tree/main/yarstack-repo-context-docs
```

Integrations:

```bash
npx skills add https://github.com/yarlson/skills/tree/main/yarstack-coderabbit-triage
```

Personal utilities:

```bash
npx skills add https://github.com/yarlson/skills/tree/main/yarstack-create-todo
```

Deprecated:

```bash
npx skills add https://github.com/yarlson/skills/tree/main/yarstack-direct-ci-push
```

`yarstack-direct-ci-push` is deprecated from the public flagship path because it stages all changes and pushes the current branch directly. Use [`yarstack-draft-pr-shipping`](yarstack-draft-pr-shipping/) for PR-based shipping; keep `yarstack-direct-ci-push` only for intentional direct-push flows.

## Current Skill Map

### Core

| Skill                                                                           | Workflow role          | What it does                                                                                                     |
| ------------------------------------------------------------------------------- | ---------------------- | ---------------------------------------------------------------------------------------------------------------- |
| [`yarstack-architecture-sparring`](yarstack-architecture-sparring/)             | Architecture decisions | Ask one focused question at a time to clarify product, system, and runtime design                                |
| [`yarstack-implementation-planning`](yarstack-implementation-planning/)         | Implementation plans   | Write implementation plans that prevent structural regressions, wrong-layer logic, file sprawl, and vague phases |
| [`yarstack-maintainable-implementation`](yarstack-maintainable-implementation/) | Maintainable coding    | Implement code under an extremely strict maintainability bar so strict review finds no structural blockers       |
| [`yarstack-code-review`](yarstack-code-review/)                                 | Code review            | Review code for security, bugs, performance, regressions, and missing tests                                      |
| [`yarstack-infra-review`](yarstack-infra-review/)                               | Infrastructure review  | Review IaC for network exposure, IAM, destructive changes, cost, and deploy risk                                 |
| [`yarstack-draft-pr-shipping`](yarstack-draft-pr-shipping/)                     | Draft PR shipping      | Check scope, branch, stage confirmed paths, commit, push, open a draft PR, wait for checks, and fix red builds   |
| [`yarstack-critical-journey-docs`](yarstack-critical-journey-docs/)             | Journey docs           | Create concise Markdown journey, flow, acceptance, telemetry, and coverage artifacts                             |
| [`yarstack-repo-context-docs`](yarstack-repo-context-docs/)                     | Repo context docs      | Create or update current-state docs under `docs/context/`                                                        |

### Integrations

| Skill                                                       | Workflow role | What it does                                                                        |
| ----------------------------------------------------------- | ------------- | ----------------------------------------------------------------------------------- |
| [`yarstack-coderabbit-triage`](yarstack-coderabbit-triage/) | Integration   | Triage CodeRabbit review comments, fix legitimate issues, push, and resolve threads |

### Personal Utilities

| Skill                                           | Workflow role    | What it does                                                         |
| ----------------------------------------------- | ---------------- | -------------------------------------------------------------------- |
| [`yarstack-create-todo`](yarstack-create-todo/) | Personal utility | Capture concrete work as committed Markdown todos grouped by project |

### Deprecated

| Skill                                                 | Workflow role | What it does                                                                                                                                                                       |
| ----------------------------------------------------- | ------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [`yarstack-direct-ci-push`](yarstack-direct-ci-push/) | Ship direct   | Deprecated public flagship path. Stages all changes and pushes the current branch directly. Use [`yarstack-draft-pr-shipping`](yarstack-draft-pr-shipping/) for PR-based shipping. |

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
