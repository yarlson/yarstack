## Code Comment Policy

Make the code clear without comments. Improve its names, structure, or control flow before adding an explanation.

Write a comment only at the user's request or when the language or repository requires one. Examples include documentation for exported symbols, docstrings required by a linter, and annotations established by the file format. Ask the user before writing any other comment.

A comment is not justified because:

- the code took effort to write or seems subtle or clever
- you learned something while doing the work
- a future maintainer might question the implementation
- an approved plan included the proposed comment

Do not use comments for:

- past decisions, implementation history, or explanations of how the code changed
- change summaries, pull request explanations, or narrative accounts
- apologies or notes about deleted code
- descriptions already clear from the syntax or names
- product opinions or guesses about future needs
- TODOs that lack an owner or a specific completion condition
- temporary notices that lack a specific removal condition

Follow the surrounding source style. Keep a comment-free file free of comments.

When changing code, correct comments that become inaccurate. Remove comments that are stale, historical, obvious, or repetitive.

List each new comment separately in the change report.

Comments document the present system, not the journey that produced it.
