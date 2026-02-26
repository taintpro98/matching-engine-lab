//! Engine v1: BTreeMap-based implementation.
//! Order book: BTreeMap<(price, timestamp, id), SellOrder> for cheapest-first, FIFO ordering.
//! id_index: HashMap<id, (price, timestamp)> for CancelSell lookup.

use std::collections::{BTreeMap, HashMap};

use engine_core::{Command, Engine, Event};
use engine_core::{AssetQty, ID, Money, Price, Timestamp};

type BookKey = (Price, Timestamp, ID);

#[derive(Clone)]
#[allow(dead_code)]
struct SellOrder {
    id: ID,
    price: Price,
    qty: AssetQty,
    timestamp: Timestamp,
}

pub struct EngineV1 {
    book: BTreeMap<BookKey, SellOrder>,
    id_index: HashMap<ID, BookKey>,
}

impl Default for EngineV1 {
    fn default() -> Self {
        Self {
            book: BTreeMap::new(),
            id_index: HashMap::new(),
        }
    }
}

impl Engine for EngineV1 {
    fn submit(&mut self, cmd: Command) -> Vec<Event> {
        match cmd {
            Command::PlaceSell {
                id,
                price,
                qty,
                timestamp,
            } => self.place_sell(id, price, qty, timestamp),
            Command::CancelSell { id } => self.cancel_sell(id),
            Command::BuyByQty { id, qty, timestamp } => self.buy_by_qty(id, qty, timestamp),
            Command::BuyByBudget {
                id,
                budget,
                timestamp,
            } => self.buy_by_budget(id, budget, timestamp),
        }
    }

    fn reset(&mut self) {
        self.book.clear();
        self.id_index.clear();
    }

    fn snapshot(&self) -> Vec<u8> {
        vec![]
    }

    fn load_snapshot(&mut self, _data: &[u8]) {}

    fn stats(&self) -> HashMap<String, String> {
        let mut m = HashMap::new();
        m.insert("engine".to_string(), "v1_btreemap".to_string());
        m.insert("book_size".to_string(), self.book.len().to_string());
        m
    }
}

impl EngineV1 {
    fn place_sell(
        &mut self,
        id: ID,
        price: Price,
        qty: AssetQty,
        timestamp: Timestamp,
    ) -> Vec<Event> {
        if price < 0 || qty <= 0 {
            return vec![Event::Rejected {
                reason: "price must be >= 0, qty must be > 0".to_string(),
            }];
        }
        if self.id_index.contains_key(&id) {
            return vec![Event::Rejected {
                reason: format!("duplicate id {}", id),
            }];
        }
        let key = (price, timestamp, id);
        let order = SellOrder {
            id,
            price,
            qty,
            timestamp,
        };
        self.book.insert(key, order);
        self.id_index.insert(id, key);
        vec![Event::Accepted]
    }

    fn cancel_sell(&mut self, id: ID) -> Vec<Event> {
        let key = match self.id_index.remove(&id) {
            Some(k) => k,
            None => {
                return vec![Event::Rejected {
                    reason: format!("sell {} not found", id),
                }]
            }
        };
        self.book.remove(&key);
        vec![Event::SellClosed]
    }

    fn buy_by_qty(&mut self, buyer_id: ID, qty: AssetQty, _timestamp: Timestamp) -> Vec<Event> {
        if qty <= 0 {
            return vec![Event::Rejected {
                reason: "qty must be > 0".to_string(),
            }];
        }
        self.match_orders(buyer_id, qty, None)
    }

    fn buy_by_budget(
        &mut self,
        buyer_id: ID,
        budget: Money,
        _timestamp: Timestamp,
    ) -> Vec<Event> {
        if budget <= 0 {
            return vec![Event::Rejected {
                reason: "budget must be > 0".to_string(),
            }];
        }
        self.match_orders(buyer_id, 0, Some(budget))
    }

    fn match_orders(
        &mut self,
        buyer_id: ID,
        max_qty: AssetQty,
        max_budget: Option<Money>,
    ) -> Vec<Event> {
        let mut events = Vec::new();
        let mut remaining_qty = max_qty;
        let mut remaining_budget = max_budget.unwrap_or(i64::MAX);
        let mut total_filled: AssetQty = 0;
        let mut total_spent: Money = 0;

        let can_continue = |rq: AssetQty, rb: Money, mb: Option<Money>| {
            (mb.is_none() && rq > 0) || (mb.is_some() && rb > 0)
        };

        while can_continue(remaining_qty, remaining_budget, max_budget) && !self.book.is_empty() {
            let (key, mut order) = match self.book.pop_first() {
                Some((k, o)) => (k, o),
                None => break,
            };

            let fill_qty = if max_budget.is_some() {
                let max_by_budget = remaining_budget / order.price;
                max_by_budget.min(order.qty)
            } else {
                remaining_qty.min(order.qty)
            };

            if fill_qty <= 0 {
                self.book.insert(key, order);
                break;
            }

            let cost = fill_qty * order.price;
            events.push(Event::Trade {
                buyer_id,
                seller_id: order.id,
                qty: fill_qty,
                price: order.price,
            });

            order.qty -= fill_qty;
            total_filled += fill_qty;
            total_spent += cost;
            if max_budget.is_some() {
                remaining_budget -= cost;
            } else {
                remaining_qty -= fill_qty;
            }

            if order.qty == 0 {
                self.id_index.remove(&order.id);
                events.push(Event::SellClosed);
            } else {
                self.book.insert(key, order);
                events.push(Event::SellUpdated);
            }
        }

        if max_budget.is_some() {
            events.push(Event::BuyResultBudget {
                spent: total_spent,
                filled: total_filled,
            });
        } else {
            events.push(Event::BuyResultQty {
                filled: total_filled,
            });
        }
        events
    }
}
