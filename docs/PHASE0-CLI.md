# Phase 0 CLI

This checkpoint adds an offline-first Go CLI around the existing App-Dir Matrix Core v0.2.0 artifact.

## Commands

```text
appdir-matrix normalize --in refs.txt --out normalized.json
appdir-matrix pin --in refs.txt --out project.json [--github-api ...]
appdir-matrix build --project project.json --template web/app-dir-matrix-core-v0.2.0.html --out matrix.html
appdir-matrix verify --artifact matrix.html [--sha matrix.html.sha256]
appdir-matrix compare --left before.json --right after.json
```

`normalize` preserves human annotation text separately from repository/subpath/external references.

`pin` resolves repository metadata, chooses a consistent explicit ref when one is supplied or the actual default branch otherwise, fetches a recursive Git tree, hard-fails on GitHub truncation, verifies explicitly named subpaths against the pinned tree, and emits typed evidence records with `authority: none`. The resolver can point at a fake/local GitHub API for deterministic tests.

`build` treats the v0.2.0 HTML as a reusable UI/runtime template. It replaces the template repository manifest, build metadata, ref pack, and embedded snapshot with the project document, then writes an external SHA-256 sidecar. It does not infer semantic dependency edges from filenames.

`verify` reopens the generated HTML, parses the embedded project document and snapshot, validates evidence classes and authority, compares repository pins and every path/SHA/type, and optionally validates the SHA-256 sidecar.

`compare` reports repository pin changes plus added/removed/changed paths between two project documents.

## Current boundary

This is a Phase 0 proof slice. It does not yet claim the complete Phase 0 exit gate because live GitHub resolution has not been exercised in this checkpoint, schema migration/version negotiation is not implemented, browser runtime smoke is environment-blocked, and the existing UI still owns its own JS rendering/export behavior.
