package enginev3_pool

import (
	"matching-engine-lab/go/internal/core"
)

// Engine implements core.Engine using object pool (placeholder).
type Engine struct {
	book map[[2]int64]struct{} // placeholder: (price, timestamp) -> {}
}

// New returns a new Engine.
func New() *Engine {
	return &Engine{
		book: make(map[[2]int64]struct{}),
	}
}

// Submit processes a command and returns events.
func (e *Engine) Submit(cmd core.Command) ([]core.Event, error) {
	_ = cmd
	return []core.Event{{Accepted: &struct{}{}}}, nil
}

// Reset clears engine state.
func (e *Engine) Reset() error {
	e.book = make(map[[2]int64]struct{})
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
	return map[string]string{"engine": "v3_pool"}, nil
}
