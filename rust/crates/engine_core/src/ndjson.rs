//! Command and Event types for NDJSON I/O.

use serde::{Deserialize, Serialize};

use crate::types::{AssetQty, ID, Money, Price, Timestamp};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum Command {
    PlaceSell {
        id: ID,
        price: Price,
        qty: AssetQty,
        timestamp: Timestamp,
    },
    CancelSell { id: ID },
    BuyByQty {
        id: ID,
        qty: AssetQty,
        timestamp: Timestamp,
    },
    BuyByBudget {
        id: ID,
        budget: Money,
        timestamp: Timestamp,
    },
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum Event {
    Accepted,
    Rejected { reason: String },
    Trade {
        buyer_id: ID,
        seller_id: ID,
        qty: AssetQty,
        price: Price,
    },
    SellUpdated,
    SellClosed,
    BuyResultQty { filled: AssetQty },
    BuyResultBudget { spent: Money, filled: AssetQty },
}

/// Parse a JSON line into a Command.
pub fn parse_command(s: &str) -> Result<Command, serde_json::Error> {
    serde_json::from_str(s)
}

/// Serialize an Event to JSON.
pub fn serialize_event(event: &Event) -> String {
    serde_json::to_string(event).unwrap()
}
