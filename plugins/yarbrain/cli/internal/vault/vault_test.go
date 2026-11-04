package vault

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func openVault(t *testing.T, dir string) *Vault {
	t.Helper()

	v, err := Open(dir)
	require.NoError(t, err)
	t.Cleanup(func() { _ = v.Close() })

	return v
}

func newVault(t *testing.T) *Vault {
	t.Helper()

	dir := filepath.Join(t.TempDir(), "vault")
	_, err := Init(dir)
	require.NoError(t, err)

	return openVault(t, dir)
}

// treeDigest hashes every path and file content under dir.
func treeDigest(t *testing.T, dir string) map[string]string {
	t.Helper()

	digests := map[string]string{}
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(dir, path)
		if entry.IsDir() {
			digests[rel] = "dir"

			return nil
		}

		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		sum := sha256.Sum256(raw)
		digests[rel] = hex.EncodeToString(sum[:])

		return nil
	})
	require.NoError(t, err)

	return digests
}

func TestInitCreatesTheLayoutOnce(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "vault")

	first, err := Init(dir)
	require.NoError(t, err)
	assert.False(t, first.Adopted)
	assert.False(t, first.Legacy)
	assert.ElementsMatch(t, []string{
		"lessons/", ".yarbrain/", ".yarbrain/ledger/", ".yarbrain/leases/",
		".yarbrain/format", ".yarbrain/.gitignore", "lessons/README.md",
	}, first.Created)

	second, err := Init(dir)
	require.NoError(t, err)
	assert.True(t, second.Adopted)
	assert.Empty(t, second.Created)
}

func TestInitAdoptsAV02VaultWithoutTouchingIt(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"episodes/2026/01", "notes/projects", "inbox/sessions", "skills/x", "archive", "reports"} {
		require.NoError(t, os.MkdirAll(filepath.Join(dir, name), 0o755))
	}

	files := map[string]string{
		"INDEX.md":                 "# Yarbrain Index\n",
		"notes/projects/a.md":      "---\nid: a\n---\n# A\n",
		"episodes/2026/01/e.md":    "---\nid: e\n---\n",
		"inbox/sessions/p.json":    "{}\n",
		"skills/x/SKILL.md":        "---\nname: x\n---\n",
		"reports/wiki-health.md":   "# Health\n",
		"lessons-not-a-dir-yet.md": "keep me\n",
	}
	for name, content := range files {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644))
	}

	before := treeDigest(t, dir)

	result, err := Init(dir)
	require.NoError(t, err)
	assert.True(t, result.Adopted)
	assert.True(t, result.Legacy)

	after := treeDigest(t, dir)
	for path, digest := range before {
		assert.Equal(t, digest, after[path], "v0.2 path %s changed", path)
	}
}

func TestOpenReportsAMissingVault(t *testing.T) {
	_, err := Open(filepath.Join(t.TempDir(), "absent"))

	var missing *MissingError
	assert.ErrorAs(t, err, &missing)
}

func TestOpenRefusesANewerFormat(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, ".yarbrain"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".yarbrain", "format"), []byte("2\n"), 0o644))

	_, err := Open(dir)

	var format *FormatError
	require.ErrorAs(t, err, &format)
	assert.Equal(t, 2, format.Found)
}

func TestAVaultWithoutMarkerReadsAsFormatOne(t *testing.T) {
	v := openVault(t, t.TempDir())

	format, err := v.Format()
	require.NoError(t, err)
	assert.Equal(t, 1, format)
}

func TestRootContainsEveryPath(t *testing.T) {
	outer := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(outer, "secret.txt"), []byte("outside"), 0o644))

	dir := filepath.Join(outer, "vault")
	_, err := Init(dir)
	require.NoError(t, err)

	v := openVault(t, dir)

	_, err = v.ReadFile("../secret.txt")
	require.Error(t, err)

	require.Error(t, v.CreateFile("../escape.md", []byte("x")))
	assert.NoFileExists(t, filepath.Join(outer, "escape.md"))
}

func TestRootDoesNotFollowASymlinkOutOfTheVault(t *testing.T) {
	outer := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(outer, "secret.txt"), []byte("outside"), 0o644))

	dir := filepath.Join(outer, "vault")
	_, err := Init(dir)
	require.NoError(t, err)

	if err := os.Symlink(filepath.Join(outer, "secret.txt"), filepath.Join(dir, "lessons", "link.md")); err != nil {
		t.Skipf("this OS user cannot create symlinks: %v", err)
	}

	_, err = openVault(t, dir).ReadFile("lessons/link.md")
	require.Error(t, err, "a symlink that points outside the vault must not be followed")
}

func TestCreateFileHasExactlyOneWinner(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "vault")
	_, err := Init(dir)
	require.NoError(t, err)

	const writers = 8

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		winners []int
	)

	for writer := range writers {
		wg.Go(func() {
			v, err := Open(dir)
			if err != nil {
				return
			}
			defer func() { _ = v.Close() }()

			if v.CreateFile("lessons/race.md", []byte{byte('a' + writer)}) == nil {
				mu.Lock()

				winners = append(winners, writer)

				mu.Unlock()
			}
		})
	}

	wg.Wait()

	require.Len(t, winners, 1)

	raw, err := os.ReadFile(filepath.Join(dir, "lessons", "race.md"))
	require.NoError(t, err)
	assert.Equal(t, []byte{byte('a' + winners[0])}, raw)
}

func TestLockIsExclusiveWhileHeld(t *testing.T) {
	v := newVault(t)
	other := openVault(t, v.Path)

	err := v.WithLock(func() error {
		return other.WithLock(func() error { return nil })
	})

	require.ErrorIs(t, err, ErrLocked)
	require.NoError(t, v.WithLock(func() error { return nil }), "the lock must be released after fn returns")
}

func TestStaleLockIsRecovered(t *testing.T) {
	v := newVault(t)
	require.NoError(t, v.CreateFile(lockFile, []byte("pid 1\n")))

	old := time.Now().Add(-time.Minute)
	require.NoError(t, os.Chtimes(filepath.Join(v.Path, ".yarbrain", "lock"), old, old))

	ran := false

	require.NoError(t, v.WithLock(func() error {
		ran = true

		return nil
	}))
	assert.True(t, ran)
	assert.NoFileExists(t, filepath.Join(v.Path, ".yarbrain", "lock"))
}

func TestFreshLockIsNotStolen(t *testing.T) {
	v := newVault(t)
	require.NoError(t, v.CreateFile(lockFile, []byte("pid 1\n")))

	err := v.WithLock(func() error { return nil })

	require.ErrorIs(t, err, ErrLocked)
	assert.FileExists(t, filepath.Join(v.Path, ".yarbrain", "lock"))
}

func TestFirstWriteCreatesTheFormatMarker(t *testing.T) {
	dir := t.TempDir()
	v := openVault(t, dir)

	require.NoError(t, v.WithLock(func() error { return nil }))

	raw, err := os.ReadFile(filepath.Join(dir, ".yarbrain", "format"))
	require.NoError(t, err)
	assert.Equal(t, "1\n", string(raw))
}

func TestReplaceFileChecksTheRevision(t *testing.T) {
	v := newVault(t)
	require.NoError(t, v.CreateFile("lessons/a.md", []byte("one")))

	err := v.ReplaceFile("lessons/a.md", Rev([]byte("stale")), []byte("two"))

	var changed *ChangedError
	require.ErrorAs(t, err, &changed)

	require.NoError(t, v.ReplaceFile("lessons/a.md", Rev([]byte("one")), []byte("two")))

	raw, err := os.ReadFile(filepath.Join(v.Path, "lessons", "a.md"))
	require.NoError(t, err)
	assert.Equal(t, "two", string(raw))

	entries, err := os.ReadDir(filepath.Join(v.Path, "lessons"))
	require.NoError(t, err)

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}

	sort.Strings(names)
	assert.Equal(t, []string{"README.md", "a.md"}, names, "no temporary file may remain")
}
