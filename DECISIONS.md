# Decisions

## D-001 — Reuse v0.2.0 HTML as the Phase 0 view/runtime template

**Context:** A working Matrix core already exists with repository cards, tree/matrix/graph views, ref-glob compilation, pinned-tree hydration, and CodeGraph overlay.

**Decision:** Treat the supplied v0.2.0 HTML as a reusable client/view. Build deterministic project truth and verification around it rather than replacing it.

**Reason:** The prompt explicitly requires recovery before rebuilding. This preserves working behavior and limits early scope.

**Alternatives considered:** rewrite UI; split HTML immediately into a framework; rebuild in Go-App now.

**Consequence:** Phase 0 can stabilize the data contract before a UI/runtime migration.

**Evidence:** source core SHA-256 `eb84b0403c240c519aa1ffb202dd916120433bc4e684ab1da0cf8c4f90d0bf7d`.

## D-002 — Go standard-library CLI is the first durable implementation surface

**Context:** Phase 0 needs deterministic packaging, hashing, GitHub resolution, verification, and future portability.

**Decision:** Implement the first durable engine in Go with no third-party runtime dependencies.

**Reason:** A small static tool fits local-first/offline use, can later support Go-App/WASM research, and keeps the portable core provider-independent without granting source-code reuse rights.

**Alternatives considered:** Node-only build scripts; Python package; browser-only continuation.

**Consequence:** Browser UI becomes a client of a testable project document rather than the only implementation of project truth.

**Evidence:** `go test ./...`, `go vet ./...`, deterministic smoke artifact.

## D-003 — Evidence never carries execution authority

**Context:** The source contract separates annotations, repository metadata, Git trees, verified subpaths, semantic graphs, and derived display levels.

**Decision:** `app-dir-matrix.project/v1` requires every evidence record to have `authority: none`; validation rejects any other value.

**Reason:** Prevents context, similarity, static analysis, or human notes from turning into permissions.

**Alternatives considered:** optional authority fields; role-like evidence flags.

**Consequence:** Future execution systems must define a separate explicit authority contract.

**Evidence:** `TestValidateRejectsAuthorityInflation`.

## D-004 — Truncated recursive Git trees are a hard error

**Context:** A truncated tree cannot prove structural completeness.

**Decision:** The resolver and project validator reject truncated trees rather than presenting them as complete.

**Reason:** Evidence-before-claims and independent verification require a clear completeness boundary.

**Alternatives considered:** partial status with rendering; silent continuation.

**Consequence:** Large repositories need a future explicit fallback/resolution strategy before they can pass Phase 0 completeness.

**Evidence:** `TestGitHubResolverRejectsTruncation`, `TestValidateRejectsTruncatedTree`.

## D-005 — SHA-256 is external sidecar integrity, not a self-hash inside the HTML

**Context:** Embedding an artifact's own final SHA changes the bytes being hashed.

**Decision:** Write `<artifact>.sha256` next to the deterministic HTML and verify it independently.

**Reason:** Simple, conventional, reproducible, and avoids recursive self-hash semantics.

**Alternatives considered:** hash-excluded metadata section; Merkle/self-describing hash envelope.

**Consequence:** Artifact + sidecar form the current integrity pair.

**Evidence:** smoke verifier confirms sidecar matches generated artifact.

## D-006 — Failed headless Chromium smoke is recorded as environment-blocked, not hidden

**Context:** Chromium 144 in this container failed to complete `--dump-dom`/screenshot startup and emitted DBus-related errors/timeouts.

**Decision:** Keep JS syntax and artifact structure checks passing, record browser smoke as blocked, and do not call Phase 0 fully verified.

**Reason:** A failed runtime validation gate must block the stronger completion claim.

**Alternatives considered:** ignore browser smoke; substitute screenshots from another artifact.

**Consequence:** Next checkpoint should rerun browser smoke in a healthy desktop/headless environment.

## D-007 — Real GitHub captures become hash-checked offline regression fixtures

**Context:** Phase 0 needed evidence that default-branch resolution, recursive-tree pinning, and explicit subpath verification work against a real public GitHub response without making future tests depend on GitHub availability or mutable branch state.

**Decision:** Preserve real API response bytes behind `app-dir-matrix.github-fixture/v1`. Each recorded request is allowlisted by method + request URI and bound to a response file by SHA-256. `pin --github-fixture-dir` replays only those responses; unrecorded requests fail and no network fallback exists.

**Reason:** The same captured source evidence can be independently replayed, audited, and hash-checked after the live source changes or becomes unavailable.

**Alternatives considered:** live GitHub on every test run; hand-written reduced fixtures only; silently falling back to live GitHub when a fixture misses.

**Consequence:** Real-source evidence is durable and deterministic, but capture freshness remains explicit rather than pretending an old fixture is current live truth.

**Evidence:** `fixtures/github-live/capture.json`, `TestCapturedGitHubFixtureReplaysOffline`, `TestGitHubFixtureHasNoNetworkFallback`, `TestGitHubFixtureRejectsTamperedResponse`, `scripts/github-live-replay.sh`.


## D-008 — Matrix becomes the observatory for a semantic runtime

**Context:** The owner-provided Phase Zero.1 semantic-tier runtime demonstrated a stronger separation than the repository-only roadmap: evidence, prediction receipts, admission, scheduling/tier movement, and outcomes can be modeled independently and audited.

**Decision:** Treat MM-manager as three cooperating planes: Go evidence/package truth, semantic runtime/policy contracts, and the Matrix observatory. The Matrix is not itself the movement policy and runtime traces are not commands.

**Reason:** This preserves structural/provenance truth while allowing multiple predictors, admission policies, Q-table experiments, and storage tiers to evolve behind stable contracts.

**Boundary:** Structural path adjacency never creates a semantic dependency by itself. Imported mutability/recomputability stay UNKNOWN until explicitly evidenced. All imported evidence remains `authority: none`.

**Evidence:** `docs/PHASE0.2-SEMANTIC-RUNTIME.md`, `contracts/runtime/`, `runtime/semantic-tier-phase-zero-2/project-to-semantic.go`.

## D-009 — Configured policy evidence is distinct from observed actor events

**Decision:** Static warmsets and declared semantic rules use `CONFIGURED_POLICY`; `ACTOR_EVENT` is reserved for actually observed actor/workflow events.

**Reason:** A configuration statement and an observation have different provenance and must not be conflated.
