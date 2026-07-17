---
name: dependency-review
description: Review dependency and supply-chain changes. Use when work touches manifests, lockfiles, generated or vendored code, build scripts, CI tool installation, containers, plugin or MCP configuration, tool versions, registries, or network-fetched scripts.
---

# Dependency Review

Review only the dependency and toolchain surface changed by the current task.

## Workflow

1. Identify every dependency, generated-state, and toolchain change.
2. Confirm each change is necessary for current behavior.
3. Verify lockfiles and equivalent generated dependency state are consistent.
4. Follow the repository's established pinning policy.
5. Check for floating versions, unpinned Git sources, curl-to-shell installation, unexpected registries, and broad install scripts.
6. Check package, tool, and CI configuration for exposed credentials or unsafe permissions.
7. Run the repository's available dependency checks.
8. Report unrelated supply-chain risks without expanding scope.

## Guardrails

- Do not upgrade unrelated dependencies.
- Do not normalize an entire lockfile unless the requested change requires it.
- Do not add scanners, services, or update automation merely to perform the review.
- Do not trust committed generated code without checking its source and expected update path.

Finish when each changed dependency or tool is necessary, consistent, appropriately pinned, and verified by available checks.
