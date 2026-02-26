# Matching Engine Benchmark Results

## Methodology

- **Dataset**: 100,000 commands (default profile: PlaceSell, CancelSell, BuyByQty, BuyByBudget)
- **Measurement**: Per-command latency (ns), throughput (ops/sec)
- **Build**: Rust `--release`, Go `go run` (no pre-build)
- **Date**: 2026-02-26

## Results Summary

| Engine | Language | Data Structure | Ops/sec | p50 (ns) | p99 (ns) | p999 (ns) |
|--------|----------|----------------|---------|----------|----------|-----------|
| v1 | Rust | BTreeMap | 479,309 | 167 | 708 | 1,167 |
| v2 | Rust | TreeMap (RB) | 424,439 | 250 | 1,084 | 1,875 |
| v1 | Go | BTree | 327,056 | 291 | 1,500 | 5,250 |
| v2 | Go | TreeMap (RB) | 317,317 | 292 | 2,917 | 8,166 |

## Findings

### Rust vs Go

- **Throughput**: Rust ~1.4–1.5× faster than Go
- **Latency**: Rust p50 ~1.7× lower, p99 ~2–3× lower
- **Tail latency**: Rust p999 is ~2–7× lower; Go shows higher variance (likely GC)

### BTree vs TreeMap (within same language)

- **Rust**: BTreeMap (v1) ~13% faster than RB TreeMap (v2)
- **Go**: BTree (v1) ~3% faster than TreeMap (v2)

### Data Structure Comparison

| Structure | Rust | Go |
|-----------|------|-----|
| BTree | `std::collections::BTreeMap` | `github.com/google/btree` |
| TreeMap | `rb_tree::RBMap` (red-black) | `github.com/emirpasic/gods` treemap (red-black) |

Both BTree implementations outperform their TreeMap counterparts in this workload, with Rust benefiting from zero-cost abstractions and no GC.

## How to Reproduce

```bash
# Generate 100K commands
python3 tools/generator/py/gen_stream.py --count 100000 > datasets/generated/load_100000.ndjson

# Run with latency
cargo run -p runner --release --manifest-path rust/Cargo.toml -- --engine v1 --latency < datasets/generated/load_100000.ndjson
go run ./cmd/runner --engine v1 --latency < datasets/generated/load_100000.ndjson

# Run full load test
./tools/bench/load_test.sh "100000" "v1 v2"
```
