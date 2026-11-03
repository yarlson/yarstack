package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/vault"
)

type setupStep struct {
	Host string `json:"host"`
	Step string `json:"step"`
}

type initResult struct {
	header

	Vault  initVault   `json:"vault"`
	Config initConfig  `json:"config"`
	Setup  []setupStep `json:"setup"`
}

type initVault struct {
	vault.InitResult

	Path string `json:"path"`
}

type initConfig struct {
	Path    string `json:"path"`
	Written bool   `json:"written"`
}

func newInitCommand(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "init <vault-path>",
		Short: "Create or adopt a vault and point the config at it",
		Long: "init creates the lesson layout in <vault-path>, or adds it to an existing\n" +
			"directory without touching any existing file, including v0.2 content. It\n" +
			"writes the config only when the config is absent or already names this vault.\n" +
			"The result lists the one-time host setup steps for the hosts it detects.",
		Example: "  yarbrain init ~/yarbrain-vault",
		Args:    exactArgs(1, "vault-path"),
		RunE: func(_ *cobra.Command, args []string) error {
			location, err := vault.Locator{FlagVault: args[0], Getenv: a.getenv, Home: a.home}.Resolve()
			if err != nil {
				return err
			}

			created, err := vault.Init(location.Vault)
			if err != nil {
				return err
			}

			configPath := a.locator().ConfigPath()

			written, err := vault.WriteConfig(configPath, location.Vault)
			if err != nil {
				return err
			}

			return writeJSON(a.stdout, initResult{
				header: okHeader(),
				Vault:  initVault{InitResult: created, Path: location.Vault},
				Config: initConfig{Path: configPath, Written: written},
				Setup:  a.setupSteps(location.Vault),
			})
		},
	}
}

// setupSteps lists the one-time host configuration that keeps vault writes
// free of prompts. The CLI never edits host configuration itself.
func (a *app) setupSteps(vaultPath string) []setupStep {
	steps := []setupStep{}

	if dirExists(filepath.Join(a.home, ".claude")) {
		steps = append(steps, setupStep{
			Host: "claude-code",
			Step: `Add "Bash(yarbrain *)" to permissions.allow in ~/.claude/settings.json.`,
		})
	}

	if dirExists(a.codexHome()) {
		steps = append(steps,
			setupStep{
				Host: "codex",
				Step: fmt.Sprintf("Add %q to writable_roots under [sandbox_workspace_write] in %s.", vaultPath, a.codexFile("config.toml")),
			},
			setupStep{
				Host: "codex",
				Step: "Pin the gate in " + a.codexFile("agents/yarbrain-gate.toml") + ": name = \"yarbrain-gate\", " +
					"description = \"Independent Yarbrain lesson gate\", developer_instructions = " +
					"\"Judge the one Yarbrain lesson id in your task with the lesson-gate skill.\", " +
					"model_reasoning_effort = \"high\", and a model different from your main session's.",
			},
		)
	}

	return steps
}

func dirExists(path string) bool {
	info, err := os.Stat(path)

	return err == nil && info.IsDir()
}
