package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/vault"
)

const (
	defaultLegacyLimit = 10
	maxLegacyLimit     = 50
)

type legacyResult struct {
	header

	Paths []string `json:"paths"`
	More  int      `json:"more"`
}

type dismissResult struct {
	header

	Path string `json:"path"`
}

func newLegacyCommand(a *app) *cobra.Command {
	var limit int

	command := &cobra.Command{
		Use:   "legacy",
		Short: "List v0.2 notes that are not converted to lessons yet",
		Long: "legacy lists the paths of v0.2 notes and vault skills that no lesson converts\n" +
			"yet. Read one with `yarbrain show <path>`, then propose a lesson with\n" +
			"converts: <path>, or dismiss it with `yarbrain legacy dismiss <path>` when it\n" +
			"holds nothing worth a lesson. The text of a v0.2 note was never gated.",
		Example: "  yarbrain legacy --limit 3",
		Args:    noArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if limit < 1 || limit > maxLegacyLimit {
				return newError("bad_input", fmt.Sprintf("--limit must be between 1 and %d", maxLegacyLimit), "Pass a smaller --limit.").
					with("field", "limit")
			}

			identity, err := a.projectIdentity(command)
			if err != nil {
				return err
			}

			v, _, err := a.openVault()
			if err != nil {
				return err
			}
			defer func() { _ = v.Close() }()

			view, err := viewFor(v, identity)
			if err != nil {
				return err
			}

			return writeJSON(a.stdout, legacyResult{
				header: okHeader(), Paths: nonNil(view.legacy[:min(limit, len(view.legacy))]), More: max(0, len(view.legacy)-limit),
			})
		},
	}
	command.Flags().IntVar(&limit, "limit", defaultLegacyLimit, fmt.Sprintf("how many paths to return, at most %d", maxLegacyLimit))
	command.AddCommand(newDismissCommand(a))

	return command
}

type dismissInput struct {
	Reason string `json:"reason"`
}

func newDismissCommand(a *app) *cobra.Command {
	var input string

	command := &cobra.Command{
		Use:   "dismiss <v0.2-path>",
		Short: "Record that a v0.2 note holds nothing worth a lesson",
		Long: "dismiss takes a v0.2 note out of the conversion backlog. The note itself stays\n" +
			"untouched. Input on stdin: {\"reason\": \"<one line, at most 200 characters>\"}",
		Example: "  echo '{\"reason\": \"status notes from one sprint\"}' | yarbrain legacy dismiss notes/projects/old.md",
		Args:    exactArgs(1, "v0.2-path"),
		RunE: func(_ *cobra.Command, args []string) error {
			var in dismissInput
			if err := a.readJSON(input, &in, `{"reason": "<why the note holds nothing worth a lesson>"}`); err != nil {
				return err
			}

			if problems := oneLine("reason", in.Reason, maxReason); len(problems) > 0 {
				return problemsError(problems)
			}

			if failure := screen(map[string]string{"reason": in.Reason}); failure != nil {
				return failure
			}

			v, _, err := a.openVault()
			if err != nil {
				return err
			}
			defer func() { _ = v.Close() }()

			if _, err := readLegacy(v, args[0]); err != nil {
				return err
			}

			err = v.WithLock(func() error {
				return v.Dismiss(vault.Dismissal{At: a.now().UTC(), Path: args[0], Session: a.session(), Reason: in.Reason})
			})
			if err != nil {
				return err
			}

			return writeJSON(a.stdout, dismissResult{header: okHeader(), Path: args[0]})
		},
	}
	command.Flags().StringVar(&input, "input", "", "read the JSON from this file instead of stdin")

	return command
}

// codexGatePinned reports whether the one-time Codex setup pinned the gate
// agent's model.
func (a *app) codexGatePinned() bool {
	_, err := os.Stat(a.home + "/.codex/agents/yarbrain-gate.toml")

	return err == nil
}
