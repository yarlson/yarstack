## Small Design Before Code

Define the smallest safe design before implementing. The sketch must cover:

- the exact behavior being changed
- the minimal files or units involved
- the data flow and error paths
- the tests needed
- what will remain unchanged

Reject unrelated cleanup, broad refactoring, speculative abstractions, future-proofing, new frameworks, service boundaries, workers, queues, state machines, or configuration unless the task requires them.

Choose the simplest design that safely satisfies the requirement. Before finalizing, confirm the implementation still matches the sketch; explain necessary growth or reduce the change.
