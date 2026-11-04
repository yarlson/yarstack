package trust

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
)

const digest = "2c9d41e07a6b"

var start = time.Date(2026, 11, 3, 10, 0, 0, 0, time.UTC)

func at(minutes int) time.Time { return start.Add(time.Duration(minutes) * time.Minute) }

func proposed(session string) Event {
	return Event{At: at(0), Kind: Proposed, Session: session, Digest: digest}
}

func trusted(session string, minutes int) Event {
	return Event{
		At: at(minutes), Kind: Verdict, Session: session, Digest: digest, Decision: lesson.DecideTrust,
		Gate: &Gate{Host: Host(session), Model: "sonnet", Effort: "high"},
	}
}

func TestATrustVerdictFromAnotherCodexThreadIsCredible(t *testing.T) {
	events := []Event{proposed("codex:author"), trusted("codex:gate", 1)}

	assert.True(t, Credible(events, digest))
}

func TestAnAuthorJudgingItsOwnLessonInCodexIsNotCredible(t *testing.T) {
	events := []Event{proposed("codex:author"), trusted("codex:author", 1)}

	assert.False(t, Credible(events, digest))
}

func TestADoubterJudgingTheLessonInCodexIsNotCredible(t *testing.T) {
	events := []Event{
		proposed("codex:author"), trusted("codex:gate", 1),
		{At: at(2), Kind: Doubt, Session: "codex:doubter", Digest: digest},
		trusted("codex:doubter", 3),
	}

	assert.False(t, Credible(events, digest))
}

func deferred(session string, minutes int) Event {
	return Event{At: at(minutes), Kind: Verdict, Session: session, Digest: digest, Decision: lesson.DecideDefer, Gate: &Gate{Host: Host(session)}}
}

func doubted(session string, minutes int) Event {
	return Event{At: at(minutes), Kind: Doubt, Session: session, Digest: digest}
}

func TestTheAuthorNeverGatesItsOwnLessonInCodex(t *testing.T) {
	laterRound := []Event{proposed("codex:author"), trusted("codex:gate", 1), doubted("codex:doubter", 2), trusted("codex:author", 3)}
	afterADefer := []Event{proposed("codex:author"), deferred("codex:gate", 1), trusted("codex:author", 2)}
	afterItsOwnDefer := []Event{proposed("codex:author"), deferred("codex:author", 1), trusted("codex:author", 2)}

	assert.False(t, Credible(laterRound, digest), "a later round still judges the author's text")
	assert.False(t, Credible(afterADefer, digest), "a defer settles nothing")
	assert.False(t, Credible(afterItsOwnDefer, digest))
}

func TestADoubterMayGateOnlyAfterAnotherTrust(t *testing.T) {
	deferredDoubt := []Event{proposed("codex:a"), trusted("codex:g1", 1), doubted("codex:x", 2), deferred("codex:g2", 3), trusted("codex:x", 4)}
	laterRound := []Event{
		proposed("codex:a"), trusted("codex:g1", 1), doubted("codex:x", 2), trusted("codex:g2", 3), doubted("codex:y", 4), trusted("codex:x", 5),
	}

	assert.False(t, Credible(deferredDoubt, digest), "a defer does not end the doubter's round")
	assert.True(t, Credible(laterRound, digest), "a trust by another gate ended the doubter's round")
}

func TestAReplacedLessonIsNotCredible(t *testing.T) {
	events := []Event{proposed("codex:a"), trusted("codex:g", 1), {At: at(2), Kind: Replaced, Session: "codex:g2", Successor: "l-0000000z"}}

	assert.False(t, Credible(events, digest), "a hand-typed trusted status must not revive a replaced lesson")
}

func TestMayJudge(t *testing.T) {
	events := []Event{proposed("codex:author"), trusted("codex:gate", 1), doubted("codex:doubter", 2)}

	assert.True(t, MayJudge(events, "codex:gate"), "the earlier gate contributed no text")
	assert.False(t, MayJudge(events, "codex:author"))
	assert.False(t, MayJudge(events, "codex:doubter"))
	assert.False(t, MayJudge(events, "manual:abc"), "a manual session is no gate")
	assert.True(t, MayJudge([]Event{proposed("claude-code:s1")}, "claude-code:s1"), "Claude Code subagents share the session id")
}

func TestClaudeCodeSessionsAreNotCompared(t *testing.T) {
	events := []Event{proposed("claude-code:s1"), trusted("claude-code:s1", 1)}

	assert.True(t, Credible(events, digest), "a Claude Code subagent shares its parent's session id; the PreToolUse hook enforces the gate instead")
}

func TestVerdictsThatAreNotCredible(t *testing.T) {
	noGate := trusted("codex:gate", 1)
	noGate.Gate = nil

	manual := trusted("manual:abc", 1)
	otherDigest := trusted("codex:gate", 1)
	otherDigest.Digest = "000000000000"
	rejected := trusted("codex:gate", 1)
	rejected.Decision = lesson.DecideReject

	cases := map[string][]Event{
		"no verdict":          {proposed("codex:author")},
		"no gate metadata":    {proposed("codex:author"), noGate},
		"manual session":      {proposed("codex:author"), manual},
		"other digest":        {proposed("codex:author"), otherDigest},
		"latest is not trust": {proposed("codex:author"), trusted("codex:gate", 1), rejected},
		"doubted after trust": {proposed("codex:author"), trusted("codex:gate", 1), {At: at(2), Kind: Doubt, Session: "codex:x", Digest: digest}},
	}
	for name, events := range cases {
		assert.False(t, Credible(events, digest), name)
	}
}

func TestConfirmationsCountDistinctNonAuthorSessionsSinceTrust(t *testing.T) {
	events := []Event{
		proposed("codex:author"),
		{At: at(1), Kind: Confirm, Session: "codex:early"},
		trusted("codex:gate", 2),
		{At: at(3), Kind: Confirm, Session: "codex:a"},
		{At: at(4), Kind: Confirm, Session: "codex:a"},
		{At: at(5), Kind: Confirm, Session: "codex:author"},
		{At: at(6), Kind: Confirm, Session: "codex:b"},
		{At: at(7), Kind: Confirm, Session: "manual:x1"},
		{At: at(8), Kind: Confirm, Session: "manual:x2"},
	}

	count, latest := Confirmations(events)

	assert.Equal(t, 2, count)
	assert.Equal(t, at(6), latest)
}

func TestRecordedAndAuthor(t *testing.T) {
	events := []Event{proposed("codex:author"), {At: at(1), Kind: Doubt, Session: "codex:d"}}

	assert.Equal(t, "codex:author", Author(events))
	assert.True(t, Recorded(events, Doubt, "codex:d"))
	assert.False(t, Recorded(events, Doubt, "codex:other"))
	assert.False(t, Recorded(events, Confirm, "codex:d"))
}

func TestDefersCountOnlyTheCurrentRun(t *testing.T) {
	deferred := func(minutes int) Event {
		return Event{At: at(minutes), Kind: Verdict, Session: "codex:g", Decision: lesson.DecideDefer}
	}

	assert.Equal(t, 0, Defers([]Event{proposed("codex:a")}))
	assert.Equal(t, 2, Defers([]Event{proposed("codex:a"), deferred(1), deferred(2)}))
	assert.Equal(t, 1, Defers([]Event{proposed("codex:a"), deferred(1), trusted("codex:g", 2), deferred(3)}))
}

func TestTrustsCountsVerdictsForOneDigest(t *testing.T) {
	other := trusted("codex:g", 3)
	other.Digest = "ffffffffffff"

	assert.Equal(t, 2, Trusts([]Event{trusted("codex:g", 1), trusted("codex:h", 2), other}, digest))
}

func TestLease(t *testing.T) {
	lease := Lease{Session: "codex:gate", Until: at(15)}

	assert.True(t, lease.HeldBy("codex:gate", at(14)))
	assert.False(t, lease.HeldBy("codex:gate", at(15)))
	assert.True(t, lease.Blocks("codex:other", at(14)))
	assert.False(t, lease.Blocks("codex:other", at(16)), "an expired lease blocks nobody")
	assert.False(t, lease.Blocks("codex:gate", at(14)))
	assert.False(t, Lease{}.Blocks("codex:other", at(1)))
}
