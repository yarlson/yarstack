package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/hook"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/project"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/vault"
)

const (
	// hookDeadline bounds the session-start hook, which runs before every
	// session.
	hookDeadline = 2 * time.Second
	// hookNotice is the hook's exit code for a configured vault that cannot
	// be read. Hosts show it to the person as a notice and never block on it.
	hookNotice = 1
)

func newHookCommand(a *app) *cobra.Command {
	var contract int

	command := &cobra.Command{
		Use:   "hook",
		Short: "Handle a host hook event (for the plugin's hooks.json)",
		Long: "hook is called by the hosts, not by agents. --contract names the CLI contract\n" +
			"the plugin expects, so an older plugin keeps working with a newer CLI.",
		Args: noArgs,
		RunE: func(*cobra.Command, []string) error {
			return newError("bad_input", "hook needs an event: session-start or pre-tool-use", "Run `yarbrain hook --help`.").
				with("field", "command")
		},
	}
	command.PersistentFlags().IntVar(&contract, "contract", contractVersion, "the CLI contract the plugin expects")

	command.AddCommand(
		&cobra.Command{
			Use:   "session-start",
			Short: "Print the counts-only session-start line",
			Args:  noArgs,
			RunE: func(command *cobra.Command, _ []string) error {
				return a.sessionStart(command.Context(), contract)
			},
		},
		&cobra.Command{
			Use:   "pre-tool-use",
			Short: "Deny gate commands outside the gate agent in Claude Code",
			Args:  noArgs,
			RunE: func(*cobra.Command, []string) error {
				return a.preToolUse(contract)
			},
		},
	)

	return command
}

// sessionStart prints one plain-text line with counts for the session's
// project, or nothing. It never writes and never blocks a session: only a
// configured vault that is missing or unreadable exits non-zero, which hosts
// show as a notice.
func (a *app) sessionStart(ctx context.Context, contract int) error {
	if contract != contractVersion {
		_, err := fmt.Fprintf(a.stdout, "yarbrain: the plugin expects CLI contract %d, but this CLI provides %d; update the yarbrain CLI.\n", contract, contractVersion)

		return err
	}

	payload, err := hook.ReadPayload(a.stdin)
	if err != nil {
		_, _ = fmt.Fprintf(a.stderr, "yarbrain: %v\n", err)

		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, hookDeadline)
	defer cancel()

	v, err := a.openHookVault()
	if err != nil {
		return err
	}

	if v == nil {
		return nil
	}
	defer func() { _ = v.Close() }()

	dir := payload.Cwd
	if dir == "" {
		dir, _ = os.Getwd()
	}

	line, err := a.countsLine(ctx, v, dir)
	if err != nil {
		_, _ = fmt.Fprintf(a.stderr, "yarbrain: %v\n", err)

		return &exitError{code: hookNotice}
	}

	if line != "" {
		_, err = fmt.Fprintln(a.stdout, line)
	}

	return err
}

// openHookVault opens the configured vault. Without a configuration the hook
// stays silent and returns no vault.
func (a *app) openHookVault() (*vault.Vault, error) {
	location, err := a.locator().Resolve()
	if errors.Is(err, vault.ErrNotConfigured) {
		return nil, nil
	}

	if err == nil {
		var v *vault.Vault
		if v, err = vault.Open(location.Vault); err == nil {
			return v, nil
		}
	}

	_, _ = fmt.Fprintf(a.stderr, "yarbrain: %v; run `yarbrain status` for the fix.\n", err)

	return nil, &exitError{code: hookNotice}
}

func (a *app) countsLine(ctx context.Context, v *vault.Vault, dir string) (string, error) {
	started := a.now()

	identity, err := project.Resolve(ctx, dir)
	if err != nil {
		return "", err
	}

	view, err := viewFor(v, identity)
	if err != nil {
		return "", err
	}

	if a.now().Sub(started) > hookDeadline {
		return "yarbrain: lesson counts are unavailable because the scan took longer than 2 seconds; run `yarbrain status`.", nil
	}

	return hook.SessionStartLine(view.counts()), nil
}

// preToolUse blocks brief and verdict outside the gate agent in Claude Code.
// It gives no opinion on anything else, and none in Codex, whose plugin hooks
// get PLUGIN_ROOT and whose threads have distinct session ids.
func (a *app) preToolUse(contract int) error {
	if contract != contractVersion || a.getenv("PLUGIN_ROOT") != "" {
		return nil
	}

	payload, err := hook.ReadPayload(a.stdin)
	if err != nil || !hook.DeniesGateCommand(payload) {
		return nil //nolint:nilerr // an unreadable payload gets no opinion; the host's own permission flow applies
	}

	output, err := hook.DenyOutput()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(a.stdout, string(output))

	return err
}
