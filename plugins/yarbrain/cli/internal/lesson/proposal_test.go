package lesson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScopesOverlap(t *testing.T) {
	assert.True(t, ScopesOverlap([]string{"origin:a", "root:1"}, []string{"root:1"}))
	assert.True(t, ScopesOverlap([]string{"all"}, []string{"root:9"}))
	assert.True(t, ScopesOverlap([]string{"root:9"}, []string{"all"}))
	assert.False(t, ScopesOverlap([]string{"origin:a"}, []string{"origin:b", "root:2"}))
}

func stored(status Status, d Draft) Lesson {
	return Lesson{Draft: d, ID: "l-00000001", Status: status}
}

func TestIdenticalProposalIsARepeat(t *testing.T) {
	d := Normalize(validDraft())

	outcome, match := ClassifyProposal([]Lesson{stored(Trusted, validDraft())}, d)

	assert.Equal(t, OutcomeRepeat, outcome)
	assert.Equal(t, "l-00000001", match.ID)
}

func TestSameTitleWithOtherContentIsTaken(t *testing.T) {
	other := validDraft()
	other.Body = "Something else entirely."

	outcome, match := ClassifyProposal([]Lesson{stored(Proposed, other)}, Normalize(validDraft()))

	assert.Equal(t, OutcomeTitleTaken, outcome)
	assert.Equal(t, Proposed, match.Status)
}

func TestTitleMatchingIgnoresCase(t *testing.T) {
	other := validDraft()
	other.Title = "RUN MIGRATIONS before seeding the test database"
	other.Body = "Different."

	outcome, _ := ClassifyProposal([]Lesson{stored(Trusted, other)}, Normalize(validDraft()))

	assert.Equal(t, OutcomeTitleTaken, outcome)
}

func TestARejectionBlocksOnlyTheIdenticalContent(t *testing.T) {
	rejected := stored(Rejected, validDraft())

	corrected := validDraft()
	corrected.Body = "Run `make migrate` once after `make db-up`."

	identical, _ := ClassifyProposal([]Lesson{rejected}, Normalize(validDraft()))
	fixed, _ := ClassifyProposal([]Lesson{rejected}, Normalize(corrected))

	assert.Equal(t, OutcomeRepeat, identical, "re-proposing rejected content returns the rejected lesson")
	assert.Equal(t, OutcomeNew, fixed, "a corrected version may be proposed")
}

func TestOtherProjectsDoNotCollide(t *testing.T) {
	other := validDraft()
	other.Scope = []string{"origin:github.com/acme/other"}
	other.Body = "Different."

	outcome, _ := ClassifyProposal([]Lesson{stored(Trusted, other)}, Normalize(validDraft()))

	assert.Equal(t, OutcomeNew, outcome)
}
