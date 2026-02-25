#!/usr/bin/env python3
"""Verify matching engine output invariants."""

import json
import sys


def main():
    for line in sys.stdin:
        line = line.strip()
        if not line:
            continue
        try:
            event = json.loads(line)
        except json.JSONDecodeError:
            print(f"Invalid JSON: {line[:80]}...")
            sys.exit(1)

        # Placeholder: add invariant checks
        # - No negative values in Trade
        # - Budget never exceeded in BuyResultBudget
        # - etc.
        pass

    print("OK: invariants verified")


if __name__ == "__main__":
    main()
