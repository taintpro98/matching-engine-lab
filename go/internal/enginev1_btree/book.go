package enginev1_btree

import (
	"github.com/google/btree"
)

// orderItem implements btree.Item for ordered sell book.
// Key: (Price, Timestamp, ID) for cheapest-first, FIFO.
type orderItem struct {
	Price     int64
	Timestamp int64
	ID        uint64
	Qty       int64
}

func (a orderItem) Less(b btree.Item) bool {
	bb := b.(orderItem)
	if a.Price != bb.Price {
		return a.Price < bb.Price
	}
	if a.Timestamp != bb.Timestamp {
		return a.Timestamp < bb.Timestamp
	}
	return a.ID < bb.ID
}
