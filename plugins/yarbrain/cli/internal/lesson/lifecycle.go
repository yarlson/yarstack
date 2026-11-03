package lesson

// Changed is the effective status of a lesson whose stored status is trusted
// but whose author parts no longer match the trusted digest, for example
// after a hand edit. It is never stored, and a changed lesson is not served.
const Changed Status = "changed"

// Decision is a gate verdict.
type Decision string

const (
	DecideTrust  Decision = "trust"
	DecideReject Decision = "reject"
	DecideRetire Decision = "retire"
	DecideDefer  Decision = "defer"
)

// Effective returns the status a lesson has right now: Changed when a
// trusted lesson's content no longer matches its trusted digest.
func Effective(l Lesson) Status {
	if l.Status == Trusted && (l.TrustedDigest == "" || l.TrustedDigest != Digest(Normalize(l.Draft))) {
		return Changed
	}

	return l.Status
}

// verdicts lists the decisions a gate may make in each effective status and
// the status each one produces. A trusted, rejected, or retired lesson takes
// no verdict: trust ends with a doubt, rejection and retirement are final.
var verdicts = map[Status]map[Decision]Status{
	Proposed: {DecideTrust: Trusted, DecideReject: Rejected, DecideDefer: Deferred},
	Deferred: {DecideTrust: Trusted, DecideReject: Rejected, DecideRetire: Retired, DecideDefer: Deferred},
	Doubted:  {DecideTrust: Trusted, DecideRetire: Retired, DecideDefer: Deferred},
	Changed:  {DecideTrust: Trusted, DecideRetire: Retired, DecideDefer: Deferred},
}

// AwaitsGate reports whether a lesson in this effective status waits for a
// gate verdict.
func AwaitsGate(effective Status) bool {
	_, ok := verdicts[effective]

	return ok
}

// AllowedDecisions lists the verdicts a gate may give in this status.
func AllowedDecisions(effective Status) []Decision {
	var allowed []Decision

	for _, decision := range []Decision{DecideTrust, DecideReject, DecideRetire, DecideDefer} {
		if _, ok := verdicts[effective][decision]; ok {
			allowed = append(allowed, decision)
		}
	}

	return allowed
}

// AfterVerdict returns the status a verdict produces, or false when the
// verdict is not allowed in this status.
func AfterVerdict(effective Status, decision Decision) (Status, bool) {
	next, ok := verdicts[effective][decision]

	return next, ok
}
