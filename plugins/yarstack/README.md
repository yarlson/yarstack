# Yarstack

Yarstack is Yar Kravtsov's portable engineering-workflow and standards plugin
for Codex and Claude Code.

Reusable skills live under `skills/`. The read-only reviewer agent that
delegating skills use in Claude Code lives under `agents/`. The canonical
guidance text is `agent-guidance/engineering-standards/standards.md`, and the
README beside it gives the evidence for each rule. A safe idempotent installer
lives at `scripts/install-engineering-standards.sh`. Run the installer with
`--print` to inspect the guidance without changing global files.

Documentation, contribution guidance, and security policy are maintained at
https://github.com/yarlson/yarstack.
