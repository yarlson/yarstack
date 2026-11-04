package search

import (
	"cmp"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
)

// Budgets for one lookup. A menu entry costs about 40 tokens, so the default
// menu stays near 200 tokens.
const (
	DefaultLimit = 5
	MaxLimit     = 20
	MaxTerms     = 8
)

// termPattern accepts terms that need no quoting in sh, cmd.exe, or
// PowerShell.
var termPattern = regexp.MustCompile(`^[A-Za-z0-9._/-]{2,40}$`)

// TermError names a search term that breaks the term rules.
type TermError struct {
	Term string
	Rule string
}

func (e *TermError) Error() string { return fmt.Sprintf("search term %q: %s", e.Term, e.Rule) }

// CheckTerms validates search terms.
func CheckTerms(terms []string) error {
	if len(terms) > MaxTerms {
		return &TermError{Term: terms[MaxTerms], Rule: fmt.Sprintf("at most %d terms", MaxTerms)}
	}

	for _, term := range terms {
		if !termPattern.MatchString(term) {
			return &TermError{Term: term, Rule: "2 to 40 characters from A-Z, a-z, 0-9, and . _ / -"}
		}
	}

	return nil
}

// Candidate is a served lesson with the facts that order the menu.
type Candidate struct {
	Lesson        lesson.Lesson
	CrossProject  bool
	Confirmations int
	LastConfirmed time.Time
	// Recent is the latest trust or confirmation time.
	Recent time.Time
}

// Rank orders candidates for a lookup. With terms, a candidate needs at least
// one match; each term scores 3 in the cue, 2 in the title, and 1 in the
// body. Ties, and a lookup without terms, go to project lessons before cross-
// project ones, then to more confirmations, then to the most recent.
func Rank(terms []string, candidates []Candidate) []Candidate {
	type scored struct {
		Candidate

		score int
	}

	var ranked []scored

	for _, candidate := range candidates {
		score := relevance(terms, candidate.Lesson)
		if len(terms) > 0 && score == 0 {
			continue
		}

		ranked = append(ranked, scored{Candidate: candidate, score: score})
	}

	slices.SortStableFunc(ranked, func(a, b scored) int {
		return cmp.Or(
			cmp.Compare(b.score, a.score),
			compareBool(a.CrossProject, b.CrossProject),
			cmp.Compare(b.Confirmations, a.Confirmations),
			b.Recent.Compare(a.Recent),
		)
	})

	result := make([]Candidate, 0, len(ranked))
	for _, r := range ranked {
		result = append(result, r.Candidate)
	}

	return result
}

func relevance(terms []string, l lesson.Lesson) int {
	cue, title, body := strings.ToLower(l.Cue), strings.ToLower(l.Title), strings.ToLower(l.Body)
	score := 0

	for _, term := range terms {
		term = strings.ToLower(term)

		switch {
		case strings.Contains(cue, term):
			score += 3
		case strings.Contains(title, term):
			score += 2
		case strings.Contains(body, term):
			score++
		}
	}

	return score
}

// compareBool orders false before true.
func compareBool(a, b bool) int {
	switch {
	case a == b:
		return 0
	case !a:
		return -1
	default:
		return 1
	}
}
