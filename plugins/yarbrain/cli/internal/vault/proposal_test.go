package vault

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
)

const proposalText = "---\ncue: 'make seed fails'\nscope: here\n---\n# Run migrations first\n\nRun `make migrate`.\n\n## Observed\n\nIt passed after.\n"

func problemsOf(t *testing.T, err error) []lesson.Problem {
	t.Helper()

	var input *InputError
	require.ErrorAs(t, err, &input)

	return input.Problems
}

func TestParseProposalReadsTheDraftAndScopeKeyword(t *testing.T) {
	draft, scope, err := ParseProposal([]byte(proposalText))

	require.NoError(t, err)
	assert.Equal(t, "here", scope)
	assert.Equal(t, "Run migrations first", draft.Title)
	assert.Equal(t, "make seed fails", draft.Cue)
	assert.Equal(t, "Run `make migrate`.", draft.Body)
	assert.Equal(t, "It passed after.", draft.Observed)
}

func TestParseProposalAcceptsAllAndEditorArtifacts(t *testing.T) {
	raw := "\xef\xbb\xbf" + strings.ReplaceAll(strings.Replace(proposalText, "scope: here", "scope: [all]", 1), "\n", "\r\n")

	_, scope, err := ParseProposal([]byte(raw))

	require.NoError(t, err)
	assert.Equal(t, "all", scope)
}

func TestParseProposalNamesTheProblemOfEveryBadKey(t *testing.T) {
	raw := strings.Replace(proposalText, "scope: here", "scope: everywhere\nid: l-00000000\nweight: 3", 1)

	_, _, err := ParseProposal([]byte(raw))

	assert.ElementsMatch(t, []lesson.Problem{
		{Field: "id", Rule: "set by the CLI; remove this key"},
		{Field: "weight", Rule: "unknown key; allowed keys are cue, scope, retires, converts"},
		{Field: "scope", Rule: "must be here or all"},
	}, problemsOf(t, err))
}

func TestParseProposalReportsAMalformedDocument(t *testing.T) {
	_, _, err := ParseProposal([]byte("# A title without frontmatter\n"))

	assert.Equal(t, []lesson.Problem{{Field: "input", Rule: "the file must start with a --- frontmatter line"}}, problemsOf(t, err))
}

func TestParseProposalReportsABodyWithoutTitle(t *testing.T) {
	_, _, err := ParseProposal([]byte("---\nscope: here\n---\nno title here\n"))

	assert.Equal(t, []lesson.Problem{{Field: "input", Rule: "the body must start with a '# Title' line"}}, problemsOf(t, err))
}

func storedLesson(id, title string) lesson.Lesson {
	return lesson.Lesson{
		Draft: lesson.Draft{Title: title, Cue: "c", Body: "b", Observed: "o", Scope: []string{"all"}},
		ID:    id, Status: lesson.Proposed, Created: "2026-11-03",
	}
}

func TestCreateLessonNamesTheFileAfterTheTitle(t *testing.T) {
	v := newVault(t)

	path, err := v.CreateLesson(storedLesson("l-0000000a", "Run migrations first"))

	require.NoError(t, err)
	assert.Equal(t, "lessons/run-migrations-first.md", path)
	assert.FileExists(t, filepath.Join(v.Path, path))
}

func TestCreateLessonAddsTheIDWhenTheNameIsTaken(t *testing.T) {
	v := newVault(t)
	require.NoError(t, os.WriteFile(filepath.Join(v.Path, "lessons", "run-migrations-first.md"), []byte("a person's file"), 0o644))

	path, err := v.CreateLesson(storedLesson("l-0000000b", "Run migrations first!"))

	require.NoError(t, err)
	assert.Equal(t, "lessons/run-migrations-first-0000000b.md", path)

	raw, err := os.ReadFile(filepath.Join(v.Path, "lessons", "run-migrations-first.md"))
	require.NoError(t, err)
	assert.Equal(t, "a person's file", string(raw), "an existing file must never be overwritten")
}

func TestCreateLessonUsesTheIDForATitleWithoutASCII(t *testing.T) {
	v := newVault(t)

	path, err := v.CreateLesson(storedLesson("l-0000000c", "Сначала миграции"))

	require.NoError(t, err)
	assert.Equal(t, "lessons/l-0000000c.md", path)
}

func TestFindLocatesAnEntryByID(t *testing.T) {
	entries := []Entry{{Path: "a", Lesson: storedLesson("l-0000000a", "A")}, {Path: "b", Lesson: storedLesson("l-0000000b", "B")}}

	found, ok := Find(entries, "l-0000000b")
	assert.True(t, ok)
	assert.Equal(t, "b", found.Path)

	_, ok = Find(entries, "l-0000000z")
	assert.False(t, ok)
}
