# B-Tree

## Use Case

Ordered map for price levels. Used in engine_v1 (Rust BTreeMap, Go BTree).

## Complexity

- Insert: O(log n)
- Delete: O(log n)
- Range scan: O(k + log n)
- Cache-friendly due to node size

## Trade-offs

- Good cache locality
- Standard library support in Rust and Go
- Baseline implementation for comparison
