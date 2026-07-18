---
name: macos-swift-review
description: Supplement code review with macOS, Swift, and SwiftUI-specific correctness analysis. Use alongside code-review whenever the requested scope contains meaningful Swift production or test code, a Swift package, an Xcode project, or a macOS Swift app; skip only when no meaningful Swift code is in scope.
---

# macOS Swift Review

Run the required Swift and macOS-specific pass alongside general review and find platform failures it is unlikely to catch. Keep the review read-only.

## Workflow

1. Confirm whether the scope is a diff, flow, module, or full app and identify the actual distribution model.
2. Trace SwiftUI state ownership, identity, navigation, presentation, windows, restoration, and model lifetime through affected flows.
3. Check actor isolation, `@MainActor`, task ownership, cancellation, async streams, and blocking main-thread work.
4. Check platform-specific persistence, atomic file access, security-scoped resources, sandbox containers, and corrupted-state recovery where used.
5. Check entitlements, hardened runtime, keychain, helpers, XPC or IPC, Apple Events, URL schemes, permissions, signing, updates, and notarization only where the selected scope depends on them.
6. Distinguish source evidence from checks requiring identities, entitlements, release artifacts, or a live macOS environment. Do not build, sign, notarize, access the network, or control the UI without separate authority.

Use `code-review` for generic concerns, `security-review` for exploit analysis, `ui-control` for direct UI evidence, `test-gap-review` for coverage, and `rollout-readiness-review` for release readiness. Follow-ups remain read-only unless separately authorized.

Report severity, affected location, realistic platform scenario, impact, smallest safe correction, verification, and confidence limits. State directly when no platform-specific finding exists.

Finish when the selected Apple-platform surface is exhausted or a named environment limitation blocks further evidence.
