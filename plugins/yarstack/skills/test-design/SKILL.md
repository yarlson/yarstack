---
name: test-design
description: Design or write deterministic tests before implementation. Use when choosing test levels, behavioral cases, fixtures, substitutes, or stable observable boundaries is non-trivial.
---

# Test Design

Own the prospective test contract and test code, not the production implementation or retrospective evidence audit.

## Workflow

1. Identify the condition, action, observable result, important side effects, and state preserved on failure.
2. Choose the lowest test level that proves the contract through a stable public or system boundary.
3. Select only cases that materially define the behavior: normal operation, important boundaries, meaningful failures, cleanup, and external effects.
4. Follow the repository's existing test structure, helpers, fixtures, and commands.
5. Write a deterministic test whose failure identifies the broken contract.
6. State what remains unverified when the behavior cannot be tested economically instead of adding a weak proxy assertion.

Use `behavior-implement` for red-green production changes, `test-gap-review` for retrospective verification sufficiency, `cli-control` for terminal evidence, and `ui-control` for graphical interface evidence.

Do not change production code or introduce broad test infrastructure unless the parent task already authorizes it. Report a testability constraint instead.

Finish with the test boundary, cases, test code when requested, command to run, and any residual gap.
