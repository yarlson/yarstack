package search

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
)

func candidate(id, title, cue, body string) Candidate {
	return Candidate{Lesson: lesson.Lesson{ID: id, Draft: lesson.Draft{Title: title, Cue: cue, Body: body}}}
}

func rankedIDs(candidates []Candidate) []string {
	out := make([]string, 0, len(candidates))
	for _, c := range candidates {
		out = append(out, c.Lesson.ID)
	}

	return out
}

func TestRankScoresCueOverTitleOverBody(t *testing.T) {
	candidates := []Candidate{
		candidate("body", "Unrelated", "unrelated", "mentions seed here"),
		candidate("cue", "Unrelated", "make seed fails", "nothing"),
		candidate("title", "Seed the database", "unrelated", "nothing"),
		candidate("none", "Unrelated", "unrelated", "nothing"),
	}

	assert.Equal(t, []string{"cue", "title", "body"}, rankedIDs(Rank([]string{"seed"}, candidates)))
}

func TestRankAddsScoresAcrossTermsAndIgnoresCase(t *testing.T) {
	candidates := []Candidate{
		candidate("one", "Seed data", "x", "y"),
		candidate("both", "Seed data", "MIGRATE first", "y"),
	}

	assert.Equal(t, []string{"both", "one"}, rankedIDs(Rank([]string{"seed", "migrate"}, candidates)))
}

func TestTiesGoToProjectThenConfirmationsThenRecency(t *testing.T) {
	now := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)
	cross := candidate("cross", "Seed", "x", "y")
	cross.CrossProject = true
	confirmed := candidate("confirmed", "Seed", "x", "y")
	confirmed.Confirmations = 3
	recent := candidate("recent", "Seed", "x", "y")
	recent.Recent = now
	old := candidate("old", "Seed", "x", "y")
	old.Recent = now.AddDate(0, -1, 0)

	assert.Equal(t, []string{"confirmed", "recent", "old", "cross"},
		rankedIDs(Rank([]string{"seed"}, []Candidate{cross, old, recent, confirmed})))
}

func TestWithoutTermsEveryCandidateIsListed(t *testing.T) {
	candidates := []Candidate{candidate("a", "A", "x", "y"), candidate("b", "B", "x", "y")}

	assert.Len(t, Rank(nil, candidates), 2)
}

func TestCheckTerms(t *testing.T) {
	require.NoError(t, CheckTerms([]string{"seed", "make", "cmd/yarbrain", "v1.2.3", "go_test"}))

	var termErr *TermError

	require.ErrorAs(t, CheckTerms([]string{"x"}), &termErr)
	require.ErrorAs(t, CheckTerms([]string{"two words"}), &termErr)
	require.ErrorAs(t, CheckTerms([]string{`"quoted"`}), &termErr)
	require.ErrorAs(t, CheckTerms([]string{"a1", "a2", "a3", "a4", "a5", "a6", "a7", "a8", "a9"}), &termErr)
	assert.Equal(t, "a9", termErr.Term)
}
