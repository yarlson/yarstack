# Yarbrain

Yarbrain is a reviewable Markdown second brain for Codex and Claude Code.

It keeps session evidence, current knowledge, and reusable procedures separate:

- episodes preserve what happened;
- notes hold the current reconciled understanding;
- skills hold approved repeatable procedures;
- inboxes keep extracted candidates reviewable;
- indexes are rebuildable from canonical Markdown.

Start with the wiki-initialize skill to choose and configure a vault. Yarbrain
does not create a vault or copy transcript contents during installation. Once a
vault is configured, lifecycle hooks load its compact index and enqueue session
locators. Use the remaining skills to decide what becomes durable knowledge.

The bundled helper uses only the Python standard library. Run its tests with:

    python3 -m unittest discover -s plugins/yarbrain/tests -p 'test_*.py'
