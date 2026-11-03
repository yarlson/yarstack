package lesson

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validDraft() Draft {
	return Draft{
		Title:    "Run migrations before seeding the test database",
		Cue:      `make seed fails with "relation users does not exist"`,
		Body:     "Run `make migrate` before `make seed`.",
		Observed: "`make seed` failed; after `make migrate` it passed.",
		Scope:    []string{"origin:github.com/acme/shop-api"},
	}
}

func TestValidDraftHasNoProblems(t *testing.T) {
	assert.Empty(t, Validate(validDraft()))
}

func TestValidateReportsEveryProblemAtOnce(t *testing.T) {
	d := Draft{Title: "two\nlines", Cue: strings.Repeat("x", 201)}

	assert.ElementsMatch(t, []Problem{
		{Field: "title", Rule: "must be one line"},
		{Field: "cue", Rule: "too long", Limit: 200, Actual: 201},
		{Field: "body", Rule: "required"},
		{Field: "observed", Rule: "required"},
		{Field: "scope", Rule: "required"},
	}, Validate(d))
}

func TestLimitsCountCharactersNotBytes(t *testing.T) {
	d := validDraft()
	d.Title = strings.Repeat("ё", 80)

	assert.Empty(t, Validate(d), "80 two-byte characters are within the 80-character limit")
}

func TestAllScopeMustStandAlone(t *testing.T) {
	d := validDraft()
	d.Scope = []string{"all", "origin:github.com/acme/shop-api"}

	assert.Contains(t, Validate(d), Problem{Field: "scope", Rule: "all must be the only identity"})
}

func TestRetiresMustBeAnID(t *testing.T) {
	d := validDraft()
	d.Retires = "the old lesson"

	assert.Contains(t, Validate(d), Problem{Field: "retires", Rule: "must be a lesson id such as l-7k3m9q2x"})
}

func TestNormalizeCanonicalizesTextAndScope(t *testing.T) {
	d := Normalize(Draft{
		Title:    "  Title  ",
		Body:     "\n\nline one   \r\nline two\t\n\n",
		Observed: "seen",
		Scope:    []string{"root:b", " origin:a ", "root:b", ""},
	})

	assert.Equal(t, "Title", d.Title)
	assert.Equal(t, "line one\nline two", d.Body)
	assert.Equal(t, []string{"origin:a", "root:b"}, d.Scope)
}

func TestDigestIgnoresFormattingButNotContent(t *testing.T) {
	base := Normalize(validDraft())

	reformatted := validDraft()
	reformatted.Body += "   \r\n\r\n"
	reformatted.Scope = []string{" origin:github.com/acme/shop-api "}

	edited := validDraft()
	edited.Body = "Run `make migrate` first."

	assert.Equal(t, Digest(base), Digest(Normalize(reformatted)))
	assert.NotEqual(t, Digest(base), Digest(Normalize(edited)))
	assert.Len(t, Digest(base), 12)
}

func TestDigestSeparatesFields(t *testing.T) {
	a := Draft{Title: "ab", Cue: "c"}
	b := Draft{Title: "a", Cue: "bc"}

	assert.NotEqual(t, Digest(a), Digest(b), "moving text between fields must change the digest")
}

func TestIntentKeyIgnoresCaseAndSpacing(t *testing.T) {
	assert.Equal(t, IntentKey("Run  Migrations first"), IntentKey(" run migrations FIRST "))
	assert.NotEqual(t, IntentKey("Run migrations first"), IntentKey("Run migrations last"))
}

func TestSlug(t *testing.T) {
	cases := map[string]string{
		"Run migrations before seeding the test database": "run-migrations-before-seeding-the-test-database",
		"`make seed` fails: relation users!":              "make-seed-fails-relation-users",
		"Всё на русском":                                  "",
		strings.Repeat("abc ", 30):                        "abc-abc-abc-abc-abc-abc-abc-abc-abc-abc-abc-abc-abc-abc-abc",
	}
	for title, want := range cases {
		assert.Equal(t, want, Slug(title), "title %q", title)
	}
}

func TestNewIDIsCrockfordAndRecognized(t *testing.T) {
	id, err := NewID(bytes.NewReader([]byte{0, 0, 0, 0, 31}))
	require.NoError(t, err)

	assert.Equal(t, "l-0000000z", id)
	assert.True(t, IsID(id))
}

func TestNewIDFailsWithoutRandomness(t *testing.T) {
	_, err := NewID(bytes.NewReader([]byte{1, 2}))

	assert.Error(t, err)
}

func TestIsIDRejectsOtherShapes(t *testing.T) {
	for _, value := range []string{"", "l-", "l-0000000", "l-0000000u", "x-00000000", "l-ABCDEFGH"} {
		assert.False(t, IsID(value), "value %q", value)
	}
}
