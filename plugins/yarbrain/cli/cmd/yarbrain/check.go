package main

import (
	"fmt"
	"slices"

	"github.com/spf13/cobra"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/vault"
)

type checkResult struct {
	header

	Checked  int             `json:"checked"`
	Findings []vault.Finding `json:"findings"`
}

func newCheckCommand(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Validate every lesson file in the vault",
		Long: "check reads every lesson and reports each problem as a finding. It exits 1\n" +
			"when there are findings. A lesson with a finding is never served.",
		Args: noArgs,
		RunE: func(*cobra.Command, []string) error {
			v, _, err := a.openVault()
			if err != nil {
				return err
			}
			defer func() { _ = v.Close() }()

			entries, findings, err := v.Lessons()
			if err != nil {
				return err
			}

			for _, entry := range entries {
				findings = append(findings, lessonFindings(entry)...)
			}

			if err := writeJSON(a.stdout, checkResult{header: okHeader(), Checked: len(entries), Findings: nonNil(findings)}); err != nil {
				return err
			}

			if len(findings) > 0 {
				return &exitError{code: exitFindings}
			}

			return nil
		},
	}
}

func lessonFindings(entry vault.Entry) []vault.Finding {
	var findings []vault.Finding

	if !slices.Contains(lesson.Statuses, entry.Lesson.Status) {
		findings = append(findings, vault.Finding{
			Path: entry.Path, Code: "unknown_status",
			Message: fmt.Sprintf("status %q is not one of %v", entry.Lesson.Status, lesson.Statuses),
		})
	}

	for _, problem := range lesson.Validate(lesson.Normalize(entry.Lesson.Draft)) {
		message := fmt.Sprintf("%s: %s", problem.Field, problem.Rule)
		if problem.Limit > 0 {
			message += fmt.Sprintf(" (limit %d, actual %d)", problem.Limit, problem.Actual)
		}

		findings = append(findings, vault.Finding{Path: entry.Path, Code: "invalid_field", Message: message})
	}

	return findings
}

func nonNil[T any](items []T) []T {
	if items == nil {
		return []T{}
	}

	return items
}
