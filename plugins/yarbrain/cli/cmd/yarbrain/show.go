package main

import (
	"github.com/spf13/cobra"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/trust"
)

// shownLesson carries lesson text only when the lesson is served. Text that no
// gate trusted never reaches an agent through show.
type shownLesson struct {
	ID            string   `json:"id"`
	Status        string   `json:"status"`
	Served        bool     `json:"served"`
	Title         string   `json:"title,omitempty"`
	Cue           string   `json:"cue,omitempty"`
	Body          string   `json:"body,omitempty"`
	Observed      string   `json:"observed,omitempty"`
	Scope         []string `json:"scope,omitempty"`
	ProjectName   string   `json:"project_name,omitempty"`
	Created       string   `json:"created,omitempty"`
	Retires       string   `json:"retires,omitempty"`
	Confirmations int      `json:"confirmations"`
	Path          string   `json:"path"`
	Rev           string   `json:"rev"`
	Digest        string   `json:"digest"`
}

type showResult struct {
	header

	Lessons []shownLesson `json:"lessons"`
}

func newShowCommand(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>...",
		Short: "Show lessons by id",
		Long: "show returns the full text of each served lesson: trusted by an independent\n" +
			"gate, unchanged since. For any other lesson it returns only its id and status.",
		Example: "  yarbrain show l-7k3m9q2x",
		Args:    minimumArgs(1, "id"),
		RunE: func(_ *cobra.Command, ids []string) error {
			v, _, err := a.openVault()
			if err != nil {
				return err
			}
			defer func() { _ = v.Close() }()

			entries, _, err := v.Lessons()
			if err != nil {
				return err
			}

			result := showResult{header: okHeader()}

			for _, id := range ids {
				target, err := loadRecord(v, entries, id)
				if err != nil {
					return err
				}

				result.Lessons = append(result.Lessons, present(target))
			}

			return writeJSON(a.stdout, result)
		},
	}
}

func present(target record) shownLesson {
	l := target.lesson()
	draft := lesson.Normalize(l.Draft)
	shown := shownLesson{
		ID: l.ID, Status: string(target.effective()), Path: target.entry.Path, Rev: target.entry.Rev, Digest: target.digest(),
	}

	if !target.served() {
		return shown
	}

	shown.Served = true
	shown.Title, shown.Cue, shown.Body, shown.Observed = draft.Title, draft.Cue, draft.Body, draft.Observed
	shown.Scope, shown.ProjectName, shown.Created, shown.Retires = draft.Scope, l.ProjectName, l.Created, draft.Retires
	shown.Confirmations, _ = trust.Confirmations(target.events)

	return shown
}
