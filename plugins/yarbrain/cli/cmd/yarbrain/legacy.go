package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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

func newLegacyCommand(a *app) *cobra.Command {
	var limit int

	command := &cobra.Command{
		Use:   "legacy",
		Short: "List v0.2 notes that are not converted to lessons yet",
		Long: "legacy lists the paths of v0.2 notes and vault skills that no lesson converts\n" +
			"yet. Read one with `yarbrain show <path>`, then propose a lesson with\n" +
			"converts: <path>. The text of a v0.2 note was never gated.",
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

	return command
}

// codexGatePinned reports whether the one-time Codex setup pinned the gate
// agent's model.
func (a *app) codexGatePinned() bool {
	_, err := os.Stat(a.home + "/.codex/agents/yarbrain-gate.toml")

	return err == nil
}
