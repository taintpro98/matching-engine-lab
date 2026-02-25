# Heap

## Use Case

Priority queue for cheapest-first matching. Min-heap keyed by (price, timestamp).

## Complexity

- Insert: O(log n)
- Extract min: O(log n)
- Peek: O(1)

## Trade-offs

- Good for extract-min heavy workloads
- Not ideal for CancelSell (requires id lookup, O(n) or auxiliary index)
