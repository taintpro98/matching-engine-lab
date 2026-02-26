package enginev2_treemap

import (
	"github.com/emirpasic/gods/utils"
)

// orderKey is the key for the TreeMap (price, timestamp, id).
type orderKey struct {
	Price     int64
	Timestamp int64
	ID        uint64
}

// OrderKeyComparator compares orderKey for TreeMap ordering (cheapest-first, FIFO).
var OrderKeyComparator utils.Comparator = func(a, b interface{}) int {
	ak := a.(orderKey)
	bk := b.(orderKey)
	if ak.Price != bk.Price {
		if ak.Price < bk.Price {
			return -1
		}
		return 1
	}
	if ak.Timestamp != bk.Timestamp {
		if ak.Timestamp < bk.Timestamp {
			return -1
		}
		return 1
	}
	if ak.ID != bk.ID {
		if ak.ID < bk.ID {
			return -1
		}
		return 1
	}
	return 0
}

// orderValue holds the sell order data.
type orderValue struct {
	Price     int64
	Timestamp int64
	ID        uint64
	Qty       int64
}
