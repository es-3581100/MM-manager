# Phase 0 CLI

This checkpoint adds an offline-first Go CLI around the existing App-Dir Matrix Core v0.2.0 artifact.

## Commands

```text
appdir-matrix normalize --in refs.txt --out normalized.json
appdir-matrix pin --in refs.txt --out project.json [--github-api ...] [--github-fixture-dir fixtures/github-live]
appdir-matrix build --project project.json --template web/app-dir-matrix-core-v0.2.0.html --out matrix.html
appdir-matrix verify --artifact matrix.html [--sha matrix.html.sha256]
appdir-matrix compare --left before.json --right after.json
```

`normalize` preserves human annotation text separately from repository/subpath/external references.

`pin` resolves repository metadata, chooses a consistent explicit ref when one is supplied or the actual default branch otherwise, fetches a recursive Git tree, hard-fails on GitHub truncation, verifies explicitly named subpaths against the pinned tree, and emits typed evidence records with `authority: none`. The resolver can point at a fake/local GitHub API for deterministic tests. With `--github-fixture-dir`, it replays a captured GitHub API session entirely offline, verifies every captured response SHA-256 before use, records capture provenance into the project ref pack, and hard-fails on any API request missing from the fixture rather than falling back to the network.

`build` treats the v0.2.0 HTML as a reusable UI/runtime template. It replaces the template repository manifest, build metadata, ref pack, and embedded snapshot with the project document, then writes an external SHA-256 sidecar. It does not infer semantic dependency edges from filenames.

`verify` reopens the generated HTML, parses the embedded project document and snapshot, validates evidence classes and authority, compares repository pins and every path/SHA/type, and optionally validates the SHA-256 sidecar.

`compare` reports repository pin changes plus added/removed/changed paths between two project documents.

## Current boundary

This is a Phase 0 proof slice. The real GitHub resolution path is now exercised through a captured public API session and deterministic offline replay. It still does not claim the complete Phase 0 exit gate because schema migration/version negotiation is not implemented, browser runtime smoke is environment-blocked, and the existing UI still owns its own JS rendering/export behavior.

## Captured real GitHub replay

The repository includes a real public GitHub capture under `fixtures/github-live/` for `octocat/Hello-World`. The fixture contains:

```text
fixtures/github-live/
├── capture.json
├── ref-glob.txt
└── octocat/Hello-World/
    ├── repo.json
    └── tree-master.json
```

`capture.json` is the request allowlist and integrity manifest. It records the exact request URI, response file, HTTP status, capture method, capture time, and SHA-256 of each raw response.

Run the complete networkless replay with:

```bash
bash scripts/github-live-replay.sh
# or
make github-live-replay
```

The script performs normalize → pin-from-capture twice → byte comparison → build → reopen/verify → project compare → JavaScript syntax check. An unrecorded API request is a hard error.
