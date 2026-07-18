---
name: ui-control
description: Verify graphical browser, desktop, or Electron behavior through the actual interface. Use when screenshot, accessibility, console, network, trace, or browser-driven evidence is required; exclude terminal-hosted interfaces.
---

# UI Control

Gather reproducible evidence from the actual graphical surface without changing production behavior.

## Workflow

1. Identify the target surface, build or session, and observable contract.
2. Reuse the authorized running or authenticated surface when it is the target. Start a local app only when necessary and authorized.
3. Prefer existing repository harnesses and locate controls by stable roles, labels, identifiers, or application markers.
4. Perform only actions needed for the requested behavior or a concrete risk.
5. Capture the smallest relevant screenshot, accessibility output, console entry, network evidence, or trace.
6. Check loading, empty, error, keyboard, focus, responsive, and cleanup states only when material to the contract.
7. Clean up only servers, profiles, sessions, and artifacts created by this workflow.

Keep verification read-only unless the parent task separately authorizes implementation. Use `behavior-implement` for fixes, `phase-validate` for contract acceptance, and `test-gap-review` for evidence sufficiency.

Do not add automation dependencies for a one-off probe, retain sensitive evidence, use stale selectors, or terminate user-owned resources.

Finish with the surface tested, actions, observed result, retained evidence, and blockers or limitations.
