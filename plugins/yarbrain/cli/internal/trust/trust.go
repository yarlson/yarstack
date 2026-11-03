// Package trust decides, from a lesson's ledger, whether its trust is
// credible and which feedback a session may still give. It does no I/O.
package trust

import (
	"slices"
	"strings"
	"time"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
)

// Kind names a ledger event.
type Kind string

const (
	Proposed Kind = "proposed"
	Verdict  Kind = "verdict"
	Doubt    Kind = "doubt"
	Confirm  Kind = "confirm"
	// Replaced records that a trusted successor retired the lesson.
	Replaced Kind = "replaced"
)

// Check is one gate check: pass or fail, with one line of reason.
type Check struct {
	Pass   bool   `json:"pass"`
	Reason string `json:"reason"`
}

// RequiredChecks are the checks every verdict carries.
var RequiredChecks = []string{"safety", "correct", "valuable", "scope"}

// Gate describes the agent that gave a verdict. Model is the gate's own
// claim; Host comes from the environment.
type Gate struct {
	Host   string `json:"host"`
	Model  string `json:"model,omitempty"`
	Effort string `json:"effort,omitempty"`
}

// Event is one line of a lesson's ledger. Notes and reasons are written by
// agents and never become served lesson text.
type Event struct {
	At          time.Time        `json:"at"`
	Kind        Kind             `json:"event"`
	Session     string           `json:"session"`
	Digest      string           `json:"digest,omitempty"`
	Decision    lesson.Decision  `json:"decision,omitempty"`
	Checks      map[string]Check `json:"checks,omitempty"`
	Reason      string           `json:"reason,omitempty"`
	Note        string           `json:"note,omitempty"`
	Gate        *Gate            `json:"gate,omitempty"`
	DuplicateOf string           `json:"duplicate_of,omitempty"`
	Successor   string           `json:"successor,omitempty"`
}

// Session prefixes identify the host that ran a command.
const (
	ClaudeCode = "claude-code"
	Codex      = "codex"
	Manual     = "manual"
)

// Host returns the host part of a session reference such as "codex:019a".
func Host(session string) string {
	host, _, _ := strings.Cut(session, ":")

	return host
}

// distinctSessions reports whether a host gives each agent its own session
// id. Codex gives each thread its own id. A Claude Code subagent shares its
// parent's id, so there the plugin's PreToolUse hook enforces that only the
// gate agent runs verdict, and the session comparison cannot apply.
func distinctSessions(host string) bool { return host == Codex }

// Credible reports whether the lesson's latest verdict is a trust verdict
// that a later session may rely on: it binds digest, carries gate metadata,
// comes from a host session, and, where the host can tell sessions apart,
// comes from a session that did not contribute to the lesson since the
// previous verdict.
func Credible(events []Event, digest string) bool {
	last := lastIndex(events, Verdict)
	if last < 0 {
		return false
	}

	verdict := events[last]
	if verdict.Decision != lesson.DecideTrust || verdict.Digest != digest {
		return false
	}

	if verdict.Gate == nil || verdict.Gate.Host == "" {
		return false
	}

	host := Host(verdict.Session)
	if host != ClaudeCode && host != Codex {
		return false
	}

	if distinctSessions(host) && slices.Contains(contributors(events, last), verdict.Session) {
		return false
	}

	return !slices.ContainsFunc(events[last+1:], func(e Event) bool { return e.Kind == Doubt })
}

// Served reports whether a lesson may reach an agent: it is trusted, its
// content still matches the trusted digest, and that trust is credible.
func Served(l lesson.Lesson, events []Event) bool {
	return lesson.Effective(l) == lesson.Trusted && Credible(events, lesson.Digest(lesson.Normalize(l.Draft)))
}

// contributors lists the sessions that proposed or doubted the lesson after
// the verdict before the one at index.
func contributors(events []Event, index int) []string {
	start := lastIndex(events[:index], Verdict) + 1

	var sessions []string

	for _, e := range events[start:index] {
		if e.Kind == Proposed || e.Kind == Doubt {
			sessions = append(sessions, e.Session)
		}
	}

	return sessions
}

// Contributors lists the sessions that proposed or doubted the lesson since
// its latest verdict: a gate from any of them would judge its own work.
func Contributors(events []Event) []string {
	return contributors(events, len(events))
}

func lastIndex(events []Event, kind Kind) int {
	for i := len(events) - 1; i >= 0; i-- {
		if events[i].Kind == kind {
			return i
		}
	}

	return -1
}

// Author returns the session that proposed the lesson.
func Author(events []Event) string {
	for _, e := range events {
		if e.Kind == Proposed {
			return e.Session
		}
	}

	return ""
}

// Recorded reports whether session already gave feedback of kind on this
// lesson. Each session doubts and confirms a lesson at most once.
func Recorded(events []Event, kind Kind, session string) bool {
	return slices.ContainsFunc(events, func(e Event) bool { return e.Kind == kind && e.Session == session })
}

// Confirmations counts the distinct sessions, other than the author, that
// confirmed the lesson since its latest trust verdict, and returns the time
// of the latest confirmation.
func Confirmations(events []Event) (int, time.Time) {
	start := lastIndex(events, Verdict) + 1
	author := Author(events)
	seen := map[string]bool{}

	var latest time.Time

	for _, e := range events[start:] {
		if e.Kind != Confirm || e.Session == author || seen[e.Session] {
			continue
		}

		seen[e.Session] = true

		if e.At.After(latest) {
			latest = e.At
		}
	}

	return len(seen), latest
}

// MaxDefers is how many defer verdicts in a row a lesson may get. After that
// the gate must decide.
const MaxDefers = 3

// Defers counts the defer verdicts since the last verdict of another kind.
func Defers(events []Event) int {
	count := 0

	for i := len(events) - 1; i >= 0; i-- {
		if events[i].Kind != Verdict {
			continue
		}

		if events[i].Decision != lesson.DecideDefer {
			break
		}

		count++
	}

	return count
}

// MaxTrusts is how many trust verdicts one digest may get. A lesson that
// keeps being doubted and trusted again is contested: the next doubt after
// that retires it.
const MaxTrusts = 3

// Trusts counts the trust verdicts that bound digest.
func Trusts(events []Event, digest string) int {
	count := 0

	for _, e := range events {
		if e.Kind == Verdict && e.Decision == lesson.DecideTrust && e.Digest == digest {
			count++
		}
	}

	return count
}

// LeaseDuration is how long a gate holds a lesson after brief.
const LeaseDuration = 15 * time.Minute

// Lease reserves a lesson for one gate session.
type Lease struct {
	Session string    `json:"session"`
	Until   time.Time `json:"until"`
}

// Blocks reports whether the lease keeps session from judging the lesson.
func (l Lease) Blocks(session string, now time.Time) bool {
	return l.Session != "" && l.Session != session && now.Before(l.Until)
}

// HeldBy reports whether session holds the lease right now.
func (l Lease) HeldBy(session string, now time.Time) bool {
	return l.Session == session && now.Before(l.Until)
}
