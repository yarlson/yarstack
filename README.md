# Yarstack

![Yarstack pirate coins](assets/yarstack.png)

Yarstack is a portable set of engineering workflows for teams using Codex and
Claude Code on real repositories. It gives agents bounded playbooks for
understanding a codebase, planning and implementing scoped changes, reviewing
risk, validating behavior, and shipping a draft pull request.

Install the same self-contained skill package in either host, then choose the
workflow that matches the job. Each skill defines its scope, evidence
requirements, checks, and guardrails so the agent does not need to invent a new
engineering process for every session.

## Install

### Codex

```sh
codex plugin marketplace add yarlson/yarstack
codex plugin add yarstack@yarstack
```

### Claude Code

```sh
claude plugin marketplace add yarlson/yarstack
claude plugin install yarstack@yarstack
```

## Upgrade

Marketplace installs are cached by version. After Yarstack publishes a new
version, refresh the marketplace and reinstall or update the plugin.

### Codex

```sh
codex plugin marketplace upgrade yarstack
codex plugin add yarstack@yarstack
```

Start a new Codex thread to load the updated plugin.

### Claude Code

```sh
claude plugin marketplace update yarstack
claude plugin update yarstack@yarstack
```

Run `/reload-plugins` or restart Claude Code. Marketplace auto-update can also
be enabled from the `/plugin` interface.

## Start with an outcome

In Codex, invoke a skill from the repository you want it to work on. Three
practical entry points are:

```text
$yarstack:go-review the whole codebase
$yarstack:pr-draft
$yarstack:coderabbit-triage
```

- `go-review` performs a read-only Go lifecycle and correctness audit.
- `pr-draft` isolates confirmed changes, pushes them, opens a draft pull
  request, and monitors its checks.
- `coderabbit-triage` judges unresolved bot feedback before any separately
  authorized remediation or thread closure.

Give the skill a concrete scope when the default is ambiguous.

## How skills work together

Use one skill when one workflow is enough. For larger changes, the skills form
a practical engineering lifecycle:

| Stage      | Skills                                                    | Result                                             |
| ---------- | --------------------------------------------------------- | -------------------------------------------------- |
| Decide     | `architecture-refine`, `technical-spike`, `plan-create`    | Settled decisions and an implementation-ready plan |
| Change     | `phase-implement`, with `behavior-implement` where needed  | One bounded implementation unit                    |
| Prove      | `phase-validate` plus focused review skills                | Criterion evidence and concrete risk findings      |
| Improve    | `phase-review`                                            | Scoped structural regressions resolved             |
| Deliver    | `phase-commit` or `pr-draft`, then `coderabbit-triage`     | Isolated local or monitored remote delivery        |

This composition is task-driven rather than a mandatory pipeline. For example,
a platform-specific Go fix can move from `go-review` to `technical-spike`, add
`dependency-review` and `security-review` where the change crosses those trust
boundaries, validate with `phase-validate`, and ship with `pr-draft`.

Yarstack supplies workflow instructions, not a hosted service or autonomous
pipeline. Skills use the tools, credentials, and repository access already
available to the host, and repository instructions still take precedence.
Reviews are read-only by default. Delegating to another skill does not expand
scope, mutation authority, network access, or permission to spawn agents.

## Skill catalogue

| Skill                                                                                 | Use it for                                                                                          |
| ------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| [`alternatives-explore`](plugins/yarstack/skills/alternatives-explore/SKILL.md)       | Recommending one non-incremental bet and a bounded falsifiable experiment.                           |
| [`architecture-refine`](plugins/yarstack/skills/architecture-refine/SKILL.md)         | Resolving only the architecture decisions needed for the next deliverable.                          |
| [`behavior-implement`](plugins/yarstack/skills/behavior-implement/SKILL.md)           | Implementing testable production behavior through red-green-refactor.                               |
| [`changes-explain`](plugins/yarstack/skills/changes-explain/SKILL.md)                 | Explaining a confirmed diff through behavior, boundaries, flow, operations, and risk.               |
| [`ci-review`](plugins/yarstack/skills/ci-review/SKILL.md)                             | Auditing CI topology and whether automation protects real contracts.                                |
| [`cli-control`](plugins/yarstack/skills/cli-control/SKILL.md)                         | Gathering repeatable evidence from terminal-hosted CLI and TUI surfaces.                            |
| [`code-review`](plugins/yarstack/skills/code-review/SKILL.md)                         | Performing general review and routing meaningful Go or Swift scopes to specialists.                 |
| [`coderabbit-triage`](plugins/yarstack/skills/coderabbit-triage/SKILL.md)             | Dispositioning CodeRabbit feedback before authorized remediation and closure.                       |
| [`critical-journey-document`](plugins/yarstack/skills/critical-journey-document/SKILL.md) | Creating one evidence-linked actor-goal journey record.                                          |
| [`dependency-review`](plugins/yarstack/skills/dependency-review/SKILL.md)             | Reviewing dependency provenance, pinning, reproducibility, and generated state.                     |
| [`docs-drift-review`](plugins/yarstack/skills/docs-drift-review/SKILL.md)             | Finding and, when authorized, repairing documentation drift from the current change.                |
| [`docs-review`](plugins/yarstack/skills/docs-review/SKILL.md)                         | Auditing or improving general repository documentation in an explicit mode.                         |
| [`go-review`](plugins/yarstack/skills/go-review/SKILL.md)                             | Running the required Go-specific pass for meaningful Go review scopes.                              |
| [`infra-review`](plugins/yarstack/skills/infra-review/SKILL.md)                       | Reviewing IaC targeting, replacement, state, availability, cost, and recovery blast radius.         |
| [`macos-swift-review`](plugins/yarstack/skills/macos-swift-review/SKILL.md)           | Running the required Swift and macOS-specific pass for meaningful Swift review scopes.              |
| [`marketing-claims-review`](plugins/yarstack/skills/marketing-claims-review/SKILL.md) | Auditing or rewriting persuasive claims against shipped-product evidence.                           |
| [`phase-commit`](plugins/yarstack/skills/phase-commit/SKILL.md)                       | Creating one explicitly requested isolated local commit without pushing.                            |
| [`phase-implement`](plugins/yarstack/skills/phase-implement/SKILL.md)                 | Implementing exactly one selected plan phase without later-phase work.                              |
| [`phase-review`](plugins/yarstack/skills/phase-review/SKILL.md)                       | Reviewing a validated phase for structural regressions introduced by the change.                    |
| [`phase-validate`](plugins/yarstack/skills/phase-validate/SKILL.md)                   | Mapping every phase criterion to current implementation evidence.                                   |
| [`plan-create`](plugins/yarstack/skills/plan-create/SKILL.md)                         | Synthesizing an implementation-ready phased plan from settled decisions.                            |
| [`plan-update`](plugins/yarstack/skills/plan-update/SKILL.md)                         | Correcting implementation-plan mechanics without changing semantic contracts.                       |
| [`pr-draft`](plugins/yarstack/skills/pr-draft/SKILL.md)                               | Publishing confirmed and validated work as a monitored draft PR.                                   |
| [`repo-context-document`](plugins/yarstack/skills/repo-context-document/SKILL.md)     | Maintaining adopted current-state documentation under `docs/context/`.                              |
| [`rollout-readiness-review`](plugins/yarstack/skills/rollout-readiness-review/SKILL.md) | Reviewing cross-cutting operational readiness before deployment.                                |
| [`security-review`](plugins/yarstack/skills/security-review/SKILL.md)                 | Tracing focused trust boundaries and concrete exploit paths.                                        |
| [`spec-update`](plugins/yarstack/skills/spec-update/SKILL.md)                         | Correcting a non-plan semantic contract when current work cannot determine correctness.             |
| [`technical-spike`](plugins/yarstack/skills/technical-spike/SKILL.md)                 | Resolving one correctness-blocking external or version-specific uncertainty.                        |
| [`test-design`](plugins/yarstack/skills/test-design/SKILL.md)                         | Designing or writing a prospective deterministic test contract.                                    |
| [`test-gap-review`](plugins/yarstack/skills/test-gap-review/SKILL.md)                 | Auditing whether current verification credibly proves a scoped behavior contract.                   |
| [`ui-control`](plugins/yarstack/skills/ui-control/SKILL.md)                           | Gathering direct evidence from graphical browser, desktop, or Electron surfaces.                    |

## Optional global engineering standards

The skills work after plugin installation. Separately, Yarstack includes shared
engineering guidance that can be installed into the global files used by Codex
and Claude Code. Preview the exact combined document without writing anything:

```sh
plugins/yarstack/scripts/install-engineering-standards.sh --print
```

Install it only when you want to replace the guidance in `~/.codex/AGENTS.md`,
`~/.agents/AGENTS.md`, and `~/.claude/CLAUDE.md`:

```sh
make install-system-prompt
```

The installer is idempotent, preserves symlinks, refuses dangling symlinks, and
creates timestamped backups before replacing differing existing guidance. The
canonical policy fragments remain under
`plugins/yarstack/agent-guidance/engineering-standards/`.

## Repository layout

```text
.agents/plugins/marketplace.json   Codex marketplace
.claude-plugin/marketplace.json    Claude Code marketplace
plugins/yarstack/                  Shared installable package
Makefile                           Local and CI entrypoints
```

The package contains separate `.codex-plugin/plugin.json` and
`.claude-plugin/plugin.json` manifests. Reusable workflows live once under
`plugins/yarstack/skills/`; host-specific adapters should be added only when a
shared implementation cannot express the required behavior.

## Validation

Run the complete read-only validation suite locally:

```sh
make validate
```

Every push and pull request runs:

- `plugin-scanner lint` and `plugin-scanner verify` for the Codex package.
- `claude plugin validate --strict` for the plugin and marketplace.

## Security

Report vulnerabilities privately as described in [SECURITY.md](SECURITY.md).

## License

[MIT](LICENSE) © 2026 Yar Kravtsov.
