package lesson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func trustedLesson() Lesson {
	l := Lesson{Draft: Normalize(validDraft()), ID: "l-00000001", Status: Trusted}
	l.TrustedDigest = Digest(l.Draft)

	return l
}

func TestEffectiveStatusDetectsEditsOfATrustedLesson(t *testing.T) {
	l := trustedLesson()
	assert.Equal(t, Trusted, Effective(l))

	l.Body = "edited by hand"
	assert.Equal(t, Changed, Effective(l))

	handTyped := trustedLesson()
	handTyped.TrustedDigest = ""
	assert.Equal(t, Changed, Effective(handTyped), "status: trusted without a digest is not trust")
}

func TestEffectiveStatusKeepsOtherStatuses(t *testing.T) {
	for _, status := range []Status{Proposed, Deferred, Rejected, Doubted, Retired} {
		l := trustedLesson()
		l.Status = status
		assert.Equal(t, status, Effective(l))
	}
}

func TestVerdictTransitions(t *testing.T) {
	cases := []struct {
		from     Status
		decision Decision
		to       Status
		ok       bool
	}{
		{Proposed, DecideTrust, Trusted, true},
		{Proposed, DecideReject, Rejected, true},
		{Proposed, DecideDefer, Deferred, true},
		{Proposed, DecideRetire, "", false},
		{Doubted, DecideTrust, Trusted, true},
		{Doubted, DecideRetire, Retired, true},
		{Doubted, DecideReject, "", false},
		{Changed, DecideTrust, Trusted, true},
		{Deferred, DecideReject, Rejected, true},
		{Trusted, DecideTrust, "", false},
		{Rejected, DecideTrust, "", false},
		{Retired, DecideTrust, "", false},
	}
	for _, tc := range cases {
		to, ok := AfterVerdict(tc.from, tc.decision)

		assert.Equal(t, tc.ok, ok, "%s + %s", tc.from, tc.decision)
		assert.Equal(t, tc.to, to, "%s + %s", tc.from, tc.decision)
	}
}

func TestOnlyUnsettledLessonsAwaitTheGate(t *testing.T) {
	assert.True(t, AwaitsGate(Proposed))
	assert.True(t, AwaitsGate(Doubted))
	assert.True(t, AwaitsGate(Changed))
	assert.True(t, AwaitsGate(Deferred))
	assert.False(t, AwaitsGate(Trusted))
	assert.False(t, AwaitsGate(Rejected))
	assert.False(t, AwaitsGate(Retired))
	assert.Equal(t, []Decision{DecideTrust, DecideRetire, DecideDefer}, AllowedDecisions(Doubted))
}
