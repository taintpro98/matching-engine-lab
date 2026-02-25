# Arena Allocator

## Use Case

Reduce allocation overhead. Pre-allocate blocks, reuse for order entries. Used in engine_v3.

## Benefits

- Fewer malloc/free calls
- Better cache locality
- Fast bulk free (reset arena)

## Trade-offs

- Memory not reclaimed until arena reset
- Fixed block size or growth strategy
