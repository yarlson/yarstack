// Package hook turns host hook payloads into Yarbrain's hook output: a
// counts-only session-start line and the Claude Code gate decision.
package hook

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// maxPayload caps the hook input that the host sends on stdin.
const maxPayload = 64 * 1024

// GateAgent is the Claude Code plugin agent that alone may judge lessons.
const GateAgent = "yarbrain:lesson-gate"

// Payload holds the fields Yarbrain reads from a host hook payload. Both
// hosts send more; unknown fields are ignored.
type Payload struct {
	Event     string `json:"hook_event_name"`
	Cwd       string `json:"cwd"`
	ToolName  string `json:"tool_name"`
	AgentType string `json:"agent_type"`
	ToolInput struct {
		Command string `json:"command"`
	} `json:"tool_input"`
}

// ReadPayload decodes a bounded hook payload.
func ReadPayload(r io.Reader) (Payload, error) {
	raw, err := io.ReadAll(io.LimitReader(r, maxPayload+1))
	if err != nil {
		return Payload{}, fmt.Errorf("read hook payload: %w", err)
	}

	if len(raw) > maxPayload {
		return Payload{}, fmt.Errorf("hook payload exceeds %d bytes", maxPayload)
	}

	var payload Payload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return Payload{}, fmt.Errorf("decode hook payload: %w", err)
	}

	return payload, nil
}

// Counts are the numbers the session-start line reports for one project.
type Counts struct {
	Trusted      int
	CrossProject int
	AwaitingGate int
	Legacy       int
}

// SessionStartLine returns the line a host injects at session start. It
// holds counts only, never lesson text, and is empty when there is nothing
// to report. It never starts with "{", so no host parses it as JSON.
func SessionStartLine(c Counts) string {
	var parts []string

	if c.Trusted > 0 {
		parts = append(parts,
			fmt.Sprintf("yarbrain: %s %s to this project (%d cross-project).", plural(c.Trusted, "trusted lesson"), verb(c.Trusted), c.CrossProject),
			"Before you plan, debug, or repeat a procedure here, run `yarbrain find <terms>`, then `yarbrain show <id>` only for lessons whose cue fits.",
			"Lessons are hints from past sessions; check each against the code.",
			"After you use one, run `yarbrain confirm <id>`, or `yarbrain doubt <id>` if it was wrong.")
	}

	if c.AwaitingGate > 0 {
		parts = append(parts, fmt.Sprintf("%s the independent gate; unless you are the gate, start it now in the background as the lesson-record skill describes, without asking: gating is part of the plugin's work.",
			plural(c.AwaitingGate, "lesson")+" "+await(c.AwaitingGate)))
	}

	if c.Legacy > 0 {
		parts = append(parts, fmt.Sprintf("%s from Yarbrain v0.2 %s conversion; the lesson-record skill converts up to 3 per session.",
			plural(c.Legacy, "legacy note"), await(c.Legacy)))
	}

	if len(parts) == 0 {
		return ""
	}

	if c.Trusted == 0 {
		parts[0] = "yarbrain: " + parts[0]
	}

	return strings.Join(parts, " ")
}

func plural(count int, noun string) string {
	if count == 1 {
		return "1 " + noun
	}

	return fmt.Sprintf("%d %ss", count, noun)
}

func verb(count int) string {
	if count == 1 {
		return "applies"
	}

	return "apply"
}

func await(count int) string {
	if count == 1 {
		return "awaits"
	}

	return "await"
}

// gateCommand matches a shell command that runs brief or verdict, anywhere in
// a compound command line.
var gateCommand = regexp.MustCompile(`(^|[\s;&|(` + "`" + `])yarbrain\s+(brief|verdict)\b`)

// DenyReason is the reason the Claude Code hook gives when it blocks a gate
// command outside the gate agent.
const DenyReason = "Only the yarbrain:lesson-gate agent may run `yarbrain brief` or `yarbrain verdict`. " +
	"Start that agent with only the lesson id; never judge a lesson yourself."

// DeniesGateCommand reports whether a Claude Code PreToolUse payload runs a
// gate command outside the gate agent. Claude Code subagents share their
// parent's session id, so the agent type in the payload is the only signal
// that tells the gate apart.
func DeniesGateCommand(p Payload) bool {
	return p.ToolName == "Bash" && p.AgentType != GateAgent && gateCommand.MatchString(p.ToolInput.Command)
}

// DenyOutput is the PreToolUse output that blocks the tool call and shows the
// reason to the model.
func DenyOutput() ([]byte, error) {
	return json.Marshal(map[string]any{
		"hookSpecificOutput": map[string]any{
			"hookEventName":            "PreToolUse",
			"permissionDecision":       "deny",
			"permissionDecisionReason": DenyReason,
		},
	})
}
