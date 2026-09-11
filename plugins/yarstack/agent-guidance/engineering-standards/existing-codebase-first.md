## Existing Codebase First

Before writing code, inspect the closest existing patterns for structure, naming, errors, logging, testing, configuration, dependency wiring, and API shape.

- Extend healthy helpers, packages, conventions, and test styles instead of creating a second way to solve the same problem.
- Before adding an abstraction, package, interface, service, helper, middleware, configuration object, or dependency, check whether an equivalent already exists.
- Before editing a function, find its callers. A defect reported on one path usually lives in the shared code; fix it once there instead of guarding each caller.
- If the existing pattern is unhealthy, explain the problem and make the smallest localized improvement. Do not silently introduce a competing pattern.

Gather only the repository context needed for the current task:

- read applicable repository instructions and linked context documents;
- identify the relevant source, tests, commands, documentation, and configuration;
- resolve cheaply discoverable facts from local evidence before asking the user;
- use current authoritative external sources only when correctness depends on unfamiliar or version-sensitive behavior;
- stop discovery when the edit surface, verification surface, constraints, and unresolved risks are clear.

Do not turn discovery into broad documentation work or inspect generated output, vendored dependencies, build artifacts, or unrelated modules without a concrete need.

State an intentional deviation from an existing pattern and why it is necessary. Otherwise let the diff show which pattern it extends.
