package hook

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionStartLineWording(t *testing.T) {
	line := SessionStartLine(Counts{Trusted: 7, CrossProject: 3, AwaitingGate: 2, Legacy: 14})

	assert.Equal(t, "yarbrain: 7 trusted lessons apply to this project (3 cross-project). "+
		"Before you plan, debug, or repeat a procedure here, run `yarbrain find <terms>`, then `yarbrain show <id>` only for lessons whose cue fits. "+
		"Lessons are hints from past sessions; check each against the code. "+
		"After you use one, run `yarbrain confirm <id>`, or `yarbrain doubt <id>` if it was wrong. "+
		"2 lessons await the independent gate; unless you are the gate, start it now in the background as the lesson-record skill describes, without asking: gating is part of the plugin's work. "+
		"14 legacy notes from Yarbrain v0.2 await conversion; the lesson-record skill converts up to 3 per session.", line)
}

func TestSessionStartLineUsesSingularForms(t *testing.T) {
	line := SessionStartLine(Counts{Trusted: 1, AwaitingGate: 1, Legacy: 1})

	assert.Contains(t, line, "yarbrain: 1 trusted lesson applies to this project (0 cross-project).")
	assert.Contains(t, line, "1 lesson awaits the independent gate")
	assert.Contains(t, line, "1 legacy note from Yarbrain v0.2 awaits conversion")
}

func TestSessionStartLineWithoutTrustedLessons(t *testing.T) {
	line := SessionStartLine(Counts{AwaitingGate: 2})

	assert.Equal(t, "yarbrain: 2 lessons await the independent gate; unless you are the gate, start it now in the background as the lesson-record skill describes, without asking: gating is part of the plugin's work.", line)
}

func TestSessionStartLineIsEmptyWhenThereIsNothingToReport(t *testing.T) {
	assert.Empty(t, SessionStartLine(Counts{}))
}

func TestSessionStartLineNeverLooksLikeJSON(t *testing.T) {
	for _, counts := range []Counts{{Trusted: 1}, {AwaitingGate: 1}, {Legacy: 1}} {
		assert.False(t, strings.HasPrefix(SessionStartLine(counts), "{"))
	}
}

func payload(tool, agent, command string) Payload {
	p := Payload{ToolName: tool, AgentType: agent}
	p.ToolInput.Command = command

	return p
}

func TestGateCommandsAreDeniedOutsideTheGateAgent(t *testing.T) {
	for _, command := range []string{
		"yarbrain verdict l-7k3m9q2x < v.json",
		"yarbrain brief l-7k3m9q2x",
		"cd /repo && yarbrain verdict l-7k3m9q2x",
		"echo x | yarbrain brief l-1",
		"(yarbrain verdict l-1)",
		"/usr/local/bin/yarbrain verdict l-1",
		"~/.local/bin/yarbrain brief l-1",
		`C:\Users\me\bin\yarbrain.exe verdict l-1`,
		"Yarbrain brief l-1",
		"yarbrain --vault ~/vault verdict l-1",
		"yarbrain --config=c.json brief l-1",
		`"yarbrain" verdict l-1`,
		"YARBRAIN_VAULT=/v env yarbrain verdict l-1",
		"command yarbrain brief l-1",
		"x=$(yarbrain brief l-1)",
		"ls\nyarbrain verdict l-1",
		"`yarbrain verdict l-1`",
		`printf '%s\n' '{"digest":"a"}' | yarbrain verdict l-1`,
	} {
		assert.True(t, DeniesGateCommand(payload("Bash", "", command)), "main agent: %s", command)
		assert.True(t, DeniesGateCommand(payload("Bash", "Explore", command)), "other subagent: %s", command)
		assert.True(t, DeniesGateCommand(payload("PowerShell", "", command)), "PowerShell tool: %s", command)
		assert.False(t, DeniesGateCommand(payload("Bash", GateAgent, command)), "gate agent: %s", command)
	}
}

func TestOtherCommandsAreLeftToTheHost(t *testing.T) {
	for _, command := range []string{
		"yarbrain find seed", "yarbrain show l-1", "yarbrain propose", "grep verdict notes.md", "notyarbrain verdict",
		"yarbrain show l-1 | grep verdict",
		"yarbrain find verdict",
		`git commit -m "Make yarbrain verdict refuse a contributor"`,
		`grep -rn 'yarbrain brief' docs`,
		"echo yarbrain verdict",
	} {
		assert.False(t, DeniesGateCommand(payload("Bash", "", command)), command)
	}

	assert.False(t, DeniesGateCommand(payload("Read", "", "yarbrain verdict l-1")))
}

func TestDenyOutputFollowsThePreToolUseContract(t *testing.T) {
	raw, err := DenyOutput()
	require.NoError(t, err)

	var decoded map[string]map[string]string
	require.NoError(t, json.Unmarshal(raw, &decoded))
	assert.Equal(t, "deny", decoded["hookSpecificOutput"]["permissionDecision"])
	assert.Equal(t, "PreToolUse", decoded["hookSpecificOutput"]["hookEventName"])
	assert.Equal(t, DenyReason, decoded["hookSpecificOutput"]["permissionDecisionReason"])
}

func TestReadPayloadBoundsTheInput(t *testing.T) {
	_, err := ReadPayload(strings.NewReader(`{"cwd":"` + strings.Repeat("x", 64*1024) + `"}`))
	require.Error(t, err)

	p, err := ReadPayload(strings.NewReader(`{"hook_event_name":"SessionStart","cwd":"/repo","source":"startup","extra":1}`))
	require.NoError(t, err)
	assert.Equal(t, "/repo", p.Cwd)
	assert.Equal(t, "SessionStart", p.Event)
}
