# Yarstack

[![Validate marketplace](https://github.com/yarlson/yarstack/actions/workflows/validate.yml/badge.svg)](https://github.com/yarlson/yarstack/actions/workflows/validate.yml)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Yarstack is Yar Kravtsov's portable plugin marketplace for Codex and Claude
Code. Both platforms install the same self-contained package while retaining
their native manifests and marketplace metadata.

The current `0.1.0` release establishes packaging, discovery, and installation.
It intentionally provides no skills, tools, hooks, agents, or other runtime
capabilities yet.

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

## Catalog

| Plugin   | Version | Codex     | Claude Code | Status               |
| -------- | ------- | --------- | ----------- | -------------------- |
| Yarstack | 0.1.0   | Supported | Supported   | Packaging foundation |

## Repository layout

```text
.agents/plugins/marketplace.json   Codex marketplace
.claude-plugin/marketplace.json    Claude Code marketplace
plugins/yarstack/                  Shared installable package
Makefile                           Local and CI entrypoints
```

The package contains separate `.codex-plugin/plugin.json` and
`.claude-plugin/plugin.json` manifests. Future portable workflows belong in
`plugins/yarstack/skills/`; host-specific adapters should be added only when a
shared implementation cannot express the required behavior.

## Validation

Run the complete read-only validation suite locally:

```sh
make validate
```

Every push and pull request runs:

- OpenAI's pinned plugin-creator validator.
- `claude plugin validate --strict` for the plugin and marketplace.
- Metadata and package-integrity checks with `jq` and `cmp`.
- A clean Codex marketplace discovery and installation smoke test.

`make smoke-codex-install` runs that stateful smoke test locally. It registers
this checkout as a Codex marketplace and installs the plugin, so it is not part
of the default validation target.

See [CONTRIBUTING.md](CONTRIBUTING.md) for the local review checklist and
[plugins/yarstack/CHANGELOG.md](plugins/yarstack/CHANGELOG.md) for releases.

## Security

Review [SECURITY.md](SECURITY.md) before reporting a vulnerability. Plugins are
a trust boundary: future scripts, hooks, MCP servers, network access, and write
actions must be explicit and narrowly scoped.

## License

[MIT](LICENSE) © 2026 Yar Kravtsov.
