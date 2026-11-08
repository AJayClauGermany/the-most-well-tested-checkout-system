"""Saving orders (rules §6): the order type, the clock, the repository interface and
the simplest repository there is - a list in memory.
"""

from __future__ import annotations

import copy
import re
from collections.abc import Callable
from dataclasses import dataclass
from datetime import datetime, timezone
from typing import Any, Protocol

from checkout.engine import checkout
from checkout.models import Basket, Catalog, Receipt

# Anything that answers "what time is it?" is a Clock: a function with no arguments that
# returns a datetime (with a time zone). Production uses `system_clock`; tests pass one that
# always says the same time - the rules (§6) demand that the time source can be replaced,
# because `created_at` is the one value that changes from run to run.
Clock = Callable[[], datetime]


def system_clock() -> datetime:
    """The real clock: the current time, in UTC."""
    return datetime.now(timezone.utc)


def format_timestamp(moment: datetime) -> str:
    """Formats a datetime the way ``created_at`` must look: whole seconds, UTC, trailing Z.

    >>> format_timestamp(datetime(2026, 1, 15, 10, 30, 0, 999_999, tzinfo=timezone.utc))
    '2026-01-15T10:30:00Z'
    """
    if moment.tzinfo is None:
        raise ValueError(f"{moment!r} has no time zone, so it cannot be converted to UTC")
    # astimezone() converts to UTC; strftime() leaves out the microseconds (no rounding).
    return moment.astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


ORDER_ID_PREFIX = "ORDER-"
# A positive whole number with no leading zero, exactly as order.schema.json spells it.
_ORDER_ID = re.compile(r"ORDER-([1-9][0-9]*)")


def order_id(n: int) -> str:
    """The id of the n-th order saved into a repository: "ORDER-1", "ORDER-2", ...

    >>> order_id(7)
    'ORDER-7'
    """
    return f"{ORDER_ID_PREFIX}{n}"


def order_number(identifier: str) -> int | None:
    """The inverse of order_id: "ORDER-7" gives 7. Anything else gives ``None``.

    >>> order_number("ORDER-7"), order_number("ORDER-07"), order_number("nonsense")
    (7, None, None)
    """
    match = _ORDER_ID.fullmatch(identifier)
    return int(match.group(1)) if match else None


@dataclass(frozen=True)
class Order:
    """A saved order, exactly as spec/schema/order.schema.json describes it."""

    id: str
    created_at: str  # UTC, ISO 8601, whole seconds, trailing Z: "2026-01-15T10:30:00Z"
    receipt: Receipt

    def to_dict(self) -> dict[str, Any]:
        """The JSON form of the order (spec/schema/order.schema.json)."""
        return {"id": self.id, "created_at": self.created_at, "receipt": self.receipt.to_dict()}

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> Order:
        """Reads an order back from its JSON form."""
        return cls(
            id=data["id"],
            created_at=data["created_at"],
            receipt=Receipt.from_dict(data["receipt"]),
        )


class OrderRepository(Protocol):
    """Where orders are kept (rules §6).

    A ``Protocol`` describes a shape: anything with these two methods *is* an
    OrderRepository, without having to inherit from it. Implemented in memory,
    as JSON files and in SQLite.
    """

    def save(self, receipt: Receipt) -> str:
        """Saves a receipt and returns the new order's id."""
        ...

    def find(self, identifier: str) -> Order | None:
        """Looks an order up. ``None`` means "not found" - a normal answer, not an error."""
        ...


class InMemoryOrderRepository:
    """Keeps orders in a plain list. Doubles as the hand-written fake for tests."""

    def __init__(self, clock: Clock | None = None) -> None:
        # Looked up at call time (not as a default argument), so a test can swap
        # `system_clock` with monkeypatch and this class picks the replacement up.
        self._clock: Clock = clock if clock is not None else system_clock
        self._orders: list[Order] = []

    def save(self, receipt: Receipt) -> str:
        identifier = order_id(len(self._orders) + 1)
        self._orders.append(
            Order(
                id=identifier,
                created_at=format_timestamp(self._clock()),
                # A copy, so whoever gave us the receipt cannot change the saved order later.
                receipt=copy.deepcopy(receipt),
            )
        )
        return identifier

    def find(self, identifier: str) -> Order | None:
        for order in self._orders:
            if order.id == identifier:
                return order
        return None

    def count(self) -> int:
        """How many orders are saved. Handy for tests ("the order repository is empty")."""
        return len(self._orders)


@dataclass(frozen=True)
class PlacedOrder:
    """What place_order gives back: the receipt and the id the repository gave it."""

    order_id: str
    receipt: Receipt


def place_order(basket: Basket, catalog: Catalog, repository: OrderRepository) -> PlacedOrder:
    """Checks out, then saves the receipt as an order.

    A failed checkout saves nothing: the CheckoutError is raised before the repository
    is touched (property P14).
    """
    receipt = checkout(basket, catalog)
    return PlacedOrder(order_id=repository.save(receipt), receipt=receipt)
