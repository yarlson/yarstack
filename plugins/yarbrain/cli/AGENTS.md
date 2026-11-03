# Yarbrain CLI instructions

`yarbrain` is the Go CLI that owns every read and write of a user's Yarbrain
vault: a folder of Markdown lessons that agents propose, an independent gate
agent judges, and later sessions confirm or doubt. The repository root
`AGENTS.md` also applies. This file adds the rules for this module.

## Boundaries

- `cmd/yarbrain` owns Cobra commands and flags, stdin and `--input` decoding,
  host session ids from the environment, the JSON envelope, the mapping from
  error codes to exit codes, dependency assembly, and the process exit. Keep
  lesson, trust, and storage policy out of command files.
- `internal/lesson` owns lesson fields, limits, normalization, the content
  digest, and lifecycle transitions. It does no I/O.
- `internal/trust` owns ledger events and the serving rules: credible verdicts,
  contributor sessions, leases, and the limits on doubt, confirm, and defer. It
  does no I/O.
- `internal/vault` owns config and vault resolution, the layout, the format
  version, frontmatter reads and owned-line edits, exclusive creates, the vault
  lock, ledger files, and read-only access to v0.2 content.
- `internal/project` owns Git identity and the committed-tree facts that the
  gate brief reports.
- `internal/textrules` owns the deterministic checks on agent-written text:
  secret patterns, text that addresses agents, risky commands, and URLs.
- `internal/search` owns term validation, ranking, and menu budgets.
- `internal/hook` owns host hook payloads and the session-start line.
- `.goreleaser.yml`, the install scripts, and the release workflow own
  distribution. Release tags are `plugins/yarbrain/cli/vX.Y.Z`, the prefix Go
  requires for a module in a subdirectory. Keep release targets aligned with
  `make cross-build`.
- Add an `internal` package only when it hides a real decision. Do not create
  generic `common`, `util`, `types`, or `models` packages.

## Product contracts

- Every command except `hook` prints exactly one JSON object on stdout:
  `{"v":1,"ok":true,...}`, or `{"v":1,"ok":false,"error":{"code","message","fix"}}`.
  Changes within `v:1` are additive. A breaking change ships as a new command
  name; an old name never gets a new meaning.
- Exit codes tell the caller what to do: 0 success, 1 `check` findings,
  2 crash, 3 fix the call, 4 retry once, 5 setup problem.
- Arguments carry only ids, paths, flags, and search terms. Free text arrives
  as UTF-8 on stdin or through `--input FILE`, capped at 32 KiB, with a BOM
  stripped and unknown keys rejected.
- Config and vault resolution keep the v0.2 order: `--vault`, `YARBRAIN_VAULT`,
  then the config from `--config`, `YARBRAIN_CONFIG`,
  `$XDG_CONFIG_HOME/yarbrain/config.json`, or `~/.config/yarbrain/config.json`
  on every OS. Never write v0.2 directories (`episodes/`, `notes/`, `inbox/`,
  `skills/`, `archive/`, `reports/`, `INDEX.md`).
- No command makes a person a gate or waits for a confirmation. Only `verdict`
  from a gate agent sets trust. `find` and `show` serve only lessons with a
  credible trust verdict; ungated lessons and legacy notes appear as counts.
- Ledger data never enters the digest or served lesson text. Every
  agent-written text passes `internal/textrules`; a finding names the field and
  the pattern, never the value, and no check ever reports "clean".
- Create lesson files with `O_EXCL`. Change them only under the vault lock,
  after rechecking `rev` and the file hash. Rewrite only single-line keys the
  CLI owns, keep every other byte, and refuse a file whose owned key is
  duplicated or spans several lines.
- `hook session-start` is read-only, prints counts only, never starts with `{`,
  and finishes within 2 seconds. It exits 0 unless a configured vault is
  missing.
- The CLI never calls the network or a model.

## Go and developer-experience rules

- Write explicit, idiomatic Go with small functions, `%w` error wrapping, and
  no panics for expected failures. Make every error name what failed, for which
  id or path, and the next call that fixes it.
- Use Cobra for commands, `go.yaml.in/yaml/v3` for frontmatter, and Testify in
  tests. Add any other dependency only when it removes more complexity or risk
  than it adds.
- Pin the Go toolchain, golangci-lint, goreleaser, and other development tools
  in `mise.toml`, not in the module graph.
- Open vault files through `os.Root` so no path can leave the vault. Run every
  `git` call through `exec.CommandContext` with a timeout, and pass contexts to
  everything that waits.
- Keep every build `CGO_ENABLED=0`. Test Windows-specific behavior, such as
  non-atomic rename, wherever it changes a contract.
- Use Testify: `require` for prerequisites and `assert` for independent
  outcomes. Use real temporary Git repositories at the Git boundary. Decode
  JSON output in tests and assert the contract, not string snippets.
- Point tests at `t.TempDir()` vaults with `YARBRAIN_CONFIG`, `HOME`, and
  `XDG_CONFIG_HOME` set through `t.Setenv`. Never read or write the real user
  configuration or vault.
- Build token-shaped test fixtures at runtime, for example `"ghp_" + body`.
  The repository's plugin scanner reads every file in this directory and
  reports literal tokens.
- Run `make check` for formatting, linting, vet, tests, builds, `go mod tidy
  -diff`, `go mod verify`, and `goreleaser check`. Run one test with
  `go test ./internal/vault -run TestName`.
- Run `make ci` before delivery. It adds race tests, CRAP, cross-builds,
  vulnerability scanning, and a release snapshot.
- New or heavily changed production functions must keep every CRAP score below
  15. Do not worsen an existing score that is above the threshold.
- Do not apply repository-wide automatic fixes except the configured formatter.
  Review every generated edit.

## Change rules

- Keep each change to one cohesive behavior. Reject speculative commands,
  storage, configuration, abstractions, and framework code.
- Preserve commands, flags, JSON keys, error codes, exit codes, the vault
  format, config resolution, and the hook output as public contracts.
- The plugin's skills in `../skills/` and its `../hooks/hooks.json` call this
  CLI. When a contract changes, update them in the same change.
- Test normal, boundary, failure, cancellation, cleanup, and concurrency
  behavior wherever the contract allows deterministic evidence.
- Update user documentation when user-visible behavior or development commands
  change.
- Before finalizing, review the full diff, remove unrelated work, run the
  applicable gates, and report any inconclusive evidence.
