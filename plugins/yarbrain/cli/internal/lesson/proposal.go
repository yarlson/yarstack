package lesson

import "slices"

// ScopesOverlap reports whether two scopes reach a common project. An all
// scope overlaps every scope.
func ScopesOverlap(a, b []string) bool {
	if slices.Contains(a, AllScope) || slices.Contains(b, AllScope) {
		return true
	}

	for _, identity := range a {
		if slices.Contains(b, identity) {
			return true
		}
	}

	return false
}

// Outcome is what a new proposal means next to the existing lessons.
type Outcome int

const (
	// OutcomeNew means the draft is a new lesson.
	OutcomeNew Outcome = iota
	// OutcomeRepeat means an identical lesson exists; proposing it again
	// changes nothing.
	OutcomeRepeat
	// OutcomeTitleTaken means an active lesson with the same title and an
	// overlapping scope says something else.
	OutcomeTitleTaken
)

// ClassifyProposal compares a normalized draft with the existing lessons. A
// lesson with the same intent key in an overlapping scope is the same lesson:
// identical content is a repeat, and different content is a conflict unless
// the existing lesson was rejected or retired. A rejection therefore blocks
// only the identical content, never a corrected version.
func ClassifyProposal(existing []Lesson, d Draft) (Outcome, Lesson) {
	key := IntentKey(d.Title)
	digest := Digest(d)

	for _, other := range existing {
		if IntentKey(other.Title) != key || !ScopesOverlap(other.Scope, d.Scope) {
			continue
		}

		if Digest(Normalize(other.Draft)) == digest {
			return OutcomeRepeat, other
		}

		if other.Status != Rejected && other.Status != Retired {
			return OutcomeTitleTaken, other
		}
	}

	return OutcomeNew, Lesson{}
}
