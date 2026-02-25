# NDJSON Format

## Rules

- Each line = one JSON command
- No JSON array wrapper
- Stream-friendly (line-delimited)
- UTF-8 encoding

## Example

```ndjson
{"PlaceSell":{"id":1,"price":100,"qty":10,"timestamp":1000}}
{"PlaceSell":{"id":2,"price":99,"qty":5,"timestamp":1001}}
{"BuyByQty":{"id":3,"qty":7,"timestamp":1002}}
{"CancelSell":{"id":1}}
```

## Parsing

- Read line by line
- Parse each line as JSON
- Skip empty lines (optional)
- Invalid JSON → Rejected event
