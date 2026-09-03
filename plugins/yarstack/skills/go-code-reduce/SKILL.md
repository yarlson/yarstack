---
name: go-code-reduce
description: Find and, when authorized, apply maintainable Go code reductions. Use when a user asks to reduce Go lines, boilerplate, duplication, or obsolete idioms without changing behavior; not for general implementation or style review.
---

# Reduce Go Code

Reduce authored maintenance burden, not raw repository lines. Prefer deleting mechanics through standard Go constructs over compressing domain logic or introducing project-specific abstractions.

## Workflow

1. Read the repository guidance, module Go version, closest source and tests, generated-file markers, and affected public consumers. Keep an assessment read-only unless edits are authorized.
2. State the behavior that must remain unchanged, including errors and wrapped identities, nil versus empty values, ordering, serialization, cancellation, cleanup, exported APIs, build variants, and relevant performance bounds.
3. Rank evidenced candidates in this order:

   - delete dead, obsolete, forwarding, compatibility, or unnecessarily exported code when consumers permit it;
   - replace hand-written mechanics with version-supported language, built-in, or standard-library operations;
   - shrink producer APIs and define small interfaces at actual consumers;
   - consolidate an algorithm only when multiple real copies differ solely by type or repeated traversal;
   - generate code only when repetition is mechanically derived from one authoritative source;
   - add a dependency only when it removes a meaningful subsystem after adapters, tooling, security, and update costs.

4. Keep explicit code when the alternative adds indirection, reflection, modes, hidden allocation, a one-use helper, or a vocabulary maintainers must learn.
5. Make one coherent reduction at a time. Preview automated rewrites such as `go fix -diff ./...`, apply only reviewed transformations, format changed Go files, and exclude unrelated modernization.
6. Run focused tests and repository-native checks. Add race tests for concurrency changes, regeneration checks for generated code, and repeated benchmarks for performance-sensitive paths.
7. Compare only relevant before-and-after evidence. Separate authored production, test, and generated code; a valid reduction lowers authored code or duplication without materially increasing complexity, API surface, dependencies, or runtime cost.

## Go guardrails

- Use only features supported by the module's declared Go version. Resolve correctness-blocking version uncertainty from local evidence or current primary sources.
- Preserve nilness, evaluation order, error identity, and partial-failure behavior when replacing loops or control flow with `slices`, `maps`, `clear`, `min`, `max`, `strings.Cut`, iterators, or other concise operations.
- Use `errors.Is` and the supported typed-error API. Wrap with `%w` only when exposing the underlying identity is intentional; never replace ordinary errors with panic-based brevity.
- Use `sync.WaitGroup.Go` only for supported launch-and-wait semantics and `errgroup` only when its error, cancellation, and concurrency behavior matches the contract. Prove cancellation and goroutine cleanup, not only the returned error.
- Do not edit generated output directly, widen an exported API through embedding, replace capability interfaces with generics, or preserve manual complexity on an unmeasured performance assumption.

## Result

Report what was reduced, the evidence that behavior remained stable, why the replacement is easier to maintain, checks and benchmarks run, authored/test/generated code impact when measured, and any rejected candidates or remaining risks. Finish when the final diff contains only reductions whose maintenance benefit exceeds their conceptual and operational cost.
