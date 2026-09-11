# Yarstack

**Give your coding agent a repeatable way to engineer changes.**

Yarstack is a plugin for **Codex and Claude Code** with reusable skills for
investigating bugs, designing features, writing and testing code, reviewing
changes, and delivering pull requests. Use it when you want to spend less time
explaining how work should be done and more time deciding what to build.

Each skill is a Markdown playbook: what the agent should inspect, what it may
change, how to check the result, and when the job is done. Install the plugin,
name a skill in your prompt, and use it in your existing repository.

<img src="assets/yarstack.png" alt="Yarstack pirate coins" width="560">

[Get started](#get-started) · [Examples](#put-it-to-work) · [Skill catalog](#skill-catalog)

## Why use it?

AI coding gets harder when a task crosses files, touches existing behavior, or
needs more proof than a passing happy-path test. Yarstack gives those tasks a
repeatable process:

- **Start from the code you have.** Investigation traces callers, tests, and
  runtime evidence before recommending a change. Design starts from the
  existing system and your actual constraints.
- **Keep the diff worth reviewing.** Implementation looks for existing code,
  standard-library solutions, and native platform features before adding more.
  Every addition needs a current reason to exist.
- **Make “done” inspectable.** Plans name acceptance evidence. Implementation
  checks behavior. Validation compares the result with the requirements and
  calls out what remains unverified.
- **Get reviews you can act on.** Code review follows concrete failure paths,
  checks findings against the source, and separates verified defects from
  unresolved concerns.

The workflows run through your agent and its available tools. Their instructions
are [readable in this repository](plugins/yarstack/skills/), so you can inspect
exactly what you are asking the agent to do.

## Get started

You need Codex or Claude Code with plugin support. Choose your host and run these
commands in your terminal.

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

The review workflow traces the affected code, examines relevant risks, and
verifies findings before reporting them. Expect file locations, concrete failure
scenarios, and suggested corrections, with missing evidence called out. Reviews
are read-only by default.

Use the same prompt pattern for any skill below: **name the skill, describe the
task, and give the constraints that matter.**

## Put it to work

### Find out why a bug happens

```text
Use Yarstack's system-investigate skill: why do retried webhook deliveries
create duplicate invoices?
```

This asks the agent to test the suspected explanation against the code and
available runtime evidence, then recommend the smallest safe next step. It stops
before editing. Once the cause and intended fix are clear:

```text
Use Yarstack's behavior-implement skill to fix the duplicate-invoice bug.
Add a regression test showing that retrying the same event creates one invoice.
```

For product behavior, `behavior-implement` starts with a failing test, implements
the fix, and reruns the relevant checks. Tooling and infrastructure changes use
native checks first when dedicated tests would add little evidence.

### Design a feature without growing the system unnecessarily

```text
Use Yarstack's system-design skill: add per-workspace usage limits to our
billing service. Ten workspaces today, one engineer, ship in two weeks.
```

The design starts inside the existing system. Each proposed component must serve
a stated requirement. The result includes the data flow, tradeoffs, excluded
scope, and the conditions that would justify a larger design.

When you accept the design, turn it into work:

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

Send these as separate requests as each step finishes. Use `architecture-refine`
when a design decision is still open. A focused bug fix can start directly with
`behavior-implement`; a phased plan is useful when the work needs one.

### Check a specific risk

Name the question you need answered:

| Your question | Skill |
| --- | --- |
| Can this change let one tenant access another tenant's data? | [`security-review`](plugins/yarstack/skills/security-review/SKILL.md) |
| Do these tests prove retries and partial failure behave correctly? | [`test-gap-review`](plugins/yarstack/skills/test-gap-review/SKILL.md) |
| Can we roll this change out and recover if it fails? | [`rollout-readiness-review`](plugins/yarstack/skills/rollout-readiness-review/SKILL.md) |
| Does CI actually check the behavior we depend on? | [`ci-review`](plugins/yarstack/skills/ci-review/SKILL.md) |
| What does this branch change, and why? | [`changes-explain`](plugins/yarstack/skills/changes-explain/SKILL.md) |

For example: “Use Yarstack's security-review skill to check tenant isolation in
this branch.”

### Take an issue through to a pull request

Use `jira-issue-refine` to turn a vague ticket into a bounded task with acceptance
criteria. With explicit delivery authority, `jira-issue-deliver` can take that
issue through implementation, validation, commits, pushes, CI fixes, and requested
automated review to a ready GitHub pull request. It stops before merging,
deploying, releasing, or changing the Jira issue.

For already completed local work, use `pr-draft` to publish a draft PR and follow
its checks. These workflows need access to the relevant services and explicit
authority for their writes. The skill does not provide accounts or credentials.

## Optional: engineering standards across your work

Skills define how to do a particular job. The bundled
[engineering standards](plugins/yarstack/agent-guidance/engineering-standards/)
define expectations across jobs: inspect existing code first, sketch a small
design, preserve contracts, handle state changes safely, test observable behavior,
and write clear names and explanations.

**Installing the plugin makes the skills available. Installing global standards
is a separate, optional step.** To preview them, clone this repository and run:

```sh
git clone https://github.com/yarlson/yarstack.git
cd yarstack
plugins/yarstack/scripts/install-engineering-standards.sh --print
```

If you want to use that text as your global agent guidance, run from the clone:

```sh
make install-system-prompt
```

This replaces the contents of `~/.codex/AGENTS.md`, `~/.agents/AGENTS.md`, and
`~/.claude/CLAUDE.md`. The installer backs up differing files before replacing
them, preserves existing symlinks, rejects dangling symlinks, and leaves matching
content unchanged.

Repository instructions still take precedence. Skills and standards guide the
agent's work; your host's permissions control its access. Invoking a review does
not authorize edits, commits, pushes, or external writes.

## Optional: carry knowledge between sessions with Yarbrain

[Yarbrain](plugins/yarbrain/README.md) is a separate plugin in the same
marketplace. Use it to keep useful findings and procedures in a Markdown vault
you control, ready to recall during later work.

It keeps session evidence, current knowledge, and approved procedures separate.
New knowledge and proposed skills stay reviewable before they become part of
what the agent relies on.

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

Initialization requires Python 3 and an explicitly chosen vault path. It
preserves existing vault files. After setup, use `wiki-capture` to save session
evidence and propose knowledge, `wiki-reconcile` to review it against existing
notes, and `wiki-recall` to retrieve relevant knowledge for a new task.

Hooks load a bounded index and queue session locators. Capture and reconciliation
remain explicit workflows; hooks do not copy transcript bodies, call a model or
network service, or change canonical notes. Yarstack works independently of
Yarbrain.

## Skill catalog

Start with the skill that matches your task. Each link opens its full
instructions, including scope, required evidence, and completion conditions.

<details>
<summary><strong>Browse all Yarstack skills</strong></summary>

### Investigate, design, and plan

| Skill | Use it to |
| --- | --- |
| [`system-investigate`](plugins/yarstack/skills/system-investigate/SKILL.md) | Explain existing behavior and recommend a next step from evidence. |
| [`system-design`](plugins/yarstack/skills/system-design/SKILL.md) | Design a system or feature around current requirements and constraints. |
| [`architecture-refine`](plugins/yarstack/skills/architecture-refine/SKILL.md) | Settle open architecture decisions before planning. |
| [`alternatives-explore`](plugins/yarstack/skills/alternatives-explore/SKILL.md) | Explore substantially different approaches and a way to test the recommendation. |
| [`technical-spike`](plugins/yarstack/skills/technical-spike/SKILL.md) | Resolve one blocking external or version-specific technical uncertainty. |
| [`plan-create`](plugins/yarstack/skills/plan-create/SKILL.md) | Turn settled decisions into ordered implementation phases. |
| [`plan-update`](plugins/yarstack/skills/plan-update/SKILL.md) | Correct a plan when implementation evidence reveals a defect. |
| [`test-design`](plugins/yarstack/skills/test-design/SKILL.md) | Design reproducible tests matched to behavioral risk. |

### Implement and verify

| Skill | Use it to |
| --- | --- |
| [`behavior-implement`](plugins/yarstack/skills/behavior-implement/SKILL.md) | Change behavior through a focused failing-test, implementation, and refactoring cycle. |
| [`phase-implement`](plugins/yarstack/skills/phase-implement/SKILL.md) | Implement one selected phase from a plan. |
| [`phase-validate`](plugins/yarstack/skills/phase-validate/SKILL.md) | Check a phase against its acceptance criteria and fix required gaps. |
| [`phase-review`](plugins/yarstack/skills/phase-review/SKILL.md) | Review a validated phase for structural regressions. |
| [`go-code-reduce`](plugins/yarstack/skills/go-code-reduce/SKILL.md) | Find or apply maintainable Go code reductions while preserving behavior. |
| [`cli-control`](plugins/yarstack/skills/cli-control/SKILL.md) | Verify CLI or TUI behavior in a real terminal. |
| [`ui-control`](plugins/yarstack/skills/ui-control/SKILL.md) | Verify browser, desktop, or Electron behavior through the interface. |

### Review a change or a specific risk

| Skill | Use it to |
| --- | --- |
| [`code-review`](plugins/yarstack/skills/code-review/SKILL.md) | Run an adaptive multi-agent review and verify its findings. |
| [`claude-review`](plugins/yarstack/skills/claude-review/SKILL.md) | Run an independent Claude Opus review, then check the findings locally. |
| [`test-gap-review`](plugins/yarstack/skills/test-gap-review/SKILL.md) | Assess whether tests and checks actually prove the intended behavior. |
| [`security-review`](plugins/yarstack/skills/security-review/SKILL.md) | Trace a focused trust boundary or exploit path. |
| [`dependency-review`](plugins/yarstack/skills/dependency-review/SKILL.md) | Check dependency provenance, version pinning, and reproducibility. |
| [`ci-review`](plugins/yarstack/skills/ci-review/SKILL.md) | Audit CI workflows and the checks they provide. |
| [`infra-review`](plugins/yarstack/skills/infra-review/SKILL.md) | Review infrastructure targeting, state, availability, cost, and recovery risks. |
| [`rollout-readiness-review`](plugins/yarstack/skills/rollout-readiness-review/SKILL.md) | Check deployment, rollback, recovery, and failure visibility. |
| [`crap-index-assess`](plugins/yarstack/skills/crap-index-assess/SKILL.md) | Assess method-level change risk using complexity and test coverage. |
| [`change-cleanup-review`](plugins/yarstack/skills/change-cleanup-review/SKILL.md) | Find unnecessary changes and weak verification before human review. |
| [`coderabbit-triage`](plugins/yarstack/skills/coderabbit-triage/SKILL.md) | Judge unresolved CodeRabbit feedback before acting on it. |

### Explain and document

| Skill | Use it to |
| --- | --- |
| [`changes-explain`](plugins/yarstack/skills/changes-explain/SKILL.md) | Explain a diff through behavior, design, and risk. |
| [`changes-report`](plugins/yarstack/skills/changes-report/SKILL.md) | Report changes that entered a branch during a date or period. |
| [`repo-context-document`](plugins/yarstack/skills/repo-context-document/SKILL.md) | Maintain implemented-system documentation under `docs/context/`. |
| [`critical-journey-document`](plugins/yarstack/skills/critical-journey-document/SKILL.md) | Document one actor pursuing one goal, linked to evidence. |
| [`docs-review`](plugins/yarstack/skills/docs-review/SKILL.md) | Audit or improve repository documentation against actual behavior. |
| [`docs-drift-review`](plugins/yarstack/skills/docs-drift-review/SKILL.md) | Find documentation made inaccurate by the current change. |
| [`text-improve`](plugins/yarstack/skills/text-improve/SKILL.md) | Rewrite technical prose for clarity while preserving meaning. |
| [`marketing-claims-review`](plugins/yarstack/skills/marketing-claims-review/SKILL.md) | Check or rewrite product claims against shipped capabilities. |
| [`spec-update`](plugins/yarstack/skills/spec-update/SKILL.md) | Correct a product, architecture, API, or behavior contract. |

### Refine issues and deliver work

| Skill | Use it to |
| --- | --- |
| [`jira-issue-refine`](plugins/yarstack/skills/jira-issue-refine/SKILL.md) | Turn a vague issue into a bounded delivery task with acceptance criteria. |
| [`jira-issue-create`](plugins/yarstack/skills/jira-issue-create/SKILL.md) | Create Jira issues from approved drafts through `acli`. |
| [`jira-issue-deliver`](plugins/yarstack/skills/jira-issue-deliver/SKILL.md) | Deliver an authorized Jira issue to a ready PR with passing checks and requested automated review complete. |
| [`roadmap-task-deliver`](plugins/yarstack/skills/roadmap-task-deliver/SKILL.md) | Deliver the next eligible roadmap task through PR merge and task-status update. |
| [`phase-commit`](plugins/yarstack/skills/phase-commit/SKILL.md) | Create one explicitly requested local commit. |
| [`pr-draft`](plugins/yarstack/skills/pr-draft/SKILL.md) | Publish completed work as an authorized draft PR and follow CI. |

</details>

<details>
<summary><strong>Browse Yarbrain skills</strong></summary>

| Skill | Use it to |
| --- | --- |
| [`wiki-initialize`](plugins/yarbrain/skills/wiki-initialize/SKILL.md) | Create or adopt a vault at a path you choose. |
| [`wiki-capture`](plugins/yarbrain/skills/wiki-capture/SKILL.md) | Preserve a session as evidence and propose knowledge or procedures. |
| [`wiki-reconcile`](plugins/yarbrain/skills/wiki-reconcile/SKILL.md) | Compare proposed knowledge with existing notes before updating them. |
| [`wiki-recall`](plugins/yarbrain/skills/wiki-recall/SKILL.md) | Retrieve up to five relevant notes, episodes, or approved procedures. |
| [`procedure-promote`](plugins/yarbrain/skills/procedure-promote/SKILL.md) | Propose a skill from a repeated, verified procedure. |
| [`wiki-maintain`](plugins/yarbrain/skills/wiki-maintain/SKILL.md) | Check vault structure, stale knowledge, duplication, and skill drift. |

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
