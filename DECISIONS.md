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