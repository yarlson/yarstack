---
name: cli-control
description: Verify terminal-hosted CLI or TUI behavior with repeatable local evidence. Use for output, prompts, interactive input, signals, hangs, resizing, startup, exit status, or terminal restoration.
---

# CLI Control

Gather reproducible evidence from terminal-hosted interfaces without changing production behavior.

## Workflow

1. Identify the command, working directory, terminal conditions, and observable contract.
2. Choose the smallest repeatable method: direct execution or repository tests for ordinary output and status; a terminal-capable harness for prompts, signals, resizing, hangs, or restoration.
3. Prefer an existing trusted harness. Otherwise use a temporary local harness outside the repository without adding dependencies.
4. Drive one action at a time and capture the smallest transcript or measurement that proves or disproves the behavior.
5. Check interrupts, cleanup, exit status, hangs, and terminal restoration only where relevant.
6. Clean up only sessions, processes, and artifacts created by this workflow.

Keep verification read-only unless the parent task separately authorizes implementation. Use `behavior-implement` and `test-design` for fixes or regression tests. Leave browser, desktop, Electron, and other graphical surfaces to `ui-control`.

Do not send credentials or destructive commands, retain sensitive transcripts, or keep one-off harness code in the repository without a product-test consumer.

Finish with the command and conditions, observed output and exit status, evidence, and any blocker or limitation.
