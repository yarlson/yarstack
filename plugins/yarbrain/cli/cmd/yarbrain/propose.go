package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/spf13/cobra"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/project"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/trust"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/vault"
)

const proposeTemplate = `---
cue: '<the observable situation: a command, error text, file, or task>'
scope: here
---
# <Imperative title, at most 80 characters>

<What to do and why, at most 1,200 characters.>

## Observed

<What happened in this session that taught it, at most 1,000 characters.>`

type proposeResult struct {
	header

	ID      string `json:"id"`
	Path    string `json:"path"`
	Status  string `json:"status"`
	Created bool   `json:"created"`
	Digest  string `json:"digest"`
}

func newProposeCommand(a *app) *cobra.Command {
	var input string

	command := &cobra.Command{
		Use:   "propose",
		Short: "Propose a lesson for the independent gate",
		Long: "propose stores a new lesson as proposed. It is not served until an\n" +
			"independent gate agent trusts it. Pass the lesson on stdin or with --input:\n\n" +
			proposeTemplate + "\n\n" +
			"scope is here (this project) or all (know-how that does not depend on one\n" +
			"project). Optional keys: retires (the id of a trusted lesson this one\n" +
			"replaces) and converts (the v0.2 note path this lesson converts).\n" +
			"Proposing identical content again is a no-op that returns the same id.",
		Example: "  yarbrain propose <<'EOF'\n" + proposeTemplate + "\n  EOF",
		Args:    noArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			raw, err := a.readInput(input)
			if err != nil {
				return err
			}

			draft, scope, err := vault.ParseProposal(raw)
			if err != nil {
				return err
			}

			identity, err := a.projectIdentity(command)
			if err != nil {
				return err
			}

			if draft, err = resolveScope(draft, scope, identity); err != nil {
				return err
			}

			draft = lesson.Normalize(draft)
			if problems := lesson.Validate(draft); len(problems) > 0 {
				return problemsError(problems)
			}

			if failure := screen(map[string]string{
				"title": draft.Title, "cue": draft.Cue, "body": draft.Body, "observed": draft.Observed,
			}); failure != nil {
				return failure
			}

			result, err := a.storeProposal(draft, identity)
			if err != nil {
				return err
			}

			return writeJSON(a.stdout, result)
		},
	}
	command.Flags().StringVar(&input, "input", "", "read the lesson from this file instead of stdin")

	return command
}

func (a *app) projectIdentity(command *cobra.Command) (project.Identity, error) {
	dir, err := os.Getwd()
	if err != nil {
		return project.Identity{}, fmt.Errorf("find the current directory: %w", err)
	}

	return project.Resolve(command.Context(), dir)
}

func resolveScope(draft lesson.Draft, scope string, identity project.Identity) (lesson.Draft, error) {
	if scope == lesson.AllScope {
		draft.Scope = []string{lesson.AllScope}

		return draft, nil
	}

	if !identity.Known() {
		return draft, newError("unknown_project", "the current directory has no project identity",
			"Run inside a git repository that has a commit, or use scope: all for know-how that applies everywhere.")
	}

	draft.Scope = identity.IDs()

	return draft, nil
}

func (a *app) storeProposal(draft lesson.Draft, identity project.Identity) (proposeResult, error) {
	v, _, err := a.openVault()
	if err != nil {
		return proposeResult{}, err
	}
	defer func() { _ = v.Close() }()

	var result proposeResult

	err = v.WithLock(func() error {
		entries, _, err := v.Lessons()
		if err != nil {
			return err
		}

		if err := checkRetires(entries, draft); err != nil {
			return err
		}

		if err := checkConverts(v, draft); err != nil {
			return err
		}

		outcome, match := lesson.ClassifyProposal(lessonsOf(entries), draft)

		switch outcome {
		case lesson.OutcomeRepeat:
			existing, _ := vault.Find(entries, match.ID)
			result = proposeResult{ID: match.ID, Path: existing.Path, Status: string(match.Status), Digest: lesson.Digest(draft)}

			return nil
		case lesson.OutcomeTitleTaken:
			return newError("title_taken", fmt.Sprintf("lesson %s already has this title with different content", match.ID),
				fmt.Sprintf("Run `yarbrain show %s`. If it covers your lesson, stop. If it is wrong, doubt it and propose a replacement with retires. Otherwise change the title.", match.ID)).
				with("id", match.ID).with("status", string(match.Status))
		}

		created, err := a.createLesson(v, entries, draft, identity)
		if err != nil {
			return err
		}

		result = created

		return v.AppendLedger(created.ID, trust.Event{
			At: a.now().UTC(), Kind: trust.Proposed, Session: a.session(), Digest: created.Digest,
		})
	})

	result.header = okHeader()

	return result, err
}

func (a *app) createLesson(v *vault.Vault, entries []vault.Entry, draft lesson.Draft, identity project.Identity) (proposeResult, error) {
	id, err := lesson.NewID(rand.Reader)
	if err != nil {
		return proposeResult{}, err
	}

	if _, taken := vault.Find(entries, id); taken {
		return proposeResult{}, errors.New("generated a lesson id that already exists; retry")
	}

	stored := lesson.Lesson{Draft: draft, ID: id, Status: lesson.Proposed, Created: a.today()}
	if !draft.IsAllScope() {
		stored.ProjectName = identity.Name
	}

	path, err := v.CreateLesson(stored)
	if err != nil {
		return proposeResult{}, err
	}

	return proposeResult{ID: id, Path: path, Status: string(lesson.Proposed), Created: true, Digest: lesson.Digest(draft)}, nil
}

func checkRetires(entries []vault.Entry, draft lesson.Draft) error {
	if draft.Retires == "" {
		return nil
	}

	target, ok := vault.Find(entries, draft.Retires)
	if !ok {
		return newError("not_found", "retires names no lesson: "+draft.Retires, "Use `yarbrain find` to get the id of the lesson to replace.").
			with("field", "retires")
	}

	if target.Lesson.Status != lesson.Trusted || !lesson.ScopesOverlap(target.Lesson.Scope, draft.Scope) {
		return newError("bad_state", fmt.Sprintf("lesson %s is %s or belongs to another project", draft.Retires, target.Lesson.Status),
			"retires may name only a trusted lesson of the same project.").with("field", "retires")
	}

	return nil
}

func checkConverts(v *vault.Vault, draft lesson.Draft) error {
	if draft.Converts == "" {
		return nil
	}

	if !vault.IsLegacyPath(draft.Converts) {
		return newError("bad_input", "converts must be a v0.2 note path under notes/ or skills/: "+draft.Converts,
			"Use a path from `yarbrain legacy`.").with("field", "converts")
	}

	if _, err := v.ReadFile(draft.Converts); err != nil {
		return newError("bad_input", "converts names no v0.2 note: "+draft.Converts, "Use a path from `yarbrain legacy`.").
			with("field", "converts")
	}

	return nil
}

func lessonsOf(entries []vault.Entry) []lesson.Lesson {
	lessons := make([]lesson.Lesson, 0, len(entries))
	for _, entry := range entries {
		lessons = append(lessons, entry.Lesson)
	}

	return slices.Clip(lessons)
}
