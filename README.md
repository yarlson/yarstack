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

The same marketplace also publishes Yarbrain as a separate plugin. Yarbrain
keeps durable agent knowledge in a user-controlled Markdown vault without
mixing that state into Yarstack's engineering workflows.

## What Yarstack is built for

### Planning that can be implemented

A useful plan settles the decisions that change the implementation path. It
names ownership, interfaces, dependencies, failure cases, compatibility needs,
rollout order, and acceptance evidence. Yarstack separates architecture
discussion, technical research, plan creation, and plan correction so each has
one clear job.

### Start with a question or a design

Most work starts in one of two places. Either something exists and you need to
understand how it behaves, or something does not exist yet and you need to
decide how to build it. Yarstack gives each its own skill.

`system-investigate` answers questions about behavior that already exists. It
treats your prompt as a hypothesis, traces the code and runtime evidence, checks
current public sources when they can change the conclusion, and ends with a
verdict and the smallest safe next step. It changes nothing. Use it when you ask
"why does this happen", "is this really how it works", or "what is the smallest
change that fixes this".

```text
Use system-investigate: why do retried webhook deliveries create duplicate invoices?
```

`system-design` proposes how to build something new: a greenfield product or a
new capability in an existing codebase. It fixes the requirements and scale
first, starts from one process and one datastore, and adds a component only
when a named requirement fails without it. It reports each component with the
requirement it serves, the concerns it deferred and what would trigger them, and
the simpler alternative it rejected.

```text
Use system-design: add per-workspace usage limits to the billing service.
Ten workspaces today, one engineer, ship in two weeks.
```

Give `system-design` the numbers you know: users, load, data volume, team size,
and deadline. Where you leave a value out, it assumes the smallest plausible one
and says so. If the result is bigger than you expected, ask which requirement
each component serves.

The two skills hand off to each other. `system-investigate` stops and points to
`system-design` when the question is about behavior that does not exist yet.
`system-design` uses `system-investigate` when it must understand the code it
extends. Both stop before `architecture-refine`, which settles the decisions
only you can make, and `plan-create`, which turns an accepted design into
implementation phases.

### The smallest change that works

Before writing code, Yarstack has the agent stop at the first option that
holds: the behavior does not need to exist, the codebase already has it, the
standard library or platform already does it, an installed dependency covers
it, or it fits in one clear line. Only then does it write new code, in the
fewest files, deleting before adding. Validation at trust boundaries,
data-loss handling, security, and accessibility are never cut to save lines.
Tests start from one check that fails when the logic breaks, and reports stay
short: what changed, what was left out and when to add it, what was checked.

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
actually sufficient. Changes to engineering tooling and infrastructure use
native checks unless their complexity or failure risk warrants focused
automated tests.

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

### Jira issues refined and delivered

`jira-issue-refine` turns a draft or existing Jira issue into a clear, bounded
delivery contract. It uses `system-investigate` to resolve discoverable facts
and keeps unresolved decisions, dependencies, safety needs, and verification
gaps visible. `jira-issue-deliver` can then own an explicitly authorized change
from Jira intake through implementation, validation, a ready GitHub pull
request, passing CI, and requested automated review. It stops before merging,
deploying, releasing, or changing the Jira issue.

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
system-design        propose the smallest design that meets stated requirements
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

## Yarbrain

[Yarbrain](plugins/yarbrain/README.md) is a separate plugin for a review-first
second brain. It separates what happened from what is currently believed and
from what an agent knows how to do:

- immutable episode notes preserve session evidence;
- canonical notes hold reconciled current knowledge;
- approved `SKILL.md` files hold repeatable procedures.

The Markdown vault is the source of truth. Search indexes and caches can be
rebuilt. Proposed memory and skill changes stay in an inbox until they are
reviewed.

| Skill                                                                     | Use it to                                                                                |
| ------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| [`wiki-initialize`](plugins/yarbrain/skills/wiki-initialize/SKILL.md)     | Create or adopt a vault at a path you approve.                                           |
| [`wiki-capture`](plugins/yarbrain/skills/wiki-capture/SKILL.md)           | Turn a completed session into one episode and reviewable candidates.                     |
| [`wiki-reconcile`](plugins/yarbrain/skills/wiki-reconcile/SKILL.md)       | Search existing notes and propose a merge, split, supersession, conflict, or new note.   |
| [`wiki-recall`](plugins/yarbrain/skills/wiki-recall/SKILL.md)             | Read at most five relevant notes, episodes, or approved procedures before a task.        |
| [`procedure-promote`](plugins/yarbrain/skills/procedure-promote/SKILL.md) | Turn a repeated, verified procedure into a proposed Agent Skill.                         |
| [`wiki-maintain`](plugins/yarbrain/skills/wiki-maintain/SKILL.md)         | Check links, provenance, duplicates, stale knowledge, pending evidence, and skill drift. |

After installation, ask the agent:

```text
Use wiki-initialize to create my Yarbrain vault at <path>.
```

Yarbrain stores the selected path in `$XDG_CONFIG_HOME/yarbrain/config.json`,
or `~/.config/yarbrain/config.json` when `XDG_CONFIG_HOME` is unset. Set
`YARBRAIN_CONFIG` to use another config file. Initialization preserves existing
vault files and creates only missing structure.

Until a vault is configured, Yarbrain hooks do nothing. After configuration:

- `SessionStart` loads at most 12,000 characters from `INDEX.md` and reports
  the pending-session count;
- `PreCompact` and `SessionEnd` enqueue session IDs and transcript locators;
- hooks do not copy transcript bodies, call a model or network service, or
  change canonical notes.

Capture and reconciliation remain explicit agent workflows. Yarbrain does not
let new evidence replace current knowledge without review, and it does not
activate a proposed procedure without approval.

## Install

Install either plugin or both. They share a marketplace but keep separate
runtime state and responsibilities.

### Codex

```sh
codex plugin marketplace add yarlson/yarstack
codex plugin add yarstack@yarstack
codex plugin add yarbrain@yarstack
```

### Claude Code

```sh
claude plugin marketplace add yarlson/yarstack
claude plugin install yarstack@yarstack
claude plugin install yarbrain@yarstack
```

Both self-contained packages keep their Codex and Claude Code manifests
separate.

## Choose a skill

### Decide and plan

- [`system-investigate`](plugins/yarstack/skills/system-investigate/SKILL.md)
  validates suspected behavior or design and recommends the smallest safe next
  step from repository, runtime, and current public evidence.
- [`system-design`](plugins/yarstack/skills/system-design/SKILL.md) proposes
  the smallest system or feature design that meets stated requirements, with
  non-goals and growth triggers.
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
  reproducible tests when behavioral risk warrants dedicated test code.

### Implement and prove

- [`behavior-implement`](plugins/yarstack/skills/behavior-implement/SKILL.md)
  changes behavior through red-green-refactor when focused tests are
  proportionate.
- [`phase-implement`](plugins/yarstack/skills/phase-implement/SKILL.md)
  implements exactly one selected plan phase.
- [`phase-validate`](plugins/yarstack/skills/phase-validate/SKILL.md) checks a
  phase against every acceptance criterion and required failure case.
- [`phase-review`](plugins/yarstack/skills/phase-review/SKILL.md) reviews a
  validated phase for structural regressions.
- [`test-gap-review`](plugins/yarstack/skills/test-gap-review/SKILL.md) asks
  whether the available tests and checks prove the contract.

### Review risk

- [`claude-review`](plugins/yarstack/skills/claude-review/SKILL.md) runs an
  independent `code-review` through Claude Opus against the repository's
  default branch, then verifies the findings locally.
- [`code-review`](plugins/yarstack/skills/code-review/SKILL.md) reviews
  changes through adaptive context tracks, risk-focused reviewer lanes,
  adversarial challenge, and verified findings.
- [`crap-index-assess`](plugins/yarstack/skills/crap-index-assess/SKILL.md)
  assesses method-level change risk from complexity and automated test coverage.
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
- [`text-improve`](plugins/yarstack/skills/text-improve/SKILL.md)
  rewrites technical prose for plain-language clarity.
- [`marketing-claims-review`](plugins/yarstack/skills/marketing-claims-review/SKILL.md)
  checks persuasive claims against shipped-product evidence.

### Verify interfaces and deliver

- [`cli-control`](plugins/yarstack/skills/cli-control/SKILL.md) verifies CLI and
  TUI behavior in a real terminal.
- [`ui-control`](plugins/yarstack/skills/ui-control/SKILL.md) gathers direct
  evidence from graphical interfaces.
- [`spec-update`](plugins/yarstack/skills/spec-update/SKILL.md) corrects a
  semantic contract when current work cannot determine correctness.
- [`jira-issue-create`](plugins/yarstack/skills/jira-issue-create/SKILL.md)
  creates user-approved Jira issues with heading-based Context, Acceptance
  criteria, and Engineering notes sections.
- [`jira-issue-refine`](plugins/yarstack/skills/jira-issue-refine/SKILL.md)
  turns a draft or existing Jira issue into a clear, evidence-backed delivery
  contract before engineering handoff.
- [`jira-issue-deliver`](plugins/yarstack/skills/jira-issue-deliver/SKILL.md)
  delivers an authorized Jira issue to a ready, green, reviewed pull request.
- [`roadmap-task-deliver`](plugins/yarstack/skills/roadmap-task-deliver/SKILL.md)
  selects the next eligible roadmap task, delivers and merges its green reviewed
  PR, then updates the task and commits its status when stored in Git.
- [`phase-commit`](plugins/yarstack/skills/phase-commit/SKILL.md) creates one
  explicitly requested local commit without pushing.
- [`pr-draft`](plugins/yarstack/skills/pr-draft/SKILL.md) publishes confirmed
  changes as a monitored draft pull request.
- [`change-cleanup-review`](plugins/yarstack/skills/change-cleanup-review/SKILL.md)
  finds unwanted scope, redundant or residual changes, low-signal prose, and
  hollow proof that shift work to reviewers.
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
repository release, then run the relevant plugin update commands below.

### Codex

```sh
codex plugin marketplace upgrade yarstack
codex plugin add yarstack@yarstack
codex plugin add yarbrain@yarstack
```

Start a new Codex thread to load the updated plugin.

### Claude Code

```sh
claude plugin marketplace update yarstack
claude plugin update yarstack@yarstack
claude plugin update yarbrain@yarstack
```

Run `/reload-plugins` or restart Claude Code.

## Validate the packages

```sh
make validate
```

The command runs `plugin-scanner lint` and `plugin-scanner verify` for both
Codex packages, then runs `claude plugin validate --strict` for both Claude Code
plugins and the marketplace.

## Repository layout

```text
.agents/plugins/marketplace.json   Codex marketplace
.claude-plugin/marketplace.json    Claude Code marketplace
plugins/yarstack/                  Engineering workflow package
plugins/yarbrain/                  Markdown second-brain package
Makefile                           Local and CI entrypoints
```

## Security

Report vulnerabilities privately as described in [SECURITY.md](SECURITY.md).

## License

[MIT](LICENSE) © 2026 Yar Kravtsov.
