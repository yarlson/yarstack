// Package lesson defines a Yarbrain lesson, its limits, and its content
// identity. It does no I/O.
package lesson

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"slices"
	"strings"
	"unicode/utf8"
)

// Status is the stored lifecycle state of a lesson.
type Status string

const (
	Proposed Status = "proposed"
	Trusted  Status = "trusted"
	Deferred Status = "deferred"
	Rejected Status = "rejected"
	Doubted  Status = "doubted"
	Retired  Status = "retired"
)

// Statuses lists every stored status in lifecycle order.
var Statuses = []Status{Proposed, Trusted, Deferred, Rejected, Doubted, Retired}

// Limits on the author parts, in characters (runes), and on the whole file in
// bytes. They keep five full lessons under about 2,000 tokens.
const (
	MaxTitle     = 80
	MaxCue       = 200
	MaxBody      = 1200
	MaxObserved  = 1000
	MaxFileBytes = 16 * 1024
)

// AllScope is the scope identity that marks a cross-project lesson.
const AllScope = "all"

// Draft holds the author parts of a lesson: everything the content digest
// covers and a gate judges.
type Draft struct {
	Title    string
	Cue      string
	Body     string
	Observed string
	// Scope holds project identities such as "origin:github.com/o/r" and
	// "root:4f1c2a9e0b7d", or exactly AllScope.
	Scope    []string
	Retires  string
	Converts string
}

// Lesson is a stored lesson: its author parts plus the keys the CLI owns.
type Lesson struct {
	Draft

	ID            string
	Status        Status
	Created       string
	ProjectName   string
	TrustedDigest string
	Successor     string
}

// IsAllScope reports whether the lesson applies to every project.
func (d Draft) IsAllScope() bool {
	return len(d.Scope) == 1 && d.Scope[0] == AllScope
}

// Problem is one violated rule in a draft.
type Problem struct {
	Field  string `json:"field"`
	Rule   string `json:"rule"`
	Limit  int    `json:"limit,omitempty"`
	Actual int    `json:"actual,omitempty"`
}

// Validate reports every rule the normalized draft breaks, so a caller can
// fix all of them in one call.
func Validate(d Draft) []Problem {
	var problems []Problem

	problems = append(problems, line("title", d.Title, MaxTitle)...)
	problems = append(problems, line("cue", d.Cue, MaxCue)...)
	problems = append(problems, text("body", d.Body, MaxBody)...)
	problems = append(problems, text("observed", d.Observed, MaxObserved)...)

	if len(d.Scope) == 0 {
		problems = append(problems, Problem{Field: "scope", Rule: "required"})
	} else if slices.Contains(d.Scope, AllScope) && len(d.Scope) > 1 {
		problems = append(problems, Problem{Field: "scope", Rule: "all must be the only identity"})
	}

	if d.Retires != "" && !IsID(d.Retires) {
		problems = append(problems, Problem{Field: "retires", Rule: "must be a lesson id such as l-7k3m9q2x"})
	}

	return problems
}

func line(field, value string, limit int) []Problem {
	if value == "" {
		return []Problem{{Field: field, Rule: "required"}}
	}

	if strings.ContainsAny(value, "\r\n") {
		return []Problem{{Field: field, Rule: "must be one line"}}
	}

	return length(field, value, limit)
}

func text(field, value string, limit int) []Problem {
	if value == "" {
		return []Problem{{Field: field, Rule: "required"}}
	}

	return length(field, value, limit)
}

func length(field, value string, limit int) []Problem {
	if count := utf8.RuneCountInString(value); count > limit {
		return []Problem{{Field: field, Rule: "too long", Limit: limit, Actual: count}}
	}

	return nil
}

// Normalize returns the canonical form of a draft: LF line endings, no
// trailing spaces, no blank lines around text, and a sorted scope without
// duplicates. Validate and Digest expect a normalized draft.
func Normalize(d Draft) Draft {
	d.Title = strings.TrimSpace(d.Title)
	d.Cue = strings.TrimSpace(d.Cue)
	d.Body = normalizeText(d.Body)
	d.Observed = normalizeText(d.Observed)
	d.Retires = strings.TrimSpace(d.Retires)
	d.Converts = strings.TrimSpace(d.Converts)

	scope := make([]string, 0, len(d.Scope))
	for _, identity := range d.Scope {
		if identity = strings.TrimSpace(identity); identity != "" {
			scope = append(scope, identity)
		}
	}

	slices.Sort(scope)
	d.Scope = slices.Compact(scope)

	return d
}

func normalizeText(value string) string {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	lines := strings.Split(value, "\n")

	for i, l := range lines {
		lines[i] = strings.TrimRight(l, " \t")
	}

	return strings.Trim(strings.Join(lines, "\n"), "\n")
}

// Digest identifies the exact author parts of a normalized draft. A trust
// verdict binds to it, so any edit of the author parts invalidates trust.
func Digest(d Draft) string {
	parts := []string{d.Title, d.Cue, strings.Join(d.Scope, "\n"), d.Body, d.Observed, d.Retires, d.Converts}
	sum := sha256.Sum256([]byte(strings.Join(parts, "\x00")))

	return hex.EncodeToString(sum[:])[:12]
}

// IntentKey compares titles: two drafts with the same key in overlapping
// scopes describe the same lesson.
func IntentKey(title string) string {
	return strings.Join(strings.Fields(strings.ToLower(title)), " ")
}

// MaxSlug bounds the file-name slug.
const MaxSlug = 60

// Slug turns a title into an ASCII file name stem. It returns "" when the
// title has no ASCII letters or digits; callers then use the id.
func Slug(title string) string {
	var b strings.Builder

	dash := false

	for _, r := range strings.ToLower(title) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)

			dash = false
		case !dash && b.Len() > 0:
			b.WriteByte('-')

			dash = true
		}

		if b.Len() >= MaxSlug {
			break
		}
	}

	return strings.TrimRight(b.String(), "-")
}

const crockford = "0123456789abcdefghjkmnpqrstvwxyz"

// NewID returns a random lesson id such as "l-7k3m9q2x": 40 bits from source,
// in lowercase Crockford base32.
func NewID(source io.Reader) (string, error) {
	var raw [5]byte
	if _, err := io.ReadFull(source, raw[:]); err != nil {
		return "", fmt.Errorf("read random id: %w", err)
	}

	bits := uint64(raw[0])<<32 | uint64(raw[1])<<24 | uint64(raw[2])<<16 | uint64(raw[3])<<8 | uint64(raw[4])

	id := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		id[i] = crockford[bits&31]
		bits >>= 5
	}

	return "l-" + string(id), nil
}

// IsID reports whether value has the shape of a lesson id.
func IsID(value string) bool {
	rest, ok := strings.CutPrefix(value, "l-")
	if !ok || len(rest) != 8 {
		return false
	}

	for _, r := range rest {
		if !strings.ContainsRune(crockford, r) {
			return false
		}
	}

	return true
}
