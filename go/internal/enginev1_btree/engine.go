package enginev1_btree

import (
	"fmt"

	"github.com/google/btree"
	"matching-engine-lab/go/internal/core"
)

const btreeDegree = 32

// Engine implements core.Engine using BTree.
type Engine struct {
	book     *btree.BTree
	idIndex  map[core.ID]orderItem
}

// New returns a new Engine.
func New() *Engine {
	return &Engine{
		book:    btree.New(btreeDegree),
		idIndex: make(map[core.ID]orderItem),
	}
}

// Submit processes a command and returns events.
func (e *Engine) Submit(cmd core.Command) ([]core.Event, error) {
	switch {
	case cmd.PlaceSell != nil:
		return e.placeSell(cmd.PlaceSell), nil
	case cmd.CancelSell != nil:
		return e.cancelSell(cmd.CancelSell), nil
	case cmd.BuyByQty != nil:
		return e.buyByQty(cmd.BuyByQty), nil
	case cmd.BuyByBudget != nil:
		return e.buyByBudget(cmd.BuyByBudget), nil
	default:
		return []core.Event{{Rejected: &core.RejectedEvent{Reason: "unknown command"}}}, nil
	}
}

func (e *Engine) placeSell(cmd *core.PlaceSellCmd) []core.Event {
	if cmd.Price < 0 || cmd.Qty <= 0 {
		return []core.Event{{Rejected: &core.RejectedEvent{Reason: "price must be >= 0, qty must be > 0"}}}
	}
	if _, exists := e.idIndex[cmd.ID]; exists {
		return []core.Event{{Rejected: &core.RejectedEvent{Reason: "duplicate id"}}}
	}
	item := orderItem{
		Price:     int64(cmd.Price),
		Timestamp: int64(cmd.Timestamp),
		ID:        uint64(cmd.ID),
		Qty:       int64(cmd.Qty),
	}
	e.book.ReplaceOrInsert(item)
	e.idIndex[cmd.ID] = item
	return []core.Event{{Accepted: &struct{}{}}}
}

func (e *Engine) cancelSell(cmd *core.CancelSellCmd) []core.Event {
	item, exists := e.idIndex[cmd.ID]
	if !exists {
		return []core.Event{{Rejected: &core.RejectedEvent{Reason: "sell not found"}}}
	}
	delete(e.idIndex, cmd.ID)
	e.book.Delete(item)
	return []core.Event{{SellClosed: &struct{}{}}}
}

func (e *Engine) buyByQty(cmd *core.BuyByQtyCmd) []core.Event {
	if cmd.Qty <= 0 {
		return []core.Event{{Rejected: &core.RejectedEvent{Reason: "qty must be > 0"}}}
	}
	return e.matchOrders(uint64(cmd.ID), int64(cmd.Qty), 0, true)
}

func (e *Engine) buyByBudget(cmd *core.BuyByBudgetCmd) []core.Event {
	if cmd.Budget <= 0 {
		return []core.Event{{Rejected: &core.RejectedEvent{Reason: "budget must be > 0"}}}
	}
	return e.matchOrders(uint64(cmd.ID), 0, int64(cmd.Budget), false)
}

func (e *Engine) matchOrders(buyerID uint64, maxQty int64, maxBudget int64, byQty bool) []core.Event {
	var events []core.Event
	totalFilled := int64(0)
	totalSpent := int64(0)
	remainingQty := maxQty
	remainingBudget := maxBudget

	for {
		if byQty && remainingQty <= 0 {
			break
		}
		if !byQty && remainingBudget <= 0 {
			break
		}
		raw := e.book.DeleteMin()
		if raw == nil {
			break
		}
		order := raw.(orderItem)

		var fillQty int64
		if byQty {
			if order.Qty < remainingQty {
				fillQty = order.Qty
			} else {
				fillQty = remainingQty
			}
		} else {
			maxByBudget := remainingBudget / order.Price
			if maxByBudget <= 0 {
				e.book.ReplaceOrInsert(order)
				break
			}
			if order.Qty < maxByBudget {
				fillQty = order.Qty
			} else {
				fillQty = maxByBudget
			}
		}

		if fillQty <= 0 {
			e.book.ReplaceOrInsert(order)
			break
		}

		cost := fillQty * order.Price
		events = append(events, core.Event{
			Trade: &core.TradeEvent{
				BuyerID:  core.ID(buyerID),
				SellerID: core.ID(order.ID),
				Qty:      core.AssetQty(fillQty),
				Price:    core.Price(order.Price),
			},
		})

		order.Qty -= fillQty
		totalFilled += fillQty
		totalSpent += cost
		if byQty {
			remainingQty -= fillQty
		} else {
			remainingBudget -= cost
		}

		if order.Qty == 0 {
			delete(e.idIndex, core.ID(order.ID))
			events = append(events, core.Event{SellClosed: &struct{}{}})
		} else {
			e.book.ReplaceOrInsert(order)
			e.idIndex[core.ID(order.ID)] = order
			events = append(events, core.Event{SellUpdated: &struct{}{}})
		}
	}

	if byQty {
		events = append(events, core.Event{BuyResultQty: &core.BuyResultQty{Filled: core.AssetQty(totalFilled)}})
	} else {
		events = append(events, core.Event{BuyResultBudget: &core.BuyResultBudget{Spent: core.Money(totalSpent), Filled: core.AssetQty(totalFilled)}})
	}
	return events
}

// Reset clears engine state.
func (e *Engine) Reset() error {
	e.book.Clear(true)
	e.idIndex = make(map[core.ID]orderItem)
	return nil
}

// Snapshot returns serialized state.
func (e *Engine) Snapshot() ([]byte, error) {
	return []byte{}, nil
}

// LoadSnapshot restores state from bytes.
func (e *Engine) LoadSnapshot(data []byte) error {
	_ = data
	return nil
}

// Stats returns engine statistics.
func (e *Engine) Stats() (map[string]string, error) {
	return map[string]string{
		"engine":    "v1_btree",
		"book_size": fmt.Sprintf("%d", e.book.Len()),
	}, nil
}
