---
name: ui-control
description: Verify browser, desktop, Electron, or other local UI behavior through the actual interface. Use when work requires screenshot, accessibility, console, network, trace, or browser-driven evidence of changed behavior.
---

# UI Control

Gather evidence from the actual UI surface using repository-native tooling when available.

## Workflow

1. Identify the UI surface and observable behavior under test.
2. Start the application with its documented local command.
3. Prefer existing Playwright, Cypress, Storybook, browser, desktop, or Electron harnesses.
4. Locate pages and controls by stable roles, labels, test identifiers, or application markers.
5. Capture before-and-after evidence only when it proves the changed behavior.
6. Inspect console, network, traces, screenshots, or accessibility output only as needed.
7. Check relevant loading, empty, error, keyboard, focus, responsive, and cleanup behavior.
8. Clean up servers, debug sessions, temporary profiles, and artifacts unless the user asks to keep them.

## Guardrails

- Do not add browser automation dependencies merely for a one-off probe.
- Do not reuse stale selectors or coordinates after navigation or layout changes.
- Avoid coordinate clicks unless a fresh screenshot immediately establishes the target.
- Do not retain sensitive screenshots, traces, HTTP bodies, profiles, or memory dumps without a concrete need.
- Do not hardcode ports, selectors, or scripts copied from another repository.

Finish when UI behavior is verified with direct evidence, fixed within scope, or blocked for a concrete reason.
