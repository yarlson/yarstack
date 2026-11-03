package main

import (
	"github.com/spf13/cobra"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/project"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/vault"
)

type statusResult struct {
	header

	Contract        int              `json:"contract"`
	Version         string           `json:"version"`
	Vault           statusVault      `json:"vault"`
	Project         project.Identity `json:"project"`
	Counts          statusCounts     `json:"counts"`
	Totals          map[string]int   `json:"totals"`
	Findings        int              `json:"findings"`
	CodexGatePinned bool             `json:"codex_gate_pinned"`
}

type statusVault struct {
	vault.Location

	Format int `json:"format"`
}

// statusCounts are the counts for the current project, as the session-start
// hook reports them.
type statusCounts struct {
	Trusted      int `json:"trusted"`
	CrossProject int `json:"cross_project"`
	AwaitingGate int `json:"awaiting_gate"`
	Legacy       int `json:"legacy"`
}

func newStatusCommand(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the contract, the vault, and this project's lesson counts",
		Long: "status is the first call an agent makes. It reports the CLI contract version,\n" +
			"where the vault is, the current project's identity, its counts of served\n" +
			"lessons, lessons awaiting the gate, and unconverted v0.2 notes, and the\n" +
			"vault-wide totals by status.",
		Args: noArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			identity, err := a.projectIdentity(command)
			if err != nil {
				return err
			}

			v, location, err := a.openVault()
			if err != nil {
				return err
			}
			defer func() { _ = v.Close() }()

			format, err := v.Format()
			if err != nil {
				return err
			}

			view, err := viewFor(v, identity)
			if err != nil {
				return err
			}

			entries, findings, err := v.Lessons()
			if err != nil {
				return err
			}

			counts := view.counts()

			return writeJSON(a.stdout, statusResult{
				header: okHeader(), Contract: contractVersion, Version: version,
				Vault:   statusVault{Location: location, Format: format},
				Project: identity,
				Counts: statusCounts{
					Trusted: counts.Trusted, CrossProject: counts.CrossProject,
					AwaitingGate: counts.AwaitingGate, Legacy: counts.Legacy,
				},
				Totals:          totals(entries),
				Findings:        len(findings),
				CodexGatePinned: a.codexGatePinned(),
			})
		},
	}
}

func totals(entries []vault.Entry) map[string]int {
	counts := map[string]int{string(lesson.Changed): 0}
	for _, status := range lesson.Statuses {
		counts[string(status)] = 0
	}

	for _, entry := range entries {
		counts[string(lesson.Effective(entry.Lesson))]++
	}

	return counts
}
