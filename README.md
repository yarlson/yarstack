# Yarstack

![Yarstack pirate coins](assets/yarstack.png)

Yarstack is Yar Kravtsov's portable plugin marketplace for Codex and Claude
Code. Both platforms install the same self-contained engineering workflows and
standards while retaining their native manifests and marketplace metadata.

The plugin packages reusable engineering skills, a shared engineering-quality
policy, and a safe installer for global Codex and Claude Code guidance.

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

## Engineering standards

The canonical policy fragments live under
`plugins/yarstack/agent-guidance/engineering-standards/`. Install their combined
form into the global guidance files for Codex and Claude Code:

```sh
make install-system-prompt
```

The installer is idempotent. It preserves symlinks, refuses dangling symlinks,
and creates timestamped backups before replacing differing existing guidance.
Use `--print` to inspect the exact combined document without writing anything.

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

Release history is in
[plugins/yarstack/CHANGELOG.md](plugins/yarstack/CHANGELOG.md).

## Security

Report vulnerabilities privately as described in [SECURITY.md](SECURITY.md).

## License

[MIT](LICENSE) © 2026 Yar Kravtsov.
