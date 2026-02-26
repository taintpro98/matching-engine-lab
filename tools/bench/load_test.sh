#!/bin/bash
# Load test matching engines across Rust and Go implementations.
# Usage: ./load_test.sh [counts] [engines]
#   counts:  space-separated command counts (default: 10000 100000 1000000)
#   engines: space-separated engine names (default: v1 v2 v3)

set -e
ROOT=$(cd "$(dirname "$0")/../.." && pwd)
GEN="$ROOT/tools/generator/py/gen_stream.py"
DATASETS="$ROOT/datasets/generated"
RUST_DIR="$ROOT/rust"
GO_DIR="$ROOT/go"

COUNTS="${1:-10000 100000 1000000}"
ENGINES="${2:-v1 v2 v3}"
WARMUP=2
RUNS=5

# Check for hyperfine
HAS_HYPERFINE=false
if command -v hyperfine &>/dev/null; then
  HAS_HYPERFINE=true
fi

mkdir -p "$DATASETS"

echo "=== Matching Engine Load Test ==="
echo "Counts: $COUNTS | Engines: $ENGINES | Hyperfine: $HAS_HYPERFINE"
echo ""

# Build once
echo "Building Rust (release)..."
(cd "$RUST_DIR" && cargo build --release -p runner 2>/dev/null) || true

for count in $COUNTS; do
  echo ""
  echo "=========================================="
  echo "  Load: $count commands"
  echo "=========================================="

  # Generate dataset
  INPUT="$DATASETS/load_${count}.ndjson"
  if [[ ! -f "$INPUT" ]] || [[ $(wc -l < "$INPUT") -ne $count ]]; then
    echo "Generating $count commands..."
    (cd "$(dirname "$GEN")" && python3 gen_stream.py --count "$count" --profile default) > "$INPUT"
  fi

  for engine in $ENGINES; do
    echo ""
    echo "--- Rust $engine ---"
    if $HAS_HYPERFINE; then
      hyperfine --warmup "$WARMUP" --runs "$RUNS" \
        "sh -c 'cd $RUST_DIR && cargo run -p runner --release -- --engine $engine < $INPUT > /dev/null'" 2>/dev/null || \
      hyperfine --warmup "$WARMUP" --runs "$RUNS" \
        "sh -c 'cd $RUST_DIR && cargo run -p runner --release -- --engine $engine < $INPUT > /dev/null'"
    else
      echo "Run 1:"
      (cd "$RUST_DIR" && cargo run -p runner --release -- --engine "$engine" < "$INPUT" > /dev/null) 2>&1 | grep -E "Processed|ops/sec" || true
      for i in 2 3 4 5; do
        echo "Run $i:"
        (cd "$RUST_DIR" && cargo run -p runner --release -- --engine "$engine" < "$INPUT" > /dev/null) 2>&1 | grep -E "Processed|ops/sec" || true
      done
    fi
  done

  for engine in $ENGINES; do
    echo ""
    echo "--- Go $engine ---"
    if $HAS_HYPERFINE; then
      hyperfine --warmup "$WARMUP" --runs "$RUNS" \
        "sh -c 'cd $GO_DIR && go run ./cmd/runner --engine $engine < $INPUT > /dev/null'" 2>/dev/null || \
      hyperfine --warmup "$WARMUP" --runs "$RUNS" \
        "sh -c 'cd $GO_DIR && go run ./cmd/runner --engine $engine < $INPUT > /dev/null'"
    else
      echo "Run 1:"
      (cd "$GO_DIR" && go run ./cmd/runner --engine "$engine" < "$INPUT" > /dev/null) 2>&1 | grep -E "Processed|ops/sec" || true
      for i in 2 3 4 5; do
        echo "Run $i:"
        (cd "$GO_DIR" && go run ./cmd/runner --engine "$engine" < "$INPUT" > /dev/null) 2>&1 | grep -E "Processed|ops/sec" || true
      done
    fi
  done
done

echo ""
echo "=== Load test complete ==="
