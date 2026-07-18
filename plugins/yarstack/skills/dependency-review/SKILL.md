---
name: dependency-review
description: Review dependency provenance and reproducibility without modifying them. Use when manifests, lockfiles, vendor state, generated dependency state, registries, fetched tools, or version pins change.
---

# Dependency Review

Own dependency necessity, provenance, pinning, reproducibility, and generated-state consistency.

## Workflow

1. Confirm the changed dependency surface, comparison target, repository pinning policy, and expected update path.
2. Establish the consumer and current need for each added, removed, or changed dependency or tool.
3. Verify manifest, lockfile, vendor, checksum, and generated-state consistency without normalizing unrelated content.
4. Check provenance, registries, Git sources, floating versions, fetched installers, lifecycle scripts, and reproducibility.
5. Prefer trusted offline and non-mutating checks. Require explicit authority before network access, installation, cache mutation, lifecycle-script execution, or regeneration.
6. Delegate exploitability and credentials to `security-review`, CI machinery to `ci-review`, infrastructure impact to `infra-review`, and runtime compatibility to `rollout-readiness-review`.

Report each finding with location, affected contract, evidence, consequence, smallest correction, verification, and uncertainty. State directly when no actionable finding exists. An evidenced inconclusive result is valid when tooling or provenance is unavailable.

Do not upgrade unrelated dependencies, churn an entire lockfile, add scanners, or trust committed generated code without checking its source.

Finish when each changed dependency is accounted for and residual provenance or reproducibility risk is explicit.
