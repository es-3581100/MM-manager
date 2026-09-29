#!/usr/bin/env bash
# Copyright © 2026 es-3581100. ALL RIGHTS RESERVED. See LICENSE and LEGAL.md.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
mkdir -p dist
rm -f dist/phase0-smoke-a.html* dist/phase0-smoke-b.html* dist/ref-glob.normalized.json dist/verify.json dist/compare-identical.json

go test ./...
go vet ./...
go run ./cmd/appdir-matrix normalize --in fixtures/ref-glob.txt --out dist/ref-glob.normalized.json
go run ./cmd/appdir-matrix build --project fixtures/project-v1.json --template web/app-dir-matrix-core-v0.2.0.html --out dist/phase0-smoke-a.html
go run ./cmd/appdir-matrix build --project fixtures/project-v1.json --template web/app-dir-matrix-core-v0.2.0.html --out dist/phase0-smoke-b.html
cmp dist/phase0-smoke-a.html dist/phase0-smoke-b.html
cmp dist/phase0-smoke-a.html.sha256 dist/phase0-smoke-b.html.sha256 || {
  # Sidecars contain the output filename; compare digest field only.
  test "$(awk '{print $1}' dist/phase0-smoke-a.html.sha256)" = "$(awk '{print $1}' dist/phase0-smoke-b.html.sha256)"
}
go run ./cmd/appdir-matrix verify --artifact dist/phase0-smoke-a.html --sha dist/phase0-smoke-a.html.sha256 > dist/verify.json
go run ./cmd/appdir-matrix compare --left fixtures/project-v1.json --right fixtures/project-v1.json > dist/compare-identical.json

grep -q '"ok": true' dist/verify.json
grep -q '"equal": true' dist/compare-identical.json
grep -q '"repo": "fixture/alpha"' dist/ref-glob.normalized.json
grep -q 'embeddedProjectDocument' dist/phase0-smoke-a.html
! grep -q 'donCannoli-burns/kol-agent-sandbox' dist/phase0-smoke-a.html

echo "SMOKE PASS"
echo "artifact_sha256=$(sha256sum dist/phase0-smoke-a.html | awk '{print $1}')"
