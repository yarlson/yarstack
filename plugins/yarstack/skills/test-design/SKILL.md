---
name: test-design
description: Design or write reproducible tests matched to behavioral risk before implementation. Use when choosing test levels, examples, properties, generators, fault cases, fixtures, substitutes, or stable observable boundaries is non-trivial.
---

# Test Design

Own the prospective test contract and test code, not the production implementation or retrospective evidence audit.

## Workflow

1. Identify the authoritative contract, condition, action, observable result, invariants, important side effects, assumptions, and state preserved on failure.
2. Choose the lowest test level that proves the contract through a stable public or system boundary.
3. Match the method to the risk: use examples for concrete cases, property or model-based tests for invariants and sequences, fuzzing for broad or hostile input, and systematic concurrency or fault-injection tests for ordering, restart, and recovery. Select only methods that address a material risk.
4. Make generated or scheduled cases reproducible with a recorded seed, useful bounds, and a minimized counterexample or retained failing input where the tooling supports them.
5. Select only cases that materially define the behavior: normal operation, important boundaries, meaningful failures, cleanup, and external effects.
6. Follow the repository's existing test structure, helpers, fixtures, and commands.
7. Write the smallest reproducible test whose failure identifies the broken contract.
8. State what remains unverified when the behavior cannot be tested economically instead of adding a weak proxy assertion.

Use `technical-spike` when a high-consequence claim may need unfamiliar model checking or proof beyond tests. Use `behavior-implement` for red-green production changes, `test-gap-review` for retrospective verification sufficiency, `cli-control` for terminal evidence, and `ui-control` for graphical interface evidence.

Prefer repository-supported tools. Do not change production code, add dependencies, or introduce broad test infrastructure unless the parent task already authorizes it. Report the missing method, expected benefit, and residual risk instead.

Finish with the test boundary, cases, test code when requested, command to run, and any residual gap.
