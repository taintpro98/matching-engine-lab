#!/usr/bin/env python3
"""Generate NDJSON command stream for matching engine benchmarks."""

import json
import random
import argparse
from profiles import get_profile


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--profile", default="default", help="Profile name from profiles.py")
    parser.add_argument("--count", type=int, default=1000, help="Number of commands")
    parser.add_argument("--seed", type=int, default=42, help="Random seed")
    args = parser.parse_args()

    random.seed(args.seed)
    profile = get_profile(args.profile)

    sell_id = 0
    buy_id = 10000
    timestamp = 0

    for _ in range(args.count):
        cmd = profile.choose_command()
        timestamp += random.randint(1, 10)

        if cmd == "PlaceSell":
            sell_id += 1
            price = random.randint(profile.min_price, profile.max_price)
            qty = random.randint(profile.min_qty, profile.max_qty)
            obj = {"PlaceSell": {"id": sell_id, "price": price, "qty": qty, "timestamp": timestamp}}
        elif cmd == "CancelSell" and sell_id > 0:
            sid = random.randint(1, sell_id)
            obj = {"CancelSell": {"id": sid}}
        elif cmd == "BuyByQty":
            buy_id += 1
            qty = random.randint(profile.min_qty, profile.max_qty)
            obj = {"BuyByQty": {"id": buy_id, "qty": qty, "timestamp": timestamp}}
        elif cmd == "BuyByBudget":
            buy_id += 1
            budget = random.randint(profile.min_budget, profile.max_budget)
            obj = {"BuyByBudget": {"id": buy_id, "budget": budget, "timestamp": timestamp}}
        else:
            continue

        print(json.dumps(obj))


if __name__ == "__main__":
    main()
