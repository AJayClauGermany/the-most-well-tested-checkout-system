"""An order repository that writes one JSON file per order into a folder."""

from __future__ import annotations

from pathlib import Path

from checkout.models import Receipt
from checkout.orders import Clock, Order


class JsonFileOrderRepository:
    """Keeps each order as ``<folder>/ORDER-<n>.json`` (spec/schema/order.schema.json).

    Not implemented yet - see tests/integration/test_json_file_repository.py.
    """

    def __init__(self, folder: Path, clock: Clock | None = None) -> None:
        raise NotImplementedError("JsonFileOrderRepository is not implemented yet")

    def save(self, receipt: Receipt) -> str:
        raise NotImplementedError("JsonFileOrderRepository.save() is not implemented yet")

    def find(self, identifier: str) -> Order | None:
        raise NotImplementedError("JsonFileOrderRepository.find() is not implemented yet")
