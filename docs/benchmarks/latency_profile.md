# Latency Profile

## Measurement

- Per-command latency (submit → events returned)
- Batch latency (full stream processing)

## Tools

- Rust: `std::time::Instant`
- Go: `time.Since`
- Percentile aggregation in runner

## Targets

- p50 < 1µs for simple commands
- p99 < 10µs
- Identify outliers (e.g., large scans)
