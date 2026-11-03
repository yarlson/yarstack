package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLegacyNotesListsNotesAndVaultSkills(t *testing.T) {
	v := newVault(t)
	files := []string{
		"notes/projects/shop.md", "notes/concepts/seeding.md", "skills/release/SKILL.md",
		"episodes/2026/01/e.md", "inbox/memory/c.md", "notes/projects/readme.txt",
	}

	for _, name := range files {
		path := filepath.Join(v.Path, name)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte("# x\n"), 0o644))
	}

	notes, err := v.LegacyNotes()

	require.NoError(t, err)
	assert.Equal(t, []string{"notes/concepts/seeding.md", "notes/projects/shop.md", "skills/release/SKILL.md"}, notes)
}

func TestLegacyNotesOfAVaultWithoutV02ContentIsEmpty(t *testing.T) {
	notes, err := newVault(t).LegacyNotes()

	require.NoError(t, err)
	assert.Empty(t, notes)
}

func TestDismissalsAreRecordedAndRead(t *testing.T) {
	v := newVault(t)

	none, err := v.Dismissed()
	require.NoError(t, err)
	assert.Empty(t, none)

	at := time.Date(2026, 11, 3, 10, 0, 0, 0, time.UTC)
	require.NoError(t, v.Dismiss(Dismissal{At: at, Path: "notes/a.md", Session: "codex:x", Reason: "nothing to keep"}))
	require.NoError(t, v.Dismiss(Dismissal{At: at, Path: "notes/b.md", Session: "codex:x", Reason: "obsolete"}))

	dismissed, err := v.Dismissed()
	require.NoError(t, err)
	assert.Equal(t, map[string]bool{"notes/a.md": true, "notes/b.md": true}, dismissed)
}

func TestACorruptDismissalFileFailsClosed(t *testing.T) {
	v := newVault(t)
	require.NoError(t, os.WriteFile(filepath.Join(v.Path, ".yarbrain", "legacy.jsonl"), []byte("{not json\n"), 0o644))

	_, err := v.Dismissed()

	var bad *BadFileError
	assert.ErrorAs(t, err, &bad)
}

func TestIsLegacyPath(t *testing.T) {
	assert.True(t, IsLegacyPath("notes/projects/shop.md"))
	assert.True(t, IsLegacyPath("skills/release/SKILL.md"))
	assert.False(t, IsLegacyPath("episodes/2026/01/e.md"))
	assert.False(t, IsLegacyPath("notes/../lessons/x.md"))
	assert.False(t, IsLegacyPath("../notes/x.md"))
	assert.False(t, IsLegacyPath("notes/x.txt"))
	assert.False(t, IsLegacyPath("lessons/x.md"))
}
