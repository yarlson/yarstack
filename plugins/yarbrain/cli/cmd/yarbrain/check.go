package main

import (
	"fmt"
	"slices"
	"time"

	"github.com/spf13/cobra"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/trust"
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

				events, err := v.Ledger(entry.Lesson.ID)
				if err != nil {
					findings = append(findings, vault.Finding{Path: entry.Path, Code: "bad_file", Message: err.Error()})

					continue
				}

				findings = append(findings, trustFindings(entry, events, a.now())...)
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

// staleAfter is how long a lesson may wait for a gate before check reports it.
const staleAfter = 7 * 24 * time.Hour

func trustFindings(entry vault.Entry, events []trust.Event, now time.Time) []vault.Finding {
	effective := lesson.Effective(entry.Lesson)

	if effective == lesson.Trusted && !trust.Served(entry.Lesson, events) {
		return []vault.Finding{{
			Path: entry.Path, Code: "not_credible",
			Message: "the lesson is marked trusted, but no credible gate verdict backs it; it is not served",
		}}
	}

	if lesson.AwaitsGate(effective) && len(events) > 0 && now.Sub(events[len(events)-1].At) > staleAfter {
		return []vault.Finding{{
			Path: entry.Path, Code: "awaiting_gate",
			Message: fmt.Sprintf("the lesson has waited for a gate since %s", events[len(events)-1].At.Format(time.DateOnly)),
		}}
	}

	return nil
}

func nonNil[T any](items []T) []T {
	if items == nil {
		return []T{}
	}

	return items
}
