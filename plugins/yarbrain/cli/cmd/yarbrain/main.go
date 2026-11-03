// Command yarbrain owns every read and write of a Yarbrain lesson vault.
package main

import (
	"context"
	"errors"
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
	a := newApp(stdin, stdout, stderr)
	root := newRootCommand(a)
	root.SetArgs(args)

	err := root.ExecuteContext(context.Background())
	if err == nil {
		return exitOK
	}

	if exit, ok := errors.AsType[*exitError](err); ok {
		return exit.code
	}

	failure := asCLIError(a.classify(err))
	if writeErr := writeError(stdout, failure); writeErr != nil {
		_, _ = fmt.Fprintf(stderr, "yarbrain: %v (while reporting %v)\n", writeErr, failure)

		return exitCrash
	}

	return failure.exitCode()
}

func newRootCommand(a *app) *cobra.Command {
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
	root.SetIn(a.stdin)
	root.SetOut(a.stdout)
	root.SetErr(a.stderr)
	root.SetVersionTemplate("{{.Version}}\n")
	root.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return newError("bad_input", err.Error(), "Run the command with --help to see its flags.").
			with("field", "flags")
	})
	root.CompletionOptions.DisableDefaultCmd = true

	root.PersistentFlags().StringVar(&a.flagVault, "vault", "", "use this vault instead of the configured one")
	root.PersistentFlags().StringVar(&a.flagConfig, "config", "", "read and write this config file")

	root.AddCommand(
		newInitCommand(a),
		newStatusCommand(a),
		newCheckCommand(a),
		newProposeCommand(a),
		newShowCommand(a),
	)

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

// exactArgs and noArgs replace Cobra's plain-text argument errors with
// bad_input, for the same reason.
func exactArgs(count int, names ...string) cobra.PositionalArgs {
	return func(_ *cobra.Command, args []string) error {
		if len(args) == count {
			return nil
		}

		return newError("bad_input", fmt.Sprintf("expected %d argument(s): %v, got %d", count, names, len(args)),
			"Run the command with --help to see its arguments.").with("field", "arguments")
	}
}

func noArgs(cmd *cobra.Command, args []string) error {
	return exactArgs(0)(cmd, args)
}

func minimumArgs(count int, name string) cobra.PositionalArgs {
	return func(_ *cobra.Command, args []string) error {
		if len(args) >= count {
			return nil
		}

		return newError("bad_input", fmt.Sprintf("expected at least %d %s argument(s), got %d", count, name, len(args)),
			"Run the command with --help to see its arguments.").with("field", "arguments")
	}
}
