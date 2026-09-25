---
name: tech-debt-audit
description: Audit a whole repository or subsystem for technical debt through parallel read-only investigation lanes, adversarial verification, and cost-based ranking, and produce an evidence-backed debt register. Use when the user wants to find, size, or prioritize existing debt across a codebase; not for reviewing one change, which belongs to code-review.
---

# Tech Debt Audit

Produce a ranked debt register a team can turn into work items. Debt is a property of existing code that makes likely future changes slower, riskier, or more expensive. Keep the audit read-only: do not edit, format, commit, push, install, upgrade, or change generated state.

## Orient

1. Confirm the target repository or subsystem, excluded paths, and the output location. Skip vendored, generated, and build output except for dependency analysis.
2. Read repository instructions, README, contributing guides, ADRs, and architecture docs. Identify languages, frameworks, runtimes, build system, test runner, CI, and deploy target.
3. Find the native quality commands and note which exist. Run only trusted read-only ones.
4. Collect global signals:
   - churn hotspots from `git log --since="12 months ago" --name-only`;
   - co-change pairs that cross module boundaries;
   - size and function count of the hottest files;
   - large files that are stale, and large files that are hot;
   - counts of `TODO`, `FIXME`, `HACK`, `@deprecated`, and lint or type suppressions per directory;
   - single-author hotspots from `git shortlog`.
5. Write a short context brief: system purpose, module map, main flows, hotspots, and available checks. Stop discovery when modules, hot paths, and verification commands are known.

## Investigate

Give each lane the context brief, its question, the read-only rules, and the finding shape below. Run lanes in parallel as subagents when the host supports them and the user permits it; otherwise run them in sequence and label the result as degraded. Ask each lane for at most 15 ranked findings, with depth on hotspots rather than breadth over cold code.

Select lanes from the repository. Skip a lane that does not apply and record why.

- **Design:** module boundaries that routine changes cross (use co-change evidence), dependency cycles, domain logic bound to frameworks or storage, catch-all packages, shallow wrappers, single-implementation interfaces, competing patterns for one concern, and duplicated domain rules.
- **Complexity:** high complexity or deep nesting in hotspots, mode flags, hidden mutable state, dead code, permanent feature flags, and drifted copy-paste clones. Use `crap-index-assess` when complexity and coverage data exist.
- **Reliability:** swallowed errors, missing timeouts or cancellation, unbounded retries, leaked resources, non-idempotent handlers on at-least-once paths, partial writes without transactions or reconciliation, and race-prone shared state.
- **Tests:** risky code with no tests, tautological or over-mocked tests, flaky patterns such as sleeps, wall-clock time, real network, and ordering dependencies, skipped tests with no condition, and slow feedback paths. Delegate sufficiency questions to `test-gap-review`.
- **Dependencies and runtime:** end-of-life runtimes, frameworks, and base images checked against current vendor dates, deprecated APIs with announced removal, unmaintained or duplicate libraries, drifted forks, and known vulnerabilities from tools already present. Delegate provenance and pinning to `dependency-review`.
- **Build and operations:** local and CI commands that disagree, non-gating checks, undocumented setup, missing observability on critical paths, unvalidated configuration sprawl, manual release steps, and migrations without rollback. Delegate CI topology to `ci-review` and infrastructure to `infra-review`.
- **Security posture:** unvalidated trust-boundary input, scattered or missing authorization, tenant isolation by convention, and broad credentials. Report only a concrete input-to-sink path, describe the class of problem rather than an exploit, and delegate deep paths to `security-review`.
- **AI integration:** agent instruction files that contradict code, commands, or each other; generated-code sprawl such as near-duplicate functions, mixed idioms in one module, unused abstractions, and narrating comments; hard-coded or retired model IDs; unversioned prompts without evaluations; model calls without timeouts, retries, or budgets; model output used as trusted input; and tool servers with broader permissions than they need.
- **Knowledge:** documentation that contradicts behavior, undocumented public contracts, single-owner hotspots, and superseded ADRs. Delegate documentation repair scope to `docs-review`.

Require each finding in this shape:

```yaml
id: <lane>-<n>
title: <one concrete line>
location: [path:line]
evidence: <short code, command output, or git statistic>
impact: <what becomes slower, riskier, or costlier, and for whom>
interest: <how the cost grows: per change, per incident, per release, or by a date>
likelihood: high | medium | low        # cost lands within six months
blast_radius: local | module | cross-module | system | external
fix_sketch: <smallest credible fix in one to three sentences>
fix_cost: S | M | L | XL
confidence: confirmed | plausible | inconclusive
```

No evidence means no finding. Style preferences and "not the latest trend" are not debt without a shown cost.

## Verify

Send each finding with high likelihood, cross-module or wider blast radius, or a security lane to an independent verifier with only the finding and repository access. Batch the remaining findings by lane. The verifier tries to disprove the claim: reproduce the evidence, look for mitigations elsewhere such as middleware, framework guarantees, database constraints, other tests, or CI gates, and check for a documented deliberate trade-off. Mark each finding `confirmed`, `downgraded` with new ratings, or `rejected` with the reason. Never treat agreement between lanes as evidence.

## Synthesize

1. Merge duplicates by root cause, keeping the strongest evidence, and group related findings into themes.
2. Score each surviving item as `(impact × likelihood × blast radius × interest) / fix cost` on a small integer scale, and show the inputs.
3. Tag each item: `quick win` (small cost, high priority), `strategic` (large cost, high priority), `watch` (low likelihood now, with a dated trigger or condition), or `accept` (carrying costs less than fixing, with the reason).
4. For each strategic item give the first safe step, the rollback point, and the check that proves progress.

## Result

Write the report to the agreed location, or return it inline when none was given. Include a summary of at most ten lines with the top themes and the first action; a hotspot map of churn × complexity; the debt register sorted by priority; one section per theme; the watch list with dates; rejected and accepted items with reasons; and the commands run, lanes skipped or failed, unavailable tools, and uncovered areas. Recommendations do not authorize implementation.

Finish only when every lane has returned or has a recorded failure, every finding has a verification disposition, and every register entry has current locations and evidence.
