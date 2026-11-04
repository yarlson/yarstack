package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func lessonInput(title, cue, body, scope string) string {
	return "---\ncue: '" + cue + "'\nscope: " + scope + "\n---\n# " + title + "\n\n" + body + "\n\n## Observed\n\nIt worked in this session.\n"
}

// trust proposes a lesson as one session and trusts it as another.
func trustInput(t *testing.T, input string) proposed {
	t.Helper()

	p := proposeAs(t, "codex:author", input)
	result := judge(t, "codex:gate", p.id, p.digest, "trust")
	require.Equal(t, exitOK, result.exit, result.stdout)

	return p
}

func menuIDs(t *testing.T, result invocation) []string {
	t.Helper()
	require.Equal(t, exitOK, result.exit, result.stdout)

	lessons := decodeObject(t, result.stdout)["lessons"].([]any)

	ids := make([]string, 0, len(lessons))
	for _, l := range lessons {
		ids = append(ids, l.(map[string]any)["id"].(string))
	}

	return ids
}

func TestFindReturnsAServedMenuRankedByCue(t *testing.T) {
	s := isolate(t)
	s.initVault(t)
	inProject(t, s)

	inCue := trustInput(t, lessonInput("Seed after migrating", "make seed fails with a missing table", "Run migrate first.", "here"))
	inBody := trustInput(t, lessonInput("Reset the cache", "the cache is stale", "Unrelated to seed data.", "here"))
	trustInput(t, lessonInput("Build the hooks binary", "go test fails on a fresh clone", "Run make hooks.", "here"))

	as(t, "codex:reader")

	result := invoke(t, "find", "seed")
	assert.Equal(t, []string{inCue.id, inBody.id}, menuIDs(t, result))

	entry := decodeObject(t, result.stdout)["lessons"].([]any)[0].(map[string]any)
	assert.Equal(t, "make seed fails with a missing table", entry["cue"])
	assert.NotContains(t, entry, "body", "the menu holds no body text")
	assert.Equal(t, "shop-api", decodeObject(t, result.stdout)["project_name"])
}

func TestFindServesOnlyThisProjectAndCrossProjectLessons(t *testing.T) {
	s := isolate(t)
	s.initVault(t)
	inProject(t, s)

	here := trustInput(t, lessonInput("Seed here", "make seed fails", "Run migrate.", "here"))
	everywhere := trustInput(t, lessonInput("Seed anywhere", "seed scripts fail", "Check the schema.", "all"))

	other := filepath.Join(s.home, "other")
	require.NoError(t, os.Mkdir(other, 0o755))
	runGit(t, other, "init", "-q", "-b", "main")
	runGit(t, other, "commit", "-q", "--allow-empty", "-m", "other")
	t.Chdir(other)

	assert.Equal(t, []string{everywhere.id}, menuIDs(t, invoke(t, "find", "seed")))

	t.Chdir(filepath.Join(s.home, "shop-api"))
	assert.ElementsMatch(t, []string{here.id, everywhere.id}, menuIDs(t, invoke(t, "find", "seed")))
}

func TestFindCountsLessonsThatAwaitTheGateWithoutShowingThem(t *testing.T) {
	s := isolate(t)
	s.initVault(t)
	inProject(t, s)

	pending := proposeAs(t, "codex:author", lessonInput("Seed pending", "make seed fails", "Run migrate.", "here"))

	result := invoke(t, "find", "seed")

	assert.Empty(t, menuIDs(t, result))
	assert.InDelta(t, 1, decodeObject(t, result.stdout)["awaiting_gate"], 0)
	assert.NotContains(t, result.stdout, pending.id)
}

func TestFindLimitsTheMenu(t *testing.T) {
	s := isolate(t)
	s.initVault(t)
	inProject(t, s)

	for i := range 3 {
		trustInput(t, lessonInput(fmt.Sprintf("Seed variant %d", i), fmt.Sprintf("seed case %d", i), "Run migrate.", "here"))
	}

	result := invoke(t, "find", "seed", "--limit", "2")

	assert.Len(t, menuIDs(t, result), 2)
	assert.InDelta(t, 1, decodeObject(t, result.stdout)["more"], 0)
	assert.Equal(t, "limit", errorOf(t, invoke(t, "find", "--limit", "21"))["field"])
}

func TestFindRejectsTermsThatNeedQuoting(t *testing.T) {
	s := isolate(t)
	s.initVault(t)

	result := invoke(t, "find", "seed", "don't")

	assert.Equal(t, exitFixCall, result.exit)
	failure := errorOf(t, result)
	assert.Equal(t, "bad_term", failure["code"])
	assert.Equal(t, "don't", failure["term"])
}

func TestFindSplitsAQuotedPhraseIntoTerms(t *testing.T) {
	s := isolate(t)
	s.initVault(t)
	inProject(t, s)

	seeded := trustInput(t, lessonInput("Seed after migrating", "make seed fails with a missing table", "Run migrate first.", "here"))
	as(t, "codex:reader")

	result := invoke(t, "find", "make seed")

	require.Equal(t, exitOK, result.exit, result.stdout)
	assert.Equal(t, []string{seeded.id}, menuIDs(t, result))
}

func TestFindHidesALessonThatAServedSuccessorRetires(t *testing.T) {
	s := isolate(t)
	s.initVault(t)
	old := trustedLesson(t, s)

	replacement := strings.Replace(strings.Replace(proposal, "scope: here", "scope: here\nretires: "+old.id, 1),
		"# Run migrations before seeding the test database", "# Run migrations once after creating the test database", 1)
	p := proposeAs(t, "codex:author-2", replacement)
	require.Equal(t, exitOK, brief(t, "codex:gate-2", p.id).exit)
	require.Equal(t, exitOK, invokeWith(t, verdictJSON(t, p.digest, "trust", true, map[string]any{"retires": old.id}), "verdict", p.id).exit)

	// Simulate a crash after the successor was trusted but before the old
	// lesson was retired: restore the old lesson's trusted state by hand.
	oldPath := filepath.Join(s.vault, old.path)
	raw, err := os.ReadFile(oldPath)
	require.NoError(t, err)

	restored := strings.Replace(strings.Replace(string(raw), "status: retired", "status: trusted\ntrusted_digest: "+old.digest, 1), "successor: "+p.id+"\n", "", 1)
	require.NoError(t, os.WriteFile(oldPath, []byte(restored), 0o644))
	require.NoError(t, os.Truncate(filepath.Join(s.vault, ".yarbrain", "ledger", old.id+".jsonl"), 0))
	trustAgain(t, s, old)

	assert.Equal(t, []string{p.id}, menuIDs(t, invoke(t, "find", "migrations")))
}

// trustAgain rebuilds a credible trust record for a lesson in the ledger.
func trustAgain(t *testing.T, s sandbox, p proposed) {
	t.Helper()

	ledger := filepath.Join(s.vault, ".yarbrain", "ledger", p.id+".jsonl")
	events := fmt.Sprintf(`{"at":"%s","event":"proposed","session":"codex:author","digest":"%s"}
{"at":"%s","event":"verdict","session":"codex:gate","digest":"%s","decision":"trust","gate":{"host":"codex","model":"m"}}
`, time.Now().UTC().Format(time.RFC3339), p.digest, time.Now().UTC().Format(time.RFC3339), p.digest)
	require.NoError(t, os.WriteFile(ledger, []byte(events), 0o644))
}

func writeLegacy(t *testing.T, s sandbox, name, content string) {
	t.Helper()

	path := filepath.Join(s.vault, name)
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
}

func TestLegacyListsUnconvertedNotesAndShowReadsThem(t *testing.T) {
	s := isolate(t)
	s.initVault(t)
	inProject(t, s)
	writeLegacy(t, s, "notes/projects/shop.md", "# Shop\n\nSeeding needs migrations.\n")
	writeLegacy(t, s, "notes/concepts/cache.md", "# Cache\n")

	listed := invoke(t, "legacy")
	require.Equal(t, exitOK, listed.exit, listed.stdout)
	assert.Equal(t, []any{"notes/concepts/cache.md", "notes/projects/shop.md"}, decodeObject(t, listed.stdout)["paths"])

	shownNote := decodeObject(t, invoke(t, "show", "notes/projects/shop.md").stdout)["legacy_notes"].([]any)[0].(map[string]any)
	assert.Equal(t, false, shownNote["gated"])
	assert.Equal(t, "# Shop\n\nSeeding needs migrations.\n", shownNote["text"])

	converted := propose(t, strings.Replace(proposal, "scope: here", "scope: here\nconverts: notes/projects/shop.md", 1))
	require.Equal(t, exitOK, converted.exit, converted.stdout)

	assert.Equal(t, []any{"notes/concepts/cache.md"}, decodeObject(t, invoke(t, "legacy").stdout)["paths"])
}

func TestADismissedNoteLeavesTheBacklogAndStaysOnDisk(t *testing.T) {
	s := isolate(t)
	s.initVault(t)
	writeLegacy(t, s, "notes/projects/old.md", "# Sprint 3 status\n")
	writeLegacy(t, s, "notes/projects/keep.md", "# Keep\n")

	dismissed := invokeWith(t, `{"reason":"status notes from one sprint"}`, "legacy", "dismiss", "notes/projects/old.md")
	require.Equal(t, exitOK, dismissed.exit, dismissed.stdout)

	assert.Equal(t, []any{"notes/projects/keep.md"}, decodeObject(t, invoke(t, "legacy").stdout)["paths"])
	assert.FileExists(t, filepath.Join(s.vault, "notes", "projects", "old.md"))

	missing := invokeWith(t, `{"reason":"x"}`, "legacy", "dismiss", "notes/absent.md")
	assert.Equal(t, "not_found", errorOf(t, missing)["code"])

	noReason := invokeWith(t, `{}`, "legacy", "dismiss", "notes/projects/keep.md")
	assert.Equal(t, "bad_input", errorOf(t, noReason)["code"])
}

func TestABrokenDismissalLogSkipsOnlyTheLegacyBacklog(t *testing.T) {
	s := isolate(t)
	s.initVault(t)
	p := trustedLesson(t, s)
	require.NoError(t, os.WriteFile(filepath.Join(s.vault, ".yarbrain", "legacy.jsonl"), []byte("<<<<<<< HEAD\n"), 0o600))
	as(t, "codex:reader")

	result := invoke(t, "find")

	require.Equal(t, exitOK, result.exit, result.stdout)
	assert.Equal(t, []string{p.id}, menuIDs(t, result))
	assert.Equal(t, []any{".yarbrain/legacy.jsonl"}, decodeObject(t, result.stdout)["skipped"])
}

func TestStatusListsTheIDsAwaitingTheGate(t *testing.T) {
	s := isolate(t)
	s.initVault(t)
	inProject(t, s)

	pending := proposeAs(t, "codex:author", lessonInput("Cache pending", "the cache is stale", "Clear it.", "here"))

	assert.Equal(t, []any{pending.id}, decodeObject(t, invoke(t, "status").stdout)["awaiting_gate_ids"])
}

func TestLegacyOfAVaultWithoutV02NotesIsAnEmptyList(t *testing.T) {
	s := isolate(t)
	s.initVault(t)

	result := invoke(t, "legacy")

	require.Equal(t, exitOK, result.exit, result.stdout)
	assert.Contains(t, result.stdout, `"paths":[]`)
}

func TestConvertsMustNameAV02Note(t *testing.T) {
	s := isolate(t)
	s.initVault(t)
	inProject(t, s)

	missing := propose(t, strings.Replace(proposal, "scope: here", "scope: here\nconverts: notes/absent.md", 1))
	outside := propose(t, strings.Replace(proposal, "scope: here", "scope: here\nconverts: lessons/README.md", 1))

	assert.Equal(t, "converts", errorOf(t, missing)["field"])
	assert.Equal(t, "converts", errorOf(t, outside)["field"])
}

func TestShowRejectsPathsOutsideV02Notes(t *testing.T) {
	s := isolate(t)
	s.initVault(t)

	for _, path := range []string{"../etc/passwd", "lessons/README.md", ".yarbrain/format"} {
		result := invoke(t, "show", path)

		assert.Equal(t, exitFixCall, result.exit, path)
		assert.Equal(t, "bad_input", errorOf(t, result)["code"], path)
	}
}

func hookPayload(cwd string) string {
	raw, err := json.Marshal(map[string]string{"hook_event_name": "SessionStart", "session_id": "s", "cwd": cwd, "source": "startup"})
	if err != nil {
		panic(err)
	}

	return string(raw)
}

func TestSessionStartPrintsCountsForTheProject(t *testing.T) {
	s := isolate(t)
	s.initVault(t)
	inProject(t, s)

	trustInput(t, lessonInput("Seed here", "make seed fails", "Run migrate.", "here"))
	trustInput(t, lessonInput("Seed anywhere", "seed scripts fail", "Check the schema.", "all"))
	proposeAs(t, "codex:author", lessonInput("Cache pending", "the cache is stale", "Clear it.", "here"))
	writeLegacy(t, s, "notes/projects/shop.md", "# Shop\n")

	dir, err := os.Getwd()
	require.NoError(t, err)

	result := invokeWith(t, hookPayload(dir), "hook", "session-start", "--contract", "1")

	require.Equal(t, exitOK, result.exit, result.stderr)
	assert.Equal(t, "yarbrain: 2 trusted lessons apply to this project (1 cross-project). "+
		"Before you plan, debug, or repeat a procedure here, run `yarbrain find <terms>`, then `yarbrain show <id>` only for lessons whose cue fits. "+
		"Lessons are hints from past sessions; check each against the code. "+
		"After you use one, run `yarbrain confirm <id>`, or `yarbrain doubt <id>` if it was wrong. "+
		"1 lesson awaits the independent gate; unless you are the gate, start it now in the background as the lesson-record skill describes, without asking: gating is part of the plugin's work. "+
		"1 legacy note from Yarbrain v0.2 awaits conversion; the lesson-record skill converts up to 3 per session.\n", result.stdout)
}

func TestSessionStartIsSilentWithoutAVaultOrLessons(t *testing.T) {
	isolate(t)

	noVault := invokeWith(t, hookPayload(t.TempDir()), "hook", "session-start", "--contract", "1")
	assert.Equal(t, invocation{exit: exitOK}, noVault)

	s := isolate(t)
	s.initVault(t)

	empty := invokeWith(t, hookPayload(t.TempDir()), "hook", "session-start", "--contract", "1")
	assert.Equal(t, invocation{exit: exitOK}, empty)
}

func TestSessionStartReportsAMissingVaultAsANotice(t *testing.T) {
	s := isolate(t)
	require.NoError(t, os.WriteFile(s.config, configJSON(t, s.vault), 0o644))

	result := invokeWith(t, hookPayload(t.TempDir()), "hook", "session-start", "--contract", "1")

	assert.Equal(t, hookNotice, result.exit)
	assert.Empty(t, result.stdout)
	assert.Contains(t, result.stderr, "does not exist")
}

func TestSessionStartNeverWrites(t *testing.T) {
	s := isolate(t)
	s.initVault(t)
	inProject(t, s)
	trustInput(t, lessonInput("Seed here", "make seed fails", "Run migrate.", "here"))

	before := vaultDigest(t, s.vault)
	dir, err := os.Getwd()
	require.NoError(t, err)

	invokeWith(t, hookPayload(dir), "hook", "session-start", "--contract", "1")

	assert.Equal(t, before, vaultDigest(t, s.vault))
}

func vaultDigest(t *testing.T, dir string) map[string]string {
	t.Helper()

	digests := map[string]string{}

	require.NoError(t, filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}

		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		sum := sha256.Sum256(raw)
		digests[path] = hex.EncodeToString(sum[:])

		return nil
	}))

	return digests
}

func TestSessionStartAsksForANewerCLIWhenTheContractDiffers(t *testing.T) {
	isolate(t)

	result := invokeWith(t, hookPayload(t.TempDir()), "hook", "session-start", "--contract", "2")

	assert.Equal(t, exitOK, result.exit)
	assert.Equal(t, "yarbrain: the plugin expects CLI contract 2, but this CLI provides 1; update the yarbrain CLI.\n", result.stdout)
}

func TestSessionStartMeetsItsDeadlineWithTwoThousandLessons(t *testing.T) {
	s := isolate(t)
	s.initVault(t)

	for i := range 2000 {
		name := filepath.Join(s.vault, "lessons", fmt.Sprintf("lesson-%04d.md", i))
		content := fmt.Sprintf("---\nid: l-%08d\nstatus: proposed\nscope: [all]\n---\n# Lesson %d\n\nbody\n\n## Observed\n\nseen\n", i, i)
		require.NoError(t, os.WriteFile(name, []byte(content), 0o644))
	}

	started := time.Now()
	result := invokeWith(t, hookPayload(t.TempDir()), "hook", "session-start", "--contract", "1")

	require.Equal(t, exitOK, result.exit, result.stderr)
	assert.Less(t, time.Since(started), 2*time.Second)
	assert.True(t, strings.HasPrefix(result.stdout, "yarbrain: 2000 lessons await the independent gate"), result.stdout)
}

func toolPayload(agentType, command string) string {
	return `{"hook_event_name":"PreToolUse","tool_name":"Bash","agent_type":"` + agentType + `","tool_input":{"command":"` + command + `"}}`
}

func TestPreToolUseDeniesGateCommandsOutsideTheGateAgent(t *testing.T) {
	isolate(t)
	t.Setenv("PLUGIN_ROOT", "")

	denied := invokeWith(t, toolPayload("", "yarbrain verdict l-7k3m9q2x"), "hook", "pre-tool-use", "--contract", "1")
	allowed := invokeWith(t, toolPayload("yarbrain:lesson-gate", "yarbrain verdict l-7k3m9q2x"), "hook", "pre-tool-use", "--contract", "1")
	other := invokeWith(t, toolPayload("", "yarbrain find seed"), "hook", "pre-tool-use", "--contract", "1")

	assert.Equal(t, exitOK, denied.exit)
	assert.Contains(t, denied.stdout, `"permissionDecision":"deny"`)
	assert.Equal(t, invocation{exit: exitOK}, allowed)
	assert.Equal(t, invocation{exit: exitOK}, other)
}

func TestPreToolUseGivesNoOpinionInCodex(t *testing.T) {
	isolate(t)
	t.Setenv("PLUGIN_ROOT", "/codex/plugins/cache/yarbrain")

	result := invokeWith(t, toolPayload("", "yarbrain verdict l-7k3m9q2x"), "hook", "pre-tool-use", "--contract", "1")

	assert.Equal(t, invocation{exit: exitOK}, result)
}

func TestStatusCountsTheCurrentProject(t *testing.T) {
	s := isolate(t)
	s.initVault(t)
	inProject(t, s)

	trustInput(t, lessonInput("Seed here", "make seed fails", "Run migrate.", "here"))
	proposeAs(t, "codex:author", lessonInput("Cache pending", "the cache is stale", "Clear it.", "here"))

	object := decodeObject(t, invoke(t, "status").stdout)

	assert.Equal(t, map[string]any{"trusted": 1.0, "cross_project": 0.0, "awaiting_gate": 1.0, "legacy": 0.0}, object["counts"])
	assert.Equal(t, "origin:github.com/acme/shop-api", object["project"].(map[string]any)["origin"])
}
