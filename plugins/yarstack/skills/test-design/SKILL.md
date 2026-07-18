---
name: test-design
description: Design and write deterministic tests as readable behavioral contracts. Use before adding or changing tests, choosing test cases or test levels, creating fixtures or test doubles, or reviewing whether a proposed test proves observable behavior.
---

# Test Contracts

Write tests that explain and prove the behavior the system promises.

## Establish the Contract

1. Read the requested behavior and the closest existing tests, implementation patterns, and repository test commands.
2. Identify the condition, action, observable outcome, important side effects, and state preserved on failure.
3. Choose the lowest test level that proves the contract through a stable boundary without coupling the test to implementation structure.
4. Select the smallest meaningful set of cases: the normal path, important boundaries, meaningful failures, cleanup or resource ownership, and externally visible side effects.

Do not add cases merely to increase coverage. If behavior cannot be tested economically, state what remains unverified and why instead of adding a weak proxy assertion.

## Keep Tests Direct

- Name tests in domain language so each reads as a behavioral contract.
- Arrange only inputs that define the scenario, perform one meaningful action, and assert explicit observable results and important side effects.
- Keep scenario-defining values visible. Hide repetitive infrastructure, fixtures, credentials, and process wiring behind focused helpers.
- Give each helper one setup responsibility, fail near broken prerequisites, and register cleanup for every resource it creates.
- Keep tests shallow, deterministic, independent, and free of ordering dependencies or hidden mutable state.
- Prefer one behavior per test. Do not use conditions or modes that change what a test means.
- Assert prerequisites separately when their failure would make the behavioral assertion misleading.

Use parameterized tests only when named cases share the same setup, action, and assertions. Keep case data limited to what varies. Split cases when combining them requires branching assertions, lifecycle modes, or materially different setup.

## Test Stable Boundaries

Assert public results and externally visible effects rather than private calls, internal ordering, or incidental representation. A behavior-preserving refactor should normally leave the test valid.

Avoid real external services, arbitrary sleeps, timing assumptions, brittle snapshots, and mocks of internal details unless the behavior requires them. Use the smallest realistic substitute at the system boundary. Prefer explicit expected values and failure messages that identify the broken contract.

For prompts and policies, test the actual contract: verify loading or assembly when that is the behavior, assert exact wording only when the text is intentionally fixed, and use an integration test or evaluation for agent behavior. Substring presence does not prove policy compliance.

## Review the Test

Before finishing, confirm:

- the name states meaningful behavior;
- the scenario is understandable without reading helpers first;
- success, failure, preserved state, and side effects are explicit where relevant;
- setup does not hide scenario inputs;
- the test fails when its promised behavior is broken;
- the suite remains deterministic and independent.

If the test is difficult to understand, simplify the product boundary or setup before adding explanatory comments.
