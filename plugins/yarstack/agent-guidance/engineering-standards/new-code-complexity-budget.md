## New Code Complexity Budget

Keep every new or heavily changed function, class, module, package, component, or service within a strict simplicity budget.

Prefer:

- one clear responsibility
- short units and shallow control flow
- guard clauses over nested conditions
- explicit data flow and lifecycle ownership
- clear names over explanatory comments
- composition over multipurpose objects
- deterministic behavior over hidden side effects

Avoid:

- modes, flags, branches, or nesting that obscure behavior
- catch-all services and premature interfaces
- generic helpers used only once
- hidden global state or implicit ownership
- mixing business logic, I/O, parsing, validation, and presentation
- concurrency without clear ownership, cancellation, bounds, and cleanup
- comments that explain complexity the code should remove

Before finalizing, confirm each changed unit is understandable without unrelated files, owns one responsibility, handles edge cases near their source, exposes a small test surface, and has an obvious place for future changes. Simplify anything that fails this check.
