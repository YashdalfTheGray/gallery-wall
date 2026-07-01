#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT="$ROOT/web/public"

mkdir -p "$OUT"

echo "building layout.wasm..."
GOOS=js GOARCH=wasm go build -o "$OUT/layout.wasm" "$ROOT/wasm"

WASM_EXEC="$(go env GOROOT)/lib/wasm/wasm_exec.js"
cp "$WASM_EXEC" "$OUT/wasm_exec.js"

echo "wrote $OUT/layout.wasm and wasm_exec.js"
