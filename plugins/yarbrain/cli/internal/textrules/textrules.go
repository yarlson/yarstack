// Package textrules holds the deterministic checks on agent-written text.
// They find only known patterns: a text without findings is not proven safe,
// so callers must never report it as clean.
package textrules

import (
	"regexp"
	"strings"
)

// Category groups patterns by the error code a caller reports.
type Category string

const (
	// Secret patterns: credentials and personal data.
	Secret Category = "secret"
	// Risky patterns: text that addresses agents, fetch-and-execute, and
	// destructive commands.
	Risky Category = "risky"
)

// Finding names the field and the pattern that matched. It never holds the
// matched text, so reporting it cannot leak a secret again.
type Finding struct {
	Field    string   `json:"field"`
	Pattern  string   `json:"pattern"`
	Category Category `json:"-"`
}

type rule struct {
	name     string
	category Category
	pattern  *regexp.Regexp
	// allow rejects a match that the pattern over-reports.
	allow func(match string) bool
}

var rules = []rule{
	{name: "private key block", category: Secret, pattern: regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----`)},
	{name: "GitHub token", category: Secret, pattern: regexp.MustCompile(`\b(gh[pousr]_[A-Za-z0-9]{30,}|github_pat_[A-Za-z0-9_]{22,})\b`)},
	{name: "AWS access key", category: Secret, pattern: regexp.MustCompile(`\b(AKIA|ASIA)[0-9A-Z]{16}\b`)},
	{name: "Slack token", category: Secret, pattern: regexp.MustCompile(`\bxox[abposr]-[A-Za-z0-9-]{10,}`)},
	{name: "API secret key", category: Secret, pattern: regexp.MustCompile(`\b(sk-[A-Za-z0-9_-]{20,}|(sk|rk)_(live|test)_[0-9A-Za-z]{16,}|AIza[0-9A-Za-z_-]{35})`)},
	{name: "JSON web token", category: Secret, pattern: regexp.MustCompile(`\beyJ[A-Za-z0-9_-]{10,}\.eyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}`)},
	{
		name:     "credential assignment",
		category: Secret,
		// A placeholder such as ${TOKEN} or <token> is not a value.
		pattern: regexp.MustCompile(`(?i)\b(password|passwd|secret|token|api[_-]?key|access[_-]?key)\s*[:=]\s*['"]?[^\s'"$<{][^\s'"]{7,}`),
	},
	{name: "credentials in a URL", category: Secret, pattern: regexp.MustCompile(`[a-z][a-z0-9+.-]*://[^/\s:@]+:[^/\s@]+@`)},
	{
		name:     "email address",
		category: Secret,
		pattern:  regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b`),
		// git@host is SSH remote syntax, not a person.
		allow: func(match string) bool { return strings.HasPrefix(match, "git@") },
	},
	{name: "home directory path", category: Secret, pattern: regexp.MustCompile(`(?i)(/Users/[^/\s]+|/home/[^/\s]+|[A-Z]:\\Users\\[^\\\s]+)`)},

	{name: "yarbrain command", category: Risky, pattern: regexp.MustCompile(`(?i)\byarbrain\s+(init|status|find|show|propose|brief|verdict|confirm|doubt|legacy|check|hook)\b`)},
	{name: "lesson id", category: Risky, pattern: regexp.MustCompile(`\bl-[0-9a-hjkmnp-tv-z]{8}\b`)},
	{name: "instruction override", category: Risky, pattern: regexp.MustCompile(`(?i)\b(ignore|disregard|forget)\s+(all\s+|any\s+)?(previous|prior|above|earlier)\s+(instructions|rules|guidance)`)},
	{
		name:     "fetch and execute",
		category: Risky,
		pattern:  regexp.MustCompile(`(?i)((curl|wget|iwr|invoke-webrequest|irm|invoke-restmethod)\b[^|\n]*\|\s*(sudo\s+)?(sh|bash|zsh|python[0-9.]*|perl|ruby|node|iex|pwsh|powershell)\b|\b(sh|bash|zsh)\s+<\(\s*(curl|wget))`),
	},
	{name: "destructive command", category: Risky, pattern: regexp.MustCompile(`(?i)(\brm\s+-[a-z]*(rf|fr)[a-z]*\s+(/|~|\$HOME|\*)(\s|$)|\bgit\s+push\b[^\n]*\s(-f|--force)\b|\bgit\s+reset\s+--hard\b|\bgit\s+clean\s+-[a-z]*f|\bchmod\s+-R\s+777\b|\bmkfs\.|\bdd\s+[^\n]*of=/dev/)`)},
}

var urlPattern = regexp.MustCompile(`(?i)(\bhttps?://|\bwww\.)`)

// LineFields are fields that must not carry a URL: they appear in every
// lookup menu, where a link would invite an agent to follow it.
var LineFields = map[string]bool{"title": true, "cue": true}

// Check returns every rule that text breaks, labeled with field.
func Check(field, text string) []Finding {
	var findings []Finding

	for _, r := range rules {
		if matches(r, text) {
			findings = append(findings, Finding{Field: field, Pattern: r.name, Category: r.category})
		}
	}

	if LineFields[field] && urlPattern.MatchString(text) {
		findings = append(findings, Finding{Field: field, Pattern: "URL", Category: Risky})
	}

	return findings
}

func matches(r rule, text string) bool {
	for _, match := range r.pattern.FindAllString(text, -1) {
		if r.allow == nil || !r.allow(match) {
			return true
		}
	}

	return false
}
