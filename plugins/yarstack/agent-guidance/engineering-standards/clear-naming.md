## Clear Naming

Choose names that make code predictable at the point of use. A reader should understand what a name refers to without opening its definition.

Match detail to scope and frequency:

- use longer, descriptive names when the definition is distant, the scope is broad, or the concept is rarely used
- use short conventional names when the scope is small and the meaning is clear from nearby code
- avoid long names that repeat context already supplied by a module, type, receiver, or package

Use consistent domain terms. Use the same verb for the same operation and parallel names for parallel interfaces. Readers should be able to predict related names.

Name functions and methods for the result they return or the action they perform. Include the affected object when context does not make it clear. Replace vague words such as `process`, `handle`, `manage`, `helper`, or `utils` with the actual operation or domain concept.

Prefer descriptive domain names over cute, mnemonic, or newly coined names. Preserve established terms when renaming would cause churn or break shared vocabulary.

Before finalizing, read each new or changed name at its use sites. Confirm that its meaning is clear, its detail fits its scope, and related code uses the same vocabulary.
