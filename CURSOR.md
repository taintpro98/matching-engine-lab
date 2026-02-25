You are Cursor.

Scaffold a new repository named matching-engine-lab.

This is a research lab project to compare multiple matching engine implementations in Rust and Go.

This is NOT a crypto exchange.

⸻

DOMAIN MODEL

There is:
	•	One Asset
	•	One Money
	•	Sellers offer Asset
	•	Buyers spend Money
	•	price = Money per Asset

There is:
	•	No bid book
	•	No limit orders
	•	No trading pairs
	•	No symbols

Only these operations are supported:
	•	PlaceSell
	•	CancelSell
	•	BuyByQty
	•	BuyByBudget

Matching rules:
	•	Cheapest-first
	•	FIFO within same price
	•	Deterministic behavior
	•	Integer arithmetic only (no floats)

All numeric types must be integer ticks:

Price     = int64
AssetQty  = int64
Money     = int64
Timestamp = int64
ID        = u64


⸻

REQUIRED DIRECTORY STRUCTURE

Create exactly this structure:

matching-engine-lab/
├── README.md
├── .gitignore
├── docs/
│   ├── interface.md
│   ├── architecture/
│   │   ├── system_overview.md
│   │   ├── matching_flow.md
│   │   ├── concurrency_model.md
│   │   └── scaling_strategy.md
│   ├── data_structures/
│   │   ├── heap.md
│   │   ├── red_black_tree.md
│   │   ├── btree.md
│   │   ├── skiplist.md
│   │   ├── arena_allocator.md
│   │   └── comparison_table.md
│   └── benchmarks/
│       ├── rust_vs_go.md
│       ├── memory_profile.md
│       └── latency_profile.md
├── spec/
│   ├── command.schema.json
│   ├── event.schema.json
│   └── ndjson_format.md
├── datasets/
│   ├── sample/
│   │   ├── tiny.ndjson
│   │   └── medium.ndjson
│   ├── generated/
│   └── README.md
├── tools/
│   ├── generator/
│   │   └── py/
│   │       ├── gen_stream.py
│   │       └── profiles.py
│   ├── verifier/
│   │   └── py/
│   │       ├── verify_invariants.py
│   │       └── check_fifo.py
│   └── bench/
│       ├── run_rust.sh
│       ├── run_go.sh
│       └── compare.sh
├── rust/
│   ├── Cargo.toml
│   ├── crates/
│   │   ├── engine_core/
│   │   │   ├── Cargo.toml
│   │   │   └── src/
│   │   │       ├── lib.rs
│   │   │       ├── types.rs
│   │   │       ├── engine_trait.rs
│   │   │       └── ndjson.rs
│   │   ├── engine_v1_btreemap/
│   │   │   ├── Cargo.toml
│   │   │   └── src/lib.rs
│   │   ├── engine_v2_skiplist/
│   │   │   ├── Cargo.toml
│   │   │   └── src/lib.rs
│   │   ├── engine_v3_arena/
│   │   │   ├── Cargo.toml
│   │   │   └── src/lib.rs
│   │   └── runner/
│   │       ├── Cargo.toml
│   │       └── src/main.rs
│   └── benches/
├── go/
│   ├── go.mod
│   ├── internal/
│   │   ├── core/
│   │   │   ├── types.go
│   │   │   ├── engine.go
│   │   │   └── ndjson.go
│   │   ├── enginev1_btree/
│   │   │   ├── engine.go
│   │   │   └── book.go
│   │   ├── enginev2_treemap/
│   │   │   ├── engine.go
│   │   │   └── book.go
│   │   └── enginev3_pool/
│   │       ├── engine.go
│   │       └── book.go
│   └── cmd/
│       └── runner/
│           └── main.go
├── outputs/
│   ├── rust/
│   └── go/


⸻

SHARED INTERFACE FILE

Create docs/interface.md containing:
	•	Domain explanation
	•	Numeric types
	•	Command definitions:
	•	PlaceSell
	•	CancelSell
	•	BuyByQty
	•	BuyByBudget
	•	Event definitions:
	•	Accepted
	•	Rejected
	•	Trade
	•	SellUpdated
	•	SellClosed
	•	BuyResultQty
	•	BuyResultBudget
	•	Matching invariants:
	•	Cheapest-first
	•	FIFO
	•	Deterministic replay
	•	No negative values
	•	Budget never exceeded
	•	Integer arithmetic only

Keep it concise and clean.

⸻

RUST REQUIREMENTS
	•	Use Cargo workspace.
	•	engine_core defines:
	•	Command enum
	•	Event enum
	•	Type aliases
	•	Engine trait:

submit(cmd) -> Vec<Event>
reset()
snapshot() -> Vec<u8>
load_snapshot(&[u8])
stats() -> HashMap<String,String>

	•	Each engine_vX crate implements the trait with placeholder logic.
	•	runner:
	•	flag --engine v1|v2|v3
	•	reads NDJSON line-by-line
	•	calls engine.submit()
	•	writes NDJSON output
	•	prints ops/sec summary

Must compile with:

cargo build


⸻

GO REQUIREMENTS
	•	Module name: matching-engine-lab/go
	•	internal/core defines:
	•	Command struct
	•	Event struct
	•	Engine interface:

Submit(cmd) ([]Event, error)
Reset() error
Snapshot() ([]byte, error)
LoadSnapshot([]byte) error
Stats() (map[string]string, error)

	•	Each enginevX implements the interface.
	•	cmd/runner:
	•	flag --engine
	•	reads NDJSON
	•	writes NDJSON
	•	prints ops/sec summary

Must compile with:

go build ./...


⸻

NDJSON RULES
	•	Each line = one JSON command
	•	No JSON array wrapper
	•	Stream-friendly

⸻

CONSTRAINTS
	•	Do NOT implement full matching logic.
	•	Provide compile-ready skeletons only.
	•	No concurrency required inside engine.
	•	Use integer arithmetic only.
	•	Ensure both Rust and Go projects compile.

⸻

When finished, return only confirmation that project is created.