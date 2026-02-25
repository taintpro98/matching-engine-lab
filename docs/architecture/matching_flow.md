# Matching Flow

## Order Book Structure

- Sell book only (no bid book)
- Keyed by price (ascending), then timestamp (ascending) for FIFO

## Matching Rules

1. **Cheapest-first** — iterate sells from lowest price
2. **FIFO** — within same price level, oldest first
3. **Deterministic** — same input sequence yields same output

## Operation Flows

### PlaceSell
- Validate id, price, qty, timestamp
- Insert into sell book at (price, timestamp)
- Emit Accepted or Rejected

### CancelSell
- Lookup sell by id
- Remove from book
- Emit SellClosed or Rejected

### BuyByQty
- Match cheapest sells until qty filled
- Emit Trade events
- Emit BuyResultQty with filled qty

### BuyByBudget
- Match cheapest sells until budget exhausted
- Emit Trade events
- Emit BuyResultBudget with spent amount and filled qty
