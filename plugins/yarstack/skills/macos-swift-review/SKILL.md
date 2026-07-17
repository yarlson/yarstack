---
name: macos-swift-review
description: Review a macOS Swift or SwiftUI codebase for platform-specific correctness and lifecycle risks. Use for full-app technical audits involving SwiftUI state, structured concurrency, persistence, sandboxing, entitlements, system integration, responsiveness, packaging, signing, or notarization.
---

# macOS Swift Review

Find concrete macOS and Swift failures that generic code review can miss. Keep the review read-only.

## Workflow

1. Map app and scene entry points, windows, commands, modules, state owners, services, persistence, networking, platform integrations, entitlements, packaging, and tests.
2. Trace important flows from user action through SwiftUI state, async work, persistence or networking, error presentation, cancellation, and cleanup.
3. Check `@State`, bindings, observable models, environment state, identity, navigation, sheets, popovers, windows, restoration, and view-model lifetime for incorrect ownership or unstable rendering.
4. Check actor isolation, `@MainActor`, unstructured tasks, task lifetime, cancellation, async streams, debounce or throttle behavior, and blocking work on the main thread.
5. Check SwiftData or Core Data migrations, atomic file writes, security-scoped bookmarks, container paths, corrupted-state recovery, caches, and backup behavior where used.
6. Check URLSession timeout and cancellation behavior, retry idempotency, offline handling, response validation, decoding errors, and user-visible recovery.
7. Check sandbox entitlements, hardened runtime, keychain access, helpers, XPC or IPC, Apple Events, URL schemes, accessibility permissions, file import or export, signing, updates, and notarization against the app's actual distribution model.
8. Check retain cycles, subscriptions, observers, timers, image or list behavior, main-thread work, accessibility, and release-only failure paths.
9. Inspect tests for state restoration, async cancellation, migrations, permissions, sandbox behavior, and platform integrations. Request focused UI, security, or test-gap review when those surfaces need deeper evidence.

## Finding Standard

Report only findings with an affected location, realistic failure scenario, impact, smallest safe fix, and verification method. State when an entitlement, signing identity, runtime environment, or release artifact is unavailable and limits confidence.
