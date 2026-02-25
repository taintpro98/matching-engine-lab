package core

// Engine is the interface for matching engine implementations.
type Engine interface {
	Submit(cmd Command) ([]Event, error)
	Reset() error
	Snapshot() ([]byte, error)
	LoadSnapshot(data []byte) error
	Stats() (map[string]string, error)
}
