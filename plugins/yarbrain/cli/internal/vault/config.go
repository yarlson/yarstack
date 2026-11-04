// Package vault owns where the vault is, how its files look on disk, and every
// read and write of them. All vault file access goes through an os.Root, so
// no path can leave the vault.
package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ErrNotConfigured means no flag, environment variable, or config file names
// a vault.
var ErrNotConfigured = errors.New("no vault is configured")

// ConfigError means a config file exists but cannot be used.
type ConfigError struct {
	Path string
	Err  error
}

func (e *ConfigError) Error() string { return fmt.Sprintf("config %s: %v", e.Path, e.Err) }
func (e *ConfigError) Unwrap() error { return e.Err }

// ConflictError means the config already names a different vault.
type ConflictError struct {
	Path     string
	Existing string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("config %s already names vault %s", e.Path, e.Existing)
}

// Locator resolves the vault in the v0.2 order: the --vault flag,
// YARBRAIN_VAULT, then the config file from --config, YARBRAIN_CONFIG,
// $XDG_CONFIG_HOME/yarbrain/config.json, or ~/.config/yarbrain/config.json.
type Locator struct {
	FlagVault  string
	FlagConfig string
	Getenv     func(string) string
	Home       string
}

// Location is a resolved vault path and where it came from.
type Location struct {
	Vault  string `json:"path"`
	Config string `json:"config,omitempty"`
	Source string `json:"source"`
}

// ConfigPath returns the config file this locator reads and writes.
func (l Locator) ConfigPath() string {
	switch {
	case l.FlagConfig != "":
		return l.expand(l.FlagConfig)
	case l.Getenv("YARBRAIN_CONFIG") != "":
		return l.expand(l.Getenv("YARBRAIN_CONFIG"))
	case l.Getenv("XDG_CONFIG_HOME") != "":
		return filepath.Join(l.expand(l.Getenv("XDG_CONFIG_HOME")), "yarbrain", "config.json")
	default:
		return filepath.Join(l.Home, ".config", "yarbrain", "config.json")
	}
}

// Resolve finds the configured vault. The vault directory may not exist.
func (l Locator) Resolve() (Location, error) {
	if l.FlagVault != "" {
		return l.absolute(l.FlagVault, Location{Source: "flag"})
	}

	if env := l.Getenv("YARBRAIN_VAULT"); env != "" {
		return l.absolute(env, Location{Source: "env"})
	}

	path := l.ConfigPath()

	configured, err := readConfig(path)
	if err != nil {
		return Location{}, err
	}

	return l.absolute(configured, Location{Config: path, Source: "config"})
}

func (l Locator) absolute(path string, location Location) (Location, error) {
	resolved, err := filepath.Abs(l.expand(path))
	if err != nil {
		return Location{}, fmt.Errorf("resolve vault path %s: %w", path, err)
	}

	location.Vault = resolved

	return location, nil
}

func (l Locator) expand(path string) string {
	if path == "~" {
		return l.Home
	}

	if rest, ok := strings.CutPrefix(path, "~/"); ok {
		return filepath.Join(l.Home, rest)
	}

	return path
}

type configFile struct {
	Vault string `json:"vault"`
}

func readConfig(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return "", fmt.Errorf("%w: %s does not exist", ErrNotConfigured, path)
	}

	if err != nil {
		return "", &ConfigError{Path: path, Err: err}
	}

	var config configFile
	if err := json.Unmarshal(raw, &config); err != nil {
		return "", &ConfigError{Path: path, Err: fmt.Errorf("invalid JSON: %w", err)}
	}

	if config.Vault == "" {
		return "", &ConfigError{Path: path, Err: errors.New(`no "vault" string`)}
	}

	return config.Vault, nil
}

// WriteConfig points the config file at vault. It writes nothing when the
// config already names the same vault and refuses a different one.
func WriteConfig(path, vault string) (written bool, err error) {
	existing, err := readConfig(path)

	switch {
	case err == nil && sameDir(existing, vault):
		return false, nil
	case err == nil:
		return false, &ConflictError{Path: path, Existing: existing}
	case !errors.Is(err, ErrNotConfigured):
		return false, err
	}

	encoded, err := json.MarshalIndent(configFile{Vault: vault}, "", "  ")
	if err != nil {
		return false, fmt.Errorf("encode config: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, fmt.Errorf("create config directory: %w", err)
	}

	if err := writeFileAtomic(path, append(encoded, '\n')); err != nil {
		return false, err
	}

	return true, nil
}

// sameDir compares directories after resolving symlinks, because v0.2 stored
// the resolved vault path.
func sameDir(a, b string) bool {
	resolvedA, resolvedB := resolvedDir(a), resolvedDir(b)

	return resolvedA != "" && resolvedA == resolvedB
}

func resolvedDir(dir string) string {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return ""
	}

	if resolved, err := filepath.EvalSymlinks(abs); err == nil {
		return resolved
	}

	return filepath.Clean(abs)
}

// writeFileAtomic replaces path through a temporary file in the same
// directory, so a reader never sees a partial config.
func writeFileAtomic(path string, data []byte) error {
	temp, err := os.CreateTemp(filepath.Dir(path), ".config-*.json")
	if err != nil {
		return fmt.Errorf("create temporary config: %w", err)
	}

	defer func() { _ = os.Remove(temp.Name()) }()

	if _, err := temp.Write(data); err != nil {
		_ = temp.Close()

		return fmt.Errorf("write temporary config: %w", err)
	}

	if err := temp.Close(); err != nil {
		return fmt.Errorf("close temporary config: %w", err)
	}

	if err := os.Rename(temp.Name(), path); err != nil {
		return fmt.Errorf("replace config %s: %w", path, err)
	}

	return nil
}
