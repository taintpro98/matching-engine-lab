#!/bin/bash
# Run Go engine benchmark

set -e
ENGINE=${1:-v1}
ROOT=$(cd "$(dirname "$0")/../.." && pwd)
cd "$ROOT/go"

go build ./cmd/runner
echo "Running Go engine_$ENGINE..."
time ./runner --engine "$ENGINE" < "$ROOT/datasets/sample/medium.ndjson" > /dev/null
