package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/vault"
)

// app holds what every command needs from the process: streams, the
// environment, and the global flags.
type app struct {
	stdin      io.Reader
	stdout     io.Writer
	stderr     io.Writer
	getenv     func(string) string
	home       string
	flagVault  string
	flagConfig string
}

func newApp(stdin io.Reader, stdout, stderr io.Writer) *app {
	home, _ := os.UserHomeDir()

	return &app{stdin: stdin, stdout: stdout, stderr: stderr, getenv: os.Getenv, home: home}
}

func (a *app) locator() vault.Locator {
	return vault.Locator{FlagVault: a.flagVault, FlagConfig: a.flagConfig, Getenv: a.getenv, Home: a.home}
}

// openVault resolves and opens the configured vault. The caller closes it.
func (a *app) openVault() (*vault.Vault, vault.Location, error) {
	location, err := a.locator().Resolve()
	if err != nil {
		return nil, vault.Location{}, err
	}

	v, err := vault.Open(location.Vault)
	if err != nil {
		return nil, location, err
	}

	return v, location, nil
}

// classify turns a domain error into a contract error with a fix. Errors it
// does not recognize stay unclassified and become internal errors.
func (a *app) classify(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := errors.AsType[*cliError](err); ok {
		return err
	}

	if classified := classifyVault(err); classified != nil {
		return classified
	}

	if errors.Is(err, fs.ErrPermission) {
		return newError("write_denied", err.Error(), a.writeDeniedFix())
	}

	if _, ok := errors.AsType[*fs.PathError](err); ok {
		return newError("io", err.Error(), "Check that the vault and config paths are readable, then retry.")
	}

	return err
}

func classifyVault(err error) *cliError {
	if errors.Is(err, vault.ErrNotConfigured) {
		return newError("no_vault", err.Error(), "Run `yarbrain init <vault-path>` to create or adopt a vault.")
	}

	if errors.Is(err, vault.ErrLocked) {
		return newError("locked", err.Error(), "Retry the same command once.")
	}

	if missing, ok := errors.AsType[*vault.MissingError](err); ok {
		return newError("no_vault", err.Error(), fmt.Sprintf("Run `yarbrain init %s`, or point the config at an existing vault.", missing.Path)).
			with("path", missing.Path)
	}

	if config, ok := errors.AsType[*vault.ConfigError](err); ok {
		return newError("bad_file", err.Error(), fmt.Sprintf("Fix %s so it holds {\"vault\": \"<path>\"}.", config.Path)).
			with("path", config.Path)
	}

	if conflict, ok := errors.AsType[*vault.ConflictError](err); ok {
		return newError("config_conflict", err.Error(), fmt.Sprintf("Pass --config <other-file>, or edit %s yourself.", conflict.Path)).
			with("path", conflict.Path).with("existing", conflict.Existing)
	}

	if _, ok := errors.AsType[*vault.FormatError](err); ok {
		return newError("contract_unsupported", err.Error(), "Install a newer yarbrain CLI.")
	}

	if bad, ok := errors.AsType[*vault.BadFileError](err); ok {
		return newError("bad_file", err.Error(), fmt.Sprintf("Fix %s by hand; the reason says what is wrong.", bad.Path)).
			with("path", bad.Path)
	}

	if changed, ok := errors.AsType[*vault.ChangedError](err); ok {
		return newError("changed_underneath", err.Error(), "Read the lesson again with `yarbrain show <id>`, then retry.").
			with("path", changed.Path)
	}

	return nil
}

// writeDeniedFix names the one-time host setup that a denied write points to.
func (a *app) writeDeniedFix() string {
	switch {
	case a.getenv("CODEX_THREAD_ID") != "":
		return "Add the vault to [sandbox_workspace_write] writable_roots in ~/.codex/config.toml, then retry."
	case a.getenv("CLAUDE_CODE_SESSION_ID") != "":
		return `Add "Bash(yarbrain *)" to permissions.allow in ~/.claude/settings.json, then retry.`
	default:
		return "Make the vault directory writable for this user, then retry."
	}
}
