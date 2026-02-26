#!/bin/bash
# Compare Rust vs Go benchmark results

set -e
ROOT=$(cd "$(dirname "$0")/../.." && pwd)
INPUT="${1:-$ROOT/datasets/sample/medium.ndjson}"

echo "=== Rust ==="
for e in v1 v2; do
  echo "--- engine_$e ---"
  cd "$ROOT/rust" && cargo run -p runner --release -- --engine "$e" < "$INPUT" 2>&1 | tail -3
done

echo ""
echo "=== Go ==="
for e in v1 v2; do
  echo "--- engine_$e ---"
  cd "$ROOT/go" && go run ./cmd/runner --engine "$e" < "$INPUT" 2>&1 | tail -3
done
