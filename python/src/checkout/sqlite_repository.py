"""An order repository backed by SQLite, using the standard library's sqlite3 module."""

from __future__ import annotations

from pathlib import Path
from types import TracebackType

from checkout.models import Receipt
from checkout.orders import Clock, Order


class SqliteOrderRepository:
    """Keeps orders in one table of an SQLite database file (or ``":memory:"``).

    Not implemented yet - see tests/integration/test_sqlite_repository.py.
    """

    def __init__(self, path: str | Path, clock: Clock | None = None) -> None:
        raise NotImplementedError("SqliteOrderRepository is not implemented yet")

    def save(self, receipt: Receipt) -> str:
        raise NotImplementedError("SqliteOrderRepository.save() is not implemented yet")

    def find(self, identifier: str) -> Order | None:
        raise NotImplementedError("SqliteOrderRepository.find() is not implemented yet")

    def close(self) -> None:
        """Closes the database file. Call it when you are done (or use a ``with`` block)."""
        raise NotImplementedError("SqliteOrderRepository.close() is not implemented yet")

    def __enter__(self) -> SqliteOrderRepository:
        raise NotImplementedError("SqliteOrderRepository.__enter__() is not implemented yet")

    def __exit__(
        self,
        exc_type: type[BaseException] | None,
        exc_value: BaseException | None,
        traceback: TracebackType | None,
    ) -> None:
        raise NotImplementedError("SqliteOrderRepository.__exit__() is not implemented yet")
