// Package project derives a project's identity from Git so that lessons
// follow a repository across clones, worktrees, and machines. Every git call
// is read-only and bounded by a timeout.
package project

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

const (
	// callTimeout bounds every ordinary git call.
	callTimeout = 2 * time.Second
	// rootBudget bounds the root-commit walk, which grows with history.
	// When it runs out, the origin identity alone matches.
	rootBudget = time.Second
)

// Identity names a project. A lesson matches the project when any of its
// scope identities equals one of IDs.
type Identity struct {
	// Origin is "origin:<host>/<path>" from remote.origin.url, without
	// scheme, credentials, port, or ".git". Empty without a usable remote.
	Origin string `json:"origin,omitempty"`
	// Root is "root:<12 hex>" from the lexicographically smallest root
	// commit. Empty without commits or when the walk exceeds its budget.
	Root string `json:"root,omitempty"`
	// Name is for display only and is never matched.
	Name string `json:"name,omitempty"`
	// Checkout is the main checkout directory; worktrees resolve to it.
	Checkout string `json:"checkout,omitempty"`
}

// IDs returns the identities a lesson scope may hold for this project.
func (i Identity) IDs() []string {
	var ids []string

	for _, id := range []string{i.Origin, i.Root} {
		if id != "" {
			ids = append(ids, id)
		}
	}

	return ids
}

// Known reports whether the project has at least one identity.
func (i Identity) Known() bool { return i.Origin != "" || i.Root != "" }

// Resolve derives the identity of the repository that contains dir. A
// directory outside any repository has an empty identity and no error.
func Resolve(ctx context.Context, dir string) (Identity, error) {
	commonDir, err := gitOutput(ctx, callTimeout, dir, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		if errors.Is(err, errNotRepository) {
			return Identity{}, nil
		}

		return Identity{}, err
	}

	identity := Identity{Checkout: checkoutOf(commonDir)}

	remote, err := gitOutput(ctx, callTimeout, dir, "config", "--get", "remote.origin.url")
	if err != nil && !errors.Is(err, errNotFound) {
		return Identity{}, err
	}

	if normalized := NormalizeOrigin(remote); normalized != "" {
		identity.Origin = "origin:" + normalized
	}

	identity.Root = rootIdentity(ctx, dir)
	identity.Name = displayName(identity)

	return identity, nil
}

func checkoutOf(commonDir string) string {
	if filepath.Base(commonDir) == ".git" {
		return filepath.Dir(commonDir)
	}

	return commonDir
}

func rootIdentity(ctx context.Context, dir string) string {
	roots, err := gitOutput(ctx, rootBudget, dir, "rev-list", "--max-parents=0", "HEAD")
	if err != nil {
		return ""
	}

	commits := strings.Fields(roots)
	if len(commits) == 0 {
		return ""
	}

	slices.Sort(commits)

	return "root:" + commits[0][:min(12, len(commits[0]))]
}

func displayName(identity Identity) string {
	if identity.Origin != "" {
		return identity.Origin[strings.LastIndex(identity.Origin, "/")+1:]
	}

	if identity.Checkout != "" {
		return filepath.Base(identity.Checkout)
	}

	return ""
}

// NormalizeOrigin reduces a remote URL to "<host>/<path>": no scheme,
// credentials, port, query, or ".git", in lowercase. It returns "" for local
// path remotes and for anything that would still carry an "@".
func NormalizeOrigin(remote string) string {
	remote = strings.TrimSpace(remote)
	if remote == "" {
		return ""
	}

	var host, path string

	switch {
	case strings.Contains(remote, "://"):
		parsed, err := url.Parse(remote)
		if err != nil || parsed.Scheme == "file" || parsed.Hostname() == "" {
			return ""
		}

		host, path = parsed.Hostname(), parsed.Path
	case strings.Contains(remote, ":") && !strings.HasPrefix(remote, "/"):
		// scp-like syntax: [user@]host:path
		hostPart, pathPart, _ := strings.Cut(remote, ":")
		host = hostPart[strings.LastIndex(hostPart, "@")+1:]
		path = pathPart
	default:
		return ""
	}

	path = strings.TrimSuffix(strings.Trim(path, "/"), ".git")
	normalized := strings.ToLower(host + "/" + path)

	if host == "" || path == "" || strings.Contains(normalized, "@") {
		return ""
	}

	return normalized
}

var (
	errNotRepository = errors.New("not a git repository")
	errNotFound      = errors.New("not found")
)

// gitOutput runs one read-only git command with a timeout. GIT_OPTIONAL_LOCKS
// keeps it from taking locks that would compete with the user's own git use.
func gitOutput(ctx context.Context, timeout time.Duration, dir string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	command := exec.CommandContext(ctx, "git", append([]string{"-C", dir}, args...)...)
	command.Env = append(os.Environ(), "GIT_OPTIONAL_LOCKS=0", "GIT_TERMINAL_PROMPT=0")

	var stdout, stderr bytes.Buffer

	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()
	if err == nil {
		return strings.TrimSpace(stdout.String()), nil
	}

	if ctx.Err() != nil {
		return "", fmt.Errorf("git %s timed out after %s: %w", args[0], timeout, ctx.Err())
	}

	if exit, ok := errors.AsType[*exec.ExitError](err); ok {
		switch {
		case exit.ExitCode() == 128 && strings.Contains(stderr.String(), "not a git repository"):
			return "", errNotRepository
		case exit.ExitCode() == 1:
			return "", errNotFound
		}
	}

	return "", fmt.Errorf("git %s: %w: %s", args[0], err, strings.TrimSpace(stderr.String()))
}
