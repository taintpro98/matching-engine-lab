# Skip List

## Use Case

Ordered structure with probabilistic balance. Used in engine_v2 (Rust).

## Complexity

- Insert: O(log n) expected
- Delete: O(log n) expected
- Range scan: O(k + log n)

## Trade-offs

- Simpler than balanced trees
- No rebalancing
- Probabilistic, but deterministic with fixed seed if needed
