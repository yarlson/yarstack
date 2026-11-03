package main

import (
	"github.com/spf13/cobra"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/vault"
)

// shownLesson carries lesson text only when the lesson is served. Text that no
// gate trusted never reaches an agent through show.
type shownLesson struct {
	ID          string   `json:"id"`
	Status      string   `json:"status"`
	Served      bool     `json:"served"`
	Title       string   `json:"title,omitempty"`
	Cue         string   `json:"cue,omitempty"`
	Body        string   `json:"body,omitempty"`
	Observed    string   `json:"observed,omitempty"`
	Scope       []string `json:"scope,omitempty"`
	ProjectName string   `json:"project_name,omitempty"`
	Created     string   `json:"created,omitempty"`
	Retires     string   `json:"retires,omitempty"`
	Path        string   `json:"path"`
	Rev         string   `json:"rev"`
	Digest      string   `json:"digest"`
}

type showResult struct {
	header

	Lessons []shownLesson `json:"lessons"`
}

func newShowCommand(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>...",
		Short: "Show lessons by id",
		Long: "show returns the full text of each trusted lesson. For a lesson that is not\n" +
			"served, such as one that awaits the gate, it returns only its id and status.",
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
				entry, ok := vault.Find(entries, id)
				if !ok {
					return newError("not_found", "no lesson has id "+id, "Use `yarbrain find <terms>` to get lesson ids.").
						with("id", id)
				}

				result.Lessons = append(result.Lessons, present(entry))
			}

			return writeJSON(a.stdout, result)
		},
	}
}

func present(entry vault.Entry) shownLesson {
	l := entry.Lesson
	draft := lesson.Normalize(l.Draft)
	shown := shownLesson{ID: l.ID, Status: string(l.Status), Path: entry.Path, Rev: entry.Rev, Digest: lesson.Digest(draft)}

	if l.Status != lesson.Trusted {
		return shown
	}

	shown.Served = true
	shown.Title, shown.Cue, shown.Body, shown.Observed = draft.Title, draft.Cue, draft.Body, draft.Observed
	shown.Scope, shown.ProjectName, shown.Created, shown.Retires = draft.Scope, l.ProjectName, l.Created, draft.Retires

	return shown
}
