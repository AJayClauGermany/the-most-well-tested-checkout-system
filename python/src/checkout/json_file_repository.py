"""An order repository that writes one JSON file per order into a folder."""

from __future__ import annotations

import json
from pathlib import Path

from checkout.models import Receipt
from checkout.orders import Clock, Order, format_timestamp, order_id, order_number, system_clock


class JsonFileOrderRepository:
    """Keeps each order as ``<folder>/ORDER-<n>.json``, in the shape of
    spec/schema/order.schema.json.

    Ids continue where the files in the folder left off, so reopening the same folder
    later keeps numbering correctly.
    """

    def __init__(self, folder: Path, clock: Clock | None = None) -> None:
        self._folder = folder
        self._clock: Clock = clock if clock is not None else system_clock
        folder.mkdir(parents=True, exist_ok=True)

    def save(self, receipt: Receipt) -> str:
        order = Order(
            id=order_id(self._highest_order_number() + 1),
            created_at=format_timestamp(self._clock()),
            receipt=receipt,
        )
        # indent=2 makes the file readable for humans; the newline at the end is good manners.
        text = json.dumps(order.to_dict(), indent=2) + "\n"
        self._path_of(order.id).write_text(text, encoding="utf-8")
        return order.id

    def find(self, identifier: str) -> Order | None:
        # Only accept real order ids: "../secret" must not read files outside the folder.
        if order_number(identifier) is None:
            return None
        path = self._path_of(identifier)
        if not path.is_file():
            return None
        return Order.from_dict(json.loads(path.read_text(encoding="utf-8")))

    def _path_of(self, identifier: str) -> Path:
        return self._folder / f"{identifier}.json"

    def _highest_order_number(self) -> int:
        """The highest n among the ORDER-<n>.json files in the folder, or 0 when there are none."""
        numbers = [
            order_number(path.stem)
            for path in self._folder.iterdir()
            if path.suffix == ".json" and path.is_file()
        ]
        return max((number for number in numbers if number is not None), default=0)
