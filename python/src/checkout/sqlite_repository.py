"""An order repository backed by SQLite, using the standard library's sqlite3 module."""

from __future__ import annotations

import json
import sqlite3
from pathlib import Path
from types import TracebackType

from checkout.models import Receipt
from checkout.orders import Clock, Order, format_timestamp, order_id, order_number, system_clock


class SqliteOrderRepository:
    """Keeps orders in one table of an SQLite database file (or ``":memory:"``).

    The database row number becomes the order number, so ids are sequential and survive
    closing and reopening the file. The receipt is stored as JSON text in one column.

    Use it in a ``with`` block to make sure the file is closed afterwards::

        with SqliteOrderRepository(path) as repository:
            repository.save(receipt)
    """

    def __init__(self, path: str | Path, clock: Clock | None = None) -> None:
        self._connection = sqlite3.connect(str(path))
        self._clock: Clock = clock if clock is not None else system_clock
        self._connection.execute(
            """
            CREATE TABLE IF NOT EXISTS orders (
                number     INTEGER PRIMARY KEY AUTOINCREMENT,
                created_at TEXT NOT NULL,
                receipt    TEXT NOT NULL
            )
            """
        )
        self._connection.commit()

    def save(self, receipt: Receipt) -> str:
        # "?" placeholders keep the values apart from the SQL, so a receipt can never be
        # mistaken for a command (no SQL injection).
        # (`pragma: no mutate`: SQL ignores upper/lower case, so mutmut's case-changing
        # mutants of this statement would behave identically - "equivalent mutants".)
        cursor = self._connection.execute(
            "INSERT INTO orders (created_at, receipt) VALUES (?, ?)",  # pragma: no mutate
            (format_timestamp(self._clock()), json.dumps(receipt.to_dict())),
        )
        self._connection.commit()
        assert cursor.lastrowid is not None  # an INSERT always has a row id
        return order_id(cursor.lastrowid)

    def find(self, identifier: str) -> Order | None:
        number = order_number(identifier)
        if number is None:
            return None
        row = self._connection.execute(
            "SELECT number, created_at, receipt FROM orders WHERE number = ?",  # pragma: no mutate
            (number,),
        ).fetchone()
        if row is None:
            return None
        found_number, created_at, receipt_json = row
        return Order(
            id=order_id(found_number),
            created_at=created_at,
            receipt=Receipt.from_dict(json.loads(receipt_json)),
        )

    def close(self) -> None:
        """Closes the database file. Call it when you are done (or use a ``with`` block)."""
        self._connection.close()

    def __enter__(self) -> SqliteOrderRepository:
        return self

    def __exit__(
        self,
        exc_type: type[BaseException] | None,
        exc_value: BaseException | None,
        traceback: TracebackType | None,
    ) -> None:
        self.close()
