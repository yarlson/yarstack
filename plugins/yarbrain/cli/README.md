# yarbrain CLI

`yarbrain` reads and writes a Yarbrain lesson vault: a folder of Markdown
lessons that agents propose, an independent gate agent trusts or rejects, and
later sessions confirm or doubt. Agents call it through the Yarbrain plugin's
skills and hooks. People can call it from a terminal.

## Install

macOS and Linux:

```sh
curl -fsSL https://raw.githubusercontent.com/yarlson/yarstack/main/plugins/yarbrain/cli/install.sh | sh
```

Windows (PowerShell):

```powershell
irm https://raw.githubusercontent.com/yarlson/yarstack/main/plugins/yarbrain/cli/install.ps1 | iex
```

The scripts download the latest release for your OS and architecture, verify
its SHA-256 against the release's `checksums.txt`, and install it to
`~/.local/bin` or `%LOCALAPPDATA%\Programs\yarbrain`. Set `YARBRAIN_VERSION`
to install a specific version and `YARBRAIN_INSTALL_DIR` to change the
directory. On macOS and Linux, the script checks that your login shell finds
`yarbrain`, because Codex runs hooks in a login shell.

`checksums.txt` catches a corrupted or swapped download. To check who built
a binary, verify its build provenance:

```sh
gh attestation verify yarbrain_<version>_<os>_<arch> --repo yarlson/yarstack \
  --signer-workflow yarlson/yarstack/.github/workflows/release-yarbrain.yml \
  --source-ref refs/tags/plugins/yarbrain/cli/v<version>
```

## Commands

Every command except `hook` prints exactly one JSON object on stdout. Run any
command with `--help` for its input format and an example.

| Command | Caller | What it does |
|---------|--------|--------------|
| `init <vault>` | setup | Create or adopt a vault and point the config at it; print the one-time host setup |
| `status` | agent | Contract version, vault, project identity, and this project's counts |
| `find [term...]` | agent | A ranked menu of served lessons: id, title, and cue |
| `show <id\|v0.2-path>...` | agent | Full text of served lessons, or of a v0.2 note for conversion |
| `propose` | author | Store a lesson as proposed; input is the lesson file format on stdin |
| `brief <id>` | gate | Lease a lesson and get it as data, with committed-tree facts |
| `verdict <id>` | gate | Trust, reject, retire, or defer, with four checks |
| `confirm <id>` | agent | Record that a served lesson worked |
| `doubt <id>` | agent | Take a served lesson out of service because it was wrong |
| `legacy` | agent | List v0.2 notes that no lesson converts yet |
| `legacy dismiss <v0.2-path>` | agent | Take a v0.2 note with nothing worth a lesson out of the backlog |
| `check` | anyone | Validate every lesson file; exit 1 on findings |
| `hook session-start\|pre-tool-use` | host | The plugin's hook handlers |

Arguments carry only ids, paths, flags, and search terms. Free text arrives on
stdin or with `--input FILE`, at most 32 KiB.

## Trust

No person gates a lesson. `find` and `show` serve a lesson only when its
latest verdict is a credible trust verdict for its current content:

- it comes from a gate agent, recorded with `brief` then `verdict`, in a host
  session (Claude Code or Codex);
- in Codex, the gate's thread did not propose the lesson and did not doubt it
  since its last trust verdict, and `brief` refuses such a thread; in Claude
  Code, whose subagents share the parent's session id, the plugin's
  PreToolUse hook lets only the `yarbrain:lesson-gate` agent run `brief` and
  `verdict`;
- the lesson text passes the same rules as `propose`;
- no edit changed the lesson's title, cue, scope, body, Observed section,
  retires, or converts since, and no doubt or replacement followed the
  verdict.

The ledger in `<vault>/.yarbrain/ledger/` records every proposal, verdict,
confirmation, and doubt. Its notes and reasons never become lesson text.

## Output contract

```json
{"v":1,"ok":true}
{"v":1,"ok":false,"error":{"code":"bad_input","message":"unknown command \"x\"","fix":"Run `yarbrain --help` to list the commands.","field":"command"}}
```

- `v` is the contract version. Changes within one version only add fields. A
  breaking change ships as a new command name.
- `error.fix` names the next call that resolves the failure. Other `error`
  fields depend on the code.
- `--help` and `--version` print plain text.

| Exit | Meaning | Error codes |
|------|---------|-------------|
| 0 | Success | |
| 1 | `check` found problems | |
| 2 | Crash: a defect in yarbrain | `internal` |
| 3 | Fix the call | `bad_input`, `bad_term`, `not_found`, `title_taken`, `bad_state`, `changed_underneath`, `unknown_project`, `project_mismatch`, `secret_suspected`, `risky_command`, `already_recorded` |
| 4 | Retry the same call once | `locked`, `busy` |
| 5 | Setup problem; tell the person | `no_vault`, `write_denied`, `bad_file`, `config_conflict`, `contract_unsupported`, `io` |

## Development

The module pins its tools in `mise.toml`. Run `mise install` once, then:

| Command | What it runs |
|---------|--------------|
| `make check` | formatting, golangci-lint, `go vet`, tests, build, `go mod tidy -diff`, `go mod verify`, `goreleaser check`, actionlint |
| `make ci` | `check`, then race tests, CRAP (every score below 15), CGO-free tests and builds, cross-builds, govulncheck, a release snapshot, and release staging |
| `go test ./cmd/yarbrain -run TestName` | one test |

The repository's `make validate` runs `make ci` for this module.

## Release

Push a tag `plugins/yarbrain/cli/vX.Y.Z` on a commit that is on `main`. Go
requires this prefix for a module in a subdirectory. The release workflow
runs `make ci`, builds the binaries, publishes them with `checksums.txt` and
the install scripts, downloads and verifies every asset, and attests their
build provenance.
