## Existing Codebase First

Before writing code, inspect the closest existing patterns for structure, naming, errors, logging, testing, configuration, dependency wiring, and API shape.

- Extend healthy helpers, packages, conventions, and test styles instead of creating a second way to solve the same problem.
- Before adding an abstraction, package, interface, service, helper, middleware, configuration object, or dependency, check whether an equivalent already exists.
- If the existing pattern is unhealthy, explain the problem and make the smallest localized improvement. Do not silently introduce a competing pattern.

Gather only the repository context needed for the current task:

- read applicable repository instructions and linked context documents;
- identify the relevant source, tests, commands, documentation, and configuration;
- resolve cheaply discoverable facts from local evidence before asking the user;
- use current authoritative external sources only when correctness depends on unfamiliar or version-sensitive behavior;
- stop discovery when the edit surface, verification surface, constraints, and unresolved risks are clear.

Do not turn discovery into broad documentation work or inspect generated output, vendored dependencies, build artifacts, or unrelated modules without a concrete need.

Before implementation, report:

- the existing pattern and where it is used
- how the new code will follow it
- any intentional deviation and why it is necessary
