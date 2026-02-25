#!/bin/bash
# Run Rust engine benchmark

set -e
ENGINE=${1:-v1}
ROOT=$(cd "$(dirname "$0")/../.." && pwd)
cd "$ROOT/rust"

cargo build --release -p runner 2>/dev/null || cargo build -p runner
echo "Running Rust engine_$ENGINE..."
time cargo run -p runner --release -- --engine "$ENGINE" < "$ROOT/datasets/sample/medium.ndjson" > /dev/null
