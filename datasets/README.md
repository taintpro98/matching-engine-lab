# Datasets

## sample/

- `tiny.ndjson` — minimal smoke test
- `medium.ndjson` — typical workload

## generated/

- Output of `tools/generator/py/gen_stream.py`
- Configurable via `profiles.py`

## Usage

```bash
# Run with sample
cargo run -p runner -- --engine v1 < datasets/sample/tiny.ndjson

# Generate and run
python tools/generator/py/gen_stream.py --profile default > datasets/generated/stream.ndjson
```
