---
name: english-text-review
description: Review English technical prose without modifying it. Use for plain-English feedback on comments, docstrings, documentation, commit or pull request text, issues, reports, runbooks, logs, errors, command help, and other human-readable technical writing rather than factual or marketing review.
---

# English Text Review

Find concrete prose problems and propose the smallest clearer wording while preserving meaning and technical precision.

## Workflow

1. Confirm the text, intended reader, purpose, and any wording that must remain exact.
2. Apply Orwell's six rules below.
3. Identify vague abstractions, repeated meaning, and unclear subjects or consequences that the rules do not make explicit.
4. Preserve exact identifiers, commands, protocol terms, legal language, quotations, externally required wording, and necessary domain vocabulary.
5. Report only findings that materially improve clarity, brevity, or comprehension.

## Orwell's Six Rules

1. Never use a metaphor, simile or other figure of speech which you are used to seeing in print.
2. Never use a long word where a short one will do.
3. If it is possible to cut a word out, always cut it out.
4. Never use the passive where you can use the active.
5. Never use a foreign phrase, a scientific word or a jargon word if you can think of an everyday English equivalent.
6. Break any of these rules sooner than say anything outright barbarous.

For each finding, include the smallest relevant excerpt, the problem, why it affects the intended reader, and a concrete replacement. Rank findings by their effect on meaning, then readability. State directly when no material finding exists.

Keep the review read-only unless editing is separately authorized. Use `docs-review` for factual accuracy and document structure, `marketing-claims-review` for claim substantiation, and the code-comment policy to decide whether a comment should exist. Do not flatten an intentional voice, alter quoted text, or trade precision for simpler-looking words.

Finish with prioritized findings, proposed wording, the reviewed scope, and any context that limits confidence.
