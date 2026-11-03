# yarbrain CLI

`yarbrain` reads and writes a Yarbrain lesson vault. Agents call it through the
Yarbrain plugin's skills and hook. People can call it from a terminal.

## Output contract

Every command except `hook` prints exactly one JSON object on stdout.

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
| 2 | Crash: a defect in yarbrain | `internal` |
| 3 | Fix the call | `bad_input`, `bad_term`, `not_found`, `title_taken`, `bad_state`, `changed_underneath`, `unknown_project`, `project_mismatch`, `secret_suspected`, `risky_command`, `already_recorded` |
| 4 | Retry the same call once | `locked`, `busy` |
| 5 | Setup problem; tell the person | `no_vault`, `write_denied`, `bad_file`, `config_conflict`, `contract_unsupported`, `io` |

## Development

The module pins its toolchain in `mise.toml`. Run `mise install` once, then:

| Command | What it runs |
|---------|--------------|
| `make check` | formatting, golangci-lint, `go vet`, tests, build, `go mod tidy -diff`, `go mod verify` |
| `make ci` | `check`, then race tests, CRAP (every score below 15), CGO-free tests and builds, cross-builds, and govulncheck |
| `go test ./cmd/yarbrain -run TestName` | one test |

The repository's `make validate` runs `make ci` for this module.
