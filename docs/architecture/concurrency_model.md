# Concurrency Model

## Current Design

- **Single-threaded** — no concurrency inside the engine
- Each engine processes commands sequentially
- Runner reads input line-by-line, submits one command at a time

## Rationale

- Research focus: compare data structure performance, not concurrency
- Deterministic replay requires sequential processing
- Simplifies benchmarking and verification

## Future Considerations

- Parallel benchmark runs (multiple processes)
- Snapshot/load for checkpointing long runs
