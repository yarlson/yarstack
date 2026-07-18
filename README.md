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
$yarstack:draft-pr
$yarstack:coderabbit-triage
```

- `go-review` performs a read-only Go lifecycle and correctness audit.
- `draft-pr` isolates confirmed changes, pushes them, opens a draft pull
  request, and monitors its checks.
- `coderabbit-triage` judges unresolved bot feedback, fixes confirmed issues,
  and resolves handled threads.

Give the skill a concrete scope when the default is ambiguous.

## How skills work together

Use one skill when one workflow is enough. For larger changes, the skills form
a practical engineering lifecycle:

| Stage      | Skills                                      | Result                                                |
| ---------- | ------------------------------------------- | ----------------------------------------------------- |
| Understand | `plan-create`                               | Relevant code, constraints, decisions, and a plan     |
| Change     | `phase-implement`                           | One bounded implementation unit                       |
| Prove      | `phase-validate` plus focused review skills | Contract checks and concrete risk findings            |
| Close      | `final-review`                              | Final diff, verification, and remaining risks checked |
| Deliver    | `draft-pr`, then `coderabbit-triage`        | A monitored draft PR and judged review feedback       |

This composition is task-driven rather than a mandatory pipeline. For example,
a platform-specific Go fix can move from `go-review` to `technical-spike`, add
`dependency-review` and `security-review` where the change crosses those trust
boundaries, close with `final-review`, and ship with `draft-pr`.

Yarstack supplies workflow instructions, not a hosted service or autonomous
pipeline. Skills use the tools, credentials, and repository access already
available to the host, and repository instructions still take precedence.

## Skill catalogue

| Skill                                                                                 | Use it for                                                                                          |
| ------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| [`alternatives-explore`](plugins/yarstack/skills/alternatives-explore/SKILL.md)       | Generating and testing a genuinely non-incremental product, workflow, or system alternative.        |
| [`architecture-refine`](plugins/yarstack/skills/architecture-refine/SKILL.md)         | Resolving underdefined architecture through one focused decision at a time.                         |
| [`behavior-implement`](plugins/yarstack/skills/behavior-implement/SKILL.md)           | Implementing production behavior through a mandatory red-green-refactor cycle.                       |
| [`changes-explain`](plugins/yarstack/skills/changes-explain/SKILL.md)                 | Explaining a diff or PR through behavior, boundaries, data flow, operations, and risk.              |
| [`ci-review`](plugins/yarstack/skills/ci-review/SKILL.md)                             | Finding stale, redundant, unsafe, or unjustified CI and delivery machinery.                         |
| [`cli-control`](plugins/yarstack/skills/cli-control/SKILL.md)                         | Verifying interactive CLI or TUI behavior with repeatable local evidence.                           |
| [`code-review`](plugins/yarstack/skills/code-review/SKILL.md)                         | Reviewing local changes for concrete correctness, security, performance, and lifecycle defects.     |
| [`coderabbit-triage`](plugins/yarstack/skills/coderabbit-triage/SKILL.md)             | Judging and resolving CodeRabbit feedback on the current pull request.                              |
| [`critical-journey-document`](plugins/yarstack/skills/critical-journey-document/SKILL.md) | Documenting traceable user journeys, acceptance scenarios, telemetry, and service flows.        |
| [`dependency-review`](plugins/yarstack/skills/dependency-review/SKILL.md)             | Reviewing dependency, lockfile, toolchain, generated-code, and supply-chain changes.                |
| [`docs-drift-review`](plugins/yarstack/skills/docs-drift-review/SKILL.md)             | Updating documentation made inaccurate by a code or configuration change.                           |
| [`docs-review`](plugins/yarstack/skills/docs-review/SKILL.md)                         | Auditing or restructuring README and repository documentation against implemented behavior.         |
| [`go-review`](plugins/yarstack/skills/go-review/SKILL.md)                             | Finding Go-specific concurrency, cancellation, resource, retry, shutdown, and race defects.         |
| [`infra-review`](plugins/yarstack/skills/infra-review/SKILL.md)                       | Reviewing infrastructure changes for blast radius, exposure, destruction, cost, and recovery risks. |
| [`macos-swift-review`](plugins/yarstack/skills/macos-swift-review/SKILL.md)           | Auditing macOS Swift and SwiftUI lifecycle, state, concurrency, sandbox, and distribution risks.    |
| [`marketing-claims-review`](plugins/yarstack/skills/marketing-claims-review/SKILL.md) | Checking product marketing claims and rewrites against shipped-product evidence.                    |
| [`phase-commit`](plugins/yarstack/skills/phase-commit/SKILL.md)                       | Creating one explicitly requested local commit for a completed, verified phase.                     |
| [`phase-implement`](plugins/yarstack/skills/phase-implement/SKILL.md)                 | Implementing exactly one selected plan phase without later-phase work.                              |
| [`phase-review`](plugins/yarstack/skills/phase-review/SKILL.md)                       | Fixing maintainability problems in a completed phase without changing behavior or scope.            |
| [`phase-validate`](plugins/yarstack/skills/phase-validate/SKILL.md)                   | Validating one implemented phase against its stated behavioral contract.                            |
| [`plan-create`](plugins/yarstack/skills/plan-create/SKILL.md)                         | Producing an implementation-ready repository plan with decisions, phases, risks, and validation.    |
| [`plan-update`](plugins/yarstack/skills/plan-update/SKILL.md)                         | Making the smallest evidence-backed correction to an implementation plan.                           |
| [`pr-draft`](plugins/yarstack/skills/pr-draft/SKILL.md)                               | Branching, committing, pushing, and shipping confirmed changes as a monitored draft PR.             |
| [`refactor-plan`](plugins/yarstack/skills/refactor-plan/SKILL.md)                     | Designing a minimal, reversible sequence for behavior-preserving structural change.                 |
| [`repo-context-document`](plugins/yarstack/skills/repo-context-document/SKILL.md)     | Maintaining verified current-state architecture and ownership documentation under `docs/context/`.  |
| [`rollout-readiness-review`](plugins/yarstack/skills/rollout-readiness-review/SKILL.md) | Reviewing deployment readiness, runtime failure modes, rollback, and recovery.                    |
| [`security-review`](plugins/yarstack/skills/security-review/SKILL.md)                 | Reviewing changed trust boundaries, dangerous inputs, permissions, secrets, and agent surfaces.     |
| [`spec-update`](plugins/yarstack/skills/spec-update/SKILL.md)                         | Correcting an authoritative product or architecture contract when implementation exposes a gap.     |
| [`technical-spike`](plugins/yarstack/skills/technical-spike/SKILL.md)                 | Resolving one bounded, implementation-blocking technical uncertainty.                               |
| [`test-design`](plugins/yarstack/skills/test-design/SKILL.md)                         | Designing deterministic tests as readable behavioral contracts.                                    |
| [`test-gap-review`](plugins/yarstack/skills/test-gap-review/SKILL.md)                 | Finding missing behavior coverage and tests that provide false confidence.                          |
| [`ui-control`](plugins/yarstack/skills/ui-control/SKILL.md)                           | Verifying browser, desktop, Electron, or other UI behavior through the actual interface.            |

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
