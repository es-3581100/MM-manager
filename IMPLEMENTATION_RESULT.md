# App-Dir Matrix — Implementation Result

Date: 2026-09-29
Checkpoint: Phase 0 deterministic core-truth slice

## Baseline

Recovered the supplied v0.2.0 reusable Matrix core rather than replacing it. The existing core already contains the browser matrix/tree/graph runtime, browser-side ref-glob compilation, pinned-tree hydration/export, and CodeGraph semantic overlay. The builder handoff and its embedded roadmap/design source were preserved byte-for-byte in this workspace.

The supplied files were recovered into an isolated build workspace and the project has since been published to `es-3581100/MM-manager`. This chunk continues from that repository state rather than rebuilding it.

## Master plan

`BUILD_PLAN.md` contains the dependency-aware Phase 0–9 plan, per-phase objective/prerequisites/current evidence/missing work/contracts/migrations/tests/risks/boundaries/deployment/capital/non-goals/exit gate/status, plus the current next executable slice.

`PHASE_STATUS.yaml` is the machine-readable status ledger.

## Implemented

### Rights-reserved publication posture

- Added top-level `LICENSE` and `LEGAL.md` asserting **Copyright © 2026 es-3581100. ALL RIGHTS RESERVED**.
- Public disclosure is explicitly separated from permission to copy, modify, redistribute, host, train on, commercialize, create derivatives, or incorporate the work into another system.
- Commercial/application, patent, trademark/branding, trade-secret/confidentiality, and recipient-contract boundaries are recorded separately.
- Source headers and generated HTML provenance carry the rights notice so exported artifacts do not silently lose it.
- Earlier roadmap language describing a `free/open core` is superseded for licensing purposes: Community/basic_4_life is a product-tier concept, not an open-source grant.


- `app-dir-matrix.project/v1` typed project document and JSON Schema.
- Go validation enforcing complete/non-truncated repository trees and `authority: none` for all evidence.
- Ref-glob normalizer preserving annotations separately from repository/subpath/external-document references.
- GitHub resolver capable of default-branch resolution, one consistent explicit-ref override, recursive tree acquisition, truncation rejection, and explicit named-subpath verification.
- Deterministic HTML packager using the existing v0.2.0 core as the view/runtime template.
- Embedded project document + pinned snapshot in emitted HTML.
- SHA-256 sidecar generation and independent verification.
- Reopen verifier that checks project schema, snapshot presence, repository pins, every path/type/SHA, evidence authority, and optional sidecar hash.
- Project comparator reporting pin changes plus added/removed/changed paths.
- Regression/unit tests and an end-to-end Phase 0 smoke script.
- Local CLI commands: `normalize`, `pin`, `build`, `verify`, `compare`.
- Hash-checked `app-dir-matrix.github-fixture/v1` replay transport for real GitHub API captures.
- `pin --github-fixture-dir` offline replay with no network fallback and capture provenance embedded into the project ref pack.
- Real public `octocat/Hello-World` metadata/tree capture plus deterministic end-to-end replay script.

## Verification

`go test ./...`: PASS.

`go vet ./...`: PASS.

`scripts/smoke.sh`: PASS.

Two independent builds of the same fixture are byte-identical. Current deterministic fixture artifact SHA-256:

`29d5638725fc2fb27d84f3791479bbbaf61d85b3fbb19dc30d7e54d82622d33b`

Generated runtime JavaScript passes `node --check`.

A non-identical drift fixture correctly reports the changed tree pin, added `docs`, and changed `README.md`.

Headless Chromium runtime smoke is blocked in the current container by Chromium/DBus startup timeouts, so the browser gate is not claimed as passed.

## Design

No new visual layer was added in this slice. The existing technical instrument-panel UI is retained as the client/view, avoiding decorative expansion before the Phase 0 data contract stabilizes.

Design authority remains the supplied Anti-AI reference. Its five gates are carried into `BUILD_PLAN.md`; no Phase 0 completion claim depends on unverified visual changes.

## Phase status

```text
Phase 0  partial
Phase 1  partial (existing structural UI only)
Phase 2  candidate (MCP-shaped ingest contract; no query server)
Phase 3  not_started
Phase 4  not_started
Phase 5  not_started
Phase 6  not_started
Phase 7  not_started
Phase 8  not_started
Phase 9  not_started
```

## Capital status

Only the no-cost Community/basic_4_life product posture is supported by current implementation; repository source remains all-rights-reserved and is not open-source licensed. Personal, Pro, Team, pack marketplace, hosted compute, and Enterprise remain roadmap metadata with no billing or entitlement implementation.

## Deviations

The prompt allowed existing equivalent bookkeeping artifacts to be updated; none existed in the supplied file set, so this checkpoint created `BUILD_PLAN.md`, `PHASE_STATUS.yaml`, `DECISIONS.md`, `VERIFICATION.md`, and `IMPLEMENTATION_RESULT.md`.

The Phase 0 tool is Go rather than a browser-only script. This is an implementation choice to make the existing HTML core a deterministic client of portable project truth rather than the sole authority implementation.

## Remaining blockers

- Rerun generated HTML in a healthy browser/headless environment and capture an actual runtime DOM/screenshot smoke.
- Add explicit project schema migration/version-negotiation tooling.
- Add import of existing v0.2.0 exported browser snapshots/manifests into `project/v1`.
- Decide the canonical Go-layer manifest/context export shape before Phase 2.

## Next executable slice

Implement explicit `project/v1` schema-version negotiation and a bounded migration framework. It must reject unknown future versions without guessing, preserve provenance through migration, and prove migrate → deterministic build → reopen/verify using fixtures. Phase 1 expansion stays blocked until that migration evidence and browser smoke pass.

## 2026-09-29 real-source replay increment

The recorded next slice is now complete. A real public GitHub API capture for `octocat/Hello-World` preserves repository metadata plus the complete recursive `master` tree as raw JSON response fixtures. The fixture manifest hashes each response and acts as a request allowlist.

`pin --github-fixture-dir fixtures/github-live` replays those exact responses entirely offline. Any unrecorded request fails; there is no automatic network fallback. Response tampering is detected before JSON parsing.

The replay independently reproduced the same `project/v1` bytes twice, verified the named `README` subpath, rendered the Matrix artifact, reopened it, validated the sidecar SHA-256, compared the duplicate projects as equal, and passed JavaScript syntax validation. Phase 0 remains `partial` only because browser-runtime smoke and schema migration/version negotiation are still unresolved.
