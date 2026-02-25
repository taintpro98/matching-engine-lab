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
- `rust/` — Rust implementations (BTreeMap, SkipList, Arena)
- `go/` — Go implementations (BTree, TreeMap, Pool)

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
```
