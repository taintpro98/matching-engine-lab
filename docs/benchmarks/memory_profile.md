# Memory Profile

## Rust

- `cargo build --release` for optimized build
- `heaptrack` or `valgrind --tool=massif` for profiling
- `#[derive(Default)]` and reuse to reduce allocations

## Go

- `go build -gcflags="-m"` to see escape analysis
- `pprof` for heap profiling
- `sync.Pool` in enginev3_pool

## Targets

- Baseline: engine_v1
- Compare: v2, v3 allocation counts
