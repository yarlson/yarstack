package main

import (
	"github.com/spf13/cobra"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/trust"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/vault"
)

// shownLesson carries lesson text only when the lesson is served, and only the
// fields its digest covers. Text that no gate trusted never reaches an agent
// through show.
type shownLesson struct {
	ID            string   `json:"id"`
	Status        string   `json:"status"`
	Served        bool     `json:"served"`
	Title         string   `json:"title,omitempty"`
	Cue           string   `json:"cue,omitempty"`
	Body          string   `json:"body,omitempty"`
	Observed      string   `json:"observed,omitempty"`
	Scope         []string `json:"scope,omitempty"`
	Retires       string   `json:"retires,omitempty"`
	Confirmations int      `json:"confirmations"`
	Path          string   `json:"path"`
	Rev           string   `json:"rev"`
	Digest        string   `json:"digest"`
}

type showResult struct {
	header

	Lessons     []shownLesson `json:"lessons"`
	LegacyNotes []legacyNote  `json:"legacy_notes,omitempty"`
}

// legacyNote is a v0.2 note read by path for conversion. No gate ever judged
// its text.
type legacyNote struct {
	Path  string `json:"path"`
	Gated bool   `json:"gated"`
	Text  string `json:"text"`
}

func newShowCommand(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "show <id|v0.2-path>...",
		Short: "Show lessons by id, or v0.2 notes by path",
		Long: "show returns the full text of each served lesson: trusted by an independent\n" +
			"gate, unchanged since. For any other lesson it returns only its id and status.\n" +
			"A v0.2 note path from `yarbrain legacy` returns that note's text, marked\n" +
			"gated: false, for conversion into a lesson.",
		Example: "  yarbrain show l-7k3m9q2x\n  yarbrain show notes/projects/shop.md",
		Args:    minimumArgs(1, "id"),
		RunE: func(_ *cobra.Command, args []string) error {
			v, _, err := a.openVault()
			if err != nil {
				return err
			}
			defer func() { _ = v.Close() }()

			entries, _, err := v.Lessons()
			if err != nil {
				return err
			}

			result := showResult{header: okHeader(), Lessons: []shownLesson{}}

			for _, arg := range args {
				if !lesson.IsID(arg) {
					note, err := readLegacy(v, arg)
					if err != nil {
						return err
					}

					result.LegacyNotes = append(result.LegacyNotes, note)

					continue
				}

				target, err := loadRecord(v, entries, arg)
				if err != nil {
					return err
				}

				result.Lessons = append(result.Lessons, present(target))
			}

			return writeJSON(a.stdout, result)
		},
	}
}

func readLegacy(v *vault.Vault, path string) (legacyNote, error) {
	if !vault.IsLegacyPath(path) {
		return legacyNote{}, newError("bad_input", path+" is neither a lesson id nor a v0.2 note path",
			"Pass a lesson id such as l-7k3m9q2x, or a path from `yarbrain legacy`.").with("field", "arguments")
	}

	raw, err := v.ReadFile(path)
	if err != nil {
		return legacyNote{}, newError("not_found", "no v0.2 note at "+path, "Use a path from `yarbrain legacy`.").with("path", path)
	}

	return legacyNote{Path: path, Gated: false, Text: string(raw[:min(len(raw), lesson.MaxFileBytes)])}, nil
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
	shown.Scope, shown.Retires = draft.Scope, draft.Retires
	shown.Confirmations, _ = trust.Confirmations(target.events)

	return shown
}
