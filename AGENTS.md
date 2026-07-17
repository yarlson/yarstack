# Repository Rules

## Commands

- Run `make validate` before finalizing changes.
- Run `make install-system-prompt` only when the user explicitly requests changes to global agent guidance; it writes outside the repository.
- Test installer changes with `YARSTACK_AGENT_HOME` pointed at a temporary directory. Do not test destructive paths against the real home directory.

## Package invariants

- Keep the installable package self-contained under `plugins/yarstack/`. Marketplace installs are cached copies; runtime paths must not traverse outside the plugin root.
- Keep Codex metadata in `.codex-plugin/plugin.json`, Claude metadata in `.claude-plugin/plugin.json`, and their marketplace schemas separate.
- Keep shared identity, author, repository, and version values aligned between manifests.
- Do not change `0.1.0` or any later explicit version unless the user requests a release/version change. When changing it, update both manifests and `plugins/yarstack/CHANGELOG.md` together.
- Do not add a plugin-root `CLAUDE.md`; Claude Code does not load it as plugin context.
- Do not declare skills, hooks, agents, MCP servers, apps, assets, or capabilities before the corresponding package content and behavior exist.

## Engineering guidance

- Treat `plugins/yarstack/agent-guidance/engineering-standards/` as the canonical source for the installed system prompt.
- Preserve policy fragment contents and assembly order unless the user explicitly requests prompt changes.
- Preserve installer idempotency, existing symlinks, preflight rejection of dangling symlinks, and backups before replacing differing files.
- Keep guidance installation explicit. Do not convert it into an automatic hook, startup action, or implicit plugin side effect.

## Change discipline

- Keep `Makefile` as direct command entrypoints. Do not put custom validators, test harnesses, dependency bootstrap frameworks, metadata reimplementations, or large shell programs in it.
- Use `plugin-scanner lint` and `plugin-scanner verify` for Codex validation and `claude plugin validate --strict` for Claude validation. Do not invent parallel validation logic.
- Keep `.github/workflows/validate.yml` calling `make validate`; do not duplicate Make targets in workflow YAML.
- Add dependencies, generated files, wrappers, configuration layers, or host-specific adapters only when a current capability requires them.
- Prefer one shared implementation. Split by host only where the platform contracts actually differ.
- Do not add placeholders, speculative abstractions, future-proofing, or unrelated repository polish.
