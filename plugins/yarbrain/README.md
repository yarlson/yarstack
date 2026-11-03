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
    PreToolUse hook blocks every other agent from judging lessons.
  - In Codex, the gate is a subagent spawned without the parent's history; a
    verdict from a thread that proposed or doubted the lesson does not count.
- **Feedback.** A session that applies a lesson runs `yarbrain confirm`; one
  that finds it wrong runs `yarbrain doubt`, which takes the lesson out of
  service until a gate judges it again. Editing a trusted lesson by hand also
  takes it out of service until a gate trusts the new text.

You can read, add, and edit lessons yourself, but no step waits for you.

## Install

1. Install the CLI:

   ```sh
   curl -fsSL https://raw.githubusercontent.com/yarlson/yarstack/main/plugins/yarbrain/cli/install.sh | sh
   ```

   On Windows: `irm https://raw.githubusercontent.com/yarlson/yarstack/main/plugins/yarbrain/cli/install.ps1 | iex`.
   See the [CLI README](cli/README.md) for verification.

2. Install the plugin: `codex plugin add yarbrain@yarstack`, or
   `claude plugin install yarbrain@yarstack`.

3. Create or adopt a vault: `yarbrain init <vault-path>`. It prints the
   one-time host setup, which keeps vault writes free of prompts:

   - Claude Code: add `Bash(yarbrain *)` to `permissions.allow`.
   - Codex: add the vault to `writable_roots` under
     `[sandbox_workspace_write]`, and create `~/.codex/agents/yarbrain-gate.toml`
     to pin the gate's model and reasoning effort.

The CLI never edits host configuration itself.

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
- Lessons are served only to projects with a Git identity, or with scope
  `all`.
- A Claude Code session started from a Codex shell inherits the Codex thread
  id, so Yarbrain counts it as that thread. Its gate verdicts then do not
  count, and its lessons wait for a gate from a normal session.
