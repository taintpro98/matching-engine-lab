# Matching Engine Interface

## Domain

- **One Asset** — what sellers offer
- **One Money** — what buyers spend
- **price** = Money per Asset (integer ticks)
- No bid book, no limit orders, no trading pairs, no symbols

## Numeric Types

| Type      | Rust   | Go     |
|-----------|--------|--------|
| Price     | i64    | int64  |
| AssetQty  | i64    | int64  |
| Money     | i64    | int64  |
| Timestamp | i64    | int64  |
| ID        | u64    | uint64 |

## Commands

| Command     | Fields                    | Description                    |
|-------------|---------------------------|--------------------------------|
| PlaceSell   | id, price, qty, timestamp | Add sell offer to book         |
| CancelSell  | id                        | Remove sell offer              |
| BuyByQty    | id, qty, timestamp        | Buy up to qty Asset            |
| BuyByBudget | id, budget, timestamp     | Buy Asset with up to budget    |

## Events

| Event          | Description                              |
|----------------|------------------------------------------|
| Accepted       | Command accepted                         |
| Rejected       | Command rejected (reason)                |
| Trade          | Execution (buyer_id, seller_id, qty, price) |
| SellUpdated    | Partial fill of sell                     |
| SellClosed     | Sell fully filled or cancelled           |
| BuyResultQty   | Result of BuyByQty                       |
| BuyResultBudget| Result of BuyByBudget                    |

## Matching Invariants

- **Cheapest-first** — lowest price sells match first
- **FIFO** — within same price, earliest timestamp first
- **Deterministic replay** — same input → same output
- **No negative values** — qty, price, budget ≥ 0
- **Budget never exceeded** — BuyByBudget spends ≤ budget
- **Integer arithmetic only** — no floats
