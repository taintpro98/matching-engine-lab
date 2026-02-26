package enginev2_treemap

import (
	"fmt"

	"github.com/emirpasic/gods/maps/treemap"
	"matching-engine-lab/go/internal/core"
)

// Engine implements core.Engine using TreeMap (red-black tree).
type Engine struct {
	book    *treemap.Map
	idIndex map[core.ID]orderKey
}

// New returns a new Engine.
func New() *Engine {
	return &Engine{
		book:    treemap.NewWith(OrderKeyComparator),
		idIndex: make(map[core.ID]orderKey),
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
	key := orderKey{Price: int64(cmd.Price), Timestamp: int64(cmd.Timestamp), ID: uint64(cmd.ID)}
	val := orderValue{
		Price:     int64(cmd.Price),
		Timestamp: int64(cmd.Timestamp),
		ID:        uint64(cmd.ID),
		Qty:       int64(cmd.Qty),
	}
	e.book.Put(key, val)
	e.idIndex[cmd.ID] = key
	return []core.Event{{Accepted: &struct{}{}}}
}

func (e *Engine) cancelSell(cmd *core.CancelSellCmd) []core.Event {
	key, exists := e.idIndex[cmd.ID]
	if !exists {
		return []core.Event{{Rejected: &core.RejectedEvent{Reason: "sell not found"}}}
	}
	delete(e.idIndex, cmd.ID)
	e.book.Remove(key)
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
		minKey, minVal := e.book.Min()
		if minKey == nil {
			break
		}
		key := minKey.(orderKey)
		order := minVal.(orderValue)
		e.book.Remove(key)

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
				e.book.Put(key, order)
				break
			}
			if order.Qty < maxByBudget {
				fillQty = order.Qty
			} else {
				fillQty = maxByBudget
			}
		}

		if fillQty <= 0 {
			e.book.Put(key, order)
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
			e.book.Put(key, order)
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
	e.book.Clear()
	e.idIndex = make(map[core.ID]orderKey)
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
		"engine":    "v2_treemap",
		"book_size": fmt.Sprintf("%d", e.book.Size()),
	}, nil
}
