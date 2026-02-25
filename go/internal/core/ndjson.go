package core

import (
	"encoding/json"
)

// Command represents a matching engine command.
type Command struct {
	PlaceSell    *PlaceSellCmd    `json:"PlaceSell,omitempty"`
	CancelSell   *CancelSellCmd   `json:"CancelSell,omitempty"`
	BuyByQty     *BuyByQtyCmd     `json:"BuyByQty,omitempty"`
	BuyByBudget  *BuyByBudgetCmd  `json:"BuyByBudget,omitempty"`
}

type PlaceSellCmd struct {
	ID        ID        `json:"id"`
	Price     Price     `json:"price"`
	Qty       AssetQty  `json:"qty"`
	Timestamp Timestamp `json:"timestamp"`
}

type CancelSellCmd struct {
	ID ID `json:"id"`
}

type BuyByQtyCmd struct {
	ID        ID        `json:"id"`
	Qty       AssetQty  `json:"qty"`
	Timestamp Timestamp `json:"timestamp"`
}

type BuyByBudgetCmd struct {
	ID        ID        `json:"id"`
	Budget    Money     `json:"budget"`
	Timestamp Timestamp `json:"timestamp"`
}

// Event represents a matching engine event.
type Event struct {
	Accepted       *struct{}       `json:"Accepted,omitempty"`
	Rejected       *RejectedEvent  `json:"Rejected,omitempty"`
	Trade          *TradeEvent     `json:"Trade,omitempty"`
	SellUpdated    *struct{}       `json:"SellUpdated,omitempty"`
	SellClosed     *struct{}       `json:"SellClosed,omitempty"`
	BuyResultQty   *BuyResultQty   `json:"BuyResultQty,omitempty"`
	BuyResultBudget *BuyResultBudget `json:"BuyResultBudget,omitempty"`
}

type RejectedEvent struct {
	Reason string `json:"reason"`
}

type TradeEvent struct {
	BuyerID  ID        `json:"buyer_id"`
	SellerID ID        `json:"seller_id"`
	Qty      AssetQty  `json:"qty"`
	Price    Price     `json:"price"`
}

type BuyResultQty struct {
	Filled AssetQty `json:"filled"`
}

type BuyResultBudget struct {
	Spent  Money    `json:"spent"`
	Filled AssetQty `json:"filled"`
}

// ParseCommand parses a JSON line into a Command.
func ParseCommand(s string) (*Command, error) {
	var cmd Command
	err := json.Unmarshal([]byte(s), &cmd)
	if err != nil {
		return nil, err
	}
	return &cmd, nil
}

// SerializeEvent serializes an Event to JSON.
func SerializeEvent(e *Event) (string, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
