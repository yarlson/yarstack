package vault

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
)

func sampleLesson() lesson.Lesson {
	return lesson.Lesson{
		Draft: lesson.Draft{
			Title:    "Run migrations before seeding the test database",
			Cue:      `make seed fails with "relation users does not exist"`,
			Body:     "Run `make migrate` before `make seed`.\nThe seed script assumes the schema.",
			Observed: "`make seed` failed; after `make migrate` it passed.",
			Scope:    []string{"origin:github.com/acme/shop-api", "root:4f1c2a9e0b7d"},
		},
		ID:          "l-7k3m9q2x",
		Status:      lesson.Proposed,
		Created:     "2026-10-28",
		ProjectName: "shop-api",
	}
}

func TestRenderThenParseRoundTrips(t *testing.T) {
	original := sampleLesson()

	raw, err := Render(original)
	require.NoError(t, err)

	parsed, err := Parse(raw, "lessons/x.md")
	require.NoError(t, err)
	assert.Equal(t, original, parsed)
}

func TestRenderedFileIsReadableMarkdown(t *testing.T) {
	raw, err := Render(sampleLesson())
	require.NoError(t, err)

	text := string(raw)
	assert.True(t, strings.HasPrefix(text, "---\nid: l-7k3m9q2x\nstatus: proposed\n"), text)
	assert.Contains(t, text, "\n---\n# Run migrations before seeding the test database\n\n")
	assert.Contains(t, text, "\n## Observed\n\n`make seed` failed")
}

func TestParseAcceptsEditorArtifacts(t *testing.T) {
	raw, err := Render(sampleLesson())
	require.NoError(t, err)

	edited := append([]byte("\xef\xbb\xbf"), []byte(strings.ReplaceAll(string(raw), "\n", "\r\n"))...)

	parsed, err := Parse(edited, "lessons/x.md")
	require.NoError(t, err)
	assert.Equal(t, "l-7k3m9q2x", parsed.ID)
	assert.Equal(t, sampleLesson().Title, parsed.Title)
}

func TestParseAcceptsABlockListScope(t *testing.T) {
	raw := "---\nid: l-0000000z\nscope:\n  - all\n---\n# T\n\nbody\n\n## Observed\n\nseen\n"

	parsed, err := Parse([]byte(raw), "lessons/x.md")
	require.NoError(t, err)
	assert.Equal(t, []string{"all"}, parsed.Scope)
}

func TestParseFindsAnObservedHeadingWithTrailingSpaces(t *testing.T) {
	raw := "---\nid: l-0000000z\nscope: [all]\n---\n# T\n\nbody\n\n## Observed  \n\nseen\n"

	parsed, err := Parse([]byte(raw), "lessons/x.md")
	require.NoError(t, err)
	assert.Equal(t, "body", parsed.Body)
	assert.Equal(t, "seen", parsed.Observed)
}

func TestParseRejectsBrokenFiles(t *testing.T) {
	cases := map[string]string{
		"no frontmatter":     "# Title\n",
		"unclosed":           "---\nid: l-0000000z\n# Title\n",
		"duplicate key":      "---\nid: l-0000000z\nstatus: proposed\nstatus: trusted\n---\n# T\n",
		"no title":           "---\nid: l-0000000z\n---\nbody only\n",
		"list as a scalar":   "---\nid: [a, b]\n---\n# T\n",
		"frontmatter scalar": "---\njust text\n---\n# T\n",
	}
	for name, raw := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := Parse([]byte(raw), "lessons/x.md")

			var bad *BadFileError
			assert.ErrorAs(t, err, &bad)
		})
	}
}

const handEdited = `---
id: l-7k3m9q2x
status: trusted   # a person's note
cue: 'make seed fails'
scope: [origin:github.com/acme/shop-api]
aliases:
  - seeding
# a comment the CLI must keep
trusted_digest: 2c9d41e07a6b
---
# Title

Body with **markdown**.

## Observed

seen
`

func TestSetOwnedKeepsEveryOtherByte(t *testing.T) {
	out, err := SetOwned([]byte(handEdited), "lessons/x.md", map[string]string{"status": "doubted", "trusted_digest": ""})
	require.NoError(t, err)

	want := strings.Replace(handEdited, "status: trusted   # a person's note\n", "status: doubted\n", 1)
	want = strings.Replace(want, "trusted_digest: 2c9d41e07a6b\n", "", 1)
	assert.Equal(t, want, string(out))
}

func TestSetOwnedAddsAMissingKeyBeforeTheClosingLine(t *testing.T) {
	out, err := SetOwned([]byte(handEdited), "lessons/x.md", map[string]string{"successor": "l-0000000z"})
	require.NoError(t, err)

	want := strings.Replace(handEdited, "trusted_digest: 2c9d41e07a6b\n---\n", "trusted_digest: 2c9d41e07a6b\nsuccessor: l-0000000z\n---\n", 1)
	assert.Equal(t, want, string(out))
}

func TestSetOwnedRefusesUnsafeFiles(t *testing.T) {
	cases := map[string]string{
		"byte-order mark": "\xef\xbb\xbf" + handEdited,
		"CRLF":            strings.ReplaceAll(handEdited, "\n", "\r\n"),
		"multi-line value": strings.Replace(handEdited, "status: trusted   # a person's note\n",
			"status: >-\n  trusted\n", 1),
		"plain continuation": strings.Replace(handEdited, "status: trusted   # a person's note\n",
			"status: trusted\n  still status\n", 1),
		"duplicate key": strings.Replace(handEdited, "cue:", "status: proposed\ncue:", 1),
		"flow mapping":  "---\n{id: l-7k3m9q2x, status: trusted, cue: c, scope: [all]}\n---\n# Title\n\nBody.\n",
		"line separator in an earlier value": strings.Replace(handEdited, "status: trusted   # a person's note\ncue: 'make seed fails'",
			"cue: 'make seed fails'\nstatus: trusted", 1),
	}
	for name, raw := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := SetOwned([]byte(raw), "lessons/x.md", map[string]string{"status": "doubted"})

			var bad *BadFileError
			assert.ErrorAs(t, err, &bad)
		})
	}
}

func TestSetOwnedQuotesValuesThatNeedIt(t *testing.T) {
	out, err := SetOwned([]byte(handEdited), "lessons/x.md", map[string]string{"project_name": "yes"})
	require.NoError(t, err)

	parsed, err := Parse(out, "lessons/x.md")
	require.NoError(t, err)
	assert.Equal(t, "yes", parsed.ProjectName, "a value YAML would read as a boolean must stay a string")
}
