# Engineering standards

[`standards.md`](standards.md) is the guidance that Yarstack can install as
global instructions for Claude Code and Codex. This file is for people who
maintain or evaluate that guidance. It explains how the guidance is built, which
first principles it follows, and which evidence supports each rule. The
installer does not copy this file.

To preview the installed text, run
`plugins/yarstack/scripts/install-engineering-standards.sh --print` from the
repository root. The [repository README](../../../../README.md#optional-engineering-standards)
describes installation.

## How the guidance is built

**One file of rules, one file of evidence.** Agents load the rules in every
session. Maintainers need the evidence only when they change a rule. Citations
in the rules would cost context in every session and would not change
behavior. Claude Code's documentation asks for files "under 200 lines" because
"longer files consume more context and reduce adherence" ([memory docs][cc-memory]).
In a 2026 study, repository context files did not generally improve task
success and raised inference cost by more than 20% ([Gloaguen et al.][gloaguen]).
So `standards.md` holds only rules that change agent behavior, and this file
holds their justification.

**One rules file, not many fragments.** The installer writes one file, and no
reader uses a subset of the rules. A single file has no assembly order to
maintain, and it makes duplicate or conflicting sections visible in one place.

**Conflict rules first, closeout last.** The section that resolves conflicts
comes first and the closeout comes last. IFScale found a primacy bias in
instruction following that peaks near 150 to 200 instructions
([Jaroslawicz et al.][ifscale]). "Lost in the Middle" found that retrieval is best
at the start and end of a context ([Liu et al.][lost-middle]). That study measured
document retrieval, not instructions, so this ordering is judgment informed by
weak evidence.

**A reason before the rules.** A section opens with one sentence of reason when
the reason is not obvious. Anthropic reports that "explaining to Claude why such
behavior is important, can help Claude better understand your goals"
([prompting best practices][claude-prompting]). It also recommends guidance at the
"right altitude", between brittle hard-coded logic and vague high-level advice
([context engineering][ctx-eng]).

**Two strengths of rule.** "Never" and "always" mark hard limits. Every other
rule is a default that a concrete reason can override. Hosts do not resolve
conflicts for the agent: in Claude Code, user and project rules do not override
each other, and "if two instructions contradict each other, Claude may pick one
arbitrarily" ([memory docs][cc-memory]). So the guidance states its own precedence.
A rule is a hard limit only when breaking it causes harm that review may not
catch, such as an unauthorized action, data loss, a leaked secret, or a false
claim of success. The marker is a plain word, not emphasis formatting, because
"if you emphasize many lines, none of them stands out"
([Claude Code best practices][cc-best]).

**The admission test for a rule.** A rule belongs in `standards.md` only when it
prevents a failure that agents show without it, or a failure that is costly
when it happens. Claude Code's guidance gives the test: "Would removing this
cause Claude to make mistakes? If not, cut it." It also excludes "self-evident
practices like 'write clean code'" ([best practices][cc-best]).

**What the file is for the host.** Claude Code delivers CLAUDE.md "as a user
message after the system prompt", and it treats the content "as context, not
enforced configuration" ([memory docs][cc-memory]). Codex gives each AGENTS.md
file to the model as a user-role message ([Codex guide][codex-agents]). The
guidance cannot enforce anything. Host permissions, hooks, and CI enforce
behavior; the guidance states intent and judgment.

## First principles

Each section of `standards.md` follows from one or more of these principles.

1. **Code is changed far longer than it is written.** "Software engineering is
   programming integrated over time" ([SWE at Google][swe-book]). The cost that
   matters is the cost of the next change. This principle drives *Add as little
   as possible*, *Design for the next change*, *Names*, and *Comments*.
2. **Every addition carries cost, and size hides defects.** Each presumptive
   feature has costs of build, delay, carry, and repair ([Fowler][yagni]). In a
   Mozilla study, larger patches made reviewers more likely to miss bugs
   ([Kononenko et al.][kononenko]). This drives *Add as little as possible* and
   *Finish and report*.
3. **Hide each decision that is likely to change.** A module should hide a
   design decision "from the others" ([Parnas 1972][parnas72]). Complexity is what
   makes a system "hard to understand and modify" ([Ousterhout][aposd-debate]).
   This drives *Design for the next change*.
4. **Observable behavior becomes a contract.** "All observable behaviors of your
   system will be depended on by somebody" ([Hyrum's Law][hyrum]). This drives
   *Preserve contracts*.
5. **Failure is normal at boundaries.** 92% of catastrophic failures in five
   distributed systems came from wrong handling of non-fatal errors that the
   software had signaled ([Yuan et al.][yuan14]). Reliability on a path needs
   checks at its endpoints ([Saltzer, Reed, Clark][e2e]). This drives *Handle
   failure explicitly* and *Change state safely*.
6. **Least privilege, fail-safe defaults, complete mediation.** These are three
   of Saltzer and Schroeder's design principles for protection
   ([1975][saltzer75]). They drive *Protect trust boundaries* and *Work within
   your authority*.
7. **A claim is evidence only if it could be shown false.** Coverage has a low
   to moderate correlation with test effectiveness once suite size is
   controlled ([Inozemtseva and Holmes][inozemtseva]). Agents game tests when
   they can ([ImpossibleBench][impossiblebench]). This drives *Prove behavior with
   checks that can fail* and *Finish and report*.
8. **An agent acts for a user, and some actions cannot be undone.** An agent
   repeats a mistake as easily as a correct step, and the user carries the
   result. This principle is judgment. It drives *Work within your authority*
   and *Understand before you change*.
9. **Reader attention is scarce.** Plain language lets readers "find what they
   need, understand what they find, and use that information" ([IPLF][iplf]).
   Agents follow fewer instructions as instruction count grows
   ([IFScale][ifscale]). This drives *Write plainly* and the size of
   `standards.md`.

## Evidence labels

- **Measured**: empirical studies with data support the rule. The entry gives the
  effect and the study's limits.
- **Practice**: recognized practitioners, organizations, or standards bodies
  document the rule from experience and argument, without controlled
  measurement.
- **Judgment**: a reasoned conclusion of this guidance. It often adapts a source
  to agents or chooses between sources that disagree.
- **Contested**: credible sources disagree. The entry names the sides and the
  choice that `standards.md` makes.

A rule can carry more than one label. Measured evidence often has a narrow
scope, and the entry states that scope.

## Evidence by section

### How to use these standards

- **Repository and user instructions take precedence.** Practice and judgment.
  Codex places closer AGENTS.md files later so that they "override earlier
  guidance" ([Codex guide][codex-agents]). Claude Code concatenates files and
  resolves contradictions arbitrarily ([memory docs][cc-memory]). The rule removes
  that arbitrary choice. "State the risk once" keeps the user informed without
  arguing.
- **Hard limits and defaults.** Judgment, based on the emphasis guidance in
  [Claude Code best practices][cc-best].
- **Conflict order.** Judgment. The order ranks errors by how hard they are to
  undo: an unauthorized external action or data loss, then wrong behavior, then
  broken consumers, then maintenance cost. TigerStyle shows the value of a
  stated order: "safety, performance, and developer experience. In that order"
  ([TigerStyle][tigerstyle]).

### Work within your authority

- **Deliver what was asked, at the scope intended.** Practice. Anthropic reports
  that Claude Opus 5 "can also expand the scope of a task, adding steps that
  weren't requested", and it gives this wording as the remedy
  ([Opus 5 prompting][opus5]).
- **Read-only reviews and explicit authority for external actions.** Judgment
  (principles 6 and 8). The host's permission system enforces access; this rule
  tells the agent where the line is when the host allows more.
- **Ask only when the answer changes behavior or scope.** Practice. Anthropic
  suggests checking in "only when different readings of the request would lead
  to materially different work" ([Opus 5 prompting][opus5]).
- **Report defects outside the task.** Practice. The Linux kernel asks for one
  problem per patch ([submitting patches][kernel-patches]). Google asks authors
  to separate refactorings from other changes ([small CLs][small-cls]).

### Understand before you change

- **Read the code, trace the path, find the callers.** Practice. Anthropic's
  guidance says "Never speculate about code you have not opened"
  ([prompting best practices][claude-prompting]). For human reviewers,
  "understanding is their main challenge" ([Bacchelli and Bird][bacchelli]).
- **Fix a defect where it starts.** Practice. DRY means that "one piece of
  knowledge in the domain changes, one piece of the system changes"
  ([Hunt and Thomas][dry]). Guards in each caller cause the change amplification
  that Ousterhout names as a symptom of complexity ([APOSD extract][aposd-extract]).
- **Confirm that an API or package exists before you depend on it.** Measured.
  Code models hallucinated package names at rates of at least 5.2% for
  commercial models and 21.7% for open-source models. 43% of hallucinated names
  repeated in all 10 runs of the same prompt ([Spracklen et al.][spracklen]).
  Repeated names are predictable, so an attacker can publish a package under
  one; that step is inference. The rule is limited to APIs the codebase does
  not already use, to avoid over-verification ([Opus 5 prompting][opus5]).
- **Follow existing patterns.** Judgment. A second way to do the same thing adds
  to the cognitive load that Ousterhout describes ([APOSD extract][aposd-extract]).
  Developers rarely choose the same name on their own (see *Names*), so
  consistency comes from reuse, not from chance.
- **Stop exploring when the change and its check are clear.** Judgment, informed
  by Anthropic's report that explicit verification instructions "cause
  over-verification" on current models ([Opus 5 prompting][opus5]).

### Add as little as possible

- **The smallest-change ladder.** Practice and local evidence. Yagni does not ban
  work that makes software easier to change, but it counts the "cost of carry"
  of anything built for a guessed future ([Fowler][yagni]). Anthropic documents a
  "tendency to overengineer by creating extra files, adding unnecessary
  abstractions, or building in flexibility that wasn't requested"
  ([prompting best practices][claude-prompting]). Agents that used an earlier
  version of this guidance "built more than tasks needed" (commit `e22616f` in
  this repository).
- **Larger changes hide more defects from review.** Measured. Patch size had a
  significant effect on missed bugs in all four models of a Mozilla study
  ([Kononenko et al.][kononenko]). The SmartBear study of Cisco reviews suggests
  200 to 400 lines per review, but its vendor author gives conclusions "without
  explanation of statistical methods", so this guidance does not use its number.
- **Never cut validation, data-loss handling, security, accessibility, or
  explicit requests.** Judgment. These protect against harm that a smaller diff
  cannot repay.
- **Outline a change of more than a few lines first.** Judgment. The outline makes
  scope growth visible before the code exists.
- **One source of truth for knowledge, not for look-alike code.** Practice. DRY
  concerns knowledge: "Most people take DRY to mean you shouldn't duplicate code.
  That's not its intention" ([Hunt and Thomas][dry]). "Duplication is far cheaper
  than the wrong abstraction", and a sign of the wrong abstraction is "passing
  parameters and adding conditional paths through shared code"
  ([Metz][metz]).
- **No helper or abstraction for a single use.** Practice. Anthropic's sample
  prompt says "Don't create helpers, utilities, or abstractions for one-time
  operations" ([prompting best practices][claude-prompting]). Fowler applies Yagni
  to "abstractions to support future flexibility" ([Yagni][yagni]).
- **Add a dependency only when it pays for itself, and confirm the package.**
  Practice and measured. Russ Cox frames dependency cost as risk and lists seven
  practices: inspect, test, abstract, isolate, avoid, upgrade, and watch
  ([Cox][cox-deps]). He also finds that "delaying an upgrade is riskier than
  upgrading quickly" when tests run, so the guidance does not discourage
  upgrades. The xz backdoor ([CVE-2024-3094][xz]) and the left-pad removal
  ([npm][leftpad]) show the cost of a bad dependency. OWASP added "Software Supply
  Chain Failures" as A03 in its 2025 Top 10 ([OWASP][owasp-top10]). Package
  hallucination makes the name check necessary ([Spracklen et al.][spracklen]).
- **No unrequired changes to interfaces or unrelated code.** Practice. Google's
  review guide asks reviewers to see that developers "solve the problem they
  know needs to be solved now" ([eng-practices][eng-practices]).
- **Structural change first, in its own commit.** Practice. "Make the change easy
  (warning: this may be hard), then make the easy change" ([Beck, 2012][beck-tweet]).
  Beck's own agent instructions say "Never mix structural and behavioral changes
  in the same commit" ([Beck, 2025][beck-claude]). The kernel asks for each
  logical change in its own patch ([submitting patches][kernel-patches]).
- **Small, separately reviewable changes.** Measured and practice. DORA finds
  that working in small batches, with other capabilities, "predicts software
  delivery performance", and that it "amplifies the positive impact of AI
  adoption" ([DORA][dora-batches]). The 2024 DORA report links AI adoption to an
  estimated 7.2% drop in delivery stability for each 25% increase in adoption,
  and it suggests that "changelists are growing in size" ([DORA 2024][dora-2024]).
  At Google, the median change modifies 24 lines ([Sadowski et al.][sadowski]).

### Design for the next change

- **Complexity is what makes change hard.** Practice. Ousterhout names three
  symptoms: change amplification, cognitive load, and unknown unknowns
  ([APOSD extract][aposd-extract]). Brooks separates essential difficulties from
  accidental ones ([No Silver Bullet][brooks]). Moseley and Marks name state as
  "the single biggest remaining cause of complexity" ([Out of the Tar Pit][tarpit]).
- **Module boundaries around decisions likely to change.** Practice. Parnas:
  "one begins with a list of difficult design decisions or design decisions
  which are likely to change. Each module is then designed to hide such a
  decision from the others" ([Parnas 1972][parnas72]). The test "name the decision
  it hides and one change that stays inside" applies that criterion.
- **Keep code that changes together in one place.** Practice. Parnas found it
  "almost always incorrect" to decompose by flowchart ([Parnas 1972][parnas72]).
  The rule bans splitting by step or layer only when an ordinary change would
  touch every piece, because a layer that hides a decision is a valid module.
- **Interfaces simpler than implementations; no pass-through layers.** Practice.
  "The best methods are those that provide a lot of functionality but have a very
  simple interface" ([Ousterhout][aposd-debate]).
- **Depth over shortness.** Contested and measured. Ousterhout: "Depth is more
  important than length: first make functions deep, then try to make them short
  enough to be easily read" ([APOSD extract][aposd-extract]). Robert Martin
  decomposes further, until a method does "One Thing", and the two agree only
  that over-decomposition is possible ([debate][aposd-debate]). The data does not
  support small units as safer:
  - In a NASA project, "the larger the module, the less error prone it was"
    ([Basili and Perricone][basili84]).
  - Fenton and Ohlsson "found no evidence to support previous claims relating
    module size to fault density" ([2000][fenton-ohlsson]).
  - After controlling for size, "none of the metrics we studied were associated
    with fault-proneness" ([El Emam et al.][elemam]).
  - Hatton's U-shaped curve of defect density against size is withdrawn by its
    author ([Hatton][hatton]).

  Holzmann's 60 lines and TigerStyle's 70 lines rest on a page or a screen, not
  on defect data ([Power of Ten][p10], [TigerStyle][tigerstyle]). The guidance
  therefore sets no length limit.
- **Invalid states impossible; parse at the boundary.** Practice. "Make illegal
  states unrepresentable" ([Minsky][minsky]). Parse "at the boundary of your
  system, before any of the data is acted upon" ([King][king]). Mixing validation
  into processing code is the "shotgun parsing" antipattern
  ([LangSec][langsec]).
- **Dependencies point one way; no cycles.** Practice. Parnas restricts the
  "uses" relation so "that its graph is loop free" ([Parnas 1979][parnas79];
  the quote was checked in the 1978 conference version).
  Hexagonal architecture keeps the application "developed and tested in
  isolation from its eventual run-time devices and databases"
  ([Cockburn][cockburn]).
- **No hidden state, side effects, or modes.** Practice. State is the main cause
  of complexity ([Out of the Tar Pit][tarpit]). Interleaving concerns is where
  "complexity comes from" ([Hickey][hickey]). Brooker warns that modal behavior is
  "difficult to test" ([AWS][aws-retries]).
- **Estimate hot-path cost; optimize where measured.** Practice. "We should forget
  about small efficiencies, say about 97% of the time... Yet we should not pass
  up our opportunities in that critical 3%" ([Knuth 1974][knuth74]).

### Names

- **One word per concept; keep established terms.** Measured. When developers
  name the same thing on their own, "the median probability was only 6.9%" that
  two choose the same name ([Feitelson et al.][feitelson]). Consistency needs
  deliberate reuse.
- **Name length grows with distance.** Practice. "The greater the distance
  between a name's declaration and its uses, the longer the name should be"
  ([Gerrand][gerrand]). The Google Go guide states the same rule in terms of
  scope ([Go style][go-style]).
- **Whole words over non-standard abbreviations.** Measured, with small effects.
  Whole words gave 19% faster comprehension (d = 0.32) ([Hofmeister et al.][hofmeister])
  and 14% faster (d = 0.27, about 2% of variance) ([Schankin et al.][schankin]).
  Other studies found no difference ([Lawrie et al. 2007][lawrie07],
  [Beniamini et al.][beniamini]). This is a default, not a hard limit.
- **No empty names such as `util` or `helper`.** Practice. "Packages named `util`,
  `common`, or `misc` provide clients with no sense of what the package
  contains" ([Go blog][go-pkg-names]).

### Comments

The comment rules are contested and judgment. The evidence:

- Ousterhout and Martin agree that "implementation code only needs comments when
  the code is nonobvious", and both accept comments for public APIs. Ousterhout
  would write "5-10x more lines of comments" than Martin
  ([debate][aposd-debate]).
- Comments drift and can mislead. iComment found 60 comment-code
  inconsistencies, 33 of them new bugs ([Tan et al.][icomment]). Only 13% to 20%
  of code changes come with a comment change, although the authors warn that this
  does not measure inconsistency ([Wen et al.][wen19]).
- The effect of comments on comprehension ranged from a 30% decrease to a 34%
  increase, depending on the code ([Abdelsalam et al.][abdelsalam]).
- The Linux kernel and the Google Go guide allow comments for what the code
  cannot carry ([kernel style][kernel-style], [Go style][go-style]).
- Agents add needless comments. Anthropic's sample prompt says "Don't add
  docstrings, comments, or type annotations to code you didn't change"
  ([prompting best practices][claude-prompting]).

The choice: no comment by default, a narrow list of facts that justify one, and
hard limits on narration, history, and TODOs without an owner. The guidance does
not require the user's approval before each comment, because an agent in an
autonomous workflow cannot ask. The report lists each new comment, so the
reviewer decides.

### Preserve contracts

- **Every observable behavior becomes a dependency.** Practice, observational
  ([Hyrum's Law][hyrum]).
- **Compiling code does not prove compatibility.** Practice. AIP-180 separates
  source, wire, and semantic compatibility, and it treats a rename as "remove and
  add" ([AIP-180][aip180]).
- **Expand, migrate, contract.** Practice. Parallel change breaks a change "into
  three distinct phases: expand, migrate, and contract" ([Sato][parallel-change]).
  Stripe migrates data with dual writes, then moves reads, then moves writes,
  then removes old data ([Stripe][stripe-migrations]).
- **Readers before writers.** Practice. "Clients and servers are never updated at
  exactly the same time" ([Protocol Buffers][protobuf-dos]). Kubernetes may not
  advance a storage version until a release supports both versions
  ([deprecation policy][k8s-deprecation]).
- **Never reuse an identifier for a different meaning.** Practice. "Never re-use a
  tag number" ([Protocol Buffers][protobuf-dos]).
- **Define behavior for missing, invalid, old, and new values.** Practice. Tolerating
  unexpected input "is no longer considered best practice in all scenarios", and
  "a flaw can become entrenched as a de facto standard" ([RFC 9413][rfc9413]).
- **Restartable migrations and a rollback path.** Practice. Kubernetes keeps both
  versions served so that users can "roll back to a previous release"
  ([deprecation policy][k8s-deprecation]).

### Handle failure explicitly

- **Handle or propagate every error; catch the specific failure.** Measured. In
  198 user-reported failures in Cassandra, HBase, HDFS, MapReduce, and Redis,
  92% of the catastrophic ones came from wrong handling of signaled errors. In
  58% of them, simple tests of the error-handling code would have found the
  fault. Ignored errors caused 25%, aborts in over-broad catch blocks 8%, and
  handlers with only a TODO 2%. "An error handler that only logs the error is
  also considered as ignoring the error" ([Yuan et al.][yuan14]).
- **Diagnosable errors without secrets.** Practice. OWASP lists data to keep out
  of logs, including session IDs, access tokens, passwords, and sensitive
  personal data ([OWASP logging][owasp-logging]). "Don't just check errors, handle
  them gracefully" ([Go proverbs][go-proverbs]).
- **Validate at boundaries; no defensive code for impossible states; check
  critical invariants.** Practice, contested, and weakly measured.
  - Anthropic: "Don't add error handling, fallbacks, or validation for scenarios
    that can't happen... Only validate at system boundaries"
    ([prompting best practices][claude-prompting]).
  - Higher assertion density correlated with lower fault density, with
    coefficients of about -0.2 to -0.38. The authors believe "enforcing the use of
    assertions would not work well" ([Kudrjavets et al.][kudrjavets]). A GitHub
    study found "a small (but significant) effect" ([Casalnuovo et al.][casalnuovo]).
  - Holzmann requires recovery when an assertion fails ([Power of Ten][p10]).
    TigerStyle says "the only correct way to handle corrupt code is to crash"
    ([TigerStyle][tigerstyle]).

  The choice: check invariants whose failure would corrupt data, stop the
  operation in the way the repository allows, and set no density target. An
  assertion fails loudly. Defensive code hides a bug by handling a state that
  should not exist.
- **Limit everything that can grow or wait.** Practice. "Every socket, process,
  pipe, or remote procedure call can and will hang", and "often, the default is
  to block forever" ([Nygard][nygard]). TigerStyle: "Put a limit on everything"
  ([TigerStyle][tigerstyle]). AWS bounds queue size and queue time
  ([load shedding][aws-shedding]). OWASP applies size limits "before buffering or
  parsing input" ([OWASP input validation][owasp-input]).
- **Retry transient failures, at one layer, with backoff and jitter.** Practice,
  with simulation evidence.
  - Independent retries in five layers can multiply database load "243x". For
    low-cost operations, AWS retries "at a single point in the stack"
    ([AWS][aws-retries]).
  - Google retries "at the layer immediately above the layer that is rejecting
    them", sets retry budgets, and says "Don't retry permanent errors"
    ([SRE overload][sre-overload], [SRE cascading][sre-cascading]).
  - In a simulation with 100 contending clients, jitter cut the call count "by more
    than half" ([Brooker][jitter]).

  The two retry-placement rules differ in detail and agree on one retry layer.
- **A visible failure over a rarely exercised fallback.** Practice. Fallback
  "isn't exercised regularly and is likely to fail or increase the scope of
  impact when it triggers during an outage" ([AWS][aws-fallback]).
- **Telemetry only when someone will use it.** Practice. "Every page should be
  actionable" ([SRE monitoring][sre-monitoring]).

### Change state safely

- **Only the endpoints can make an operation take effect once.** Practice. A
  function "can completely and correctly be implemented only with the knowledge
  and help of the application standing at the end points"
  ([Saltzer, Reed, Clark][e2e]).
- **Atomic, or record progress.** Practice. Atomic phases commit "before initiating
  any foreign state mutation", and recovery points let a retry resume
  ([Brandur][brandur]). The transactional outbox stores a message "as part of the
  transaction that updates the business entities" ([Richardson][outbox]).
- **Idempotent processing, a key for the intent, reject mismatched parameters.**
  Practice. "Any message ever sent may be delivered multiple times"
  ([Helland][helland07]). AWS uses a "caller-provided client request identifier"
  and returns "a validation error indicating a parameter mismatch"
  ([AWS][aws-idempotent]).
- **Unknown outcomes and no exactly-once claims.** Practice and measured.
  "Exactly-once semantics: retry + idempotence or deduplication"
  ([Kleppmann][kleppmann-notes]). "You cannot have exactly-once message delivery"
  ([Treat][treat]). Documented guarantees also fail: PostgreSQL 12.3 serializable
  transactions "were not, in fact, serializable" ([Jepsen][jepsen-pg]).
- **Owned concurrency.** Practice. In most systems, "unhandled errors in background
  tasks are simply discarded" ([N. J. Smith][structured-concurrency]).
  "Goroutines are not garbage collected; they must exit on their own"
  ([Go blog][go-pipelines]). Failing to cancel a context "leaks the child and its
  children" ([context package][go-context]).

### Protect trust boundaries

- **Least privilege, fail-safe defaults, complete mediation.** Practice. "Every
  access to every object must be checked for authority" ([Saltzer and
  Schroeder][saltzer75]).
- **Agents need explicit security rules.** Measured. About 40% of 1,689 Copilot
  programs were vulnerable ([Pearce et al.][pearce]). Users with an AI assistant
  "wrote significantly less secure code" and were "more likely to believe they
  wrote secure code" ([Perry et al.][perry]). Both studies predate current agents.
- **Validate untrusted input where it enters.** Practice. OWASP: "Define what the
  application accepts and reject values outside those rules", and it names
  "trusting internal sources" as a pitfall ([OWASP input validation][owasp-input]).
  The rule adds model output as an untrusted source. That addition is judgment,
  for prompt injection.
- **Keep data apart from code at every sink.** Measured frequency and practice.
  Cross-site scripting and SQL injection rank first and second in the 2025 CWE Top
  25, which draws on 39,080 CVE records ([MITRE][cwe-top25]).
- **Restrict the destinations of fetched URLs.** Practice. OWASP uses an allowlist
  when the targets are known, and notes that "domain allowlisting alone does not
  prevent DNS rebinding" ([OWASP SSRF][owasp-ssrf]).
- **Authorization on every access.** Measured frequency. Missing Authorization
  (CWE-862) ranks fourth in 2025 ([MITRE][cwe-top25]). Broken Access Control is A01
  in the OWASP Top 10:2025 ([OWASP][owasp-top10]).
- **No secrets in logs, output, commits, or code.** Practice ([OWASP
  logging][owasp-logging]).

### Prove behavior with checks that can fail

- **A check is evidence only when it can fail.** Measured. Mutants were "coupled
  with 70% of high-priority bugs" at Google, and developers who used mutation
  testing wrote more tests ([Petrović et al.][petrovic21]).
- **One deterministic check; add cases only for what the change alters.** Practice
  and judgment. Google's Beyoncé Rule: "If you liked it, then you shoulda put a
  test on it" ([SWE at Google, ch. 11][swe-ch11]). The limit on extra cases
  answers the over-building that agents show ([Claude Code best
  practices][cc-best]: "tests for cases that can't happen").
- **Confirm that a new test can fail.** Judgment. TDD's red step and mutation
  testing share this point. TDD research finds that small steps matter more than
  writing the test first (see *Work in small steps*).
- **Test through the public interface; assert state, not interactions.** Practice.
  "Test via public APIs", "Test state, not interactions", and "Test behaviors, not
  methods" ([SWE at Google, ch. 12][swe-ch12]). "Test for observable behaviour
  instead" ([Vocke][vocke]).
- **Real implementations, then fakes; stubs and mocks only for what the system
  does not own.** Practice. "A real implementation should be preferred over a test
  double" and "a fake is often the ideal solution" ([SWE at Google,
  ch. 13][swe-ch13]). "Only unmanaged dependencies should be replaced with mocks"
  ([Khorikov][khorikov-mock]). Fowler presents mockist testing as a valid style
  ([Mocks Aren't Stubs][mocks-stubs]). The guidance chooses the classicist style
  because it survives refactoring.
- **Independent expected values; no logic in tests.** Practice. "Don't put logic in
  tests" ([SWE at Google, ch. 12][swe-ch12]). An expected value that comes from the
  code under test only proves that the code agrees with itself.
- **Deterministic, independent tests.** Measured. Of 161 classified flaky-test
  fixes, async waits caused 45%, concurrency 20%, and test-order dependency 12%
  ([Luo et al.][luo14]). At Google, about 1.5% of test runs gave a flaky result,
  and almost 16% of tests had some flakiness ([Micco][micco]).
- **Readable tests and clear failure messages.** Practice. A failure message
  should "clearly express the desired outcome, the actual outcome, and any
  relevant parameters", and tests follow "DAMP over DRY" ([SWE at Google,
  ch. 12][swe-ch12]). Beck lists readable and specific among twelve test
  properties that trade off ([Test Desiderata][desiderata]).
- **Work in small steps.** Measured. In a study of 39 professionals, quality and
  productivity gains went with granularity and uniformity of steps; test-first
  order "had no important influence" ([Fucci et al.][fucci]). A meta-analysis of
  27 studies found that TDD has "a small positive effect on quality but little to
  no discernible effect on productivity" ([Rafique and Mišić][rafique]).
- **Never weaken or special-case a test; stop and explain.** Measured.
  - Claude models cheat on impossible tasks mainly "through modifying test cases"
    (more than 79% of their cheating) ([ImpossibleBench][impossiblebench]).
  - A prompt that says to stop and explain incorrect tests cut cheating by GPT-5
    to 1%. A prompt with only "DO NOT MODIFY THE TESTS" left it above 85%
    ([ImpossibleBench][impossiblebench]).
  - In Anthropic's system card, an anti-hard-coding prompt cut Claude Opus 4 from
    47% to 5% on impossible tasks ([system cards][system-cards]).
  - For o3, METR found that such instructions had "a nearly negligible effect"
    ([METR][metr-rh]).

  The wording "stop and explain why" follows the prompts that worked. Read-only
  tests in the host are a stronger control.
- **Coverage is not proof.** Measured. Coverage "should not be used as a quality
  target" ([Inozemtseva and Holmes][inozemtseva]). At Google, engineers treat an
  80% bar "like a ceiling" ([SWE at Google, ch. 11][swe-ch11]).
- **Native checks first for tooling and infrastructure.** Judgment. Linters,
  validators, dry runs, and plans already prove most tooling behavior. A
  dedicated test earns its maintenance cost only where they cannot.
- **Prompt text does not prove agent behavior.** Measured and practice. CLAUDE.md
  is "context, not enforced configuration" ([memory docs][cc-memory]). METR's
  result above shows that an instruction can have almost no effect.

### Write plainly

- **Plain language.** Practice. The IPLF definition, on which ISO 24495-1 is
  based, asks that readers "can easily find what they need, understand what they
  find, and use that information" ([IPLF][iplf], [ISO 24495-1][iso24495]).
- **About 80% of the way to ASD-STE100.** Judgment, with narrow measured support.
  - STE's maintainers say it "is not intended for general-purpose writing", but
    that "many of its principles... can be effectively adopted" ([ASD-STE100
    FAQ][ste]).
  - With 175 aircraft maintenance technicians, "comprehension was significantly
    improved with Simplified English, particularly... for non-native English
    speakers" ([Chervak et al.][chervak]).
  - Google and Microsoft recommend a conversational register, which a strict STE
    requirement would conflict with ([Google style][google-style]).
  - STE's maintainers warn that AI text "can appear clear, authoritative, and
    consistent with STE, even when it does not correctly apply the rules"
    ([STEMG][ste-ai]).

  The partial adoption keeps the sentence-level rules that the sources share and
  leaves out the controlled dictionary.
- **Match length to need.** Practice. Anthropic reports that files Claude Opus 5
  writes "are often longer than on prior models" ([Opus 5 prompting][opus5]).
- **Commit messages state the problem, the change, and why.** Practice. "Describe
  your problem" ([kernel][kernel-patches]). "Use the body to explain what and why
  vs. how" ([Beams][beams]). A change description says what changes and why
  ([Google][cl-descriptions]).

### Finish and report

- **State the finished outcome, not extra verification steps.** Practice.
  Anthropic: "If your prompt contains explicit verification instructions...
  remove them: instructions like these cause over-verification"
  ([Opus 5 prompting][opus5]). The closeout states what must be true of the final
  diff and the reported checks.
- **Lead with what changed; keep small reports short.** Judgment and local
  evidence. Agents that used an earlier version of this guidance "wrapped small
  changes in long reports" (commit `e22616f`).
- **Report results as they are.** Measured. Models are sycophantic: Claude 1.3
  "wrongly admits mistakes on 98% of questions" when the user asks "Are you
  sure?" ([Sharma et al.][sharma]). Claude Code's guidance asks Claude to "show
  evidence rather than asserting success" ([best practices][cc-best]).

## Advice not adopted

These popular rules were considered and left out.

- **Fixed limits on function length or cyclomatic complexity.** The size data
  above does not support them. McCabe called his limit of 10 "reasonable, but not
  magical" ([McCabe][mccabe]). At method level, cyclomatic complexity explains only
  part of the variance in lines of code ([Landman et al.][landman]), so the
  metric is not redundant, but no study here ties a threshold to defects.
- **Coverage targets.** See *Coverage is not proof*.
- **A required assertion density.** The authors of the main study argue against
  it ([Kudrjavets et al.][kudrjavets]).
- **Test-first as a rule for all work.** The order of test and code did not matter
  in the TDD studies above ([Fucci et al.][fucci]), and Ousterhout objects that TDD
  "discourages design thinking" ([debate][aposd-debate]). The `behavior-implement`
  skill keeps red-green-refactor as its own workflow.
- **"Be liberal in what you accept."** RFC 9413 rejects that reading of the
  robustness principle ([RFC 9413][rfc9413]).
- **"A comment is always a failure."** Martin's position. Even he accepts
  comments for public APIs ([debate][aposd-debate]).
- **Circuit breakers by default.** They "introduce modal behavior into systems
  that can be difficult to test" ([AWS][aws-retries]). The guidance limits retries
  instead.
- **Fixed test-pyramid ratios.** Google's own sources give 80/15/5
  ([SWE at Google, ch. 11][swe-ch11]) and 70/20/10 ([Wacker][wacker]). Fowler's
  pyramid depends on the cost of high-level tests ([TestPyramid][fowler-pyramid]).
- **CRAP thresholds as a gate.** The authors call the metric "highly experimental"
  ([Savoia][crap4j]).
- **Emphasis markers and generic "verify your work" steps.** See *How the guidance
  is built* and *Finish and report*.
- **Zero technical debt and zero dependencies.** These TigerStyle policies fit a
  safety-critical database with its own toolchain ([TigerStyle][tigerstyle]). They
  do not fit general work.

## Change the guidance

- Before you add a rule, name the failure it prevents, and confirm that agents
  show that failure without the rule. Remove a rule when agents follow it
  without being told.
- Keep `standards.md` under 200 lines.
- Use "never" and "always" only for hard limits.
- When you change a rule, update its entry here with its label and sources.
- Text in a prompt does not prove behavior. To test a change, run the same tasks
  with and without it and compare the results.

## References

Agent behavior and instruction files

- Anthropic, [Claude Code memory][cc-memory] and [best practices][cc-best]
- Anthropic, [Prompting best practices][claude-prompting] and [Prompting Claude Opus 5][opus5]
- Anthropic, [Effective context engineering for AI agents][ctx-eng], 2025
- Anthropic, [System cards][system-cards]
- OpenAI, [Custom instructions with AGENTS.md][codex-agents]
- Gloaguen et al., [Evaluating AGENTS.md][gloaguen], 2026
- Jaroslawicz et al., [How Many Instructions Can LLMs Follow at Once?][ifscale], 2025
- Liu et al., [Lost in the Middle][lost-middle], TACL 2024
- Spracklen et al., [We Have a Package for You!][spracklen], USENIX Security 2025
- Pearce et al., [Asleep at the Keyboard?][pearce], IEEE S&P 2022
- Perry et al., [Do Users Write More Insecure Code with AI Assistants?][perry], CCS 2023
- Zhong, Raghunathan, Carlini, [ImpossibleBench][impossiblebench], 2025
- METR, [Recent Frontier Models Are Reward Hacking][metr-rh], 2025
- Sharma et al., [Towards Understanding Sycophancy in Language Models][sharma], ICLR 2024
- DORA, [2024 report][dora-2024], [2025 report][dora-2025], [Working in small batches][dora-batches]

Design and complexity

- Parnas, [On the Criteria To Be Used in Decomposing Systems into Modules][parnas72], CACM 1972
- Parnas, [Designing Software for Ease of Extension and Contraction][parnas79], IEEE TSE 1979
- Brooks, [No Silver Bullet: Essence and Accidents of Software Engineering][brooks], 1986
- Moseley and Marks, [Out of the Tar Pit][tarpit], 2006
- Ousterhout, [A Philosophy of Software Design][aposd], 2nd ed., 2021; [book extract][aposd-extract]; [discussion with Robert Martin][aposd-debate]
- Hickey, [Simple Made Easy][hickey] (transcript), 2011
- Fowler, [Yagni][yagni], 2015
- Metz, [The Wrong Abstraction][metz], 2016
- Hunt and Thomas, [Orthogonality and the DRY Principle][dry] (interview), 2003
- Beck, [make the change easy][beck-tweet], 2012; [BPlusTree3 CLAUDE.md][beck-claude], 2025
- Cockburn, [Hexagonal Architecture][cockburn], 2005
- Knuth, [Structured Programming with go to Statements][knuth74], ACM Computing Surveys 1974
- Minsky, [Effective ML Revisited][minsky], 2011
- King, [Parse, don't validate][king], 2019
- Momot et al., [The Seven Turrets of Babel][langsec], IEEE SecDev 2016
- Ajmani, [Package names][go-pkg-names], Go blog, 2015
- Pike, [Go Proverbs][go-proverbs], 2015
- TigerBeetle, [TigerStyle][tigerstyle]
- Holzmann, [The Power of Ten][p10], IEEE Computer 2006
- Cox, [Our Software Dependency Problem][cox-deps], 2019

Size, metrics, and defects

- Basili and Perricone, [Software Errors and Complexity][basili84], CACM 1984
- Fenton and Ohlsson, [Quantitative analysis of faults and failures in a complex software system][fenton-ohlsson], IEEE TSE 2000
- El Emam et al., [The Confounding Effect of Class Size][elemam], IEEE TSE 2001
- Hatton, [Reexamining the Fault Density-Component Size Connection][hatton], IEEE Software 1997
- McCabe, [A Complexity Measure][mccabe], IEEE TSE 1976
- Landman et al., [Empirical analysis of the relationship between CC and SLOC][landman], JSEP 2016
- Kudrjavets, Nagappan, Ball, [Assessing the Relationship between Software Assertions and Faults][kudrjavets], ISSRE 2006
- Casalnuovo et al., [Assert Use in GitHub Projects][casalnuovo], ICSE 2015
- Savoia, [crap4j][crap4j], 2007
- NVD, [CVE-2024-3094][xz]; npm, [kik, left-pad, and npm][leftpad], 2016

Failure, state, and concurrency

- Yuan et al., [Simple Testing Can Prevent Most Critical Failures][yuan14], OSDI 2014
- Saltzer, Reed, Clark, [End-to-End Arguments in System Design][e2e], ACM TOCS 1984
- Amazon Builders' Library: [Timeouts, retries, and backoff with jitter][aws-retries], [Making retries safe with idempotent APIs][aws-idempotent], [Avoiding fallback][aws-fallback], [Using load shedding][aws-shedding]
- Brooker, [Exponential Backoff and Jitter][jitter], 2015
- Google, SRE book: [Addressing Cascading Failures][sre-cascading], [Handling Overload][sre-overload], [Monitoring Distributed Systems][sre-monitoring]
- Nygard, [Release It!, Integration Points excerpt][nygard], 2018
- Leach, [Implementing Stripe-like Idempotency Keys in Postgres][brandur], 2017
- Helland, [Life beyond Distributed Transactions][helland07], CIDR 2007
- Treat, [You Cannot Have Exactly-Once Delivery][treat], 2015
- Kleppmann, [Distributed Systems lecture notes][kleppmann-notes], 2021
- Richardson, [Transactional outbox][outbox]
- Kingsbury, [Jepsen: PostgreSQL 12.3][jepsen-pg], 2020
- Smith, [Notes on structured concurrency][structured-concurrency], 2018
- Go: [Pipelines and cancellation][go-pipelines], [context package][go-context]

Security

- Saltzer and Schroeder, [The Protection of Information in Computer Systems][saltzer75], Proc. IEEE 1975
- OWASP, [Input Validation][owasp-input], [Logging][owasp-logging], [SSRF Prevention][owasp-ssrf], [Top 10:2025][owasp-top10]
- MITRE, [2025 CWE Top 25][cwe-top25]

Testing, review, and change size

- Winters, Manshreck, Wright, [Software Engineering at Google][swe-book], 2020: [ch. 11][swe-ch11], [ch. 12][swe-ch12], [ch. 13][swe-ch13]
- Beck, [Test Desiderata][desiderata], 2019
- Khorikov, [When to Mock][khorikov-mock], 2020
- Fowler, [Mocks Aren't Stubs][mocks-stubs], 2007; [TestPyramid][fowler-pyramid], 2012
- Vocke, [The Practical Test Pyramid][vocke], 2018
- Wacker, [Just Say No to More End-to-End Tests][wacker], 2015
- Inozemtseva and Holmes, [Coverage Is Not Strongly Correlated with Test Suite Effectiveness][inozemtseva], ICSE 2014
- Petrović et al., [Does Mutation Testing Improve Testing Practices?][petrovic21], ICSE 2021
- Luo et al., [An Empirical Analysis of Flaky Tests][luo14], FSE 2014
- Micco, [Flaky Tests at Google and How We Mitigate Them][micco], 2016
- Fucci et al., [A Dissection of the Test-Driven Development Process][fucci], IEEE TSE 2017
- Rafique and Mišić, [The Effects of Test-Driven Development][rafique], IEEE TSE 2013
- Google, [Engineering Practices][eng-practices]: [Small CLs][small-cls], [CL descriptions][cl-descriptions]
- Sadowski et al., [Modern Code Review: A Case Study at Google][sadowski], ICSE-SEIP 2018
- Bacchelli and Bird, [Expectations, Outcomes, and Challenges of Modern Code Review][bacchelli], ICSE 2013
- Kononenko et al., [Investigating Code Review Quality][kononenko], ICSME 2015
- Linux kernel, [Submitting patches][kernel-patches]; Beams, [How to Write a Git Commit Message][beams], 2014

Contracts

- Wright, [Hyrum's Law][hyrum]
- Sato, [ParallelChange][parallel-change], 2014
- Xu, [Online migrations at scale][stripe-migrations], Stripe, 2017
- Thomson and Schinazi, [RFC 9413: Maintaining Robust Protocols][rfc9413], 2023
- [Protocol Buffers dos and don'ts][protobuf-dos]; [AIP-180 Backwards compatibility][aip180]
- [Go 1 compatibility][go1compat]; [Kubernetes deprecation policy][k8s-deprecation]

Naming, comments, and writing

- Feitelson et al., [How Developers Choose Names][feitelson], IEEE TSE 2022
- Gerrand, [What's in a name?][gerrand], 2014; Google, [Go style decisions][go-style]
- Hofmeister, Siegmund, Holt, [Shorter Identifier Names Take Longer to Comprehend][hofmeister], SANER 2017
- Schankin et al., [Descriptive Compound Identifier Names Improve Source Code Comprehension][schankin], ICPC 2018
- Lawrie et al., [Effective identifier names for comprehension and memory][lawrie07], 2007
- Beniamini et al., [Meaningful Identifier Names: The Case of Single-Letter Variables][beniamini], ICPC 2017
- Tan et al., [/* iComment: Bugs or Bad Comments? */][icomment], SOSP 2007
- Wen et al., [A Large-Scale Empirical Study on Code-Comment Inconsistencies][wen19], ICPC 2019
- Abdelsalam et al., [The Effect of Comments on Program Comprehension][abdelsalam], EMSE 2026
- Linux kernel, [Coding style][kernel-style]
- ASD-STE100, [FAQ][ste] and [STE and AI white paper][ste-ai], 2026
- Chervak, Drury, Ouellette, [Simplified English for aircraft workcards][chervak], HFES 1996
- IPLF, [Plain language definitions][iplf]; ISO, [24495-1:2023][iso24495]
- Google, [Developer documentation style guide][google-style]

[cc-memory]: https://code.claude.com/docs/en/memory
[cc-best]: https://code.claude.com/docs/en/best-practices
[claude-prompting]: https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/claude-prompting-best-practices
[opus5]: https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/prompting-claude-opus-5
[ctx-eng]: https://www.anthropic.com/engineering/effective-context-engineering-for-ai-agents
[system-cards]: https://www.anthropic.com/system-cards
[codex-agents]: https://learn.chatgpt.com/docs/agent-configuration/agents-md
[gloaguen]: https://arxiv.org/abs/2602.11988
[ifscale]: https://arxiv.org/abs/2507.11538
[lost-middle]: https://aclanthology.org/2024.tacl-1.9/
[spracklen]: https://www.usenix.org/conference/usenixsecurity25/presentation/spracklen
[pearce]: https://arxiv.org/abs/2108.09293
[perry]: https://doi.org/10.1145/3576915.3623157
[impossiblebench]: https://arxiv.org/abs/2510.20270
[metr-rh]: https://metr.org/blog/2025-06-05-recent-reward-hacking/
[sharma]: https://arxiv.org/abs/2310.13548
[dora-2024]: https://services.google.com/fh/files/misc/2024_final_dora_report.pdf
[dora-2025]: https://services.google.com/fh/files/misc/2025_state_of_ai_assisted_software_development.pdf
[dora-batches]: https://dora.dev/capabilities/working-in-small-batches/
[parnas72]: https://doi.org/10.1145/361598.361623
[parnas79]: https://doi.org/10.1109/TSE.1979.234169
[brooks]: https://www.cs.unc.edu/techreports/86-020.pdf
[tarpit]: https://curtclifton.net/papers/MoseleyMarks06a.pdf
[aposd]: https://web.stanford.edu/~ouster/cgi-bin/aposd.php
[aposd-extract]: https://web.stanford.edu/~ouster/cgi-bin/aposd2ndEdExtract.pdf
[aposd-debate]: https://github.com/johnousterhout/aposd-vs-clean-code
[hickey]: https://github.com/matthiasn/talk-transcripts/blob/master/Hickey_Rich/SimpleMadeEasy-mostly-text.md
[yagni]: https://martinfowler.com/bliki/Yagni.html
[metz]: https://sandimetz.com/blog/2016/1/20/the-wrong-abstraction
[dry]: https://www.artima.com/articles/orthogonality-and-the-dry-principle
[beck-tweet]: https://x.com/KentBeck/status/250733358307500032
[beck-claude]: https://github.com/KentBeck/BPlusTree3/blob/main/rust/docs/CLAUDE.md
[cockburn]: https://alistair.cockburn.us/hexagonal-architecture/
[knuth74]: https://doi.org/10.1145/356635.356640
[minsky]: https://blog.janestreet.com/effective-ml-revisited/
[king]: https://lexi-lambda.github.io/blog/2019/11/05/parse-don-t-validate/
[langsec]: https://langsec.org/papers/langsec-cwes-secdev2016.pdf
[go-pkg-names]: https://go.dev/blog/package-names
[go-proverbs]: https://go-proverbs.github.io/
[tigerstyle]: https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/TIGER_STYLE.md
[p10]: https://spinroot.com/gerard/pdf/P10.pdf
[cox-deps]: https://research.swtch.com/deps
[basili84]: https://doi.org/10.1145/69605.2085
[fenton-ohlsson]: https://doi.org/10.1109/32.879815
[elemam]: https://doi.org/10.1109/32.935855
[hatton]: https://leshatton.org/IEEE_Soft_97b.html
[mccabe]: https://doi.org/10.1109/TSE.1976.233837
[landman]: https://doi.org/10.1002/smr.1760
[kudrjavets]: https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/tr-2006-54.pdf
[casalnuovo]: https://www.cs.ucdavis.edu/~filkov/papers/assert-main.pdf
[crap4j]: https://www.artima.com/weblogs/viewpost.jsp?thread=215899
[xz]: https://nvd.nist.gov/vuln/detail/CVE-2024-3094
[leftpad]: https://blog.npmjs.org/post/141577284765/kik-left-pad-and-npm
[yuan14]: https://www.usenix.org/conference/osdi14/technical-sessions/presentation/yuan
[e2e]: https://web.mit.edu/Saltzer/www/publications/endtoend/endtoend.pdf
[aws-retries]: https://builder.aws.com/content/3EumjoZascWd1oZiEgL8ORlv3qE/timeouts-retries-and-backoff-with-jitter
[aws-idempotent]: https://aws.amazon.com/builders-library/making-retries-safe-with-idempotent-APIs/
[aws-fallback]: https://builder.aws.com/content/3EuS9Sakq7L3VLQIF3qzfMfke1Y/avoiding-fallback-in-distributed-systems
[aws-shedding]: https://builder.aws.com/content/3Eun1EEyX6p2e3VYNyRLSJzLuMV/using-load-shedding-to-avoid-overload
[jitter]: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
[sre-cascading]: https://sre.google/sre-book/addressing-cascading-failures/
[sre-overload]: https://sre.google/sre-book/handling-overload/
[sre-monitoring]: https://sre.google/sre-book/monitoring-distributed-systems/
[nygard]: https://media.pragprog.com/titles/mnee2/antipatterns.pdf
[brandur]: https://brandur.org/idempotency-keys
[helland07]: https://www.cidrdb.org/cidr2007/papers/cidr07p15.pdf
[treat]: https://bravenewgeek.com/you-cannot-have-exactly-once-delivery/
[kleppmann-notes]: https://www.cl.cam.ac.uk/teaching/2122/ConcDisSys/dist-sys-notes.pdf
[outbox]: https://microservices.io/patterns/data/transactional-outbox.html
[jepsen-pg]: https://jepsen.io/analyses/postgresql-12.3
[structured-concurrency]: https://vorpus.org/blog/notes-on-structured-concurrency-or-go-statement-considered-harmful/
[go-pipelines]: https://go.dev/blog/pipelines
[go-context]: https://pkg.go.dev/context
[saltzer75]: https://web.mit.edu/Saltzer/www/publications/protection/Basic.html
[owasp-input]: https://cheatsheetseries.owasp.org/cheatsheets/Input_Validation_Cheat_Sheet.html
[owasp-logging]: https://cheatsheetseries.owasp.org/cheatsheets/Logging_Cheat_Sheet.html#data-to-exclude
[owasp-ssrf]: https://cheatsheetseries.owasp.org/cheatsheets/Server_Side_Request_Forgery_Prevention_Cheat_Sheet.html
[owasp-top10]: https://top10.owasp.org/2025/
[cwe-top25]: https://cwe.mitre.org/top25/archive/2025/2025_cwe_top25.html
[swe-book]: https://abseil.io/resources/swe-book
[swe-ch11]: https://abseil.io/resources/swe-book/html/ch11.html
[swe-ch12]: https://abseil.io/resources/swe-book/html/ch12.html
[swe-ch13]: https://abseil.io/resources/swe-book/html/ch13.html
[desiderata]: https://testdesiderata.com/
[khorikov-mock]: https://enterprisecraftsmanship.com/posts/when-to-mock/
[mocks-stubs]: https://martinfowler.com/articles/mocksArentStubs.html
[fowler-pyramid]: https://martinfowler.com/bliki/TestPyramid.html
[vocke]: https://martinfowler.com/articles/practical-test-pyramid.html
[wacker]: https://testing.googleblog.com/2015/04/just-say-no-to-more-end-to-end-tests.html
[inozemtseva]: https://doi.org/10.1145/2568225.2568271
[petrovic21]: https://arxiv.org/abs/2103.07189
[luo14]: https://mir.cs.illinois.edu/marinov/publications/LuoETAL14FlakyTestsAnalysis.pdf
[micco]: https://testing.googleblog.com/2016/05/flaky-tests-at-google-and-how-we.html
[fucci]: https://doi.org/10.1109/TSE.2016.2616877
[rafique]: https://doi.org/10.1109/TSE.2012.28
[eng-practices]: https://google.github.io/eng-practices/
[small-cls]: https://google.github.io/eng-practices/review/developer/small-cls.html
[cl-descriptions]: https://google.github.io/eng-practices/review/developer/cl-descriptions.html
[sadowski]: https://research.google/pubs/modern-code-review-a-case-study-at-google/
[bacchelli]: https://www.microsoft.com/en-us/research/publication/expectations-outcomes-and-challenges-of-modern-code-review/
[kononenko]: https://plg.uwaterloo.ca/~migod/papers/2015/icsme15-OleksiiOlgaLatifa.pdf
[kernel-patches]: https://docs.kernel.org/process/submitting-patches.html
[beams]: https://cbea.ms/git-commit/
[hyrum]: https://www.hyrumslaw.com/
[parallel-change]: https://martinfowler.com/bliki/ParallelChange.html
[stripe-migrations]: https://stripe.com/blog/online-migrations
[rfc9413]: https://www.rfc-editor.org/rfc/rfc9413
[protobuf-dos]: https://protobuf.dev/best-practices/dos-donts/
[aip180]: https://google.aip.dev/180
[go1compat]: https://go.dev/doc/go1compat
[k8s-deprecation]: https://kubernetes.io/docs/reference/deprecation-policy/
[feitelson]: https://doi.org/10.1109/TSE.2020.2976920
[gerrand]: https://go.dev/talks/2014/names.slide
[go-style]: https://google.github.io/styleguide/go/decisions
[hofmeister]: https://doi.org/10.1109/SANER.2017.7884623
[schankin]: https://doi.org/10.1145/3196321.3196332
[lawrie07]: https://doi.org/10.1007/s11334-007-0031-2
[beniamini]: https://doi.org/10.1109/ICPC.2017.18
[icomment]: https://doi.org/10.1145/1294261.1294276
[wen19]: https://www.inf.usi.ch/faculty/lanza/PUBS/P/Wen2019a.pdf
[abdelsalam]: https://doi.org/10.1007/s10664-025-10721-2
[kernel-style]: https://docs.kernel.org/process/coding-style.html
[ste]: https://www.asd-ste100.org/STE_faq.html
[ste-ai]: https://www.asd-ste100.org/assets/files/WhitePaper-ASD-STE100_and_AI.pdf
[chervak]: https://researchconnect.buffalo.edu/en/publications/simplified-english-for-aircraft-workcards/
[iplf]: https://www.iplfederation.org/plain-language-definitions/
[iso24495]: https://www.iso.org/standard/78907.html
[google-style]: https://developers.google.com/style/highlights
