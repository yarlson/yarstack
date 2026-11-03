package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/trust"
)

const (
	maxConfirmNote = 200
	maxDoubtNote   = 500
)

type confirmInput struct {
	Note string `json:"note"`
}

type doubtInput struct {
	Digest string `json:"digest"`
	Note   string `json:"note"`
}

type confirmResult struct {
	header

	ID            string `json:"id"`
	Confirmations int    `json:"confirmations"`
}

type doubtResult struct {
	header

	ID     string `json:"id"`
	Status string `json:"status"`
}

func newConfirmCommand(a *app) *cobra.Command {
	var input string

	command := &cobra.Command{
		Use:   "confirm <id>",
		Short: "Record that a served lesson worked",
		Long: "confirm records that applying a served lesson worked. Each session confirms a\n" +
			"lesson at most once, and the author's confirmations do not count. Input on\n" +
			"stdin: {\"note\": \"<one line, at most 200 characters>\"}",
		Example: "  echo '{\"note\": \"seed passed after migrate\"}' | yarbrain confirm l-7k3m9q2x",
		Args:    exactArgs(1, "id"),
		RunE: func(_ *cobra.Command, args []string) error {
			var in confirmInput
			if err := a.readJSON(input, &in); err != nil {
				return err
			}

			if problems := oneLine("note", in.Note, maxConfirmNote); len(problems) > 0 {
				return problemsError(problems)
			}

			if failure := screen(map[string]string{"note": in.Note}); failure != nil {
				return failure
			}

			result, err := a.recordConfirm(args[0], in)
			if err != nil {
				return err
			}

			return writeJSON(a.stdout, result)
		},
	}
	command.Flags().StringVar(&input, "input", "", "read the JSON from this file instead of stdin")

	return command
}

func newDoubtCommand(a *app) *cobra.Command {
	var input string

	command := &cobra.Command{
		Use:   "doubt <id>",
		Short: "Take a served lesson out of service because it was wrong",
		Long: "doubt records that a served lesson failed or contradicts the code. The lesson\n" +
			"is no longer served until the independent gate judges it again. Each session\n" +
			"doubts a lesson at most once. Input on stdin:\n" +
			"{\"digest\": \"<the digest from show>\", \"note\": \"<what went wrong, at most 500 characters>\"}",
		Example: "  yarbrain doubt l-7k3m9q2x < doubt.json",
		Args:    exactArgs(1, "id"),
		RunE: func(_ *cobra.Command, args []string) error {
			var in doubtInput
			if err := a.readJSON(input, &in); err != nil {
				return err
			}

			problems := oneLine("note", in.Note, maxDoubtNote)
			if in.Digest == "" {
				problems = append(problems, lesson.Problem{Field: "digest", Rule: "required"})
			}

			if len(problems) > 0 {
				return problemsError(problems)
			}

			if failure := screen(map[string]string{"note": in.Note}); failure != nil {
				return failure
			}

			result, err := a.recordDoubt(args[0], in)
			if err != nil {
				return err
			}

			return writeJSON(a.stdout, result)
		},
	}
	command.Flags().StringVar(&input, "input", "", "read the JSON from this file instead of stdin")

	return command
}

func (a *app) readJSON(path string, target any) error {
	raw, err := a.readInput(path)
	if err != nil {
		return err
	}

	return decodeStrict(raw, target)
}

func (a *app) recordConfirm(id string, in confirmInput) (confirmResult, error) {
	v, _, err := a.openVault()
	if err != nil {
		return confirmResult{}, err
	}
	defer func() { _ = v.Close() }()

	session := a.session()

	var result confirmResult

	err = v.WithLock(func() error {
		entries, _, err := v.Lessons()
		if err != nil {
			return err
		}

		target, err := loadRecord(v, entries, id)
		if err != nil {
			return err
		}

		switch {
		case !target.served():
			return badState(target, "it is not served", "Confirm only a lesson that find or show served to you.")
		case session == trust.Author(target.events):
			return badState(target, "this session proposed it", "Only another session's confirmation counts.")
		case trust.Recorded(target.events, trust.Confirm, session):
			return newError("already_recorded", "this session already confirmed lesson "+id, "Nothing to do.").with("id", id)
		}

		event := trust.Event{At: a.now().UTC(), Kind: trust.Confirm, Session: session, Digest: target.digest(), Note: in.Note}
		if err := v.AppendLedger(id, event); err != nil {
			return err
		}

		count, _ := trust.Confirmations(append(target.events, event))
		result = confirmResult{header: okHeader(), ID: id, Confirmations: count}

		return nil
	})

	return result, err
}

func (a *app) recordDoubt(id string, in doubtInput) (doubtResult, error) {
	v, _, err := a.openVault()
	if err != nil {
		return doubtResult{}, err
	}
	defer func() { _ = v.Close() }()

	session := a.session()

	var result doubtResult

	err = v.WithLock(func() error {
		entries, _, err := v.Lessons()
		if err != nil {
			return err
		}

		target, err := loadRecord(v, entries, id)
		if err != nil {
			return err
		}

		switch {
		case !target.served():
			return badState(target, "it is not served", "Doubt only a lesson that find or show served to you.")
		case in.Digest != target.digest():
			return newError("changed_underneath", "the lesson changed after you read it", fmt.Sprintf("Run `yarbrain show %s` and doubt the current text.", id)).
				with("id", id).with("digest", target.digest())
		case trust.Recorded(target.events, trust.Doubt, session):
			return newError("already_recorded", "this session already doubted lesson "+id, "Nothing to do.").with("id", id)
		}

		// A lesson trusted again and again on the same content is contested:
		// the next doubt retires it instead of sending it back to the gate.
		next := lesson.Doubted
		if trust.Trusts(target.events, target.digest()) >= trust.MaxTrusts {
			next = lesson.Retired
		}

		if err := v.SetLessonKeys(target.entry, map[string]string{"status": string(next), "trusted_digest": ""}); err != nil {
			return err
		}

		event := trust.Event{At: a.now().UTC(), Kind: trust.Doubt, Session: session, Digest: in.Digest, Note: in.Note}
		if err := v.AppendLedger(id, event); err != nil {
			return err
		}

		result = doubtResult{header: okHeader(), ID: id, Status: string(next)}

		return nil
	})

	return result, err
}
