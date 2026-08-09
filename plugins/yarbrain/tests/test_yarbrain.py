from __future__ import annotations

import json
import os
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


PLUGIN_ROOT = Path(__file__).resolve().parents[1]
SCRIPT = PLUGIN_ROOT / "scripts" / "yarbrain.py"


class YarbrainCommandTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temporary_directory = tempfile.TemporaryDirectory()
        self.addCleanup(self.temporary_directory.cleanup)
        self.root = Path(self.temporary_directory.name)
        self.vault = self.root / "vault"
        self.config = self.root / "config" / "config.json"
        self.environment = {
            **os.environ,
            "YARBRAIN_CONFIG": str(self.config),
            "PLUGIN_ROOT": str(PLUGIN_ROOT),
        }

    def run_command(
        self,
        *arguments: str,
        input_text: str | None = None,
        environment: dict[str, str] | None = None,
    ) -> subprocess.CompletedProcess[str]:
        return subprocess.run(
            [sys.executable, str(SCRIPT), *arguments],
            input=input_text,
            text=True,
            capture_output=True,
            env={**self.environment, **(environment or {})},
            check=False,
        )

    def initialize(self) -> None:
        result = self.run_command("init", str(self.vault))
        self.assertEqual(result.returncode, 0, result.stderr)

    def write_note(self, relative_path: str, content: str) -> None:
        path = self.vault / relative_path
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(content, encoding="utf-8")

    def test_initialize_creates_contract_and_is_idempotent(self) -> None:
        self.initialize()
        second = self.run_command("init", str(self.vault))

        self.assertEqual(second.returncode, 0, second.stderr)
        self.assertTrue((self.vault / "INDEX.md").is_file())
        self.assertTrue((self.vault / "episodes").is_dir())
        self.assertTrue((self.vault / "inbox" / "sessions").is_dir())
        self.assertEqual(
            json.loads(self.config.read_text(encoding="utf-8"))["vault"],
            str(self.vault.resolve()),
        )
        lint = self.run_command("lint")
        self.assertEqual(lint.returncode, 0, lint.stdout + lint.stderr)

    def test_hook_is_noop_until_a_vault_is_configured(self) -> None:
        payload = json.dumps(
            {"hook_event_name": "SessionEnd", "session_id": "session-1"}
        )

        result = self.run_command("hook", input_text=payload)

        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertFalse(self.vault.exists())

    def test_hook_enqueues_only_session_locator_and_is_idempotent(self) -> None:
        self.initialize()
        payload = json.dumps(
            {
                "hook_event_name": "SessionEnd",
                "session_id": "session-1",
                "transcript_path": "/tmp/transcript.jsonl",
                "cwd": "/tmp/project",
                "secret": "must-not-be-copied",
            }
        )

        first = self.run_command("hook", input_text=payload)
        second = self.run_command("hook", input_text=payload)

        self.assertEqual(first.returncode, 0, first.stderr)
        self.assertEqual(second.returncode, 0, second.stderr)
        records = list((self.vault / "inbox" / "sessions").glob("pending-*.json"))
        self.assertEqual(len(records), 1)
        record_text = records[0].read_text(encoding="utf-8")
        record = json.loads(record_text)
        self.assertEqual(record["session_id"], "session-1")
        self.assertEqual(record["transcript_path"], "/tmp/transcript.jsonl")
        self.assertNotIn("must-not-be-copied", record_text)

    def test_precompact_enqueues_a_checkpoint_locator(self) -> None:
        self.initialize()
        payload = json.dumps(
            {
                "hook_event_name": "PreCompact",
                "session_id": "session-2",
                "transcript_path": "/tmp/transcript.jsonl",
                "trigger": "auto",
            }
        )

        result = self.run_command("hook", input_text=payload)

        self.assertEqual(result.returncode, 0, result.stderr)
        records = list((self.vault / "inbox" / "sessions").glob("pending-*.json"))
        self.assertEqual(len(records), 1)
        record = json.loads(records[0].read_text(encoding="utf-8"))
        self.assertEqual(record["event"], "PreCompact")
        self.assertEqual(record["trigger"], "auto")

    def test_claude_hook_records_the_host(self) -> None:
        self.initialize()
        payload = json.dumps(
            {"hook_event_name": "SessionEnd", "session_id": "claude-session"}
        )

        result = self.run_command(
            "hook",
            input_text=payload,
            environment={
                "PLUGIN_ROOT": "",
                "CLAUDE_PLUGIN_ROOT": str(PLUGIN_ROOT),
            },
        )

        self.assertEqual(result.returncode, 0, result.stderr)
        records = list((self.vault / "inbox" / "sessions").glob("pending-*.json"))
        self.assertEqual(len(records), 1)
        record = json.loads(records[0].read_text(encoding="utf-8"))
        self.assertEqual(record["host"], "claude-code")

    def test_hook_surfaces_an_invalid_config(self) -> None:
        self.config.parent.mkdir(parents=True)
        self.config.write_text("not json", encoding="utf-8")

        result = self.run_command(
            "hook",
            input_text=json.dumps(
                {"hook_event_name": "SessionEnd", "session_id": "session-1"}
            ),
        )

        self.assertEqual(result.returncode, 2)
        self.assertIn("Cannot read Yarbrain config", result.stderr)

    def test_session_start_loads_compact_index_and_pending_count(self) -> None:
        self.initialize()
        end_payload = json.dumps(
            {"hook_event_name": "SessionEnd", "session_id": "session-1"}
        )
        self.run_command("hook", input_text=end_payload)

        result = self.run_command(
            "hook", input_text=json.dumps({"hook_event_name": "SessionStart"})
        )

        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("Yarbrain index (reference only", result.stdout)
        self.assertIn("Pending session records: 1", result.stdout)

    def test_search_prefers_matching_canonical_note(self) -> None:
        self.initialize()
        self.write_note(
            "notes/concepts/refresh-token-rotation.md",
            "---\n"
            "id: refresh-token-rotation\n"
            "kind: semantic\n"
            "scope: repo:app\n"
            "status: current\n"
            "created: 2026-08-09\n"
            "updated: 2026-08-09\n"
            "verified_at: 2026-08-09\n"
            "volatility: medium\n"
            "aliases:\n"
            "  - token rotation\n"
            "sources:\n"
            "  - session://codex/app/one\n"
            "---\n\n"
            "# Refresh-token rotation\n\n"
            "## Current understanding\n\n"
            "Rotate refresh tokens after successful use.\n",
        )

        result = self.run_command("search", "refresh token rotation")

        self.assertEqual(result.returncode, 0, result.stderr)
        matches = json.loads(result.stdout)
        self.assertEqual(matches[0]["id"], "refresh-token-rotation")
        self.assertEqual(matches[0]["kind"], "semantic")
        lint = self.run_command("lint")
        self.assertEqual(lint.returncode, 0, lint.stdout + lint.stderr)

    def test_search_rejects_a_nonpositive_result_limit(self) -> None:
        self.initialize()

        result = self.run_command("search", "anything", "--limit", "0")

        self.assertEqual(result.returncode, 2)
        self.assertIn("Search limit must be at least 1", result.stderr)

    def test_lint_reports_duplicate_ids_and_broken_links(self) -> None:
        self.initialize()
        note = (
            "---\n"
            "id: duplicate\n"
            "kind: semantic\n"
            "scope: repo:app\n"
            "status: current\n"
            "created: 2026-08-09\n"
            "updated: 2026-08-09\n"
            "verified_at: 2026-08-09\n"
            "volatility: low\n"
            "aliases: []\n"
            "sources:\n"
            "  - session://codex/app/one\n"
            "---\n\n"
            "# Duplicate\n\n"
            "## Current understanding\n\n"
            "See [[missing-note]].\n"
        )
        self.write_note("notes/concepts/one.md", note)
        self.write_note("notes/concepts/two.md", note)

        result = self.run_command("lint")

        self.assertEqual(result.returncode, 1)
        issue_types = {issue["type"] for issue in json.loads(result.stdout)["issues"]}
        self.assertIn("duplicate-id", issue_types)
        self.assertIn("broken-link", issue_types)

    def test_lint_reports_missing_note_metadata(self) -> None:
        self.initialize()
        self.write_note(
            "notes/concepts/incomplete.md",
            "---\nid: incomplete\nkind: semantic\nsources:\n---\n\n# Incomplete\n",
        )

        result = self.run_command("lint")

        self.assertEqual(result.returncode, 1)
        messages = [issue["message"] for issue in json.loads(result.stdout)["issues"]]
        self.assertIn("Required frontmatter is missing: scope", messages)
        self.assertIn("Required frontmatter is missing: verified_at", messages)
        self.assertIn(
            "Required section is missing: ## Current understanding", messages
        )
        self.assertIn("Required frontmatter has no value: sources", messages)

    def test_lint_reports_empty_current_understanding(self) -> None:
        self.initialize()
        self.write_note(
            "notes/concepts/empty.md",
            "---\n"
            "id: empty\n"
            "kind: semantic\n"
            "scope: repo:app\n"
            "status: current\n"
            "created: 2026-08-09\n"
            "updated: 2026-08-09\n"
            "verified_at: 2026-08-09\n"
            "volatility: low\n"
            "aliases: []\n"
            "sources:\n"
            "  - session://codex/app/one\n"
            "---\n\n"
            "# Empty\n\n"
            "## Current understanding\n",
        )

        result = self.run_command("lint")

        self.assertEqual(result.returncode, 1)
        messages = [issue["message"] for issue in json.loads(result.stdout)["issues"]]
        self.assertIn(
            "Required section has no content: ## Current understanding", messages
        )

    def test_lint_reports_duplicate_episode_sources(self) -> None:
        self.initialize()
        episode = (
            "---\n"
            "id: {identifier}\n"
            "kind: episode\n"
            "agent: codex\n"
            "project: app\n"
            "started: 2026-08-09T10:00:00Z\n"
            "ended: 2026-08-09T11:00:00Z\n"
            "source: session://codex/app/session-1\n"
            "topics:\n"
            "  - auth\n"
            "---\n\n"
            "# Episode\n"
        )
        self.write_note("episodes/2026/08/one.md", episode.format(identifier="one"))
        self.write_note("episodes/2026/08/two.md", episode.format(identifier="two"))

        result = self.run_command("lint")

        self.assertEqual(result.returncode, 1)
        issue_types = {issue["type"] for issue in json.loads(result.stdout)["issues"]}
        self.assertIn("duplicate-episode-source", issue_types)

    def test_rebuild_index_uses_canonical_notes(self) -> None:
        self.initialize()
        self.write_note(
            "notes/concepts/retries.md",
            "---\n"
            "id: bounded-retries\n"
            "kind: semantic\n"
            "sources:\n"
            "  - session://codex/app/one\n"
            "---\n\n"
            "# Bounded retries\n",
        )

        result = self.run_command("rebuild-index")

        self.assertEqual(result.returncode, 0, result.stderr)
        index = (self.vault / "INDEX.md").read_text(encoding="utf-8")
        self.assertIn("[Bounded retries](notes/concepts/retries.md)", index)


if __name__ == "__main__":
    unittest.main()
