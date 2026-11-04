package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/trust"
)

func TestLedgerKeepsEventsInOrder(t *testing.T) {
	v := newVault(t)
	first := trust.Event{At: time.Date(2026, 11, 3, 10, 0, 0, 0, time.UTC), Kind: trust.Proposed, Session: "codex:a", Digest: "d1"}
	second := trust.Event{
		At: first.At.Add(time.Minute), Kind: trust.Verdict, Session: "codex:g", Digest: "d1",
		Decision: lesson.DecideTrust, Gate: &trust.Gate{Host: "codex", Model: "m"},
		Checks: map[string]trust.Check{"safety": {Pass: true, Reason: "no secrets"}},
	}

	require.NoError(t, v.AppendLedger("l-0000000a", first))
	require.NoError(t, v.AppendLedger("l-0000000a", second))

	events, err := v.Ledger("l-0000000a")
	require.NoError(t, err)
	assert.Equal(t, []trust.Event{first, second}, events)
}

func TestLedgerOfALessonWithoutEventsIsEmpty(t *testing.T) {
	events, err := newVault(t).Ledger("l-0000000a")

	require.NoError(t, err)
	assert.Empty(t, events)
}

func TestACorruptLedgerFailsClosed(t *testing.T) {
	v := newVault(t)
	require.NoError(t, os.WriteFile(filepath.Join(v.Path, ".yarbrain", "ledger", "l-0000000a.jsonl"), []byte("{\"event\":\n"), 0o644))

	_, err := v.Ledger("l-0000000a")

	var bad *BadFileError
	assert.ErrorAs(t, err, &bad)
}

func TestLedgerRejectsIDsThatAreNotLessonIDs(t *testing.T) {
	v := newVault(t)

	require.Error(t, v.AppendLedger("../escape", trust.Event{}))

	_, err := v.Ledger("../../etc/passwd")
	assert.Error(t, err)
}

func TestLeaseRoundTrip(t *testing.T) {
	v := newVault(t)
	lease := trust.Lease{Session: "codex:g", Until: time.Date(2026, 11, 3, 10, 15, 0, 0, time.UTC)}

	empty, err := v.Lease("l-0000000a")
	require.NoError(t, err)
	assert.Equal(t, trust.Lease{}, empty)

	require.NoError(t, v.SetLease("l-0000000a", lease))

	read, err := v.Lease("l-0000000a")
	require.NoError(t, err)
	assert.Equal(t, lease, read)

	require.NoError(t, v.ClearLease("l-0000000a"))
	require.NoError(t, v.ClearLease("l-0000000a"), "clearing a missing lease is not an error")

	cleared, err := v.Lease("l-0000000a")
	require.NoError(t, err)
	assert.Equal(t, trust.Lease{}, cleared)
}

func TestSetLessonKeysRefusesAFileChangedSinceItWasRead(t *testing.T) {
	v := newVault(t)
	writeLesson(t, v, "a.md", lessonText("l-0000000a"))

	entries, _, err := v.Lessons()
	require.NoError(t, err)

	writeLesson(t, v, "a.md", lessonText("l-0000000a")+"\nedited\n")

	err = v.SetLessonKeys(entries[0], map[string]string{"status": "trusted"})

	var changed *ChangedError
	require.ErrorAs(t, err, &changed)
}

func TestSetLessonKeysRewritesOwnedKeys(t *testing.T) {
	v := newVault(t)
	writeLesson(t, v, "a.md", lessonText("l-0000000a"))

	entries, _, err := v.Lessons()
	require.NoError(t, err)

	require.NoError(t, v.SetLessonKeys(entries[0], map[string]string{"status": "trusted", "trusted_digest": "abc"}))

	updated, _, err := v.Lessons()
	require.NoError(t, err)
	assert.Equal(t, lesson.Trusted, updated[0].Lesson.Status)
	assert.Equal(t, "abc", updated[0].Lesson.TrustedDigest)
}
