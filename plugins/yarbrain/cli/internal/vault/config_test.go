package vault

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func envOf(values map[string]string) func(string) string {
	return func(key string) string { return values[key] }
}

func writeConfigFile(t *testing.T, path, content string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
}

// configNaming encodes a config that names vault, escaped for every OS.
func configNaming(t *testing.T, vault string) string {
	t.Helper()

	raw, err := json.Marshal(map[string]string{"vault": vault})
	require.NoError(t, err)

	return string(raw)
}

func TestResolutionKeepsTheV02Order(t *testing.T) {
	home := t.TempDir()
	xdg := t.TempDir()
	vaults := t.TempDir()
	explicit := filepath.Join(t.TempDir(), "explicit.json")
	fromEnv := filepath.Join(t.TempDir(), "env.json")

	vault := func(name string) string { return filepath.Join(vaults, name) }

	writeConfigFile(t, filepath.Join(home, ".config", "yarbrain", "config.json"), configNaming(t, vault("home")))
	writeConfigFile(t, filepath.Join(xdg, "yarbrain", "config.json"), configNaming(t, vault("xdg")))
	writeConfigFile(t, fromEnv, configNaming(t, vault("env-config")))
	writeConfigFile(t, explicit, configNaming(t, vault("flag-config")))

	all := map[string]string{"YARBRAIN_VAULT": vault("env"), "YARBRAIN_CONFIG": fromEnv, "XDG_CONFIG_HOME": xdg}
	cases := []struct {
		name    string
		locator Locator
		want    string
		source  string
	}{
		{"flag vault wins", Locator{FlagVault: vault("flag"), FlagConfig: explicit, Getenv: envOf(all), Home: home}, vault("flag"), "flag"},
		{"YARBRAIN_VAULT next", Locator{FlagConfig: explicit, Getenv: envOf(all), Home: home}, vault("env"), "env"},
		{"--config next", Locator{FlagConfig: explicit, Getenv: envOf(map[string]string{"YARBRAIN_CONFIG": fromEnv}), Home: home}, vault("flag-config"), "config"},
		{"YARBRAIN_CONFIG next", Locator{Getenv: envOf(map[string]string{"YARBRAIN_CONFIG": fromEnv, "XDG_CONFIG_HOME": xdg}), Home: home}, vault("env-config"), "config"},
		{"XDG_CONFIG_HOME next", Locator{Getenv: envOf(map[string]string{"XDG_CONFIG_HOME": xdg}), Home: home}, vault("xdg"), "config"},
		{"~/.config last", Locator{Getenv: envOf(nil), Home: home}, vault("home"), "config"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			location, err := tc.locator.Resolve()
			require.NoError(t, err)
			assert.Equal(t, tc.want, location.Vault)
			assert.Equal(t, tc.source, location.Source)
		})
	}
}

func TestResolutionExpandsHome(t *testing.T) {
	home := t.TempDir()

	location, err := Locator{FlagVault: "~/brain", Getenv: envOf(nil), Home: home}.Resolve()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(home, "brain"), location.Vault)
}

func TestMissingConfigMeansNotConfigured(t *testing.T) {
	_, err := Locator{Getenv: envOf(nil), Home: t.TempDir()}.Resolve()

	assert.ErrorIs(t, err, ErrNotConfigured)
}

func TestInvalidConfigIsAConfigError(t *testing.T) {
	for name, content := range map[string]string{"bad json": "{", "no vault": `{"other": 1}`} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.json")
			writeConfigFile(t, path, content)

			_, err := Locator{FlagConfig: path, Getenv: envOf(nil), Home: t.TempDir()}.Resolve()

			var configErr *ConfigError
			require.ErrorAs(t, err, &configErr)
			assert.Equal(t, path, configErr.Path)
		})
	}
}

func TestWriteConfigUsesTheV02FormatOnce(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "config.json")
	vault := filepath.Join(t.TempDir(), "vault")

	written, err := WriteConfig(path, vault)
	require.NoError(t, err)
	assert.True(t, written)

	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	quoted, err := json.Marshal(vault)
	require.NoError(t, err)
	assert.Equal(t, "{\n  \"vault\": "+string(quoted)+"\n}\n", string(raw))

	again, err := WriteConfig(path, vault)
	require.NoError(t, err)
	assert.False(t, again, "the same vault must not rewrite the config")
}

func TestWriteConfigTreatsASymlinkToTheSameVaultAsTheSameVault(t *testing.T) {
	real := filepath.Join(t.TempDir(), "vault")
	require.NoError(t, os.Mkdir(real, 0o755))

	link := filepath.Join(t.TempDir(), "vault-link")
	if err := os.Symlink(real, link); err != nil {
		t.Skipf("symlinks are not available: %v", err)
	}

	path := filepath.Join(t.TempDir(), "config.json")
	quoted, err := json.Marshal(real)
	require.NoError(t, err)
	writeConfigFile(t, path, `{"vault": `+string(quoted)+`}`)

	written, err := WriteConfig(path, link)

	require.NoError(t, err, "v0.2 stored the resolved path, so init through a link names the same vault")
	assert.False(t, written)
}

func TestWriteConfigRefusesADifferentVault(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	writeConfigFile(t, path, `{"vault": "/vaults/existing"}`)

	_, err := WriteConfig(path, "/vaults/other")

	var conflict *ConflictError
	require.ErrorAs(t, err, &conflict)
	assert.Equal(t, "/vaults/existing", conflict.Existing)

	raw, readErr := os.ReadFile(path)
	require.NoError(t, readErr)
	assert.JSONEq(t, `{"vault": "/vaults/existing"}`, string(raw))
}
