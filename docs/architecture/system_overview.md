# System Overview

The matching-engine-lab is a research project comparing multiple matching engine implementations across Rust and Go.

## Components

- **Engine Core** — Shared types, commands, events, engine trait/interface
- **Engine Implementations** — v1 (BTreeMap/BTree), v2 (SkipList/TreeMap), v3 (Arena/Pool)
- **Runner** — CLI that reads NDJSON, submits to engine, writes NDJSON output

## Data Flow

```
NDJSON input → Parser → Engine.submit(cmd) → Events → NDJSON output
```

## Design Principles

- Single-threaded engines (no concurrency inside)
- Integer arithmetic only
- Deterministic behavior for reproducible benchmarks
