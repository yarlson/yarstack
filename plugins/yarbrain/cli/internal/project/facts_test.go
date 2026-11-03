package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func repoWithFiles(t *testing.T, files map[string]string) string {
	t.Helper()

	dir := newRepo(t)
	for name, content := range files {
		path := filepath.Join(dir, name)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	}

	git(t, dir, "add", "-A")
	git(t, dir, "commit", "-q", "-m", "files")

	return dir
}

func TestGatherChecksClaimsAgainstTheCommittedTree(t *testing.T) {
	dir := repoWithFiles(t, map[string]string{
		"Makefile":         ".PHONY: seed\nmigrate:\n\tgo run ./cmd/migrate\nseed: migrate\n",
		"scripts/seed.sql": "insert into users values (1);\n",
		"README.md":        "# Shop API\n\n## Testing\n\nRun make test.\n",
	})
	require.NoError(t, os.WriteFile(filepath.Join(dir, "draft.sql"), []byte("uncommitted"), 0o644))

	text := "Run `make migrate` before `make seed`; `scripts/seed.sql` assumes the schema. " +
		"`draft.sql` and `docs/missing.md` and `make deploy` are not committed. `--force` is a flag."

	facts, err := Gather(t.Context(), dir, text)
	require.NoError(t, err)

	assert.NotEmpty(t, facts.Head)
	assert.Equal(t, "main", facts.DefaultBranch)
	assert.True(t, facts.OnDefaultBranch)
	assert.Equal(t, []string{"draft.sql"}, facts.Dirty)
	assert.Equal(t, []PathFact{
		{Path: "scripts/seed.sql", Exists: true},
		{Path: "draft.sql", Exists: false},
		{Path: "docs/missing.md", Exists: false},
	}, facts.Paths)
	assert.Equal(t, []TargetFact{
		{Target: "migrate", Exists: true},
		{Target: "seed", Exists: true},
		{Target: "deploy", Exists: false},
	}, facts.MakeTargets)
	assert.Equal(t, []DocFact{{File: "README.md", Headings: []string{"# Shop API", "## Testing"}}}, facts.Docs)
}

func TestGatherSeesABranchThatIsNotMerged(t *testing.T) {
	dir := repoWithFiles(t, map[string]string{"README.md": "# A\n"})
	git(t, dir, "checkout", "-q", "-b", "feature")
	git(t, dir, "commit", "-q", "--allow-empty", "-m", "feature work")

	facts, err := Gather(t.Context(), dir, "")
	require.NoError(t, err)

	assert.Equal(t, "main", facts.DefaultBranch)
	assert.False(t, facts.OnDefaultBranch, "work only on a feature branch is not true for every clone yet")
}

func TestGatherOutsideARepositoryIsEmpty(t *testing.T) {
	facts, err := Gather(t.Context(), t.TempDir(), "`a/b.go` and make x")

	require.NoError(t, err)
	assert.Empty(t, facts.Head)
	assert.Empty(t, facts.Paths)
	assert.Empty(t, facts.MakeTargets)
}
