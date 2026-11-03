// Package hook turns host hook payloads into Yarbrain's hook output: a
// counts-only session-start line and the Claude Code gate decision.
package hook

import (
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"
	"unicode"
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

// commandWrappers run the command that follows them.
var commandWrappers = []string{"env", "command", "exec", "sudo", "nohup", "time"}

// runsGateCommand reports whether a Bash or PowerShell command line calls
// yarbrain brief or verdict in any of its simple commands.
func runsGateCommand(line string) bool {
	return slices.ContainsFunc(simpleCommands(line), isGateCall)
}

// simpleCommands splits a command line into simple commands and their words.
// It honors single and double quotes and ends a simple command at ; & | ( )
// a newline, or a backtick. A backslash stays literal, because PowerShell
// and Windows paths use it. It recognizes gate calls; it is no shell parser.
func simpleCommands(line string) [][]string {
	var (
		commands [][]string
		words    []string
		word     strings.Builder
		inWord   bool
		quote    rune
	)

	endWord := func() {
		if inWord {
			words, inWord = append(words, word.String()), false
			word.Reset()
		}
	}

	for _, r := range line {
		switch {
		case quote != 0 && r == quote:
			quote = 0
		case quote != 0:
			word.WriteRune(r)
		case r == '\'' || r == '"':
			quote, inWord = r, true
		case strings.ContainsRune(";&|()\n`", r):
			endWord()

			if len(words) > 0 {
				commands, words = append(commands, words), nil
			}
		case unicode.IsSpace(r):
			endWord()
		default:
			word.WriteRune(r)

			inWord = true
		}
	}

	endWord()

	return append(commands, words)
}

// isGateCall reports whether a simple command runs yarbrain brief or verdict:
// after assignments and wrappers, a program named yarbrain in any directory
// and case, then global flags, then brief or verdict.
func isGateCall(words []string) bool {
	i := 0
	for i < len(words) && (isAssignment(words[i]) || slices.Contains(commandWrappers, words[i]) || strings.HasPrefix(words[i], "-")) {
		i++
	}

	if i == len(words) || !isYarbrain(words[i]) {
		return false
	}

	for i++; i < len(words); i++ {
		switch word := words[i]; {
		case word == "--vault" || word == "--config":
			i++
		case strings.HasPrefix(word, "-"):
		default:
			return word == "brief" || word == "verdict"
		}
	}

	return false
}

func isAssignment(word string) bool {
	name, _, ok := strings.Cut(word, "=")

	return ok && name != "" && strings.IndexFunc(name, func(r rune) bool { return r != '_' && !unicode.IsLetter(r) && !unicode.IsDigit(r) }) < 0
}

func isYarbrain(word string) bool {
	name := strings.ToLower(word[strings.LastIndexAny(word, `/\`)+1:])

	return name == "yarbrain" || name == "yarbrain.exe"
}

// DenyReason is the reason the Claude Code hook gives when it blocks a gate
// command outside the gate agent.
const DenyReason = "Only the yarbrain:lesson-gate agent may run `yarbrain brief` or `yarbrain verdict`. " +
	"Start that agent with only the lesson id; never judge a lesson yourself."

// DeniesGateCommand reports whether a Claude Code PreToolUse payload runs a
// gate command outside the gate agent. Claude Code subagents share their
// parent's session id, so the agent type in the payload is the only signal
// that tells the gate apart.
func DeniesGateCommand(p Payload) bool {
	return (p.ToolName == "Bash" || p.ToolName == "PowerShell") && p.AgentType != GateAgent && runsGateCommand(p.ToolInput.Command)
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
