"""Saving orders (rules §6): the order type, the clock, the repository interface and
the simplest repository there is - a list in memory.
"""

from __future__ import annotations

from collections.abc import Callable
from dataclasses import dataclass
from datetime import datetime, timezone
from typing import Any, Protocol

from checkout.models import Basket, Catalog, Receipt

# Anything that answers "what time is it?" is a Clock: a function with no arguments that
# returns a datetime (with a time zone). Tests pass one that always says the same time.
Clock = Callable[[], datetime]


def system_clock() -> datetime:
    """The real clock: the current time, in UTC."""
    return datetime.now(timezone.utc)


def format_timestamp(moment: datetime) -> str:
    """Formats a datetime the way ``created_at`` must look: whole seconds, UTC, trailing Z.

    Not implemented yet - see tests/unit/test_orders.py.
    """
    raise NotImplementedError(f"format_timestamp({moment}) is not implemented yet")


ORDER_ID_PREFIX = "ORDER-"


def order_id(n: int) -> str:
    """The id of the n-th order saved into a repository: "ORDER-1", "ORDER-2", ...

    Not implemented yet - see tests/unit/test_orders.py.
    """
    raise NotImplementedError(f"order_id({n}) is not implemented yet")


def order_number(identifier: str) -> int | None:
    """The inverse of order_id: "ORDER-7" gives 7. Anything else gives ``None``.

    Not implemented yet - see tests/unit/test_orders.py.
    """
    raise NotImplementedError(f"order_number({identifier!r}) is not implemented yet")


@dataclass(frozen=True)
class Order:
    """A saved order, exactly as spec/schema/order.schema.json describes it."""

    id: str
    created_at: str  # UTC, ISO 8601, whole seconds, trailing Z: "2026-01-15T10:30:00Z"
    receipt: Receipt

    def to_dict(self) -> dict[str, Any]:
        """The JSON form of the order (spec/schema/order.schema.json)."""
        raise NotImplementedError("Order.to_dict() is not implemented yet")

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> Order:
        """Reads an order back from its JSON form."""
        raise NotImplementedError("Order.from_dict() is not implemented yet")


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
        self._clock = clock

    def save(self, receipt: Receipt) -> str:
        raise NotImplementedError("InMemoryOrderRepository.save() is not implemented yet")

    def find(self, identifier: str) -> Order | None:
        raise NotImplementedError("InMemoryOrderRepository.find() is not implemented yet")

    def count(self) -> int:
        """How many orders are saved. Handy for tests ("the order repository is empty")."""
        raise NotImplementedError("InMemoryOrderRepository.count() is not implemented yet")


@dataclass(frozen=True)
class PlacedOrder:
    """What place_order gives back: the receipt and the id the repository gave it."""

    order_id: str
    receipt: Receipt


def place_order(basket: Basket, catalog: Catalog, repository: OrderRepository) -> PlacedOrder:
    """Checks out, then saves the receipt as an order. A failed checkout saves nothing.

    Not implemented yet - see tests/unit/test_orders.py.
    """
    raise NotImplementedError("place_order() is not implemented yet")
