# Yarstack

![Yarstack pirate coins](assets/yarstack.png)

Generating a diff is not the same as engineering a change. The hard parts are
deciding what to change, keeping the design coherent, proving the behavior, and
leaving the repository easier to work in.

Yarstack packages those engineering disciplines as reusable workflows for
Codex and Claude Code. It helps an agent plan before coding, implement in small
steps, test the contract, review concrete risks, and deliver only when the
evidence supports it.

This is not a prompt that asks an agent to "be a senior engineer." It is a set
of explicit workflows and engineering standards that define what the agent must
inspect, what it may change, how it should verify the result, and where it must
stop.

## What Yarstack is built for

### Planning that can be implemented

A useful plan settles the decisions that change the implementation path. It
names ownership, interfaces, dependencies, failure cases, compatibility needs,
rollout order, and acceptance evidence. Yarstack separates architecture
discussion, technical research, plan creation, and plan correction so each has
one clear job.

### Correctness before completion

Code that works on the happy path is not enough. Yarstack makes normal paths,
boundaries, partial failure, cleanup, cancellation, ordering, retries, stale
state, and recovery part of the work when they apply. Validation checks the
implemented change against its contract instead of trusting the implementation
summary.

### Tests as product contracts

Tests should describe promised behavior, not mirror internal code. Yarstack
supports test design before implementation, red-green-refactor while changing
behavior, and a separate test-gap review to ask whether the final evidence is
actually sufficient.

### Modules that contain decisions

Yarstack treats module boundaries as an engineering tool, not a directory
layout exercise. A module should own a coherent policy, representation, or
protocol, hide the details callers do not need, and keep likely changes from
spreading through unrelated code.

### Names that make code predictable

Clear names reduce the amount of code a reader must open. Yarstack asks agents
to use consistent domain terms, match detail to scope, and name operations for
the result they return or the action they perform.

### English that says what happened

Engineering work ends in prose: comments, errors, documentation, commits, pull
requests, runbooks, and reports. Yarstack includes a plain-English standard and
a dedicated review workflow because vague writing hides vague thinking.

## One engineering system, two parts

Yarstack contains:

- **Skills** for specific jobs such as planning, implementation, validation,
  review, documentation, and delivery. Each skill has a defined scope,
  evidence standard, authority boundary, and completion condition.
- **Engineering standards** for the rules that should apply throughout the
  work: small design, codebase-first decisions, correctness, tests, module
  boundaries, contract evolution, state safety, naming, comments, and plain
  English.

Use a single skill for a focused task or compose several skills into a larger
workflow. Repository instructions still take precedence. Reviews are read-only
by default, and invoking one workflow does not grant another workflow permission
to edit, commit, push, or change external systems.

## A typical change

```text
architecture-refine  settle the decisions that change the design
plan-create          turn settled decisions into ordered implementation phases
phase-implement      implement one phase and stop
phase-validate       prove each acceptance criterion against current evidence
phase-review         find structural regressions introduced by the change
phase-commit         create one validated local commit when explicitly requested
pr-draft             publish and monitor a draft pull request when authorized
```

That sequence is available, not mandatory. A small bug may need only
`behavior-implement`. A risky infrastructure change may add `infra-review`,
`security-review`, and `rollout-readiness-review`. The task determines the
workflow.

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

The same self-contained package supplies the shared skills to both hosts while
keeping their native manifests separate.

## Choose a skill

### Decide and plan

- [`alternatives-explore`](plugins/yarstack/skills/alternatives-explore/SKILL.md)
  tests one non-incremental alternative.
- [`architecture-refine`](plugins/yarstack/skills/architecture-refine/SKILL.md)
  settles the architecture decisions needed for the next deliverable.
- [`technical-spike`](plugins/yarstack/skills/technical-spike/SKILL.md) resolves
  one blocking external or version-specific uncertainty.
- [`plan-create`](plugins/yarstack/skills/plan-create/SKILL.md) creates an
  implementation-ready phased plan.
- [`plan-update`](plugins/yarstack/skills/plan-update/SKILL.md) corrects plan
  mechanics when implementation evidence proves them wrong.
- [`test-design`](plugins/yarstack/skills/test-design/SKILL.md) designs
  reproducible tests around behavioral risk.

### Implement and prove

- [`behavior-implement`](plugins/yarstack/skills/behavior-implement/SKILL.md)
  changes production behavior through red-green-refactor.
- [`phase-implement`](plugins/yarstack/skills/phase-implement/SKILL.md)
  implements exactly one selected plan phase.
- [`phase-validate`](plugins/yarstack/skills/phase-validate/SKILL.md) checks a
  phase against every acceptance criterion and required failure case.
- [`phase-review`](plugins/yarstack/skills/phase-review/SKILL.md) reviews a
  validated phase for structural regressions.
- [`test-gap-review`](plugins/yarstack/skills/test-gap-review/SKILL.md) asks
  whether the available tests and checks prove the contract.

### Review risk

- [`code-review`](plugins/yarstack/skills/code-review/SKILL.md) reviews
  correctness, performance, compatibility, lifecycle, and maintainability.
- [`go-review`](plugins/yarstack/skills/go-review/SKILL.md) adds Go-specific
  correctness and lifecycle analysis.
- [`macos-swift-review`](plugins/yarstack/skills/macos-swift-review/SKILL.md)
  adds Swift, SwiftUI, and macOS-specific analysis.
- [`security-review`](plugins/yarstack/skills/security-review/SKILL.md) traces a
  focused trust boundary or exploit path.
- [`dependency-review`](plugins/yarstack/skills/dependency-review/SKILL.md)
  checks provenance, pinning, and reproducibility.
- [`ci-review`](plugins/yarstack/skills/ci-review/SKILL.md) audits whether CI
  protects real contracts.
- [`infra-review`](plugins/yarstack/skills/infra-review/SKILL.md) reviews
  infrastructure targeting, state, cost, availability, and recovery risk.
- [`rollout-readiness-review`](plugins/yarstack/skills/rollout-readiness-review/SKILL.md)
  checks deployment, rollback, recovery, and production failure visibility.

### Understand and document

- [`changes-explain`](plugins/yarstack/skills/changes-explain/SKILL.md) explains
  a confirmed diff through behavior, boundaries, flow, and risk.
- [`changes-report`](plugins/yarstack/skills/changes-report/SKILL.md) reports a
  period of repository changes by product outcome.
- [`repo-context-document`](plugins/yarstack/skills/repo-context-document/SKILL.md)
  maintains adopted current-state documentation under `docs/context/`.
- [`critical-journey-document`](plugins/yarstack/skills/critical-journey-document/SKILL.md)
  records one actor pursuing one goal, linked to current evidence.
- [`docs-review`](plugins/yarstack/skills/docs-review/SKILL.md) audits general
  repository documentation against implemented behavior.
- [`docs-drift-review`](plugins/yarstack/skills/docs-drift-review/SKILL.md) finds
  documentation made stale by a code or configuration change.
- [`english-text-review`](plugins/yarstack/skills/english-text-review/SKILL.md)
  reviews technical prose for plain-English clarity.
- [`marketing-claims-review`](plugins/yarstack/skills/marketing-claims-review/SKILL.md)
  checks persuasive claims against shipped-product evidence.

### Verify interfaces and deliver

- [`cli-control`](plugins/yarstack/skills/cli-control/SKILL.md) verifies CLI and
  TUI behavior in a real terminal.
- [`ui-control`](plugins/yarstack/skills/ui-control/SKILL.md) gathers direct
  evidence from graphical interfaces.
- [`spec-update`](plugins/yarstack/skills/spec-update/SKILL.md) corrects a
  semantic contract when current work cannot determine correctness.
- [`phase-commit`](plugins/yarstack/skills/phase-commit/SKILL.md) creates one
  explicitly requested local commit without pushing.
- [`pr-draft`](plugins/yarstack/skills/pr-draft/SKILL.md) publishes confirmed
  changes as a monitored draft pull request.
- [`coderabbit-triage`](plugins/yarstack/skills/coderabbit-triage/SKILL.md)
  judges unresolved review feedback before authorized remediation.

## Install the engineering standards

Skills work as soon as the plugin is installed. The shared engineering
standards are optional global guidance. Preview the exact assembled text first:

```sh
plugins/yarstack/scripts/install-engineering-standards.sh --print
```

Install it only when you want Yarstack to replace the global guidance used by
Codex and Claude Code:

```sh
make install-system-prompt
```

The installer writes to `~/.codex/AGENTS.md`, `~/.agents/AGENTS.md`, and
`~/.claude/CLAUDE.md`. It preserves symlinks, rejects dangling symlinks, creates
timestamped backups before replacing differing files, and makes no change when
the installed content is already current.

## Upgrade

Marketplace installs are cached by version. Refresh the marketplace after a new
Yarstack release.

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

Run `/reload-plugins` or restart Claude Code.

## Validate the package

```sh
make validate
```

The command runs `plugin-scanner lint` and `plugin-scanner verify` for the Codex
package, then runs `claude plugin validate --strict` for the Claude Code plugin
and marketplace.

## Repository layout

```text
.agents/plugins/marketplace.json   Codex marketplace
.claude-plugin/marketplace.json    Claude Code marketplace
plugins/yarstack/                  Shared installable package
Makefile                           Local and CI entrypoints
```

## Security

Report vulnerabilities privately as described in [SECURITY.md](SECURITY.md).

## License

[MIT](LICENSE) © 2026 Yar Kravtsov.
