"""Command generation profiles for benchmark streams."""

import random
from dataclasses import dataclass
from typing import List


@dataclass
class Profile:
    name: str
    weights: dict  # PlaceSell, CancelSell, BuyByQty, BuyByBudget
    min_price: int
    max_price: int
    min_qty: int
    max_qty: int
    min_budget: int
    max_budget: int

    def choose_command(self) -> str:
        commands = list(self.weights.keys())
        weights = list(self.weights.values())
        return random.choices(commands, weights=weights)[0]


PROFILES = {
    "default": Profile(
        name="default",
        weights={"PlaceSell": 4, "CancelSell": 1, "BuyByQty": 2, "BuyByBudget": 2},
        min_price=90,
        max_price=110,
        min_qty=1,
        max_qty=100,
        min_budget=100,
        max_budget=10000,
    ),
    "sell_heavy": Profile(
        name="sell_heavy",
        weights={"PlaceSell": 8, "CancelSell": 1, "BuyByQty": 1, "BuyByBudget": 1},
        min_price=95,
        max_price=105,
        min_qty=10,
        max_qty=50,
        min_budget=500,
        max_budget=5000,
    ),
    "buy_heavy": Profile(
        name="buy_heavy",
        weights={"PlaceSell": 2, "CancelSell": 1, "BuyByQty": 4, "BuyByBudget": 4},
        min_price=98,
        max_price=102,
        min_qty=5,
        max_qty=200,
        min_budget=1000,
        max_budget=20000,
    ),
}


def get_profile(name: str) -> Profile:
    if name not in PROFILES:
        raise ValueError(f"Unknown profile: {name}. Available: {list(PROFILES.keys())}")
    return PROFILES[name]
