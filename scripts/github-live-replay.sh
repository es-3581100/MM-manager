#!/usr/bin/env bash
# Copyright © 2026 es-3581100. ALL RIGHTS RESERVED. See LICENSE and LEGAL.md.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
OUT="dist/github-live-replay"
FIXTURE="fixtures/github-live"
mkdir -p "$OUT"

# This replay is intentionally networkless. The fixture transport hard-fails
# on any request that was not captured in capture.json.
go test ./...
go vet ./...

go run ./cmd/appdir-matrix normalize \
  --in "$FIXTURE/ref-glob.txt" \
  --out "$OUT/normalized.json"

go run ./cmd/appdir-matrix pin \
  --in "$FIXTURE/ref-glob.txt" \
  --out "$OUT/project-a.json" \
  --project-id github-live-replay \
  --build-id capture-20260929 \
  --github-fixture-dir "$FIXTURE" \
  2>"$OUT/pin-a.stderr"

go run ./cmd/appdir-matrix pin \
  --in "$FIXTURE/ref-glob.txt" \
  --out "$OUT/project-b.json" \
  --project-id github-live-replay \
  --build-id capture-20260929 \
  --github-fixture-dir "$FIXTURE" \
  2>"$OUT/pin-b.stderr"

cmp "$OUT/project-a.json" "$OUT/project-b.json"

python - <<'PY2'
import json
from pathlib import Path
p = json.loads(Path("dist/github-live-replay/project-a.json").read_text())
f = p.get("ref_pack", {}).get("github_fixture", {})
assert f.get("schema") == "app-dir-matrix.github-fixture/v1", f
assert f.get("repository") == "octocat/Hello-World", f
assert f.get("capture_method") == "GitHub API connector", f
PY2

go run ./cmd/appdir-matrix build \
  --project "$OUT/project-a.json" \
  --template web/app-dir-matrix-core-v0.2.0.html \
  --out "$OUT/mm-live-replay.html" \
  >"$OUT/build.stdout"

go run ./cmd/appdir-matrix verify \
  --artifact "$OUT/mm-live-replay.html" \
  --sha "$OUT/mm-live-replay.html.sha256" \
  >"$OUT/verify.json"

go run ./cmd/appdir-matrix compare \
  --left "$OUT/project-a.json" \
  --right "$OUT/project-b.json" \
  >"$OUT/compare.json"

python - <<'PY'
from pathlib import Path
src = Path("dist/github-live-replay/mm-live-replay.html").read_text()
start = src.rfind("<script>")
end = src.rfind("</script>")
if start < 0 or end <= start:
    raise SystemExit("runtime script not found")
Path("dist/github-live-replay/runtime.js").write_text(src[start + len("<script>"):end])
PY
node --check "$OUT/runtime.js"

sha256sum \
  "$FIXTURE/capture.json" \
  "$FIXTURE/octocat/Hello-World/repo.json" \
  "$FIXTURE/octocat/Hello-World/tree-master.json" \
  "$OUT/project-a.json" \
  "$OUT/mm-live-replay.html" \
  "$OUT/verify.json" \
  >"$OUT/checksums.txt"

printf 'GitHub live-capture offline replay PASS\n'
cat "$OUT/verify.json"
