package main

import (
	"slices"
	"time"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/hook"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/project"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/search"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/trust"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/vault"
)

// maxSkipped bounds the unreadable paths a lookup reports.
const maxSkipped = 10

// scopedView is what one project sees of the vault. find, status, and the
// session-start hook all read it, so they always agree.
type scopedView struct {
	candidates []search.Candidate
	// awaiting holds the ids of lessons in scope that await the gate, oldest
	// file first. Ids carry no lesson text.
	awaiting []string
	legacy   []string
	skipped  []string
}

func (s scopedView) counts() hook.Counts {
	counts := hook.Counts{Trusted: len(s.candidates), AwaitingGate: len(s.awaiting), Legacy: len(s.legacy)}

	for _, candidate := range s.candidates {
		if candidate.CrossProject {
			counts.CrossProject++
		}
	}

	return counts
}

// viewFor reads the vault as the project with identity sees it: served
// lessons in its scope or cross-project, minus any lesson that a served
// successor retires, plus the lessons awaiting the gate and the unconverted
// v0.2 notes. It only reads.
func viewFor(v *vault.Vault, identity project.Identity) (scopedView, error) {
	entries, findings, err := v.Lessons()
	if err != nil {
		return scopedView{}, err
	}

	view := scopedView{skipped: []string{}}
	for _, finding := range findings {
		view.addSkipped(finding.Path)
	}

	view.collect(v, entries, identity)

	if view.legacy, err = unconvertedNotes(v, entries); err != nil {
		return scopedView{}, err
	}

	return view, nil
}

// collect gathers the served lessons and the lessons awaiting the gate in the
// project's scope. A lesson that a served lesson retires stays hidden, so a
// crash between trusting a successor and retiring its predecessor never
// serves both.
func (s *scopedView) collect(v *vault.Vault, entries []vault.Entry, identity project.Identity) {
	retiredBySuccessor := map[string]bool{}

	for _, entry := range entries {
		l := entry.Lesson
		if !inScope(l, identity) {
			continue
		}

		events, err := v.Ledger(l.ID)
		if err != nil {
			s.addSkipped(entry.Path)

			continue
		}

		switch {
		case trust.Served(l, events):
			s.candidates = append(s.candidates, candidateOf(l, events))
			retiredBySuccessor[l.Retires] = true
		case lesson.AwaitsGate(lesson.Effective(l)):
			s.awaiting = append(s.awaiting, l.ID)
		}
	}

	s.candidates = slices.DeleteFunc(s.candidates, func(c search.Candidate) bool { return retiredBySuccessor[c.Lesson.ID] })
}

// unconvertedNotes lists the v0.2 notes that await conversion: no lesson
// converts them, except through a rejected conversion, and no session
// dismissed them.
func unconvertedNotes(v *vault.Vault, entries []vault.Entry) ([]string, error) {
	handled, err := v.Dismissed()
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.Lesson.Converts != "" && entry.Lesson.Status != lesson.Rejected {
			handled[entry.Lesson.Converts] = true
		}
	}

	notes, err := v.LegacyNotes()
	if err != nil {
		return nil, err
	}

	return slices.DeleteFunc(notes, func(note string) bool { return handled[note] }), nil
}

func (s *scopedView) addSkipped(path string) {
	if len(s.skipped) < maxSkipped {
		s.skipped = append(s.skipped, path)
	}
}

// inScope reports whether a lesson applies to the project. Outside a known
// project only cross-project lessons apply.
func inScope(l lesson.Lesson, identity project.Identity) bool {
	if l.IsAllScope() {
		return true
	}

	ids := identity.IDs()

	return len(ids) > 0 && lesson.ScopesOverlap(l.Scope, ids)
}

func candidateOf(l lesson.Lesson, events []trust.Event) search.Candidate {
	confirmations, lastConfirm := trust.Confirmations(events)

	recent := lastConfirm
	for _, e := range events {
		if e.Kind == trust.Verdict && e.At.After(recent) {
			recent = e.At
		}
	}

	return search.Candidate{
		Lesson:        lesson.Lesson{Draft: lesson.Normalize(l.Draft), ID: l.ID, Status: l.Status, ProjectName: l.ProjectName, Created: l.Created},
		CrossProject:  l.IsAllScope(),
		Confirmations: confirmations,
		LastConfirmed: lastConfirm,
		Recent:        recent,
	}
}

// dateOf formats a time as a date, or "" for the zero time.
func dateOf(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format(time.DateOnly)
}
