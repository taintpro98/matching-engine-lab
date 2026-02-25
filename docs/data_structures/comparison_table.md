# Data Structure Comparison

| Structure    | Insert | Delete | Scan   | Cancel by ID | Cache |
|-------------|--------|--------|--------|--------------|-------|
| Heap        | O(log n) | O(n)* | N/A  | O(n)         | Medium |
| Red-Black   | O(log n) | O(log n) | O(k+log n) | O(log n)** | Medium |
| B-Tree      | O(log n) | O(log n) | O(k+log n) | O(log n)** | Good |
| Skip List   | O(log n) | O(log n) | O(k+log n) | O(log n)** | Medium |
| Arena+BTree | O(log n) | O(log n) | O(k+log n) | O(log n)** | Best |

\* Heap delete by value requires search  
\** With id→(price,ts) index
