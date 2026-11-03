package project

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()

	command := exec.CommandContext(t.Context(), "git", append([]string{"-C", dir}, args...)...)
	command.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.com",
		"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.com",
		"GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1")

	out, err := command.CombinedOutput()
	require.NoError(t, err, "git %v: %s", args, out)

	return strings.TrimSpace(string(out))
}

func newRepo(t *testing.T) string {
	t.Helper()

	dir := filepath.Join(t.TempDir(), "shop-api")
	require.NoError(t, os.Mkdir(dir, 0o755))
	git(t, dir, "init", "-q", "-b", "main")

	return dir
}

func commit(t *testing.T, dir string) string {
	t.Helper()
	git(t, dir, "commit", "-q", "--allow-empty", "-m", "c")

	return git(t, dir, "rev-parse", "HEAD")
}

func TestNormalizeOrigin(t *testing.T) {
	cases := map[string]string{
		"git@github.com:Acme/Shop-API.git":                       "github.com/acme/shop-api",
		"https://github.com/acme/shop-api.git":                   "github.com/acme/shop-api",
		"https://x-access-token:secret@github.com/acme/shop-api": "github.com/acme/shop-api",
		"ssh://git@gitlab.example.com:2222/team/app.git":         "gitlab.example.com/team/app",
		"https://github.com/acme/shop-api/":                      "github.com/acme/shop-api",
		"/srv/git/shop-api.git":                                  "",
		"file:///srv/git/shop-api.git":                           "",
		"":                                                       "",
	}
	for remote, want := range cases {
		assert.Equal(t, want, NormalizeOrigin(remote), "remote %q", remote)
	}
}

func TestOutsideARepositoryTheIdentityIsUnknown(t *testing.T) {
	identity, err := Resolve(t.Context(), t.TempDir())

	require.NoError(t, err)
	assert.False(t, identity.Known())
	assert.Empty(t, identity.IDs())
}

func TestARepositoryWithoutCommitsHasNoIdentity(t *testing.T) {
	identity, err := Resolve(t.Context(), newRepo(t))

	require.NoError(t, err)
	assert.False(t, identity.Known())
	assert.Equal(t, "shop-api", identity.Name)
}

func TestRootIdentityWithoutARemote(t *testing.T) {
	dir := newRepo(t)
	head := commit(t, dir)

	identity, err := Resolve(t.Context(), dir)

	require.NoError(t, err)
	assert.Equal(t, []string{"root:" + head[:12]}, identity.IDs())
	assert.Equal(t, "shop-api", identity.Name)
}

func TestOriginAddedLaterKeepsTheRootIdentity(t *testing.T) {
	dir := newRepo(t)
	head := commit(t, dir)

	before, err := Resolve(t.Context(), dir)
	require.NoError(t, err)

	git(t, dir, "remote", "add", "origin", "git@github.com:acme/shop-api.git")

	after, err := Resolve(t.Context(), dir)
	require.NoError(t, err)

	assert.Equal(t, []string{"origin:github.com/acme/shop-api", "root:" + head[:12]}, after.IDs())
	assert.Contains(t, after.IDs(), before.Root, "an old lesson scoped to the root identity must still match")
}

func TestWorktreesShareTheIdentity(t *testing.T) {
	dir := newRepo(t)
	commit(t, dir)
	git(t, dir, "remote", "add", "origin", "https://github.com/acme/shop-api.git")

	worktree := filepath.Join(t.TempDir(), "feature")
	git(t, dir, "worktree", "add", "-q", "-b", "feature", worktree)

	main, err := Resolve(t.Context(), dir)
	require.NoError(t, err)

	linked, err := Resolve(t.Context(), filepath.Join(worktree))
	require.NoError(t, err)

	assert.Equal(t, main.IDs(), linked.IDs())

	mainCheckout, err := filepath.EvalSymlinks(dir)
	require.NoError(t, err)

	linkedCheckout, err := filepath.EvalSymlinks(linked.Checkout)
	require.NoError(t, err)
	assert.Equal(t, mainCheckout, linkedCheckout, "a worktree resolves to its main checkout")
}

func TestSubdirectoriesShareTheIdentity(t *testing.T) {
	dir := newRepo(t)
	commit(t, dir)

	sub := filepath.Join(dir, "a", "b")
	require.NoError(t, os.MkdirAll(sub, 0o755))

	top, err := Resolve(t.Context(), dir)
	require.NoError(t, err)

	nested, err := Resolve(t.Context(), sub)
	require.NoError(t, err)
	assert.Equal(t, top.IDs(), nested.IDs())
}

func TestMultipleRootsUseTheSmallest(t *testing.T) {
	dir := newRepo(t)
	first := commit(t, dir)

	git(t, dir, "checkout", "-q", "--orphan", "other")
	second := commit(t, dir)
	git(t, dir, "checkout", "-q", "main")
	git(t, dir, "merge", "-q", "--allow-unrelated-histories", "-m", "merge", "other")

	identity, err := Resolve(t.Context(), dir)
	require.NoError(t, err)

	smallest := min(first, second)
	assert.Equal(t, "root:"+smallest[:12], identity.Root)
}

func TestCancelledContextFailsInsteadOfGuessing(t *testing.T) {
	dir := newRepo(t)
	commit(t, dir)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	_, err := Resolve(ctx, dir)
	assert.Error(t, err)
}
