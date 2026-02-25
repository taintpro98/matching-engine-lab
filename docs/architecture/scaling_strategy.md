# Scaling Strategy

## Horizontal Scaling

- Run multiple engine instances in parallel
- Partition by workload (e.g., different input files)
- Aggregate benchmark results

## Vertical Scaling

- Optimize data structures (heap, BTree, SkipList, arena)
- Reduce allocations (object pools, arena allocators)
- Profile memory and latency

## Benchmark Scaling

- tiny.ndjson — smoke tests
- medium.ndjson — typical workload
- generated/ — stress tests with configurable profiles
