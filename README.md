# matching-engine-lab

A research lab project to compare multiple matching engine implementations in Rust and Go.

**This is NOT a crypto exchange.**

## Domain

- One Asset, One Money
- Sellers offer Asset, Buyers spend Money
- price = Money per Asset
- Operations: PlaceSell, CancelSell, BuyByQty, BuyByBudget
- Matching: Cheapest-first, FIFO, deterministic, integer arithmetic only

## Structure

- `docs/` — Interface, architecture, data structures, benchmarks
- `spec/` — JSON schemas and NDJSON format
- `datasets/` — Sample and generated test data
- `tools/` — Generator, verifier, benchmark scripts
- `rust/` — Rust implementations (v1 BTreeMap, v2 TreeMap/RB)
- `go/` — Go implementations (v1 BTree, v2 TreeMap/RB)

## Quick Start

```bash
# Rust
cd rust && cargo build

# Go
cd go && go build ./...
```

## Run

```bash
# Rust runner
cargo run -p runner -- --engine v1 < input.ndjson > output.ndjson

# Go runner
go run ./cmd/runner --engine v1 < input.ndjson > output.ndjson

# With latency percentiles (p50, p99, p999)
cargo run -p runner -- --engine v1 --latency < input.ndjson > output.ndjson
go run ./cmd/runner --engine v1 --latency < input.ndjson > output.ndjson
```

## Load Test

```bash
# Full load test (10K, 100K, 1M commands across all engines)
./tools/bench/load_test.sh

# Custom: 50K commands, engines v1 and v2
./tools/bench/load_test.sh "50000" "v1 v2"

# Install hyperfine for better stats: brew install hyperfine
```
