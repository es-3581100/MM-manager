# Verification

Date: 2026-09-29

## Source recovery / provenance

- Builder prompt SHA-256: `ef7c90be67bc49c48e75852bf23dc29c88625c27b73fb17f9a107deccd25a609`
- App-Dir Matrix Core v0.2.0 source SHA-256: `eb84b0403c240c519aa1ffb202dd916120433bc4e684ab1da0cf8c4f90d0bf7d`
- Extracted roadmap SHA-256: `873d02e7c4608fdda8e5b9803c104cceefa96321eafb5a87282be1384356e864`
- Extracted Anti-AI design reference SHA-256: `944e1d2a4c0fd6711155f1ba63fcc13db1d0dbfb88ee1696799981a0636fdac4`

These match the hashes declared by the supplied handoff.

## Go tests

Command:

```bash
go test ./...
```

Result: PASS.

Covered behaviors include:

- ref-glob Markdown/escaped URL normalization;
- duplicate repository collapse while preserving annotations;
- evidence authority inflation rejection;
- truncated project tree rejection;
- derived path-depth level contract;
- deterministic artifact generation;
- generated artifact reopen/verification;
- snapshot drift detection;
- default-branch resolution with fake GitHub API;
- recursive tree request with `recursive=1`;
- truncation rejection from resolver;
- named subpath verification and missing-path failure.

## Go vet

Command:

```bash
go vet ./...
```

Result: PASS.

## Phase 0 smoke

Command:

```bash
bash scripts/smoke.sh
```

Result: PASS.

The script:

1. runs tests and vet;
2. normalizes the fixture ref glob;
3. builds the same project twice;
4. confirms generated HTML bytes are identical;
5. confirms artifact SHA-256 digests are identical;
6. reopens/verifies the generated HTML + sidecar;
7. confirms identical-project compare is equal;
8. confirms no supplied template repository data leaked into the generated fixture artifact.

Deterministic smoke artifact SHA-256:

```text
29d5638725fc2fb27d84f3791479bbbaf61d85b3fbb19dc30d7e54d82622d33b
```

Verifier report:

```json
{
  "ok": true,
  "project_id": "phase0-smoke",
  "schema_version": "app-dir-matrix.project/v1",
  "repository_count": 1,
  "entry_count": 6,
  "artifact_sha256": "29d5638725fc2fb27d84f3791479bbbaf61d85b3fbb19dc30d7e54d82622d33b",
  "checks": [
    "embedded project document validates",
    "embedded snapshot parses",
    "repository pins match snapshot",
    "tree entries match project document",
    "evidence classes preserve authority=none",
    "sha256 sidecar matches artifact"
  ]
}
```

## Drift comparison

Command:

```bash
go run ./cmd/appdir-matrix compare \
  --left fixtures/project-v1.json \
  --right fixtures/project-v1-drift.json
```

Result: PASS; expected non-equality detected.

Observed drift:

- tree pin `tree-alpha-001` -> `tree-alpha-002`;
- added path `docs`;
- changed path `README.md`.

## Runtime JavaScript syntax

The executable JavaScript block was extracted from the generated HTML and checked with:

```bash
node --check dist/runtime.js
```

Result: PASS.

This check caught no JavaScript syntax errors after the generator replacement bug was repaired.

## Generator defect found during artifact inspection

The first implementation replaced a `const REPOS = ...;` value by searching for the first semicolon. A semicolon inside a repository boundary string caused stale template repository JSON to remain after the new value.

Artifact inspection exposed the defect. The replacement now operates on the complete JavaScript constant line, and `TestBuildVerifyAndDeterminism` fails if a known source-template repository leaks into the generated fixture artifact.

Status: FIXED + regression-tested.

## Browser runtime smoke

Attempted with Chromium 144 in headless mode using both DOM dump and screenshot paths.

Result: BLOCKED BY ENVIRONMENT.

Observed behavior: Chromium failed to complete startup before timeout and emitted repeated DBus connection errors. No screenshot or DOM artifact was produced. Because the browser gate did not pass, Phase 0 remains `partial` rather than `verified`.

This is not treated as evidence that the generated HTML is browser-correct; only the static HTML structure, embedded JSON, deterministic packaging, verifier, and JavaScript syntax have passed here.

## Build artifact hashes

```text
e672b82c15cda9a40f770efe2fc983d72b7a68ee8fce0f8832c5a4844aacb0b9  bin/appdir-matrix
29d5638725fc2fb27d84f3791479bbbaf61d85b3fbb19dc30d7e54d82622d33b  dist/phase0-smoke-a.html
781af06696e356ae8f1c81a543c7504a1fe9a6bb118c8ab047cd4845e00b7fb1  dist/ref-glob.normalized.json
d2422b42e1950ebfa27061b2596eb5fda1b96963f816a5d05432086708f09938  schema/app-dir-matrix-project.schema.json
```

## Rights notice verification

Verified after legal-posture integration:

```text
LICENSE sha256 34f316241d9745f988de176bb9dcecc02167a5bbf3c159fcf31c912d93136871
LEGAL.md sha256 637d9e09a2635965c04536a06ad55e3157c50cab17656476857d7c511fb730d0
README.md sha256 0b3ccc53ba50aa8ed11410c495a152f209f9ce4e91d824a248f2d3e10869ebb3
phase0 HTML sha256 29d5638725fc2fb27d84f3791479bbbaf61d85b3fbb19dc30d7e54d82622d33b
```

The generated HTML begins with an all-rights-reserved provenance comment and the reusable web template carries the same notice. `go test ./...`, `go vet ./...`, the deterministic smoke, reopen verification, and `node --check` remained passing after the change.


## Real public GitHub capture + offline replay

A live GitHub API capture was recorded on 2026-09-29 for the public repository `octocat/Hello-World`. GitHub reported default branch `master`; the recursive tree response was complete (`truncated: false`) at tree SHA `7fd1a60b01f91b314f59955a4e4d4e80d8edf11d` and contained the named `README` blob SHA `980a0d5f19a64b4b30a87d4206aade58726b60e3`.

Captured response integrity:

```text
607d2ca1df1d71a64fd45553c8a1360de7a52c35b8dbb2f6e9e69f1c345bb572  fixtures/github-live/octocat/Hello-World/repo.json
23f3f4cce7631d59dc56deaec224736de7813861f3de28ea6db9569fe6d4ea83  fixtures/github-live/octocat/Hello-World/tree-master.json
de810ed5b2386aa06505ca795115d3abe6ae55af4584538f89bdad2132520c87  fixtures/github-live/capture.json
```

Command:

```bash
bash scripts/github-live-replay.sh
```

Result: PASS.

The replay ran without network fallback and performed normalize → pin twice → byte comparison → deterministic HTML build → reopen verification → identical-project comparison → JavaScript syntax check. The generated project preserved GitHub capture provenance in `ref_pack.github_fixture`. The verifier reported one repository, one structural entry, matching pins/tree entries, `authority: none`, and a valid SHA-256 sidecar.

Additional regression tests:

- `TestCapturedGitHubFixtureReplaysOffline`: PASS.
- `TestGitHubFixtureHasNoNetworkFallback`: PASS.
- `TestGitHubFixtureRejectsTamperedResponse`: PASS.

This resolves the previous “live GitHub pin round-trip not yet recorded” blocker. It does not resolve browser-runtime smoke or schema migration/version-negotiation.
