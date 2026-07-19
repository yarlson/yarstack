## Module Boundaries

Draw module boundaries around coherent knowledge and changes likely to occur together. A module should own a stable domain concept, policy, representation, protocol, or other design decision and keep its changes from spreading to unrelated code.

Do not split modules by execution order, workflow step, technical layer, CRUD resource, or arbitrary size. Put code together when the same likely change affects it; separate code when it uses a different model, policy, vocabulary, or source of change.

Expose the smallest useful interface from the client's point of view. Hide storage, representation, algorithms, and integration details that clients do not need. Make errors, consistency, latency, resource ownership, and other behavior clients must rely on explicit in the contract.

Keep dependencies explicit and directional. Stable domain policy should not depend directly on volatile frameworks, storage, transport, or global registries. Translate external types and terms at the boundary. Avoid cycles, shared mutable state, duplicate knowledge, and paths that bypass the owning module.

Do not create generic buckets such as `common`, `util`, `types`, `interfaces`, or `api`. Avoid shallow wrappers that only forward calls, interfaces added only for symmetry or mocking, and modules that do not hide a real decision or contain change.

Use observed changes to test the design. If routine work crosses a boundary or requires unrelated domain knowledge, reconsider ownership or the contract. Before finalizing a boundary, name what it owns, what it hides, which callers need it, and one expected change that should remain inside it.
