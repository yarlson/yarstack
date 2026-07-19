## Contract Evolution

Treat public APIs, schemas, serialized data, configuration, protocols, and other published behavior as contracts. Before changing one, identify its consumers, stored data, deployed versions, and compatibility requirements.

Preserve established inputs, outputs, defaults, errors, and semantics unless the task explicitly changes them. Do not assume a change is safe because known source code still compiles; unknown clients and persisted data may depend on observable behavior.

For an incompatible change, use expand, migrate, then contract:

1. Add support for the new contract without removing the old one.
2. Move consumers and data, and verify the new path in current evidence.
3. Remove the old contract only after its consumers and data are gone.

When versions may overlap, deploy readers that accept both forms before writers produce the new form. Keep migrations restartable and safe after partial completion. Preserve a working rollback path until old and new versions no longer need to coexist.

Define how invalid, missing, old, and new values behave throughout the transition. Do not silently reinterpret existing data or reuse an established name for incompatible behavior.

Before finalizing, confirm the change order, mixed-version behavior, migration evidence, removal condition, and rollback point.
