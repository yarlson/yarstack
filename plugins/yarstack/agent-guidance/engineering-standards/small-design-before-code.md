## Small Design Before Code

Find the smallest change that satisfies the requirement before writing any code. Stop at the first option that holds:

1. It does not need to exist. The behavior is already there, or the need is speculative. Skip it and say so in one line.
2. The codebase already has it. Extend the existing helper, type, or pattern instead of writing a second one.
3. The standard library already does it. Call it.
4. The platform already does it. Prefer a native control, a database constraint, or a runtime or framework guarantee to application code.
5. An installed dependency already does it. Use it, and do not add a dependency for what a few lines cover.
6. It fits in one clear line. Write one line.
7. Otherwise write the minimum code that works, in the fewest files, deleting code before adding it.

Choose the option only after understanding the problem: read the code the change touches and trace the real flow end to end. A small change in the wrong place is a second defect, not a small change. When two options are the same size, take the one that is correct at the boundaries.

Never cut to save lines: validation at trust boundaries, error handling that prevents data loss, security controls, accessibility, or anything the user explicitly asked for. When the user insists on the larger version, build it without re-arguing.

When the change is more than a few lines, sketch it before editing: the exact behavior being changed, the minimal files or units, the data flow and error paths, the checks needed, and what remains unchanged. Reject unrelated cleanup, broad refactoring, speculative abstractions, future-proofing, new frameworks, service boundaries, workers, queues, state machines, or configuration unless the task requires them. Before finalizing, confirm the implementation still matches the sketch; explain necessary growth or reduce the change.

For behavior-preserving structural work, identify the behavior that must remain unchanged, keep the edit sequence mechanical and reversible, name rollback points when partial application is risky, and specify focused regression checks before editing.
