# Contributing

Yarstack accepts focused improvements to its marketplace packaging and plugin
capabilities.

## Development principles

- Keep `plugins/yarstack/` self-contained; installed plugins cannot rely on
  files outside their package.
- Prefer shared Agent Skills under `plugins/yarstack/skills/` over duplicated
  host-specific instructions.
- Add hooks, MCP servers, agents, scripts, or dependencies only for a concrete
  capability and document their permissions and failure behavior.
- Keep the Codex and Claude Code plugin names, versions, authors, and repository
  URLs aligned.
- Never commit credentials, tokens, private endpoints, or generated caches.

## Pull requests

1. Explain the user-visible behavior and why the change is necessary.
2. Update both manifest versions when installable package content changes.
3. Add contract-focused tests or evaluations for new behavior.
4. Update `plugins/yarstack/CHANGELOG.md`.
5. Run the complete validation suite locally:

   ```sh
   make validate
   ```

6. Confirm the repository validation workflow is green.

Marketplace metadata, documentation, and release notes should remain accurate
for both supported hosts. Platform-specific behavior must be labeled clearly.
