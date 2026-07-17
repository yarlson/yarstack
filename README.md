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

## Upgrade

Marketplace installs are cached by version. After Yarstack publishes a new
version, refresh the marketplace and reinstall or update the plugin.

### Codex

```sh
codex plugin marketplace upgrade yarstack
codex plugin add yarstack@yarstack
```

Start a new Codex thread to load the updated plugin.

### Claude Code

```sh
claude plugin marketplace update yarstack
claude plugin update yarstack@yarstack
```

Run `/reload-plugins` or restart Claude Code. Marketplace auto-update can also
be enabled from the `/plugin` interface.

## Skill catalogue

| Skill                                                                                 | Use it for                                                                                          |
| ------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| [`alternative-ideation`](plugins/yarstack/skills/alternative-ideation/SKILL.md)       | Generating and testing a genuinely non-incremental product, workflow, or system alternative.        |
| [`architecture-sparring`](plugins/yarstack/skills/architecture-sparring/SKILL.md)     | Resolving underdefined architecture through one focused decision at a time.                         |
| [`ci-review`](plugins/yarstack/skills/ci-review/SKILL.md)                             | Finding stale, redundant, unsafe, or unjustified CI and delivery machinery.                         |
| [`cli-control`](plugins/yarstack/skills/cli-control/SKILL.md)                         | Verifying interactive CLI or TUI behavior with repeatable local evidence.                           |
| [`code-review`](plugins/yarstack/skills/code-review/SKILL.md)                         | Reviewing local changes for concrete correctness, security, performance, and lifecycle defects.     |
| [`coderabbit-triage`](plugins/yarstack/skills/coderabbit-triage/SKILL.md)             | Judging and resolving CodeRabbit feedback on the current pull request.                              |
| [`context-map`](plugins/yarstack/skills/context-map/SKILL.md)                         | Mapping the files, tests, conventions, commands, and risks relevant to a change.                    |
| [`critical-journey-docs`](plugins/yarstack/skills/critical-journey-docs/SKILL.md)     | Documenting traceable user journeys, acceptance scenarios, telemetry, and service flows.            |
| [`dependency-review`](plugins/yarstack/skills/dependency-review/SKILL.md)             | Reviewing dependency, lockfile, toolchain, generated-code, and supply-chain changes.                |
| [`docs-drift-review`](plugins/yarstack/skills/docs-drift-review/SKILL.md)             | Updating documentation made inaccurate by a code or configuration change.                           |
| [`docs-review`](plugins/yarstack/skills/docs-review/SKILL.md)                         | Auditing or restructuring README and repository documentation against implemented behavior.         |
| [`draft-pr`](plugins/yarstack/skills/draft-pr/SKILL.md)                               | Branching, committing, pushing, and shipping confirmed changes as a monitored draft PR.             |
| [`explain-changes`](plugins/yarstack/skills/explain-changes/SKILL.md)                 | Explaining a diff or PR through behavior, boundaries, data flow, operations, and risk.              |
| [`final-review`](plugins/yarstack/skills/final-review/SKILL.md)                       | Closing completed work against its contract, final diff, and verification evidence.                 |
| [`go-review`](plugins/yarstack/skills/go-review/SKILL.md)                             | Finding Go-specific concurrency, cancellation, resource, retry, shutdown, and race defects.         |
| [`infra-review`](plugins/yarstack/skills/infra-review/SKILL.md)                       | Reviewing infrastructure changes for blast radius, exposure, destruction, cost, and recovery risks. |
| [`macos-swift-review`](plugins/yarstack/skills/macos-swift-review/SKILL.md)           | Auditing macOS Swift and SwiftUI lifecycle, state, concurrency, sandbox, and distribution risks.    |
| [`marketing-claims-review`](plugins/yarstack/skills/marketing-claims-review/SKILL.md) | Checking product marketing claims and rewrites against shipped-product evidence.                    |
| [`phase-commit`](plugins/yarstack/skills/phase-commit/SKILL.md)                       | Creating one explicitly requested local commit for a completed, verified phase.                     |
| [`phase-implement`](plugins/yarstack/skills/phase-implement/SKILL.md)                 | Implementing exactly one selected plan phase without later-phase work.                              |
| [`phase-review`](plugins/yarstack/skills/phase-review/SKILL.md)                       | Fixing maintainability problems in a completed phase without changing behavior or scope.            |
| [`phase-validate`](plugins/yarstack/skills/phase-validate/SKILL.md)                   | Validating one implemented phase against its stated behavioral contract.                            |
| [`plan-create`](plugins/yarstack/skills/plan-create/SKILL.md)                         | Producing an implementation-ready repository plan with decisions, phases, risks, and validation.    |
| [`plan-update`](plugins/yarstack/skills/plan-update/SKILL.md)                         | Making the smallest evidence-backed correction to an implementation plan.                           |
| [`refactor-plan`](plugins/yarstack/skills/refactor-plan/SKILL.md)                     | Designing a minimal, reversible sequence for behavior-preserving structural change.                 |
| [`repo-context-docs`](plugins/yarstack/skills/repo-context-docs/SKILL.md)             | Maintaining verified current-state architecture and ownership documentation under `docs/context/`.  |
| [`rollout-review`](plugins/yarstack/skills/rollout-review/SKILL.md)                   | Reviewing deployment readiness, runtime failure modes, rollback, and recovery.                      |
| [`security-review`](plugins/yarstack/skills/security-review/SKILL.md)                 | Reviewing changed trust boundaries, dangerous inputs, permissions, secrets, and agent surfaces.     |
| [`spec-update`](plugins/yarstack/skills/spec-update/SKILL.md)                         | Correcting an authoritative product or architecture contract when implementation exposes a gap.     |
| [`technical-spike`](plugins/yarstack/skills/technical-spike/SKILL.md)                 | Resolving one bounded, implementation-blocking technical uncertainty.                               |
| [`test-gap-review`](plugins/yarstack/skills/test-gap-review/SKILL.md)                 | Finding missing behavior coverage and tests that provide false confidence.                          |
| [`ui-control`](plugins/yarstack/skills/ui-control/SKILL.md)                           | Verifying browser, desktop, Electron, or other UI behavior through the actual interface.            |

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

## Security

Report vulnerabilities privately as described in [SECURITY.md](SECURITY.md).

## License

[MIT](LICENSE) © 2026 Yar Kravtsov.
