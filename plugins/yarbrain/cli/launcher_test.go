package install_test

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeCLI stands in for a release binary: it prints its arguments, then
// its stdin.
var fakeCLI = []byte("#!/bin/sh\necho \"ran $*\"\ncat\n")

type launcherRun struct {
	stdout string
	stderr string
	code   int
}

type envelope struct {
	OK    bool `json:"ok"`
	Error struct {
		Code string `json:"code"`
		Fix  string `json:"fix"`
	} `json:"error"`
}

type launcherSetup struct {
	launcher string
	home     string
	cache    string
	config   string
	server   *httptest.Server
}

func launcherAsset(t *testing.T) string {
	t.Helper()

	if runtime.GOOS == "windows" {
		t.Skip("the launcher is a POSIX shell script")
	}

	return fmt.Sprintf("yarbrain_0.3.0_%s_%s", runtime.GOOS, runtime.GOARCH)
}

// newLauncher lays out a plugin with the launcher and a pin file whose
// checksum is pinned, and a release server that serves served.
func newLauncher(t *testing.T, pinned string, served []byte) launcherSetup {
	t.Helper()

	asset := launcherAsset(t)
	plugin := t.TempDir()

	script, err := os.ReadFile(filepath.Join("..", "bin", "yarbrain"))
	require.NoError(t, err)
	require.NoError(t, os.MkdirAll(filepath.Join(plugin, "bin"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(plugin, "bin", "yarbrain"), script, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(plugin, "cli-checksums.txt"),
		fmt.Appendf(nil, "%s  yarbrain_0.3.0_other_arch\n%s  %s\n", digest([]byte("other")), pinned, asset), 0o644))

	mux := http.NewServeMux()
	mux.HandleFunc("/plugins/yarbrain/cli/v0.3.0/"+asset, func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write(served) })
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return launcherSetup{
		launcher: filepath.Join(plugin, "bin", "yarbrain"),
		home:     t.TempDir(),
		cache:    t.TempDir(),
		config:   t.TempDir(),
		server:   server,
	}
}

func (s launcherSetup) run(t *testing.T, env []string, args ...string) launcherRun {
	t.Helper()

	return s.runInput(t, env, "", args...)
}

func (s launcherSetup) runInput(t *testing.T, env []string, input string, args ...string) launcherRun {
	t.Helper()

	// A launcher that ran itself in a loop would otherwise hang the test.
	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
	t.Cleanup(cancel)

	command := exec.CommandContext(ctx, s.launcher, args...)
	// A fixed system PATH keeps any yarbrain that this machine installed out
	// of the run; a later PATH in env replaces it.
	command.Env = append([]string{
		"PATH=/usr/bin:/bin:/usr/sbin:/sbin",
		"HOME=" + s.home,
		"XDG_CACHE_HOME=" + s.cache,
		"XDG_CONFIG_HOME=" + s.config,
		"YARBRAIN_RELEASE_BASE=" + s.server.URL,
	}, env...)

	var stdout, stderr strings.Builder

	command.Stdin = strings.NewReader(input)
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()

	var exit *exec.ExitError
	if err != nil && !errors.As(err, &exit) {
		require.NoError(t, err)
	}

	return launcherRun{stdout: stdout.String(), stderr: stderr.String(), code: command.ProcessState.ExitCode()}
}

func decodeEnvelope(t *testing.T, out string) envelope {
	t.Helper()

	var result envelope
	require.NoError(t, json.Unmarshal([]byte(out), &result), "stdout must be one JSON object: %q", out)

	return result
}

func TestLauncherDownloadsThePinnedBuildAndRunsItOnTheNextCall(t *testing.T) {
	setup := newLauncher(t, digest(fakeCLI), fakeCLI)

	first := setup.run(t, nil, "status")

	assert.Equal(t, 4, first.code, first.stderr)
	result := decodeEnvelope(t, first.stdout)
	assert.False(t, result.OK)
	assert.Equal(t, "busy", result.Error.Code)
	assert.Equal(t, "Run the same command again.", result.Error.Fix)

	cached, err := os.ReadFile(filepath.Join(setup.cache, "yarbrain", digest(fakeCLI), "yarbrain"))
	require.NoError(t, err)
	assert.Equal(t, fakeCLI, cached)

	setup.server.Close()

	second := setup.run(t, nil, "find", "make", "seed")

	assert.Equal(t, 0, second.code, second.stderr)
	assert.Equal(t, "ran find make seed\n", second.stdout)
}

func TestLauncherRefusesADownloadThatDoesNotMatchThePin(t *testing.T) {
	setup := newLauncher(t, digest(fakeCLI), []byte("#!/bin/sh\necho swapped\n"))

	run := setup.run(t, nil, "status")

	assert.Equal(t, 5, run.code)
	assert.Equal(t, "io", decodeEnvelope(t, run.stdout).Error.Code)
	assert.Contains(t, run.stderr, "does not match the plugin's checksum; nothing was installed")
	assert.NoFileExists(t, filepath.Join(setup.cache, "yarbrain", digest(fakeCLI), "yarbrain"))
}

func TestLauncherNamesTheSandboxFixWhenTheDownloadFails(t *testing.T) {
	setup := newLauncher(t, digest(fakeCLI), fakeCLI)
	setup.server.Close()

	run := setup.run(t, nil, "status")

	assert.Equal(t, 5, run.code)
	result := decodeEnvelope(t, run.stdout)
	assert.Equal(t, "io", result.Error.Code)
	assert.Equal(t, "If a sandbox blocks network access, run the same command again with network access allowed.", result.Error.Fix)
}

func TestLauncherHookStaysSilentWithoutAVault(t *testing.T) {
	setup := newLauncher(t, digest(fakeCLI), fakeCLI)
	setup.server.Close()

	run := setup.run(t, []string{"YARBRAIN_NO_DOWNLOAD=1"}, "hook", "session-start", "--contract", "1")

	assert.Equal(t, 0, run.code)
	assert.Empty(t, run.stdout)
	assert.NoDirExists(t, filepath.Join(setup.cache, "yarbrain"))
}

func (s launcherSetup) writeConfig(t *testing.T) {
	t.Helper()

	require.NoError(t, os.MkdirAll(filepath.Join(s.config, "yarbrain"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(s.config, "yarbrain", "config.json"), []byte(`{"vault":"/v"}`), 0o644))
}

func TestLauncherHookAsksClaudeCodeToDownloadWhenAVaultIsSetUp(t *testing.T) {
	setup := newLauncher(t, digest(fakeCLI), fakeCLI)
	setup.server.Close()
	setup.writeConfig(t)

	run := setup.run(t, []string{"YARBRAIN_NO_DOWNLOAD=1"}, "hook", "session-start", "--contract", "1")

	assert.Equal(t, 0, run.code)
	assert.Contains(t, run.stdout, "Before you start the task, run `yarbrain status` now, without asking")
	assert.Contains(t, run.stdout, "run `yarbrain find <terms>` before you plan")
	assert.NoDirExists(t, filepath.Join(setup.cache, "yarbrain"))
}

func TestLauncherHookGivesCodexItsPathWhenAVaultIsSetUp(t *testing.T) {
	setup := newLauncher(t, digest(fakeCLI), fakeCLI)
	setup.server.Close()
	setup.writeConfig(t)

	run := setup.run(t, []string{"YARBRAIN_NO_DOWNLOAD=1", "PLUGIN_ROOT=/plugin"}, "hook", "session-start", "--contract", "1")

	assert.Equal(t, 0, run.code)
	assert.Contains(t, run.stdout, "Before you start the task, run `\""+setup.launcher+"\" status` now, without asking")
	assert.Contains(t, run.stdout, "then `\""+setup.launcher+"\" show <id>` only for lessons whose cue fits.")
}

func TestLauncherHookPassesItsPayloadToTheCachedBuild(t *testing.T) {
	setup := newLauncher(t, digest(fakeCLI), fakeCLI)
	setup.run(t, nil, "status")
	setup.server.Close()

	run := setup.runInput(t, []string{"YARBRAIN_NO_DOWNLOAD=1"}, `{"tool_name":"Bash"}`, "hook", "pre-tool-use", "--contract", "1")

	assert.Equal(t, 0, run.code)
	assert.Equal(t, "ran hook pre-tool-use --contract 1\n{\"tool_name\":\"Bash\"}", run.stdout)
}

func TestLauncherPreToolUseRefusesAGateCallBeforeTheDownload(t *testing.T) {
	setup := newLauncher(t, digest(fakeCLI), fakeCLI)
	setup.server.Close()

	run := setup.runInput(t, []string{"YARBRAIN_NO_DOWNLOAD=1"},
		`{"tool_name":"Bash","tool_input":{"command":"yarbrain status; yarbrain brief l-abc"}}`, "hook", "pre-tool-use", "--contract", "1")

	assert.Equal(t, 2, run.code)
	assert.Empty(t, run.stdout)
	assert.Contains(t, run.stderr, "Run `yarbrain status` alone first; it downloads the CLI.")
}

func TestLauncherPreToolUseAllowsOtherCallsBeforeTheDownload(t *testing.T) {
	for _, command := range []string{"yarbrain status", "git commit -m 'brief fix'"} {
		t.Run(command, func(t *testing.T) {
			setup := newLauncher(t, digest(fakeCLI), fakeCLI)
			setup.server.Close()

			run := setup.runInput(t, []string{"YARBRAIN_NO_DOWNLOAD=1"},
				`{"tool_name":"Bash","tool_input":{"command":"`+command+`"}}`, "hook", "pre-tool-use", "--contract", "1")

			assert.Equal(t, 0, run.code, run.stderr)
			assert.Empty(t, run.stdout)
		})
	}
}

func TestLauncherHookUsesAYarbrainOnPathBeforeTheDownload(t *testing.T) {
	setup := newLauncher(t, digest(fakeCLI), fakeCLI)
	setup.server.Close()

	installed := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(installed, "yarbrain"), []byte("#!/bin/sh\necho \"installed $*\"\n"), 0o755))

	run := setup.run(t, []string{"YARBRAIN_NO_DOWNLOAD=1", "PATH=" + installed + ":/usr/bin:/bin:/usr/sbin:/sbin"},
		"hook", "session-start", "--contract", "1")

	assert.Equal(t, 0, run.code, run.stderr)
	assert.Equal(t, "installed hook session-start --contract 1\n", run.stdout)
}

func TestLauncherHookDoesNotRunItselfFromPath(t *testing.T) {
	setup := newLauncher(t, digest(fakeCLI), fakeCLI)
	setup.server.Close()
	setup.writeConfig(t)

	run := setup.run(t, []string{"YARBRAIN_NO_DOWNLOAD=1", "PATH=" + filepath.Dir(setup.launcher) + ":/usr/bin:/bin:/usr/sbin:/sbin"},
		"hook", "session-start", "--contract", "1")

	assert.Equal(t, 0, run.code, run.stderr)
	assert.Contains(t, run.stdout, "has not downloaded its CLI yet")
}

func TestLauncherHookExpandsAHomeRelativeConfig(t *testing.T) {
	setup := newLauncher(t, digest(fakeCLI), fakeCLI)
	setup.server.Close()
	require.NoError(t, os.WriteFile(filepath.Join(setup.home, "yarbrain.json"), []byte(`{"vault":"/v"}`), 0o644))

	run := setup.run(t, []string{"YARBRAIN_NO_DOWNLOAD=1", "YARBRAIN_CONFIG=~/yarbrain.json"}, "hook", "session-start", "--contract", "1")

	assert.Equal(t, 0, run.code, run.stderr)
	assert.Contains(t, run.stdout, "has not downloaded its CLI yet")
}

func TestLauncherRefusesToRunOfflineOutsideAHook(t *testing.T) {
	setup := newLauncher(t, digest(fakeCLI), fakeCLI)
	setup.server.Close()

	run := setup.run(t, []string{"YARBRAIN_NO_DOWNLOAD=1"}, "status")

	assert.Equal(t, 5, run.code)
	result := decodeEnvelope(t, run.stdout)
	assert.Equal(t, "io", result.Error.Code)
	assert.Equal(t, "Run the same command without YARBRAIN_NO_DOWNLOAD.", result.Error.Fix)
}

func TestLauncherSessionStartGivesCodexItsPathWhenYarbrainIsNotOnPath(t *testing.T) {
	setup := newLauncher(t, digest(fakeCLI), fakeCLI)
	setup.run(t, nil, "status")

	run := setup.run(t, []string{"YARBRAIN_NO_DOWNLOAD=1", "PLUGIN_ROOT=/plugin"}, "hook", "session-start", "--contract", "1")

	assert.Equal(t, 0, run.code, run.stderr)
	assert.Equal(t, "ran hook session-start --contract 1 The yarbrain command is not on PATH in this session; run it as \""+setup.launcher+"\".\n", run.stdout)
}

func TestLauncherSessionStartPrintsOnlyTheCLILineInCodexWhenYarbrainIsOnPath(t *testing.T) {
	setup := newLauncher(t, digest(fakeCLI), fakeCLI)
	setup.run(t, nil, "status")

	run := setup.run(t, []string{"YARBRAIN_NO_DOWNLOAD=1", "PLUGIN_ROOT=/plugin", "PATH=" + filepath.Dir(setup.launcher) + ":/usr/bin:/bin:/usr/sbin:/sbin"},
		"hook", "session-start", "--contract", "1")

	assert.Equal(t, 0, run.code, run.stderr)
	assert.Equal(t, "ran hook session-start --contract 1\n", run.stdout)
}

func TestLauncherSessionStartPrintsOnlyTheCLILineInClaudeCode(t *testing.T) {
	setup := newLauncher(t, digest(fakeCLI), fakeCLI)
	setup.run(t, nil, "status")

	run := setup.run(t, []string{"YARBRAIN_NO_DOWNLOAD=1"}, "hook", "session-start", "--contract", "1")

	assert.Equal(t, 0, run.code, run.stderr)
	assert.Equal(t, "ran hook session-start --contract 1\n", run.stdout)
}

func TestPluginPinsOneVersionForEveryReleaseTarget(t *testing.T) {
	pins, err := os.Open(filepath.Join("..", "cli-checksums.txt"))
	require.NoError(t, err)

	t.Cleanup(func() { _ = pins.Close() })

	versions := map[string]bool{}
	targets := map[string]bool{}

	lines := bufio.NewScanner(pins)
	for lines.Scan() {
		fields := strings.Fields(lines.Text())
		require.Len(t, fields, 2, "line %q", lines.Text())
		assert.Len(t, fields[0], 64, "sha256 of %s", fields[1])

		parts := strings.SplitN(strings.TrimPrefix(fields[1], "yarbrain_"), "_", 2)
		require.Len(t, parts, 2, "asset %q", fields[1])
		versions[parts[0]] = true
		targets[parts[1]] = true
	}

	require.NoError(t, lines.Err())
	assert.Len(t, versions, 1, "the plugin must pin one CLI version: %v", versions)
	assert.Equal(t, map[string]bool{
		"darwin_amd64": true, "darwin_arm64": true, "linux_amd64": true, "linux_arm64": true, "windows_amd64.exe": true,
	}, targets)
}
