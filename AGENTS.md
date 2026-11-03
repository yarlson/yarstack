# Repository Rules

## Commands

- Run `make validate` before finalizing changes. It also runs `make ci` for the Go CLI in `plugins/yarbrain/cli`, which needs the tools pinned in its `mise.toml`.
- Run `make install-system-prompt` only when the user explicitly requests changes to global agent guidance; it writes outside the repository.
- Test installer changes with `YARSTACK_AGENT_HOME` pointed at a temporary directory. Do not test destructive paths against the real home directory.

## Package invariants

- Keep each installable package self-contained under `plugins/<plugin-name>/`. Marketplace installs are cached copies; runtime paths must not traverse outside the owning plugin root.
- Follow [`SKILL_FORMAT.md`](SKILL_FORMAT.md) when adding or changing skills.
- Keep Codex metadata in `.codex-plugin/plugin.json`, Claude metadata in `.claude-plugin/plugin.json`, and their marketplace schemas separate.
- Keep each plugin's shared identity, author, repository, and version values aligned between its manifests.
- Do not change `0.1.0` or any later explicit version unless the user requests a release/version change. When changing it, update both manifests for that plugin together.
- Do not add a plugin-root `CLAUDE.md`; Claude Code does not load it as plugin context.
- Do not declare skills, hooks, agents, MCP servers, apps, assets, or capabilities before the corresponding package content and behavior exist.
- Keep every procedure in a skill. Agents under `plugins/<plugin-name>/agents/` load only in Claude Code, because Codex plugins cannot ship agents. Use an agent as a worker that applies skills, and keep each skill usable without it.

## Engineering guidance

- Treat `plugins/yarstack/agent-guidance/engineering-standards/standards.md` as the canonical source for the installed global guidance.
- Preserve the guidance text unless the user explicitly requests prompt changes. When a rule changes, update its evidence entry in `plugins/yarstack/agent-guidance/engineering-standards/README.md`.
- Preserve installer idempotency, existing symlinks, preflight rejection of dangling symlinks, and backups before replacing differing files.
- Keep guidance installation explicit. Do not convert it into an automatic hook, startup action, or implicit plugin side effect.

## Yarbrain invariants

- Keep the vault and all mutable knowledge outside the installed or cached plugin root.
- Do not choose, create, or configure a vault during plugin installation. Create or adopt a vault only with `yarbrain init` and a path the user gave, and preserve existing vault contents.
- The `yarbrain` CLI in `plugins/yarbrain/cli` owns every vault read and write. Skills and hooks call it; they never edit vault files directly.
- Keep lifecycle hooks read-only and bounded. They may print counts and decide on gate commands, but must not write, print lesson text, copy transcripts, or call a model or network service.
- No person gates a lesson. Trust comes only from a credible verdict that a fresh-context gate agent records through the CLI. Never weaken the credibility rules to make a lesson served.
- Never modify or delete v0.2 vault content (`episodes/`, `notes/`, `inbox/`, `skills/`, `archive/`, `reports/`, `INDEX.md`).
- Test Yarbrain with a temporary `YARBRAIN_CONFIG` and vault. Never use the real user configuration or vault.

## Change discipline

- Keep `Makefile` as direct command entrypoints. Do not put custom validators, test harnesses, dependency bootstrap frameworks, metadata reimplementations, or large shell programs in it.
- Use `plugin-scanner lint` and `plugin-scanner verify` for Codex validation and `claude plugin validate --strict` for Claude validation. Do not invent parallel validation logic.
- Keep `.github/workflows/validate.yml` calling `make validate`, plus `make -C plugins/yarbrain/cli test` on macOS and Windows; do not duplicate Make targets in workflow YAML.
- Add dependencies, generated files, wrappers, configuration layers, or host-specific adapters only when a current capability requires them.
- Prefer one shared implementation. Split by host only where the platform contracts actually differ.
- Do not add placeholders, speculative abstractions, future-proofing, or unrelated repository polish.
