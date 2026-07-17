---
name: cli-control
description: Verify user-visible CLI or TUI behavior with repeatable local evidence. Use when work changes terminal output, prompts, interactive input, interrupts, hangs, resizing, startup behavior, or terminal demonstrations.
---

# CLI Control

Use a repeatable local harness instead of manual terminal poking.

## Workflow

1. Identify the command, working directory, and observable behavior under test.
2. Prefer repository-native integration tests, end-to-end tests, demo scripts, PTY helpers, or expect scripts.
3. If none exists, use a temporary local harness outside the repository.
4. Drive one action at a time and wait for concrete output before continuing.
5. Capture the smallest transcript or measurement that proves or disproves the behavior.
6. Check interrupts, cleanup, exit status, hangs, and terminal restoration where relevant.
7. Clean up temporary sessions, processes, and artifacts unless the user asks to keep them.
8. Turn findings into a scoped fix or explicit blocker.

## Guardrails

- Do not send credentials or destructive commands through the harness.
- Do not hardcode paths or assumptions from another repository.
- Do not keep one-off harness code in the repository unless it belongs in the product test suite.
- Do not add a testing dependency for a one-off probe without a current need.
- Do not treat screenshots or transcripts as sufficient when a stable automated test is practical.

Finish when the CLI or TUI behavior is verified with local evidence, fixed within scope, or blocked for a concrete reason.
