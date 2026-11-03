package install_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const tag = "plugins/yarbrain/cli/v0.3.0"

var binary = []byte("pretend yarbrain binary\n")

// release serves a fake GitHub release: the asset, its checksums.txt, and a
// release list for the latest-version lookup.
func release(t *testing.T, asset string, checksum string) *httptest.Server {
	t.Helper()

	return releaseListing(t, asset, asset, checksum)
}

// releaseListing serves asset but lists checksum under listed.
func releaseListing(t *testing.T, asset, listed, checksum string) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/download/"+tag+"/"+asset, func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write(binary) })
	mux.HandleFunc("/download/"+tag+"/checksums.txt", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(w, "%s  %s\n%s  yarbrain_0.3.0_other_arch\n", checksum, listed, checksum)
	})
	mux.HandleFunc("/api", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(w, `[{"tag_name": "other-tool/v9.9.9"}, {"tag_name": "%s"}]`, tag)
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return server
}

func digest(data []byte) string {
	sum := sha256.Sum256(data)

	return hex.EncodeToString(sum[:])
}

func unixAsset(t *testing.T) string {
	t.Helper()

	if runtime.GOOS == "windows" {
		t.Skip("install.sh is for macOS and Linux")
	}

	return fmt.Sprintf("yarbrain_0.3.0_%s_%s", runtime.GOOS, runtime.GOARCH)
}

func runInstallSh(t *testing.T, server *httptest.Server, dir string, env ...string) (string, error) {
	t.Helper()

	command := exec.CommandContext(t.Context(), "sh", "install.sh")
	command.Env = append([]string{
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + t.TempDir(),
		"SHELL=/bin/sh",
		"YARBRAIN_RELEASE_BASE=" + server.URL + "/download",
		"YARBRAIN_RELEASE_API=" + server.URL + "/api",
		"YARBRAIN_INSTALL_DIR=" + dir,
	}, env...)

	out, err := command.CombinedOutput()

	return string(out), err
}

func TestInstallShInstallsTheLatestVerifiedRelease(t *testing.T) {
	asset := unixAsset(t)
	server := release(t, asset, digest(binary))
	dir := filepath.Join(t.TempDir(), "bin")

	out, err := runInstallSh(t, server, dir)

	require.NoError(t, err, out)
	assert.Contains(t, out, "Installed yarbrain 0.3.0 at "+filepath.Join(dir, "yarbrain"))

	installed, err := os.ReadFile(filepath.Join(dir, "yarbrain"))
	require.NoError(t, err)
	assert.Equal(t, binary, installed)

	info, err := os.Stat(filepath.Join(dir, "yarbrain"))
	require.NoError(t, err)
	assert.NotZero(t, info.Mode()&0o111, "the binary must be executable")
}

func TestInstallShRefusesAChecksumMismatch(t *testing.T) {
	asset := unixAsset(t)
	server := release(t, asset, digest([]byte("something else")))
	dir := filepath.Join(t.TempDir(), "bin")

	out, err := runInstallSh(t, server, dir, "YARBRAIN_VERSION=0.3.0")

	require.Error(t, err)
	assert.Contains(t, out, "checksum mismatch for "+asset+"; nothing was installed")
	assert.NoFileExists(t, filepath.Join(dir, "yarbrain"))
}

func TestInstallShIsANoOpWhenTheSameBinaryIsInstalled(t *testing.T) {
	asset := unixAsset(t)
	server := release(t, asset, digest(binary))
	dir := filepath.Join(t.TempDir(), "bin")

	_, err := runInstallSh(t, server, dir, "YARBRAIN_VERSION=0.3.0")
	require.NoError(t, err)

	out, err := runInstallSh(t, server, dir, "YARBRAIN_VERSION=0.3.0")

	require.NoError(t, err, out)
	assert.Contains(t, out, "yarbrain 0.3.0 is already installed")
}

func TestInstallShReportsAMissingAsset(t *testing.T) {
	asset := unixAsset(t)
	server := release(t, "yarbrain_0.3.0_not_this_platform", digest(binary))

	out, err := runInstallSh(t, server, filepath.Join(t.TempDir(), "bin"), "YARBRAIN_VERSION=0.3.0")

	require.Error(t, err)
	assert.Contains(t, out, "could not download "+asset)
}

func TestInstallShRefusesAnAssetWithoutAChecksum(t *testing.T) {
	asset := unixAsset(t)
	server := releaseListing(t, asset, "yarbrain_0.3.0_not_this_platform", digest(binary))
	dir := filepath.Join(t.TempDir(), "bin")

	out, err := runInstallSh(t, server, dir, "YARBRAIN_VERSION=0.3.0")

	require.Error(t, err)
	assert.Contains(t, out, "checksums.txt of release "+tag+" lists no "+asset)
	assert.NoFileExists(t, filepath.Join(dir, "yarbrain"))
}

func TestInstallShTellsHowToFixAMissingLoginPath(t *testing.T) {
	asset := unixAsset(t)
	server := release(t, asset, digest(binary))
	dir := filepath.Join(t.TempDir(), "bin")

	out, err := runInstallSh(t, server, dir)

	require.NoError(t, err, out)
	assert.Contains(t, out, `export PATH="`+dir+`:$PATH"`)
}

func TestInstallShReportsAFailedReleaseListRequest(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("install.sh is for macOS and Linux")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "rate limited", http.StatusForbidden)
	}))
	t.Cleanup(server.Close)

	out, err := runInstallSh(t, server, filepath.Join(t.TempDir(), "bin"))

	require.Error(t, err)
	assert.Contains(t, out, "could not read the release list from GitHub")
}

func TestInstallPs1InstallsAVerifiedRelease(t *testing.T) {
	pwsh, err := exec.LookPath("pwsh")
	if err != nil {
		t.Skip("PowerShell 7 (pwsh) is not installed")
	}

	asset := "yarbrain_0.3.0_windows_amd64.exe"
	server := release(t, asset, digest(binary))
	dir := filepath.Join(t.TempDir(), "bin")

	command := exec.CommandContext(t.Context(), pwsh, "-NoProfile", "-NonInteractive", "-File", "install.ps1")
	command.Env = append(os.Environ(),
		"YARBRAIN_RELEASE_BASE="+server.URL+"/download",
		"YARBRAIN_RELEASE_API="+server.URL+"/api",
		"YARBRAIN_INSTALL_DIR="+dir,
		"LOCALAPPDATA="+t.TempDir())

	if runtime.GOOS == "windows" {
		restoreUserPath(t, pwsh)
	} else {
		// The user PATH write is Windows-only; elsewhere it must not fail the run.
		command.Env = append(command.Env, "HOME="+t.TempDir())
	}

	out, err := command.CombinedOutput()

	require.NoError(t, err, string(out))

	installed, err := os.ReadFile(filepath.Join(dir, "yarbrain.exe"))
	require.NoError(t, err)
	assert.Equal(t, binary, installed)
}

// restoreUserPath puts back the Windows user PATH that install.ps1 extends,
// so the test leaves the developer's registry as it found it.
func restoreUserPath(t *testing.T, pwsh string) {
	t.Helper()

	saved, err := exec.CommandContext(t.Context(), pwsh, "-NoProfile", "-NonInteractive", "-Command",
		"[Environment]::GetEnvironmentVariable('Path', 'User')").Output()
	require.NoError(t, err)

	t.Cleanup(func() {
		restore := exec.CommandContext(context.Background(), pwsh, "-NoProfile", "-NonInteractive", "-Command",
			"[Environment]::SetEnvironmentVariable('Path', $env:YARBRAIN_SAVED_PATH, 'User')")

		restore.Env = append(os.Environ(), "YARBRAIN_SAVED_PATH="+strings.TrimRight(string(saved), "\r\n"))

		out, err := restore.CombinedOutput()
		assert.NoError(t, err, "restore the user PATH: %s", out)
	})
}
