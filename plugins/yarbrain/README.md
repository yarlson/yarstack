# Yarbrain

Yarbrain lets coding agents keep what a session taught them. An agent records
a small lesson; an independent gate agent trusts or rejects it; later
sessions in Claude Code or Codex find it, apply it, and report whether it
worked. Lessons are Markdown files in a vault you own.

## How it works

- **Lesson.** A title, a cue (the situation that makes it relevant, such as
  an error message), what to do, and what was observed. One file per lesson
  in `<vault>/lessons/`.
- **Lookup.** At session start, a hook prints one line with counts for the
  current project. Agents run `yarbrain find <terms>` for a short menu and
  `yarbrain show <id>` for the full text.
- **Trust without a person.** A proposed lesson is not served. The
  `lesson-record` skill starts a gate agent with fresh context and only the
  lesson id. The gate checks the lesson against the committed code, its docs,
  and similar lessons, and records trust, reject, retire, or defer. Agents
  start the gate without asking you, in the session that proposed or doubted
  the lesson or in a later one. Each gate is one subagent run: Sonnet at high
  effort in Claude Code, and the pinned role or the session's model in Codex.
  - In Claude Code, the gate is the plugin's `yarbrain:lesson-gate` agent, and a
    PreToolUse hook denies `yarbrain brief` and `yarbrain verdict` in every
    other agent's Bash and PowerShell calls.
  - In Codex, the gate is a subagent spawned without the parent's history;
    `brief` refuses a thread that proposed the lesson or doubted it since its
    last trust.
- **Feedback.** A session that applies a lesson runs `yarbrain confirm`; one
  that finds it wrong runs `yarbrain doubt`, which takes the lesson out of
  service until a gate judges it again. Editing a trusted lesson by hand also
  takes it out of service until a gate trusts the new text.

You can read, add, and edit lessons yourself, but no step waits for you.

## Install

1. Install the plugin: `codex plugin add yarbrain@yarstack`, or
   `claude plugin install yarbrain@yarstack`.
2. Ask your agent to set up Yarbrain, for example: "Set up Yarbrain with its
   vault at ~/yarbrain."

The agent does the rest:

- It runs the plugin's launcher, `bin/yarbrain`. The first call downloads the
  `yarbrain` CLI build that the plugin pins, checks its SHA-256 against
  [`cli-checksums.txt`](cli-checksums.txt), and caches it in
  `~/.cache/yarbrain` (or `$XDG_CACHE_HOME/yarbrain`).
- It creates or adopts the vault with `yarbrain init <vault-path>`.
- It applies the one-time host setup that `init` prints, which keeps vault
  writes free of prompts:
  - Claude Code: add `Bash(yarbrain *)` to `permissions.allow`.
  - Codex: add the vault to `writable_roots` under
    `[sandbox_workspace_write]`, and create `~/.codex/agents/yarbrain-gate.toml`
    to pin the gate's model and reasoning effort.

Your host asks you to approve each settings edit. In Codex, the sandbox also
asks before the first download. The CLI never edits host configuration itself.

A plugin update also updates the CLI: the new version pins a new build, and
the launcher downloads it on first use. Hooks never download. Until the first
download, the PreToolUse hook refuses every `yarbrain brief` and
`yarbrain verdict` call, and the agent runs `yarbrain status` alone first.

The launcher is a POSIX shell script; on Windows it runs in Git Bash. Codex
on Windows, and Claude Code on Windows without Git Bash, need the CLI on
`PATH` instead:
`irm https://raw.githubusercontent.com/yarlson/yarstack/main/plugins/yarbrain/cli/install.ps1 | iex`.
The [CLI README](cli/README.md) covers manual installs and verification.

## Upgrading from v0.2

Run `yarbrain init` on your existing vault path; the config file does not
change. Yarbrain never writes the v0.2 content (`episodes/`, `notes/`,
`inbox/`, `skills/`, `archive/`, `reports/`, `INDEX.md`). Unconverted notes
appear as a count, and `lesson-record` converts or dismisses up to three per
session, each through the gate.

## Limits

- In Codex, the gate's model pin needs the role file from setup, and the
  gate's read-only rule is an instruction, not a sandbox boundary.
- A verdict is a record in the vault's ledger. The CLI refuses verdicts that
  are not credible, but any process that can write the vault could forge one.
- The Claude Code hook reads the command line. A gate call hidden in a
  variable or in a nested shell string, such as `sh -c '...'`, gets past it,
  and so does the launcher run through a shell, as in
  `sh <plugin>/bin/yarbrain brief <id>`.
- Claude Code on Windows without Git Bash runs hooks in PowerShell, where the
  plugin's hooks fail: no counts line and no gate check.
- claude.ai and Cowork do not install a plugin with a top-level `bin/`
  directory, so Yarbrain works only in Claude Code and Codex.
- A `yarbrain` that you installed on `PATH` yourself can run instead of the
  pinned build. Codex agents run it whenever they call `yarbrain`, and the
  hooks run it until the pinned build is downloaded. Remove it to always run
  the pinned build.
- Lessons are served only to projects with a Git identity, or with scope
  `all`.
- A Claude Code session started from a Codex shell inherits the Codex thread
  id, so Yarbrain counts it as that thread. `brief` refuses its gate agent,
  and its lessons wait for a gate from a normal session.
