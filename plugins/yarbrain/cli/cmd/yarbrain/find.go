package main

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/search"
)

type menuEntry struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Cue           string `json:"cue"`
	CrossProject  bool   `json:"cross_project"`
	Confirmations int    `json:"confirmations"`
	LastConfirmed string `json:"last_confirmed,omitempty"`
}

type findResult struct {
	header

	ProjectName  string      `json:"project_name,omitempty"`
	Lessons      []menuEntry `json:"lessons"`
	More         int         `json:"more"`
	AwaitingGate int         `json:"awaiting_gate"`
	Legacy       int         `json:"legacy"`
	Skipped      []string    `json:"skipped"`
}

func newFindCommand(a *app) *cobra.Command {
	var limit int

	command := &cobra.Command{
		Use:   "find [term...]",
		Short: "List served lessons for this project, best match first",
		Long: "find returns a short menu of served lessons for the current project and\n" +
			"cross-project lessons: id, title, and cue. Pick the lessons whose cue fits\n" +
			"your situation and read them with show. Terms are ranked, not all required:\n" +
			"a term scores 3 in the cue, 2 in the title, 1 in the body. Without terms,\n" +
			"find lists every served lesson. Lessons that await the gate and v0.2 notes\n" +
			"appear only as counts.",
		Example: "  yarbrain find seed migrate\n  yarbrain find --limit 10",
		RunE: func(command *cobra.Command, terms []string) error {
			if err := search.CheckTerms(terms); err != nil {
				return termError(err)
			}

			if limit < 1 || limit > search.MaxLimit {
				return newError("bad_input", fmt.Sprintf("--limit must be between 1 and %d", search.MaxLimit), "Pass a smaller --limit.").
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

			ranked := search.Rank(terms, view.candidates)
			result := findResult{
				header: okHeader(), ProjectName: identity.Name, Lessons: []menuEntry{},
				More: max(0, len(ranked)-limit), AwaitingGate: len(view.awaiting), Legacy: len(view.legacy), Skipped: view.skipped,
			}

			for _, candidate := range ranked[:min(limit, len(ranked))] {
				result.Lessons = append(result.Lessons, menuEntry{
					ID: candidate.Lesson.ID, Title: candidate.Lesson.Title, Cue: candidate.Lesson.Cue,
					CrossProject: candidate.CrossProject, Confirmations: candidate.Confirmations,
					LastConfirmed: dateOf(candidate.LastConfirmed),
				})
			}

			return writeJSON(a.stdout, result)
		},
	}
	command.Flags().IntVar(&limit, "limit", search.DefaultLimit, fmt.Sprintf("how many lessons to return, at most %d", search.MaxLimit))

	return command
}

func termError(err error) error {
	if term, ok := errors.AsType[*search.TermError](err); ok {
		return newError("bad_term", err.Error(), "Use 2 or 3 distinctive words such as a command, a file name, or an error token.").
			with("term", term.Term).with("rule", term.Rule)
	}

	return err
}
