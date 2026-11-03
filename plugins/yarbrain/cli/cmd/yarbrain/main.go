// Command yarbrain owns every read and write of a Yarbrain lesson vault.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// version is set at build time with -ldflags "-X main.version=...".
var version = "dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// run executes one command and returns the process exit code. Every outcome
// except --help, --version, and the hook is one JSON object on stdout.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	root := newRootCommand(stdin, stdout, stderr)
	root.SetArgs(args)

	err := root.Execute()
	if err == nil {
		return exitOK
	}

	failure := asCLIError(err)
	if writeErr := writeError(stdout, failure); writeErr != nil {
		_, _ = fmt.Fprintf(stderr, "yarbrain: %v (while reporting %v)\n", writeErr, failure)

		return exitCrash
	}

	return failure.exitCode()
}

func newRootCommand(stdin io.Reader, stdout, stderr io.Writer) *cobra.Command {
	root := &cobra.Command{
		Use:   "yarbrain",
		Short: "Read and write a Yarbrain lesson vault",
		Long: "yarbrain records reusable know-how as lessons in a Markdown vault.\n" +
			"Every command prints one JSON object on stdout. Exit codes: 0 success,\n" +
			"1 check findings, 2 crash, 3 fix the call, 4 retry once, 5 setup problem.",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          rejectUnknownCommand,
		RunE: func(*cobra.Command, []string) error {
			return newError("bad_input", "no command given", "Run `yarbrain --help` to list the commands.").
				with("field", "command")
		},
	}
	root.SetIn(stdin)
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SetVersionTemplate("{{.Version}}\n")
	root.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return newError("bad_input", err.Error(), "Run the command with --help to see its flags.").
			with("field", "flags")
	})
	root.CompletionOptions.DisableDefaultCmd = true

	return root
}

// rejectUnknownCommand replaces Cobra's plain-text "unknown command" error,
// which would otherwise reach the caller as an internal error.
func rejectUnknownCommand(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}

	return newError("bad_input", fmt.Sprintf("unknown command %q", args[0]), "Run `yarbrain --help` to list the commands.").
		with("field", "command")
}
