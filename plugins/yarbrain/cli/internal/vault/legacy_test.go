package vault

import (
	"os"
	"path/filepath"
	"testing"

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

func TestIsLegacyPath(t *testing.T) {
	assert.True(t, IsLegacyPath("notes/projects/shop.md"))
	assert.True(t, IsLegacyPath("skills/release/SKILL.md"))
	assert.False(t, IsLegacyPath("episodes/2026/01/e.md"))
	assert.False(t, IsLegacyPath("notes/../lessons/x.md"))
	assert.False(t, IsLegacyPath("../notes/x.md"))
	assert.False(t, IsLegacyPath("notes/x.txt"))
	assert.False(t, IsLegacyPath("lessons/x.md"))
}
