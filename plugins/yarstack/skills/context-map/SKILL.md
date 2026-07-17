---
name: context-map
description: Map the repository context needed for a requested change or plan phase. Use before implementation or validation when relevant files, tests, commands, conventions, risks, or unknowns are not yet clear.
---

# Context Map

Gather enough local context to work safely without reading the entire repository.

## Workflow

1. Read the request or selected plan phase.
2. Read the applicable `AGENTS.md` and linked context documents.
3. Search for relevant source, tests, commands, documentation, and configuration.
4. Inspect nearby implementation and test patterns before choosing an approach.
5. Identify unknowns that affect correctness.
6. Resolve discoverable unknowns from repository evidence and installed tool versions.
7. Use current official documentation only when local evidence is insufficient for unfamiliar or version-sensitive behavior.

## Required Context

Before proceeding, know:

- which files may need changes;
- which files and behaviors need verification;
- which commands provide relevant evidence;
- which local conventions constrain the change;
- which risks and assumptions remain.

## Guardrails

- Do not turn context mapping into broad documentation work.
- Do not inspect generated output, build artifacts, vendored dependencies, or unrelated modules unless required.
- Do not ask the user for locations or facts that are cheaply discoverable locally.

Finish when the relevant edit surface, verification surface, constraints, and remaining unknowns are clear.
