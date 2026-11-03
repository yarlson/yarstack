package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type invocation struct {
	exit   int
	stdout string
	stderr string
}

func invoke(t *testing.T, args ...string) invocation {
	t.Helper()

	return invokeWith(t, "", args...)
}

func invokeWith(t *testing.T, stdin string, args ...string) invocation {
	t.Helper()

	var stdout, stderr bytes.Buffer

	exit := run(args, strings.NewReader(stdin), &stdout, &stderr)

	return invocation{exit: exit, stdout: stdout.String(), stderr: stderr.String()}
}

func decodeObject(t *testing.T, output string) map[string]any {
	t.Helper()

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	require.Len(t, lines, 1, "stdout must hold exactly one JSON line: %q", output)

	var object map[string]any
	require.NoError(t, json.Unmarshal([]byte(lines[0]), &object), "stdout: %q", output)

	return object
}

func TestBareInvocationIsACallerError(t *testing.T) {
	result := invoke(t)

	assert.Equal(t, exitFixCall, result.exit)
	object := decodeObject(t, result.stdout)
	assert.InDelta(t, 1, object["v"], 0)
	assert.Equal(t, false, object["ok"])
	failure, ok := object["error"].(map[string]any)
	require.True(t, ok, "error must be an object: %v", object)
	assert.Equal(t, "bad_input", failure["code"])
	assert.Equal(t, "command", failure["field"])
	assert.NotEmpty(t, failure["fix"])
}

func TestUnknownCommandIsACallerErrorWithJSON(t *testing.T) {
	result := invoke(t, "no-such-command")

	assert.Equal(t, exitFixCall, result.exit)
	failure, ok := decodeObject(t, result.stdout)["error"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "bad_input", failure["code"])
	assert.Empty(t, result.stderr)
}

func TestUnknownFlagIsACallerErrorWithJSON(t *testing.T) {
	result := invoke(t, "--no-such-flag")

	assert.Equal(t, exitFixCall, result.exit)
	failure, ok := decodeObject(t, result.stdout)["error"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "bad_input", failure["code"])
	assert.Equal(t, "flags", failure["field"])
}

func TestVersionFlagPrintsOnlyTheVersion(t *testing.T) {
	result := invoke(t, "--version")

	assert.Equal(t, exitOK, result.exit)
	assert.Equal(t, version+"\n", result.stdout)
}

func TestExitCodeClassesFollowTheContract(t *testing.T) {
	cases := map[string]int{
		"internal":     exitCrash,
		"bad_input":    exitFixCall,
		"title_taken":  exitFixCall,
		"locked":       exitRetry,
		"busy":         exitRetry,
		"no_vault":     exitSetup,
		"write_denied": exitSetup,
	}
	for code, want := range cases {
		assert.Equal(t, want, newError(code, "m", "f").exitCode(), "code %s", code)
	}
}

func TestUnclassifiedErrorBecomesInternal(t *testing.T) {
	failure := asCLIError(errors.New("boom"))

	assert.Equal(t, "internal", failure.Code)
	assert.Equal(t, exitCrash, failure.exitCode())
}

func TestErrorEnvelopeCarriesCodeSpecificFields(t *testing.T) {
	var stdout bytes.Buffer

	failure := newError("title_taken", "taken", "Run `yarbrain show l-x`.").with("id", "l-x")
	require.NoError(t, writeError(&stdout, failure))

	object := decodeObject(t, stdout.String())
	body, ok := object["error"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "l-x", body["id"])
	assert.Equal(t, "title_taken", body["code"])
}
