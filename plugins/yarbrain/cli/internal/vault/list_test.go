package vault

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeLesson(t *testing.T, v *Vault, name, raw string) {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(v.Path, "lessons", name), []byte(raw), 0o644))
}

func lessonText(id string) string {
	return "---\nid: " + id + "\nstatus: proposed\nscope: [all]\n---\n# Title " + id + "\n\nbody\n\n## Observed\n\nseen\n"
}

func codesByPath(findings []Finding) map[string]string {
	codes := map[string]string{}
	for _, finding := range findings {
		codes[finding.Path] = finding.Code
	}

	return codes
}

func TestLessonsServesOnlyUsableFiles(t *testing.T) {
	v := newVault(t)

	writeLesson(t, v, "good.md", lessonText("l-0000000a"))
	writeLesson(t, v, "dropped-in.md", "---\nstatus: proposed\n---\n# Mine\n")
	writeLesson(t, v, "broken.md", "no frontmatter\n")
	writeLesson(t, v, "huge.md", lessonText("l-0000000b")+strings.Repeat("x", 16*1024))
	writeLesson(t, v, "twin-1.md", lessonText("l-0000000c"))
	writeLesson(t, v, "twin-2.md", lessonText("l-0000000c"))
	writeLesson(t, v, ".tmp-123", "ignored")
	writeLesson(t, v, "notes.txt", "ignored")

	entries, findings, err := v.Lessons()
	require.NoError(t, err)

	require.Len(t, entries, 1)
	assert.Equal(t, "lessons/good.md", entries[0].Path)
	assert.Equal(t, "l-0000000a", entries[0].Lesson.ID)

	assert.Equal(t, map[string]string{
		"lessons/dropped-in.md": "not_adopted",
		"lessons/broken.md":     "bad_file",
		"lessons/huge.md":       "file_too_large",
		"lessons/twin-1.md":     "duplicate_id",
		"lessons/twin-2.md":     "duplicate_id",
	}, codesByPath(findings))
}

func TestLessonsOfAnUnadoptedVaultIsEmpty(t *testing.T) {
	v := openVault(t, t.TempDir())

	entries, findings, err := v.Lessons()
	require.NoError(t, err)
	assert.Empty(t, entries)
	assert.Empty(t, findings)
}

func TestEntryRevisionMatchesTheFileBytes(t *testing.T) {
	v := newVault(t)
	writeLesson(t, v, "good.md", lessonText("l-0000000a"))

	entries, _, err := v.Lessons()
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, Rev([]byte(lessonText("l-0000000a"))), entries[0].Rev)
}
