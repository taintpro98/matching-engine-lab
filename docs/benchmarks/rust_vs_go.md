# Rust vs Go Benchmarks

## Metrics

- **Ops/sec** — commands processed per second
- **Latency** — p50, p99, p999 (use `--latency` flag)
- **Memory** — peak RSS, allocations

## Methodology

- Same NDJSON input for both
- Warm-up run, then timed run
- Multiple iterations, report median

## Scripts

- `tools/bench/run_rust.sh` — run Rust engines
- `tools/bench/run_go.sh` — run Go engines
- `tools/bench/compare.sh` — compare results
- `tools/bench/load_test.sh` — full load test (generates data, runs all engines)

## Load Test

```bash
./tools/bench/load_test.sh                    # 10K, 100K, 1M × v1,v2
./tools/bench/load_test.sh "100000" "v1 v2"   # 100K × v1,v2
```

Uses `hyperfine` if installed (`brew install hyperfine`) for mean/stddev; otherwise runs 5 iterations manually.

## Latency Percentiles

Add `--latency` to the runner for per-command latency (ns):

```bash
cargo run -p runner -- --engine v1 --latency < input.ndjson
go run ./cmd/runner --engine v1 --latency < input.ndjson
```
