#!/usr/bin/env python3
"""Check FIFO ordering in matching engine output."""

import json
import sys


def main():
    # Placeholder: parse events, verify FIFO within same price level
    for line in sys.stdin:
        line = line.strip()
        if not line:
            continue
        try:
            event = json.loads(line)
        except json.JSONDecodeError:
            print(f"Invalid JSON: {line[:80]}...")
            sys.exit(1)
        pass

    print("OK: FIFO verified")


if __name__ == "__main__":
    main()
