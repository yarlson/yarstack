package project

import (
	"context"
	"errors"
	"regexp"
	"slices"
	"strings"
)

// Bounds on the facts one brief collects.
const (
	maxPaths    = 20
	maxTargets  = 10
	maxDirty    = 50
	maxHeadings = 40
)

// Facts are committed-tree facts a gate uses to check a lesson's claims
// without running repository code. A claim that holds only in uncommitted or
// unmerged work is not yet true for every clone.
type Facts struct {
	Head            string       `json:"head,omitempty"`
	DefaultBranch   string       `json:"default_branch,omitempty"`
	OnDefaultBranch bool         `json:"on_default_branch"`
	Dirty           []string     `json:"dirty_paths"`
	Paths           []PathFact   `json:"paths"`
	MakeTargets     []TargetFact `json:"make_targets"`
	Docs            []DocFact    `json:"docs"`
}

// PathFact says whether a path the lesson mentions exists at HEAD.
type PathFact struct {
	Path   string `json:"path"`
	Exists bool   `json:"exists"`
}

// TargetFact says whether a make target the lesson mentions exists in the
// tracked Makefile at HEAD.
type TargetFact struct {
	Target string `json:"target"`
	Exists bool   `json:"exists"`
}

// DocFact lists the headings of a tracked instruction or README file, so a
// gate can tell whether a lesson restates them.
type DocFact struct {
	File     string   `json:"file"`
	Headings []string `json:"headings"`
}

var (
	backticked = regexp.MustCompile("`([^`\\s]+)`")
	makeTarget = regexp.MustCompile(`\bmake\s+([A-Za-z0-9][A-Za-z0-9_.-]*)`)
	fileLike   = regexp.MustCompile(`(/|\.[A-Za-z0-9]{1,8}$)`)
)

// Gather collects the facts for the lesson text from the repository that
// contains dir. Outside a repository it returns empty facts.
func Gather(ctx context.Context, dir, text string) (Facts, error) {
	facts := Facts{Dirty: []string{}, Paths: []PathFact{}, MakeTargets: []TargetFact{}, Docs: []DocFact{}}

	head, err := gitOutput(ctx, callTimeout, dir, "rev-parse", "--verify", "--quiet", "HEAD")
	if errors.Is(err, errNotRepository) || errors.Is(err, errNotFound) {
		return facts, nil
	}

	if err != nil {
		return facts, err
	}

	facts.Head = head
	facts.DefaultBranch = defaultBranch(ctx, dir)

	if facts.DefaultBranch != "" {
		_, err := gitOutput(ctx, callTimeout, dir, "merge-base", "--is-ancestor", "HEAD", facts.DefaultBranch)
		facts.OnDefaultBranch = err == nil
	}

	if facts.Dirty, err = dirtyPaths(ctx, dir); err != nil {
		return facts, err
	}

	for _, path := range mentionedPaths(text) {
		_, err := gitOutput(ctx, callTimeout, dir, "cat-file", "-e", "HEAD:"+path)
		facts.Paths = append(facts.Paths, PathFact{Path: path, Exists: err == nil})
	}

	facts.MakeTargets = targetFacts(ctx, dir, text)
	facts.Docs = docFacts(ctx, dir)

	return facts, nil
}

func defaultBranch(ctx context.Context, dir string) string {
	if ref, err := gitOutput(ctx, callTimeout, dir, "symbolic-ref", "--quiet", "--short", "refs/remotes/origin/HEAD"); err == nil {
		return ref
	}

	for _, branch := range []string{"main", "master"} {
		if _, err := gitOutput(ctx, callTimeout, dir, "rev-parse", "--verify", "--quiet", "refs/heads/"+branch); err == nil {
			return branch
		}
	}

	return ""
}

func dirtyPaths(ctx context.Context, dir string) ([]string, error) {
	status, err := gitOutput(ctx, callTimeout, dir, "status", "--porcelain=v1", "--untracked-files=normal")
	if err != nil {
		return nil, err
	}

	paths := []string{}

	for line := range strings.Lines(status) {
		if len(paths) == maxDirty {
			break
		}

		if path := strings.TrimSpace(line); len(path) > 3 {
			paths = append(paths, strings.TrimSpace(path[3:]))
		}
	}

	return paths, nil
}

func mentionedPaths(text string) []string {
	var paths []string

	for _, match := range backticked.FindAllStringSubmatch(text, -1) {
		path := strings.TrimPrefix(match[1], "./")
		if strings.Contains(path, "://") || strings.HasPrefix(path, "-") || strings.HasPrefix(path, "/") || strings.HasPrefix(path, "~") || !fileLike.MatchString(path) {
			continue
		}

		if !slices.Contains(paths, path) {
			paths = append(paths, path)
		}

		if len(paths) == maxPaths {
			break
		}
	}

	return paths
}

func targetFacts(ctx context.Context, dir, text string) []TargetFact {
	facts := []TargetFact{}

	makefile, err := gitOutput(ctx, callTimeout, dir, "show", "HEAD:Makefile")
	if err != nil {
		makefile = ""
	}

	for _, match := range makeTarget.FindAllStringSubmatch(text, -1) {
		target := match[1]

		if slices.ContainsFunc(facts, func(f TargetFact) bool { return f.Target == target }) {
			continue
		}

		pattern := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(target) + `\s*:`)
		facts = append(facts, TargetFact{Target: target, Exists: pattern.MatchString(makefile)})

		if len(facts) == maxTargets {
			break
		}
	}

	return facts
}

func docFacts(ctx context.Context, dir string) []DocFact {
	docs := []DocFact{}

	for _, file := range []string{"README.md", "AGENTS.md", "CLAUDE.md"} {
		content, err := gitOutput(ctx, callTimeout, dir, "show", "HEAD:"+file)
		if err != nil {
			continue
		}

		var headings []string

		for line := range strings.Lines(content) {
			if strings.HasPrefix(line, "#") && len(headings) < maxHeadings {
				headings = append(headings, strings.TrimSpace(line))
			}
		}

		docs = append(docs, DocFact{File: file, Headings: headings})
	}

	return docs
}
