package vault

import (
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

func TestResolutionKeepsTheV02Order(t *testing.T) {
	home := t.TempDir()
	xdg := t.TempDir()
	explicit := filepath.Join(t.TempDir(), "explicit.json")
	fromEnv := filepath.Join(t.TempDir(), "env.json")

	writeConfigFile(t, filepath.Join(home, ".config", "yarbrain", "config.json"), `{"vault": "/vaults/home"}`)
	writeConfigFile(t, filepath.Join(xdg, "yarbrain", "config.json"), `{"vault": "/vaults/xdg"}`)
	writeConfigFile(t, fromEnv, `{"vault": "/vaults/env-config"}`)
	writeConfigFile(t, explicit, `{"vault": "/vaults/flag-config"}`)

	all := map[string]string{"YARBRAIN_VAULT": "/vaults/env", "YARBRAIN_CONFIG": fromEnv, "XDG_CONFIG_HOME": xdg}
	cases := []struct {
		name    string
		locator Locator
		want    string
		source  string
	}{
		{"flag vault wins", Locator{FlagVault: "/vaults/flag", FlagConfig: explicit, Getenv: envOf(all), Home: home}, "/vaults/flag", "flag"},
		{"YARBRAIN_VAULT next", Locator{FlagConfig: explicit, Getenv: envOf(all), Home: home}, "/vaults/env", "env"},
		{"--config next", Locator{FlagConfig: explicit, Getenv: envOf(map[string]string{"YARBRAIN_CONFIG": fromEnv}), Home: home}, "/vaults/flag-config", "config"},
		{"YARBRAIN_CONFIG next", Locator{Getenv: envOf(map[string]string{"YARBRAIN_CONFIG": fromEnv, "XDG_CONFIG_HOME": xdg}), Home: home}, "/vaults/env-config", "config"},
		{"XDG_CONFIG_HOME next", Locator{Getenv: envOf(map[string]string{"XDG_CONFIG_HOME": xdg}), Home: home}, "/vaults/xdg", "config"},
		{"~/.config last", Locator{Getenv: envOf(nil), Home: home}, "/vaults/home", "config"},
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
	assert.Equal(t, "{\n  \"vault\": \""+vault+"\"\n}\n", string(raw))

	again, err := WriteConfig(path, vault)
	require.NoError(t, err)
	assert.False(t, again, "the same vault must not rewrite the config")
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
