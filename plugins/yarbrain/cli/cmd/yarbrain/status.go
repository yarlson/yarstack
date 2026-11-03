package main

import (
	"github.com/spf13/cobra"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/vault"
)

type statusResult struct {
	header

	Contract int            `json:"contract"`
	Version  string         `json:"version"`
	Vault    statusVault    `json:"vault"`
	Counts   map[string]int `json:"counts"`
	Findings int            `json:"findings"`
}

type statusVault struct {
	vault.Location

	Format int  `json:"format"`
	Legacy bool `json:"legacy"`
}

func newStatusCommand(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the contract version, the vault, and lesson counts",
		Long: "status is the first call an agent makes: it reports the CLI contract\n" +
			"version, where the vault is, and how many lessons it holds by status.",
		Args: noArgs,
		RunE: func(*cobra.Command, []string) error {
			v, location, err := a.openVault()
			if err != nil {
				return err
			}
			defer func() { _ = v.Close() }()

			format, err := v.Format()
			if err != nil {
				return err
			}

			entries, findings, err := v.Lessons()
			if err != nil {
				return err
			}

			counts := map[string]int{}
			for _, status := range lesson.Statuses {
				counts[string(status)] = 0
			}

			for _, entry := range entries {
				counts[string(entry.Lesson.Status)]++
			}

			return writeJSON(a.stdout, statusResult{
				header:   okHeader(),
				Contract: contractVersion,
				Version:  version,
				Vault:    statusVault{Location: location, Format: format, Legacy: v.HasLegacy()},
				Counts:   counts,
				Findings: len(findings),
			})
		},
	}
}
