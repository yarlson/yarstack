package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/trust"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/vault"
)

// session identifies the host session that runs this command. A command run
// outside a host gets a fresh manual id, so it never matches a gate session.
func (a *app) session() string {
	if id := a.getenv("CODEX_THREAD_ID"); id != "" {
		return trust.Codex + ":" + id
	}

	if id := a.getenv("CLAUDE_CODE_SESSION_ID"); id != "" {
		return trust.ClaudeCode + ":" + id
	}

	return trust.Manual + ":" + strings.ToLower(rand.Text()[:12])
}

// record is a lesson with its ledger, as read under the vault lock.
type record struct {
	entry  vault.Entry
	events []trust.Event
}

func (r record) lesson() lesson.Lesson { return r.entry.Lesson }

func (r record) digest() string { return lesson.Digest(lesson.Normalize(r.entry.Lesson.Draft)) }

func (r record) effective() lesson.Status { return lesson.Effective(r.entry.Lesson) }

func (r record) served() bool { return trust.Served(r.entry.Lesson, r.events) }

func loadRecord(v *vault.Vault, entries []vault.Entry, id string) (record, error) {
	entry, ok := vault.Find(entries, id)
	if !ok {
		return record{}, newError("not_found", "no lesson has id "+id, "Use `yarbrain find <terms>` to get lesson ids.").
			with("id", id)
	}

	events, err := v.Ledger(id)
	if err != nil {
		return record{}, err
	}

	return record{entry: entry, events: events}, nil
}

// decodeStrict reads one JSON object and rejects unknown keys and trailing
// data, so a misspelled field is an error instead of a silent default.
func decodeStrict(raw []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(target); err != nil {
		return newError("bad_input", "the input is not the expected JSON: "+err.Error(),
			"Run the command with --help to see the input format.").with("field", "input")
	}

	if decoder.More() {
		return newError("bad_input", "the input holds more than one JSON value",
			"Pass exactly one JSON object.").with("field", "input")
	}

	return nil
}

// oneLine checks an agent-written line such as a reason or a note.
func oneLine(field, value string, limit int) []lesson.Problem {
	switch {
	case strings.TrimSpace(value) == "":
		return []lesson.Problem{{Field: field, Rule: "required"}}
	case strings.ContainsAny(value, "\r\n"):
		return []lesson.Problem{{Field: field, Rule: "must be one line"}}
	case utf8.RuneCountInString(value) > limit:
		return []lesson.Problem{{Field: field, Rule: "too long", Limit: limit, Actual: utf8.RuneCountInString(value)}}
	}

	return nil
}

// gateStep names the step that makes a lesson served when it awaits the
// gate, and is empty otherwise. It points to the skill, which owns the
// host-specific procedure.
func gateStep(id string, status lesson.Status) string {
	if !lesson.AwaitsGate(status) {
		return ""
	}

	return "Start the independent gate for " + id + " now, in the background, as the lesson-record skill describes, without asking: " +
		"gating is part of the plugin's work. Never run brief or verdict yourself."
}

func badState(r record, message, fix string) *cliError {
	return newError("bad_state", fmt.Sprintf("lesson %s is %s: %s", r.lesson().ID, r.effective(), message), fix).
		with("id", r.lesson().ID).with("status", string(r.effective()))
}
