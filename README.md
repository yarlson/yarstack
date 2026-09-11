# Yarstack

![Yarstack pirate coins](assets/yarstack.png)

Yarstack is a plugin for Codex and Claude Code. It adds skills for investigating
bugs, designing features, writing and testing code, reviewing changes, and
delivering pull requests.

Each skill is a Markdown file that tells the agent what to inspect, what it may
change, how to check the result, and when to stop. Name a skill in your prompt
to use its instructions with your repository and available tools.

[Get started](#get-started) · [Examples](#examples) · [Skill catalog](#skill-catalog)

## Why use it?

- `system-investigate` traces the relevant code, tests, and available runtime
  evidence before recommending a fix. It reports when the evidence does not
  support a conclusion.
- `system-design` starts with your requirements and current system. It adds a
  component only when it can name the requirement that needs it.
- `behavior-implement` looks for existing code, standard-library solutions, and
  native platform features before writing new code. It starts product changes
  with a failing test.
- `phase-validate` checks each acceptance criterion against the implementation
  and test results. It identifies criteria that remain unverified.
- `code-review` checks suspected defects against the source before reporting
  them as findings. Each finding must describe a concrete failure scenario.

You can [read every skill's instructions](plugins/yarstack/skills/) before using it.

## Get started

You need Codex or Claude Code with plugin support. Run the commands for the one
you use in your terminal.

### Codex

```sh
codex plugin marketplace add yarlson/yarstack
codex plugin add yarstack@yarstack
```

Start a new Codex thread after installation. See the
[Codex plugin command reference](https://learn.chatgpt.com/docs/developer-commands#codex-plugin)
for host setup and command details.

### Claude Code

```sh
claude plugin marketplace add yarlson/yarstack
claude plugin install yarstack@yarstack
```

Run `/reload-plugins` in Claude Code or start a new session. See
[Claude Code's plugin guide](https://code.claude.com/docs/en/discover-plugins)
for installation options.

### Try your first skill

Open a repository with local changes and ask your agent:

```text
Use Yarstack's code-review skill to review my uncommitted changes.
```

The agent reads the changed code and its callers, then checks suspected bugs
against the source. The report lists verified defects with file locations,
failure scenarios, and suggested corrections. It separates those findings from
concerns it could not verify and notes missing evidence. Reviews are read-only
unless you authorize edits.

To use another skill, replace `code-review` with its name and describe the task.
Include constraints such as the deadline or behavior that must stay unchanged.

## Examples

### Find out why a bug happens

```text
Use Yarstack's system-investigate skill to find out why retried webhook
deliveries create duplicate invoices.
```

The agent checks whether retries cause duplicate invoices and recommends a fix
if the evidence supports one. This skill does not edit code. Once you agree on
the cause and fix, ask:

```text
Use Yarstack's behavior-implement skill to fix the duplicate-invoice bug.
Add a regression test showing that retrying the same event creates one invoice.
```

For product behavior, `behavior-implement` starts with a failing test, implements
the fix, and reruns the relevant checks. For tooling and infrastructure, use the
existing validation commands first. Add dedicated tests when the remaining risk
warrants them.

### Design and implement a feature

```text
Use Yarstack's system-design skill to design per-workspace usage limits for
our billing service. We have ten workspaces. One engineer will implement
this in two weeks.
```

The agent proposes a design within the current system. It explains which
requirement each component satisfies and what would require a larger design.
It describes how data moves through the system and records tradeoffs and work
left out of scope.

When you accept the design, ask the agent to write a plan:

```text
Use Yarstack's plan-create skill to turn the agreed design into a plan at
docs/usage-limits-plan.md, with ordered phases and acceptance checks.
```

Then work through one phase at a time:

```text
Use Yarstack's phase-implement skill to implement phase 1 of docs/usage-limits-plan.md.
Use Yarstack's phase-validate skill to check phase 1 against the plan.
Use Yarstack's phase-review skill to review phase 1 for structural regressions.
```

Send one request at a time. After each step finishes, send the next one. Use
`architecture-refine` if you still need to settle a design decision.

Use a phased plan when the work needs separate implementation steps. For a
focused bug fix, start with `behavior-implement`.

### Check a specific risk

| Your question | Skill |
| --- | --- |
| Can this change let one tenant access another tenant's data? | [`security-review`](plugins/yarstack/skills/security-review/SKILL.md) |
| Do these tests prove retries and partial failure behave correctly? | [`test-gap-review`](plugins/yarstack/skills/test-gap-review/SKILL.md) |
| Can we roll this change out and recover if it fails? | [`rollout-readiness-review`](plugins/yarstack/skills/rollout-readiness-review/SKILL.md) |
| Does CI actually check the behavior we depend on? | [`ci-review`](plugins/yarstack/skills/ci-review/SKILL.md) |
| What does this branch change, and why? | [`changes-explain`](plugins/yarstack/skills/changes-explain/SKILL.md) |

For example, ask:

```text
Use Yarstack's security-review skill to check tenant isolation in this branch.
```

### Deliver a Jira issue

`jira-issue-refine` helps the agent define an issue's scope and acceptance
criteria. When you authorize delivery, `jira-issue-deliver` implements the issue,
validates the change, commits it, and pushes it to GitHub. It opens a pull request
and fixes CI failures.

It finishes when the pull request is ready for review, its required checks pass,
and requested automated review feedback is resolved. It does not merge, deploy,
release, or change the Jira issue.

For completed local changes, `pr-draft` publishes a draft pull request and
monitors its CI checks. These workflows require authenticated access to the
relevant services. You must authorize commits, pushes, and other external writes.

## Optional engineering standards

Skills apply to individual tasks. You can also install the bundled
[engineering standards](plugins/yarstack/agent-guidance/engineering-standards/)
as global guidance. They tell the agent to inspect existing code, sketch a design
before editing, preserve contracts, and test observable behavior. They also cover
safe state changes, naming, and technical writing.

The plugin does not install these global instructions. To preview the exact
text, clone this repository and run:

```sh
git clone https://github.com/yarlson/yarstack.git
cd yarstack
plugins/yarstack/scripts/install-engineering-standards.sh --print
```

If you want to use that text as your global agent guidance, run from the clone:

```sh
make install-system-prompt
```

The installer replaces the contents of `~/.codex/AGENTS.md`, `~/.agents/AGENTS.md`,
and `~/.claude/CLAUDE.md`. It first backs up files whose contents differ. It
preserves symlinks and refuses to write through dangling ones. If a file already
matches, it leaves it alone.

Repository instructions still take precedence. Your host's permissions control
the agent's access. Asking for a review does not authorize edits or external
writes.

## Optional session memory with Yarbrain

[Yarbrain](plugins/yarbrain/README.md) is a separate plugin in the same
marketplace. It stores session evidence and notes in a Markdown vault at a path
you choose. You can recall those notes during later sessions.

Episode notes preserve what happened in a session. Other notes record the current
understanding. You review proposed knowledge before adding it to those notes and
approve proposed skills before activating them.

After adding the marketplace, install Yarbrain for your host:

```sh
codex plugin add yarbrain@yarstack
```

Or, for Claude Code:

```sh
claude plugin install yarbrain@yarstack
```

Start a new Codex thread or reload Claude Code's plugins, then ask:

```text
Use Yarbrain's wiki-initialize skill to create my vault at <path>.
```

Initialization requires Python 3 and a vault path you choose. It preserves
existing vault files. Use `wiki-capture` after a session to save evidence and
propose notes or procedures. Use `wiki-reconcile` to compare those proposals with
existing knowledge. Use `wiki-recall` to find relevant notes for your next task.

After you configure a vault, hooks load up to 12,000 characters from its index
at session start. Other hooks queue session IDs and transcript paths for later
capture. Hooks do not copy transcript bodies, call models or network services,
or edit current notes. Capture and reconciliation happen when you invoke those
skills. You can use Yarstack without installing Yarbrain.

## Skill catalog

Each link opens the skill's instructions. They define what the agent may change,
which checks it must run, and when it should stop.

<details>
<summary><strong>Browse all Yarstack skills</strong></summary>

### Investigate, design, and plan

| Skill | Use it to |
| --- | --- |
| [`system-investigate`](plugins/yarstack/skills/system-investigate/SKILL.md) | Explain existing behavior and recommend a change when the evidence supports it. |
| [`system-design`](plugins/yarstack/skills/system-design/SKILL.md) | Design a system or feature around current requirements and constraints. |
| [`architecture-refine`](plugins/yarstack/skills/architecture-refine/SKILL.md) | Settle open architecture decisions before planning. |
| [`alternatives-explore`](plugins/yarstack/skills/alternatives-explore/SKILL.md) | Compare different ways to solve the problem and propose a test for the recommendation. |
| [`technical-spike`](plugins/yarstack/skills/technical-spike/SKILL.md) | Resolve one technical uncertainty about an external tool or specific version that blocks the work. |
| [`plan-create`](plugins/yarstack/skills/plan-create/SKILL.md) | Turn settled decisions into ordered implementation phases. |
| [`plan-update`](plugins/yarstack/skills/plan-update/SKILL.md) | Fix incorrect phase order, prerequisites, or checks when implementation exposes them. |
| [`test-design`](plugins/yarstack/skills/test-design/SKILL.md) | Design repeatable tests for the behavior and failure cases that matter to the change. |

### Implement and verify

| Skill | Use it to |
| --- | --- |
| [`behavior-implement`](plugins/yarstack/skills/behavior-implement/SKILL.md) | Write a failing test, implement the behavior, and refactor with checks passing. |
| [`phase-implement`](plugins/yarstack/skills/phase-implement/SKILL.md) | Implement one selected phase from a plan. |
| [`phase-validate`](plugins/yarstack/skills/phase-validate/SKILL.md) | Check a phase against its acceptance criteria and fix what fails to meet them. |
| [`phase-review`](plugins/yarstack/skills/phase-review/SKILL.md) | Check a validated phase for new problems in code structure. |
| [`go-code-reduce`](plugins/yarstack/skills/go-code-reduce/SKILL.md) | Reduce Go code without changing its behavior or making it harder to maintain. |
| [`cli-control`](plugins/yarstack/skills/cli-control/SKILL.md) | Verify CLI or TUI behavior in a real terminal. |
| [`ui-control`](plugins/yarstack/skills/ui-control/SKILL.md) | Verify browser, desktop, or Electron behavior through the interface. |

### Review a change or a specific risk

| Skill | Use it to |
| --- | --- |
| [`code-review`](plugins/yarstack/skills/code-review/SKILL.md) | Assign agents to review specific risks, then verify their findings. |
| [`claude-review`](plugins/yarstack/skills/claude-review/SKILL.md) | Run an independent Claude Opus review, then check the findings locally. |
| [`test-gap-review`](plugins/yarstack/skills/test-gap-review/SKILL.md) | Check whether tests and other evidence prove the intended behavior. |
| [`security-review`](plugins/yarstack/skills/security-review/SKILL.md) | Trace a focused trust boundary or exploit path. |
| [`dependency-review`](plugins/yarstack/skills/dependency-review/SKILL.md) | Check dependency provenance, version pinning, and reproducibility. |
| [`ci-review`](plugins/yarstack/skills/ci-review/SKILL.md) | Audit CI workflows and the checks they provide. |
| [`infra-review`](plugins/yarstack/skills/infra-review/SKILL.md) | Review infrastructure changes for mistakes in deployment targets, state, cost, availability, or recovery. |
| [`rollout-readiness-review`](plugins/yarstack/skills/rollout-readiness-review/SKILL.md) | Check deployment and recovery procedures, including rollback and how operators detect failures. |
| [`crap-index-assess`](plugins/yarstack/skills/crap-index-assess/SKILL.md) | Assess which methods are risky to change using code complexity and test coverage. |
| [`change-cleanup-review`](plugins/yarstack/skills/change-cleanup-review/SKILL.md) | Find unnecessary changes and checks that do not prove their claims before human review. |
| [`coderabbit-triage`](plugins/yarstack/skills/coderabbit-triage/SKILL.md) | Judge unresolved CodeRabbit feedback before acting on it. |

### Explain and document

| Skill | Use it to |
| --- | --- |
| [`changes-explain`](plugins/yarstack/skills/changes-explain/SKILL.md) | Explain a change's behavior, design decisions, and risks. |
| [`changes-report`](plugins/yarstack/skills/changes-report/SKILL.md) | Report changes that entered a branch during a date or period. |
| [`repo-context-document`](plugins/yarstack/skills/repo-context-document/SKILL.md) | Document how the system currently works under `docs/context/`. |
| [`critical-journey-document`](plugins/yarstack/skills/critical-journey-document/SKILL.md) | Document how one actor completes a task, with links to the code and other evidence. |
| [`docs-review`](plugins/yarstack/skills/docs-review/SKILL.md) | Audit or improve repository documentation against actual behavior. |
| [`docs-drift-review`](plugins/yarstack/skills/docs-drift-review/SKILL.md) | Find documentation made inaccurate by the current change. |
| [`text-improve`](plugins/yarstack/skills/text-improve/SKILL.md) | Rewrite technical prose for clarity while preserving meaning. |
| [`marketing-claims-review`](plugins/yarstack/skills/marketing-claims-review/SKILL.md) | Check or rewrite product claims against what the product does. |
| [`spec-update`](plugins/yarstack/skills/spec-update/SKILL.md) | Correct an inaccurate product, architecture, API, or behavior specification. |

### Refine issues and deliver work

| Skill | Use it to |
| --- | --- |
| [`jira-issue-refine`](plugins/yarstack/skills/jira-issue-refine/SKILL.md) | Clarify an issue's scope and acceptance criteria before implementation. |
| [`jira-issue-create`](plugins/yarstack/skills/jira-issue-create/SKILL.md) | Create Jira issues from approved drafts through `acli`. |
| [`jira-issue-deliver`](plugins/yarstack/skills/jira-issue-deliver/SKILL.md) | Implement a Jira issue, resolve requested automated review feedback, and open a ready pull request with passing checks. |
| [`roadmap-task-deliver`](plugins/yarstack/skills/roadmap-task-deliver/SKILL.md) | Implement the next eligible roadmap task, merge its reviewed pull request after checks pass, and update the task status. |
| [`phase-commit`](plugins/yarstack/skills/phase-commit/SKILL.md) | Create one explicitly requested local commit. |
| [`pr-draft`](plugins/yarstack/skills/pr-draft/SKILL.md) | Commit and push completed work, open a draft pull request, and follow its CI checks. |

</details>

<details>
<summary><strong>Browse Yarbrain skills</strong></summary>

| Skill | Use it to |
| --- | --- |
| [`wiki-initialize`](plugins/yarbrain/skills/wiki-initialize/SKILL.md) | Create or adopt a vault at a path you choose. |
| [`wiki-capture`](plugins/yarbrain/skills/wiki-capture/SKILL.md) | Save a session as evidence and propose notes or reusable procedures. |
| [`wiki-reconcile`](plugins/yarbrain/skills/wiki-reconcile/SKILL.md) | Compare proposed knowledge with existing notes before updating them. |
| [`wiki-recall`](plugins/yarbrain/skills/wiki-recall/SKILL.md) | Retrieve up to five relevant notes, episodes, or approved procedures. |
| [`procedure-promote`](plugins/yarbrain/skills/procedure-promote/SKILL.md) | Propose a skill for a procedure verified in repeated use. |
| [`wiki-maintain`](plugins/yarbrain/skills/wiki-maintain/SKILL.md) | Check vault links and structure. Find stale, duplicate, or conflicting notes and outdated skills. |

</details>

## Update

Refresh the marketplace and update the plugins you have installed.

### Codex

```sh
codex plugin marketplace upgrade yarstack
codex plugin add yarstack@yarstack
```

If you also use Yarbrain, run `codex plugin add yarbrain@yarstack`. Start a new
Codex thread to load the updated plugins.

### Claude Code

```sh
claude plugin marketplace update yarstack
claude plugin update yarstack@yarstack
```

If you also use Yarbrain, run `claude plugin update yarbrain@yarstack`. Run
`/reload-plugins` or restart Claude Code.

## Work on Yarstack

Each plugin is self-contained under `plugins/yarstack/` or `plugins/yarbrain/`.
Codex and Claude Code use separate manifests and marketplace catalogs. See
[AGENTS.md](AGENTS.md) for repository rules and [SKILL_FORMAT.md](SKILL_FORMAT.md)
for skill conventions.

Run the repository checks before submitting changes:

```sh
make validate
```

This runs the Yarbrain Python tests, `plugin-scanner lint` and
`plugin-scanner verify` for both Codex packages, and
`claude plugin validate --strict` for both Claude Code plugins and the marketplace.
The command requires Python 3, `pipx`, and the Claude Code CLI; the Makefile pins
the scanner version.

Report vulnerabilities privately as described in [SECURITY.md](SECURITY.md).

[MIT](LICENSE) © 2026 Yar Kravtsov.
