# Red-Black Tree

## Use Case

Ordered map for price levels. Self-balancing BST.

## Complexity

- Insert: O(log n)
- Delete: O(log n)
- Range scan: O(k + log n) for k elements

## Trade-offs

- Balanced height
- Good for ordered iteration (cheapest-first)
- Cancel by id requires secondary index
