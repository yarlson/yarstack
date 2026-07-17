---
name: docs-drift-review
description: Review repository documentation for drift caused by a code or configuration change. Use when changed behavior may affect README files, agent instructions, plans, API or schema docs, operational docs, examples, configuration templates, or command help.
---

# Documentation Drift Review

Update only documentation made inaccurate by the current change.

## Workflow

1. Identify changed behavior, commands, configuration, APIs, workflows, and contracts.
2. Search repository documentation that describes those surfaces.
3. Compare the documentation with the implementation.
4. Update only directly affected canonical documents.
5. Follow the repository's existing documentation style.
6. Update examples, templates, and command help when they are part of the affected contract.

## Guardrails

- Do not document unchanged behavior.
- Do not rewrite documentation for style alone.
- Do not create a parallel documentation system.
- Do not put implementation-plan details into user-facing documentation unless that is the established convention.
- Do not leave contract changes only in chat when a repository document owns them.

Finish when affected documentation is accurate, or when inspection provides concrete evidence that no documentation change is needed.
