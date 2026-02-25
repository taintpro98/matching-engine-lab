# Rust vs Go Benchmarks

## Metrics

- **Ops/sec** — commands processed per second
- **Latency** — p50, p99, p999
- **Memory** — peak RSS, allocations

## Methodology

- Same NDJSON input for both
- Warm-up run, then timed run
- Multiple iterations, report median

## Scripts

- `tools/bench/run_rust.sh` — run Rust engines
- `tools/bench/run_go.sh` — run Go engines
- `tools/bench/compare.sh` — compare results
