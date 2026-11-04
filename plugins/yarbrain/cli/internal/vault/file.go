package vault

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"

	yaml "go.yaml.in/yaml/v3"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
)

// Frontmatter keys. The CLI owns ownedKeys and may rewrite their lines; the
// author keys and any key a person adds are never rewritten.
var ownedKeys = []string{"id", "status", "created", "project_name", "trusted_digest", "successor"}

const observedHeading = "## Observed"

var bom = []byte("\xef\xbb\xbf")

// Rev identifies the exact bytes of a file. Writers compare it before they
// replace a file, so a concurrent edit is never lost.
func Rev(raw []byte) string {
	sum := sha256.Sum256(raw)

	return hex.EncodeToString(sum[:])[:12]
}

type document struct {
	lines []string
	// closing is the index of the closing "---" line. Frontmatter YAML line n
	// is file line n, because the opening "---" is line 0.
	closing int
	mapping *yaml.Node
}

func splitDocument(text, name string) (*document, error) {
	lines := strings.Split(text, "\n")
	if lines[0] != "---" {
		return nil, &BadFileError{Path: name, Reason: "the file must start with a --- frontmatter line"}
	}

	closing := slices.Index(lines[1:], "---") + 1
	if closing == 0 {
		return nil, &BadFileError{Path: name, Reason: "the frontmatter has no closing --- line"}
	}

	var root yaml.Node
	if err := yaml.Unmarshal([]byte(strings.Join(lines[1:closing], "\n")), &root); err != nil {
		return nil, &BadFileError{Path: name, Reason: "the frontmatter is not valid YAML: " + err.Error()}
	}

	mapping := &yaml.Node{Kind: yaml.MappingNode}
	if len(root.Content) > 0 {
		mapping = root.Content[0]
	}

	if mapping.Kind != yaml.MappingNode {
		return nil, &BadFileError{Path: name, Reason: "the frontmatter must be a mapping of keys to values"}
	}

	seen := map[string]bool{}

	for i := 0; i < len(mapping.Content); i += 2 {
		key := mapping.Content[i].Value
		if seen[key] {
			return nil, &BadFileError{Path: name, Reason: fmt.Sprintf("the frontmatter key %q appears twice", key)}
		}

		seen[key] = true
	}

	return &document{lines: lines, closing: closing, mapping: mapping}, nil
}

func (d *document) pair(key string) (index int, keyNode, valueNode *yaml.Node) {
	for i := 0; i < len(d.mapping.Content); i += 2 {
		if d.mapping.Content[i].Value == key {
			return i / 2, d.mapping.Content[i], d.mapping.Content[i+1]
		}
	}

	return -1, nil, nil
}

// Parse reads a lesson file. It accepts a byte-order mark and CRLF line
// endings, which people's editors may add.
func Parse(raw []byte, name string) (lesson.Lesson, error) {
	text := strings.ReplaceAll(string(bytes.TrimPrefix(raw, bom)), "\r\n", "\n")

	doc, err := splitDocument(text, name)
	if err != nil {
		return lesson.Lesson{}, err
	}

	var l lesson.Lesson

	scalars := map[string]*string{
		"id": &l.ID, "cue": &l.Cue, "retires": &l.Retires, "converts": &l.Converts,
		"created": &l.Created, "project_name": &l.ProjectName,
		"trusted_digest": &l.TrustedDigest, "successor": &l.Successor,
	}
	for key, target := range scalars {
		if _, _, value := doc.pair(key); value != nil {
			if value.Kind != yaml.ScalarNode {
				return lesson.Lesson{}, &BadFileError{Path: name, Reason: fmt.Sprintf("the key %q must hold one value", key)}
			}

			*target = value.Value
		}
	}

	if _, _, value := doc.pair("status"); value != nil {
		l.Status = lesson.Status(value.Value)
	}

	if _, _, value := doc.pair("scope"); value != nil {
		scope, err := scopeValues(value, name)
		if err != nil {
			return lesson.Lesson{}, err
		}

		l.Scope = scope
	}

	title, body, observed, err := splitBody(doc.lines[doc.closing+1:], name)
	if err != nil {
		return lesson.Lesson{}, err
	}

	l.Title, l.Body, l.Observed = title, body, observed

	return l, nil
}

// proposalKeys are the frontmatter keys a proposal may carry.
var proposalKeys = []string{"cue", "scope", "retires", "converts"}

// InputError lists every problem in a proposal.
type InputError struct{ Problems []lesson.Problem }

func (e *InputError) Error() string {
	return fmt.Sprintf("the input has %d problem(s), first: %s %s", len(e.Problems), e.Problems[0].Field, e.Problems[0].Rule)
}

// ParseProposal reads a lesson proposal in the file format. Its scope is the
// keyword "here" or "all", which the caller resolves. It rejects keys the CLI
// owns and keys it does not know.
func ParseProposal(raw []byte) (lesson.Draft, string, error) {
	text := strings.ReplaceAll(string(bytes.TrimPrefix(raw, bom)), "\r\n", "\n")

	doc, err := splitDocument(text, "input")
	if err != nil {
		return lesson.Draft{}, "", &InputError{Problems: []lesson.Problem{{Field: "input", Rule: inputReason(err)}}}
	}

	problems := proposalKeyProblems(doc.mapping)

	parsed, err := Parse([]byte(text), "input")
	if err != nil {
		return lesson.Draft{}, "", &InputError{Problems: append(problems, lesson.Problem{Field: "input", Rule: inputReason(err)})}
	}

	scope := strings.Join(parsed.Scope, ",")
	if scope != "here" && scope != lesson.AllScope {
		problems = append(problems, lesson.Problem{Field: "scope", Rule: "must be here or all"})
	}

	if len(problems) > 0 {
		return lesson.Draft{}, "", &InputError{Problems: problems}
	}

	return parsed.Draft, scope, nil
}

func proposalKeyProblems(mapping *yaml.Node) []lesson.Problem {
	var problems []lesson.Problem

	for i := 0; i < len(mapping.Content); i += 2 {
		key := mapping.Content[i].Value

		switch {
		case slices.Contains(ownedKeys, key):
			problems = append(problems, lesson.Problem{Field: key, Rule: "set by the CLI; remove this key"})
		case !slices.Contains(proposalKeys, key):
			problems = append(problems, lesson.Problem{Field: key, Rule: "unknown key; allowed keys are cue, scope, retires, converts"})
		}
	}

	return problems
}

func inputReason(err error) string {
	if bad, ok := errors.AsType[*BadFileError](err); ok {
		return bad.Reason
	}

	return err.Error()
}

func scopeValues(node *yaml.Node, name string) ([]string, error) {
	if node.Kind == yaml.ScalarNode {
		return []string{node.Value}, nil
	}

	if node.Kind != yaml.SequenceNode {
		return nil, &BadFileError{Path: name, Reason: "scope must be a list such as [all]"}
	}

	scope := make([]string, 0, len(node.Content))
	for _, item := range node.Content {
		if item.Kind != yaml.ScalarNode {
			return nil, &BadFileError{Path: name, Reason: "every scope entry must be a single value"}
		}

		scope = append(scope, item.Value)
	}

	return scope, nil
}

func splitBody(lines []string, name string) (title, body, observed string, err error) {
	start := 0
	for start < len(lines) && strings.TrimSpace(lines[start]) == "" {
		start++
	}

	if start == len(lines) || !strings.HasPrefix(lines[start], "# ") {
		return "", "", "", &BadFileError{Path: name, Reason: "the body must start with a '# Title' line"}
	}

	title = strings.TrimPrefix(lines[start], "# ")
	rest := lines[start+1:]

	section := func(lines []string) string { return strings.Trim(strings.Join(lines, "\n"), "\n") }

	split := slices.IndexFunc(rest, func(line string) bool { return strings.TrimRight(line, " \t") == observedHeading })
	if split < 0 {
		return title, section(rest), "", nil
	}

	return title, section(rest[:split]), section(rest[split+1:]), nil
}

// Render writes a lesson in the canonical file format.
func Render(l lesson.Lesson) ([]byte, error) {
	mapping := &yaml.Node{Kind: yaml.MappingNode}
	add := func(key string, value *yaml.Node) {
		mapping.Content = append(mapping.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: key}, value)
	}

	text := func(key, value string) {
		if value != "" {
			add(key, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value})
		}
	}

	text("id", l.ID)
	text("status", string(l.Status))
	text("cue", l.Cue)

	scope := &yaml.Node{Kind: yaml.SequenceNode, Style: yaml.FlowStyle}
	for _, identity := range l.Scope {
		scope.Content = append(scope.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: identity})
	}

	add("scope", scope)
	text("retires", l.Retires)
	text("converts", l.Converts)
	text("project_name", l.ProjectName)
	text("created", l.Created)
	text("trusted_digest", l.TrustedDigest)
	text("successor", l.Successor)

	var frontmatter bytes.Buffer

	encoder := yaml.NewEncoder(&frontmatter)
	encoder.SetIndent(2)

	if err := encoder.Encode(mapping); err != nil {
		return nil, fmt.Errorf("encode frontmatter: %w", err)
	}

	if err := encoder.Close(); err != nil {
		return nil, fmt.Errorf("encode frontmatter: %w", err)
	}

	var out bytes.Buffer

	fmt.Fprintf(&out, "---\n%s---\n# %s\n\n", frontmatter.String(), l.Title)

	if l.Body != "" {
		fmt.Fprintf(&out, "%s\n\n", l.Body)
	}

	fmt.Fprintf(&out, "%s\n\n%s\n", observedHeading, l.Observed)

	return out.Bytes(), nil
}

// SetOwned rewrites only the lines of CLI-owned keys and keeps every other
// byte. An empty value removes the key. It refuses files it cannot change
// safely: a byte-order mark, CRLF line endings, or an owned key that spans
// several lines.
func SetOwned(raw []byte, name string, values map[string]string) ([]byte, error) {
	if bytes.HasPrefix(raw, bom) {
		return nil, &BadFileError{Path: name, Reason: "the file starts with a byte-order mark; save it as UTF-8 without one"}
	}

	if bytes.Contains(raw, []byte("\r\n")) {
		return nil, &BadFileError{Path: name, Reason: "the file uses CRLF line endings; save it with LF line endings"}
	}

	text := string(raw)

	for _, key := range ownedKeys {
		value, ok := values[key]
		if !ok {
			continue
		}

		updated, err := setLine(text, name, key, value)
		if err != nil {
			return nil, err
		}

		text = updated
	}

	if _, err := Parse([]byte(text), name); err != nil {
		return nil, fmt.Errorf("rewritten %s does not parse: %w", name, err)
	}

	return []byte(text), nil
}

func setLine(text, name, key, value string) (string, error) {
	doc, err := splitDocument(text, name)
	if err != nil {
		return "", err
	}

	rendered, err := renderLine(key, value)
	if err != nil {
		return "", err
	}

	index, keyNode, valueNode := doc.pair(key)
	if keyNode == nil {
		if value == "" {
			return text, nil
		}

		lines := slices.Insert(doc.lines, doc.closing, rendered)

		return strings.Join(lines, "\n"), nil
	}

	// YAML line numbers match the file's "\n" lines only when every key starts
	// its own line and no value holds a Unicode line break.
	if keyNode.Line >= doc.closing || !strings.HasPrefix(doc.lines[keyNode.Line], key+":") {
		return "", &BadFileError{Path: name, Reason: fmt.Sprintf("the key %q is not at the start of its own line; write one key per line", key)}
	}

	if !doc.singleLine(index, keyNode, valueNode) {
		return "", &BadFileError{Path: name, Reason: fmt.Sprintf("the key %q spans several lines; put it on one line", key)}
	}

	lines := doc.lines
	if value == "" {
		lines = slices.Delete(lines, keyNode.Line, keyNode.Line+1)
	} else {
		lines[keyNode.Line] = rendered
	}

	return strings.Join(lines, "\n"), nil
}

// singleLine reports whether a value occupies only its key's line: every line
// before the next key may hold only blank space or a comment.
func (d *document) singleLine(index int, keyNode, valueNode *yaml.Node) bool {
	if valueNode.Kind != yaml.ScalarNode || valueNode.Line != keyNode.Line {
		return false
	}

	if valueNode.Style&(yaml.LiteralStyle|yaml.FoldedStyle) != 0 {
		return false
	}

	next := d.closing
	if pairStart := 2 * (index + 1); pairStart < len(d.mapping.Content) {
		next = d.mapping.Content[pairStart].Line
	}

	if next <= keyNode.Line || next > d.closing {
		return false
	}

	for _, between := range d.lines[keyNode.Line+1 : next] {
		trimmed := strings.TrimSpace(between)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			return false
		}
	}

	return true
}

func renderLine(key, value string) (string, error) {
	node := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: key},
		{Kind: yaml.ScalarNode, Tag: "!!str", Value: value},
	}}

	encoded, err := yaml.Marshal(node)
	if err != nil {
		return "", fmt.Errorf("encode %s: %w", key, err)
	}

	line := strings.TrimSuffix(string(encoded), "\n")
	if strings.Contains(line, "\n") {
		return "", fmt.Errorf("value for %s does not fit on one line", key)
	}

	return line, nil
}
