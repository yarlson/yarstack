## Tests as Product Contracts

Tests must explain the behavior the system promises, not the mechanics used to arrange the test. A reader should quickly understand the condition, action, observable outcome, and what remains true on failure.

### Keep tests direct

- Arrange only inputs that define the scenario, perform one meaningful action, and assert the observable result and important side effects.
- Name tests in domain language so they read like behavioral contracts.
- Hide temporary directories, fake servers, process wiring, credentials, fixtures, and repetitive construction behind focused setup helpers. Keep scenario-defining values visible.
- Make each setup helper do one cohesive job, fail near a broken prerequisite, and register cleanup where it creates the resource.
- Keep test code shallow, deterministic, independent, and free of hidden mutable state or ordering dependencies.
- Prefer one behavior per test. Avoid conditionals that change a test's meaning and loops outside simple parameterized cases.

### Use tables only when they clarify

Use parameterized or table-based tests when named cases share the same setup, action, and assertions, especially for validation, parsing, mappings, and boundaries.

Keep case data focused on what varies. Separate materially different behaviors when combining them requires mode flags, branching assertions, different lifecycle expectations, or substantially different setup.

Do not force a table when a few direct tests make the contracts easier to understand.

### Verify observable behavior

Cover the smallest useful set of contracts:

- the normal path
- important input boundaries
- meaningful failure paths and preserved state
- cleanup and resource ownership
- externally visible side effects

Assert prerequisites separately when their failure would make later assertions misleading. Prefer explicit expected values over clever generation and failure messages that identify the broken contract.

Avoid real external services, arbitrary sleeps, timing assumptions, shared global state, and brittle snapshots unless the behavior requires them. Use the smallest realistic substitute at the system boundary; do not mock internal details merely to increase isolation.

### Test prompts and policies honestly

Sentence or substring presence does not prove that an agent follows a prompt or policy. Test the actual contract:

- test loading or assembly when that is the contract
- assert exact text only when wording is intentionally external and fixed
- use an integration test or evaluation when the contract is agent behavior

If behavior cannot be tested cheaply, state the remaining risk instead of adding a false-confidence string assertion.

### Before finalizing

Confirm that test names describe meaningful behavior, each test is understandable without reading its helpers first, setup does not hide the scenario, table cases share one execution path, success and failure contracts are clear, behavior-preserving refactors would keep tests valid, and the suite is deterministic and independent.

If a test is difficult to read, simplify the product boundary or setup before adding explanatory comments. Simplify any test that fails this review.
