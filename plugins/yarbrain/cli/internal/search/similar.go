// Package search ranks lessons for lookups and finds lessons that resemble a
// draft, within fixed budgets.
package search

import (
	"slices"
	"strings"
	"unicode"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
)

// minShared is how many distinctive words two lessons must share in their
// title and cue to count as similar.
const minShared = 2

// Similar returns up to limit candidates whose title and cue share the most
// distinctive words with the draft, most similar first. A gate uses them to
// spot duplicates and contradictions.
func Similar(d lesson.Draft, candidates []lesson.Lesson, limit int) []lesson.Lesson {
	target := words(d.Title + " " + d.Cue)

	type scored struct {
		lesson lesson.Lesson
		shared int
	}

	var matches []scored

	for _, candidate := range candidates {
		shared := 0

		for word := range words(candidate.Title + " " + candidate.Cue) {
			if target[word] {
				shared++
			}
		}

		if shared >= minShared {
			matches = append(matches, scored{lesson: candidate, shared: shared})
		}
	}

	slices.SortStableFunc(matches, func(a, b scored) int { return b.shared - a.shared })

	result := make([]lesson.Lesson, 0, min(limit, len(matches)))
	for _, match := range matches[:min(limit, len(matches))] {
		result = append(result, match.lesson)
	}

	return result
}

// stopWords are too common to signal similarity.
var stopWords = map[string]bool{
	"the": true, "and": true, "for": true, "with": true, "when": true, "before": true,
	"after": true, "from": true, "into": true, "that": true, "this": true, "run": true,
	"use": true, "not": true, "are": true, "was": true, "fails": true,
}

func words(text string) map[string]bool {
	set := map[string]bool{}

	for _, word := range strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	}) {
		if len(word) >= 3 && !stopWords[word] {
			set[word] = true
		}
	}

	return set
}
