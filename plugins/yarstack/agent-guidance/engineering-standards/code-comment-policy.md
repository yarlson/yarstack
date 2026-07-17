## Code Comment Policy

Use comments only to clarify the current code: invariants, constraints, non-obvious behavior, edge cases, safety or protocol rules, ownership, lifecycle, concurrency, ordering, or reasons an obvious approach is unsafe.

Prefer clear names and simple control flow over comments. Do not use comments as:

- decision logs or implementation history
- change summaries, PR explanations, or prose narratives
- apologies or commentary about removed code
- restatements of syntax or function names
- product commentary or future speculation
- TODOs without an owner or concrete completion condition
- “temporary” notes without a clear removal condition

A useful comment should explain what must remain true, what external contract forces the behavior, what edge case is intentional, or what would break if the code changed.

Before finalizing, review every touched comment. Remove stale, historical, obvious, or duplicated text; update comments that no longer match the code; and add a comment only when simplification cannot make the behavior clear.

Comments document the present system, not the journey that produced it.
