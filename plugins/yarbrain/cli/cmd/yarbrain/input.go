package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/textrules"
)

// maxInput caps every free-text input; a lesson file is at most 16 KiB.
const maxInput = 32 * 1024

// readInput reads free text from --input FILE or stdin. Free text never
// arrives in arguments, because Windows shells mangle quoted arguments.
func (a *app) readInput(path string) ([]byte, error) {
	reader := a.stdin

	if path != "" {
		file, err := os.Open(path)
		if err != nil {
			return nil, newError("bad_input", err.Error(), "Pass a readable file to --input, or pipe the input on stdin.").
				with("field", "input")
		}
		defer func() { _ = file.Close() }()

		reader = file
	}

	data, err := io.ReadAll(io.LimitReader(reader, maxInput+1))
	if err != nil {
		return nil, fmt.Errorf("read input: %w", err)
	}

	if len(data) > maxInput {
		return nil, newError("bad_input", fmt.Sprintf("the input exceeds %d bytes", maxInput), "Shorten the input.").
			with("field", "input").with("limit", maxInput)
	}

	if len(bytes.TrimSpace(data)) == 0 {
		return nil, newError("bad_input", "the input is empty", "Pass the input on stdin or with --input FILE; see --help for the format.").
			with("field", "input")
	}

	return data, nil
}

func problemsError(problems []lesson.Problem) *cliError {
	return newError("bad_input", fmt.Sprintf("%d field(s) break the lesson rules", len(problems)),
		"Fix every listed field and run the same command again.").with("problems", problems)
}

// screen runs the deterministic text rules on agent-written fields. A secret
// finding wins over a risky one, because it must never be stored.
func screen(fields map[string]string) *cliError {
	var findings []textrules.Finding

	for _, field := range []string{"title", "cue", "body", "observed", "reason", "note"} {
		if text, ok := fields[field]; ok {
			findings = append(findings, textrules.Check(field, text)...)
		}
	}

	if len(findings) == 0 {
		return nil
	}

	code, fix := "risky_command", "Describe the step without the flagged text; a person can run risky steps themselves."

	for _, finding := range findings {
		if finding.Category == textrules.Secret {
			code, fix = "secret_suspected", "Remove the secret or personal data, or replace it with a placeholder such as <token>."

			break
		}
	}

	return newError(code, fmt.Sprintf("%d field(s) match a known pattern", len(findings)), fix).with("findings", findings)
}
