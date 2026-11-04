package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// sandbox isolates a test from the real user configuration, vault, and host
// session: HOME, the config path, and every variable the CLI reads.
type sandbox struct {
	home   string
	config string
	vault  string
}

func isolate(t *testing.T) sandbox {
	t.Helper()

	s := sandbox{home: t.TempDir()}
	s.config = filepath.Join(s.home, "yarbrain-config.json")
	s.vault = filepath.Join(s.home, "vault")

	t.Setenv("HOME", s.home)
	t.Setenv("USERPROFILE", s.home)
	t.Setenv("YARBRAIN_CONFIG", s.config)
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("YARBRAIN_VAULT", "")
	t.Setenv("CLAUDE_CODE_SESSION_ID", "")
	t.Setenv("CODEX_THREAD_ID", "")
	t.Setenv("CLAUDE_EFFORT", "")
	t.Setenv("CODEX_HOME", "")

	return s
}

// configJSON encodes a config file that names vault, escaping the path for
// every OS.
func configJSON(t *testing.T, vault string) []byte {
	t.Helper()

	raw, err := json.Marshal(map[string]string{"vault": vault})
	require.NoError(t, err)

	return raw
}

func (s sandbox) initVault(t *testing.T) {
	t.Helper()

	result := invoke(t, "init", s.vault)
	require.Equal(t, exitOK, result.exit, result.stdout)
}

func errorOf(t *testing.T, result invocation) map[string]any {
	t.Helper()

	failure, ok := decodeObject(t, result.stdout)["error"].(map[string]any)
	require.True(t, ok, "stdout must hold an error object: %s", result.stdout)

	return failure
}

func TestInitCreatesAVaultAndWritesTheConfig(t *testing.T) {
	s := isolate(t)

	result := invoke(t, "init", s.vault)

	require.Equal(t, exitOK, result.exit, result.stdout)
	object := decodeObject(t, result.stdout)
	assert.Equal(t, true, object["ok"])
	assert.Equal(t, map[string]any{"path": s.config, "written": true}, object["config"])

	vault, ok := object["vault"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, s.vault, vault["path"])
	assert.Equal(t, false, vault["adopted"])
	assert.DirExists(t, filepath.Join(s.vault, "lessons"))

	raw, err := os.ReadFile(s.config)
	require.NoError(t, err)
	assert.JSONEq(t, string(configJSON(t, s.vault)), string(raw))
}

func TestInitIsIdempotent(t *testing.T) {
	s := isolate(t)
	s.initVault(t)

	result := invoke(t, "init", s.vault)

	require.Equal(t, exitOK, result.exit, result.stdout)
	object := decodeObject(t, result.stdout)
	assert.Equal(t, map[string]any{"path": s.config, "written": false}, object["config"])
	assert.Equal(t, true, object["vault"].(map[string]any)["adopted"])
}

func TestInitRefusesToRepointAConfiguredVault(t *testing.T) {
	s := isolate(t)
	s.initVault(t)

	result := invoke(t, "init", filepath.Join(s.home, "other-vault"))

	assert.Equal(t, exitSetup, result.exit)
	failure := errorOf(t, result)
	assert.Equal(t, "config_conflict", failure["code"])
	assert.Equal(t, s.vault, failure["existing"])
}

func TestInitListsSetupForDetectedHosts(t *testing.T) {
	s := isolate(t)
	require.NoError(t, os.Mkdir(filepath.Join(s.home, ".claude"), 0o755))
	require.NoError(t, os.Mkdir(filepath.Join(s.home, ".codex"), 0o755))

	result := invoke(t, "init", s.vault)

	require.Equal(t, exitOK, result.exit, result.stdout)

	steps, ok := decodeObject(t, result.stdout)["setup"].([]any)
	require.True(t, ok)

	hosts := map[string]int{}
	for _, step := range steps {
		hosts[step.(map[string]any)["host"].(string)]++
	}

	assert.Equal(t, map[string]int{"claude-code": 1, "codex": 2}, hosts)
}

func TestCodexSetupFollowsCodexHome(t *testing.T) {
	s := isolate(t)
	codexHome := filepath.Join(t.TempDir(), "codex")
	require.NoError(t, os.MkdirAll(filepath.Join(codexHome, "agents"), 0o755))
	t.Setenv("CODEX_HOME", codexHome)

	result := invoke(t, "init", s.vault)

	require.Equal(t, exitOK, result.exit, result.stdout)
	assert.Contains(t, result.stdout, filepath.ToSlash(filepath.Join(codexHome, "config.toml")))

	require.NoError(t, os.WriteFile(filepath.Join(codexHome, "agents", "yarbrain-gate.toml"), nil, 0o600))
	assert.Equal(t, true, decodeObject(t, invoke(t, "status").stdout)["codex_gate_pinned"])
}

func TestInitWithoutAnArgumentIsACallerError(t *testing.T) {
	isolate(t)

	result := invoke(t, "init")

	assert.Equal(t, exitFixCall, result.exit)
	assert.Equal(t, "bad_input", errorOf(t, result)["code"])
}

func TestInitIntoAReadOnlyDirectoryIsWriteDenied(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows ignores read-only permission bits on directories")
	}

	s := isolate(t)
	locked := filepath.Join(s.home, "locked")
	require.NoError(t, os.Mkdir(locked, 0o555))
	t.Cleanup(func() { _ = os.Chmod(locked, 0o755) })

	t.Setenv("CODEX_THREAD_ID", "thread-1")

	result := invoke(t, "init", filepath.Join(locked, "vault"))

	assert.Equal(t, exitSetup, result.exit)
	failure := errorOf(t, result)
	assert.Equal(t, "write_denied", failure["code"])
	assert.Contains(t, failure["fix"], "writable_roots")
}

func TestStatusWithoutAVaultTellsHowToCreateOne(t *testing.T) {
	isolate(t)

	result := invoke(t, "status")

	assert.Equal(t, exitSetup, result.exit)
	failure := errorOf(t, result)
	assert.Equal(t, "no_vault", failure["code"])
	assert.Contains(t, failure["fix"], "yarbrain init")
}

func TestStatusReportsAMissingConfiguredVault(t *testing.T) {
	s := isolate(t)
	require.NoError(t, os.WriteFile(s.config, configJSON(t, s.vault), 0o644))

	result := invoke(t, "status")

	assert.Equal(t, exitSetup, result.exit)
	failure := errorOf(t, result)
	assert.Equal(t, "no_vault", failure["code"])
	assert.Equal(t, s.vault, failure["path"])
}

func TestStatusReportsTheContractAndCounts(t *testing.T) {
	s := isolate(t)
	s.initVault(t)

	result := invoke(t, "status")

	require.Equal(t, exitOK, result.exit, result.stdout)
	object := decodeObject(t, result.stdout)
	assert.InDelta(t, 1, object["contract"], 0)
	assert.Equal(t, map[string]any{"trusted": 0.0, "cross_project": 0.0, "awaiting_gate": 0.0, "legacy": 0.0}, object["counts"])
	assert.Equal(t, map[string]any{
		"proposed": 0.0, "trusted": 0.0, "deferred": 0.0, "rejected": 0.0, "doubted": 0.0, "retired": 0.0, "changed": 0.0,
	}, object["totals"])
	assert.Equal(t, false, object["codex_gate_pinned"])

	vault, ok := object["vault"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, s.vault, vault["path"])
	assert.Equal(t, "config", vault["source"])
}

func TestBrokenConfigIsABadFile(t *testing.T) {
	s := isolate(t)
	require.NoError(t, os.WriteFile(s.config, []byte("{"), 0o644))

	result := invoke(t, "status")

	assert.Equal(t, exitSetup, result.exit)
	assert.Equal(t, "bad_file", errorOf(t, result)["code"])
}

func TestCheckOfACleanVaultHasNoFindings(t *testing.T) {
	s := isolate(t)
	s.initVault(t)

	result := invoke(t, "check")

	assert.Equal(t, exitOK, result.exit, result.stdout)
	assert.Equal(t, []any{}, decodeObject(t, result.stdout)["findings"])
}

func TestCheckReportsFindingsAndExitsOne(t *testing.T) {
	s := isolate(t)
	s.initVault(t)

	lesson := "---\nid: l-0000000a\nstatus: maybe\nscope: [all]\n---\n# Title\n\nbody\n"
	require.NoError(t, os.WriteFile(filepath.Join(s.vault, "lessons", "a.md"), []byte(lesson), 0o644))

	result := invoke(t, "check")

	assert.Equal(t, exitFindings, result.exit)
	findings, ok := decodeObject(t, result.stdout)["findings"].([]any)
	require.True(t, ok)

	codes := map[string]bool{}
	for _, finding := range findings {
		codes[finding.(map[string]any)["code"].(string)] = true
	}

	assert.Equal(t, map[string]bool{"unknown_status": true, "invalid_field": true}, codes)
}
