package main

import (
	"fmt"
	"slices"

	"github.com/spf13/cobra"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/trust"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/vault"
)

const verdictTemplate = `{
  "digest": "<the digest from brief>",
  "decision": "trust | reject | retire | defer",
  "checks": {
    "safety":   {"pass": true, "reason": "<one line>"},
    "correct":  {"pass": true, "reason": "<one line>"},
    "valuable": {"pass": true, "reason": "<one line>"},
    "scope":    {"pass": true, "reason": "<one line>"}
  },
  "reason": "<one line, at most 200 characters>",
  "gate": {"model": "<your model>", "effort": "<your reasoning effort>"}
}`

type verdictInput struct {
	Digest      string                 `json:"digest"`
	Decision    lesson.Decision        `json:"decision"`
	Checks      map[string]trust.Check `json:"checks"`
	Reason      string                 `json:"reason"`
	DuplicateOf string                 `json:"duplicate_of"`
	Retires     string                 `json:"retires"`
	Gate        struct {
		Model  string `json:"model"`
		Effort string `json:"effort"`
	} `json:"gate"`
}

type verdictResult struct {
	header

	ID      string `json:"id"`
	Status  string `json:"status"`
	Served  bool   `json:"served"`
	Retired string `json:"retired,omitempty"`
}

func newVerdictCommand(a *app) *cobra.Command {
	var input string

	command := &cobra.Command{
		Use:   "verdict <id>",
		Short: "Record the gate's decision on one lesson",
		Long: "verdict is for the independent gate agent only, after brief. It records\n" +
			"trust, reject, retire, or defer with the four checks. trust needs every\n" +
			"check to pass. Optional keys: duplicate_of (with reject) and retires (repeat\n" +
			"the lesson's own retires to confirm the replacement). Input on stdin:\n\n" + verdictTemplate,
		Example: "  yarbrain verdict l-7k3m9q2x < verdict.json",
		Args:    exactArgs(1, "id"),
		RunE: func(_ *cobra.Command, args []string) error {
			raw, err := a.readInput(input)
			if err != nil {
				return err
			}

			var in verdictInput
			if err := decodeStrict(raw, &in); err != nil {
				return err
			}

			if err := validateVerdict(in); err != nil {
				return err
			}

			result, err := a.recordVerdict(args[0], in)
			if err != nil {
				return err
			}

			return writeJSON(a.stdout, result)
		},
	}
	command.Flags().StringVar(&input, "input", "", "read the verdict JSON from this file instead of stdin")

	return command
}

func validateVerdict(in verdictInput) error {
	var problems []lesson.Problem

	if !slices.Contains([]lesson.Decision{lesson.DecideTrust, lesson.DecideReject, lesson.DecideRetire, lesson.DecideDefer}, in.Decision) {
		problems = append(problems, lesson.Problem{Field: "decision", Rule: "must be trust, reject, retire, or defer"})
	}

	fields := map[string]string{"reason": in.Reason}
	problems = append(problems, oneLine("reason", in.Reason, maxReason)...)
	problems = append(problems, oneLine("gate.model", in.Gate.Model, 80)...)

	allPass := true

	for _, name := range trust.RequiredChecks {
		check, ok := in.Checks[name]
		if !ok {
			problems = append(problems, lesson.Problem{Field: "checks." + name, Rule: "required"})

			continue
		}

		problems = append(problems, oneLine("checks."+name+".reason", check.Reason, maxReason)...)
		fields["reason"] += "\n" + check.Reason
		allPass = allPass && check.Pass
	}

	if len(in.Checks) != len(trust.RequiredChecks) {
		problems = append(problems, lesson.Problem{Field: "checks", Rule: "must hold exactly safety, correct, valuable, and scope"})
	}

	if in.Decision == lesson.DecideTrust && !allPass {
		problems = append(problems, lesson.Problem{Field: "decision", Rule: "trust needs every check to pass"})
	}

	if in.DuplicateOf != "" && (in.Decision != lesson.DecideReject || !lesson.IsID(in.DuplicateOf)) {
		problems = append(problems, lesson.Problem{Field: "duplicate_of", Rule: "only with reject, and must be a lesson id"})
	}

	if len(problems) > 0 {
		return problemsError(problems)
	}

	if failure := screen(fields); failure != nil {
		return failure
	}

	return nil
}

func (a *app) recordVerdict(id string, in verdictInput) (verdictResult, error) {
	v, _, err := a.openVault()
	if err != nil {
		return verdictResult{}, err
	}
	defer func() { _ = v.Close() }()

	session := a.session()

	var result verdictResult

	err = v.WithLock(func() error {
		entries, _, err := v.Lessons()
		if err != nil {
			return err
		}

		target, err := loadRecord(v, entries, id)
		if err != nil {
			return err
		}

		next, err := a.checkVerdict(v, entries, target, in, session)
		if err != nil {
			return err
		}

		event := trust.Event{
			At: a.now().UTC(), Kind: trust.Verdict, Session: session, Digest: in.Digest, Decision: in.Decision,
			Checks: in.Checks, Reason: in.Reason, DuplicateOf: in.DuplicateOf,
			Gate: &trust.Gate{Host: trust.Host(session), Model: in.Gate.Model, Effort: a.gateEffort(in.Gate.Effort)},
		}

		if err := a.applyVerdict(v, entries, target, next, event); err != nil {
			return err
		}

		result = verdictResult{
			header: okHeader(), ID: id, Status: string(next),
			Served: next == lesson.Trusted && trust.Credible(append(target.events, event), in.Digest),
		}

		if next == lesson.Trusted {
			result.Retired = target.lesson().Retires
		}

		return nil
	})

	return result, err
}

// gateEffort prefers the effort the host reports for this agent over the
// gate's own claim.
func (a *app) gateEffort(claimed string) string {
	if effort := a.getenv("CLAUDE_EFFORT"); effort != "" {
		return effort
	}

	return claimed
}

func (a *app) checkVerdict(v *vault.Vault, entries []vault.Entry, target record, in verdictInput, session string) (lesson.Status, error) {
	lease, err := v.Lease(target.lesson().ID)
	if err != nil {
		return "", err
	}

	if !lease.HeldBy(session, a.now()) {
		return "", badState(target, "this session holds no lease on it", fmt.Sprintf("Run `yarbrain brief %s` first; the lease lasts 15 minutes.", target.lesson().ID))
	}

	if in.Digest != target.digest() {
		return "", newError("changed_underneath", "the lesson changed after brief", fmt.Sprintf("Run `yarbrain brief %s` again and judge the current text.", target.lesson().ID)).
			with("id", target.lesson().ID).with("digest", target.digest())
	}

	next, ok := lesson.AfterVerdict(target.effective(), in.Decision)
	if !ok || !slices.Contains(allowedNow(target), in.Decision) {
		return "", badState(target, "the decision "+string(in.Decision)+" is not allowed", "Choose one of the allowed decisions.").
			with("allowed", allowedNow(target))
	}

	if in.Retires != target.lesson().Retires {
		return "", newError("bad_input", "retires must repeat the lesson's own retires value", "Copy retires from brief, or leave it out when the lesson has none.").
			with("field", "retires")
	}

	if in.DuplicateOf != "" {
		if _, ok := vault.Find(entries, in.DuplicateOf); !ok {
			return "", newError("not_found", "duplicate_of names no lesson: "+in.DuplicateOf, "Use an id from the similar lessons in brief.").
				with("field", "duplicate_of")
		}
	}

	if next == lesson.Trusted && target.lesson().Retires != "" {
		if err := checkRetireTarget(v, entries, target); err != nil {
			return "", err
		}
	}

	return next, nil
}

func checkRetireTarget(v *vault.Vault, entries []vault.Entry, target record) error {
	old, err := loadRecord(v, entries, target.lesson().Retires)
	if err != nil {
		return err
	}

	if !old.served() || !lesson.ScopesOverlap(old.lesson().Scope, target.lesson().Scope) {
		return badState(old, "it is not a served lesson of the same project", "Reject the replacement, or propose it without retires.")
	}

	return nil
}

// applyVerdict writes the new status, then the ledger event, then retires a
// replaced lesson. Lookups hide a lesson that a served lesson retires, so a
// crash between the steps never serves both.
func (a *app) applyVerdict(v *vault.Vault, entries []vault.Entry, target record, next lesson.Status, event trust.Event) error {
	keys := map[string]string{"status": string(next), "trusted_digest": ""}
	if next == lesson.Trusted {
		keys["trusted_digest"] = event.Digest
	}

	if err := v.SetLessonKeys(target.entry, keys); err != nil {
		return err
	}

	if err := v.AppendLedger(target.lesson().ID, event); err != nil {
		return err
	}

	if err := v.ClearLease(target.lesson().ID); err != nil {
		return err
	}

	if next != lesson.Trusted || target.lesson().Retires == "" {
		return nil
	}

	old, _ := vault.Find(entries, target.lesson().Retires)
	if err := v.SetLessonKeys(old, map[string]string{"status": string(lesson.Retired), "trusted_digest": "", "successor": target.lesson().ID}); err != nil {
		return err
	}

	return v.AppendLedger(old.Lesson.ID, trust.Event{At: event.At, Kind: trust.Replaced, Session: event.Session, Successor: target.lesson().ID})
}
