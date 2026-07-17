---
name: technical-spike
description: Run a bounded technical investigation when implementation or validation depends on unknown, unfamiliar, or version-sensitive framework, library, tool, runtime, protocol, or platform behavior. Produce only conclusions that affect the current task.
---

# Technical Spike

Resolve one implementation-blocking technical question without drifting into broad research.

## Workflow

1. State the exact technical question.
2. Check local evidence first: manifests, lockfiles, configuration, code, tests, installed versions, and command help.
3. If local evidence is insufficient, consult current official documentation for the specific tool or library.
4. Record only conclusions that affect implementation or validation.
5. Convert the findings into concrete design, compatibility, or verification implications.
6. Stop when the current task has a safe path forward.

## Result Shape

```md
## Technical Spike Result

Question: ...
Local evidence: ...
External evidence, if used: ...
Conclusion: ...
Impact: ...
Remaining uncertainty: ...
```

Do not research adjacent features, upgrade dependencies without authorization, add abstractions to hide uncertainty, or use unofficial sources when authoritative documentation is available.
