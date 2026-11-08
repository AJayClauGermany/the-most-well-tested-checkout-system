"""Test helper: reads the shared, language-agnostic spec files.

The files live in the repository's top-level ``spec/`` folder and are shared by
every language, so we read them from there instead of copying them.
"""

from __future__ import annotations

import json
from dataclasses import dataclass
from pathlib import Path
from typing import Any

from checkout import Basket, Catalog, ErrorCode


def find_spec_folder() -> Path:
    """Finds the repository's ``spec/`` folder.

    We start from *this file* (not from the current working directory), so it
    works no matter which folder you run ``pytest`` from.
    """
    here = Path(__file__).resolve()
    for folder in here.parents:
        if (folder / "spec" / "checkout-rules.md").is_file():
            return folder / "spec"
    raise FileNotFoundError(f"Could not find the spec/ folder above {here}")


SPEC_FOLDER = find_spec_folder()


@dataclass(frozen=True)
class ExpectedTotals:
    """The six totals a spec case expects on the receipt (all in cents)."""

    subtotal: int
    promotion_discount: int
    coupon_discount: int
    shipping: int
    tax: int
    total: int


@dataclass(frozen=True)
class SpecCase:
    """One case from ``spec/test-cases/*.json``.

    Exactly one of ``expected`` and ``expected_error`` is set.
    """

    file: str  # e.g. "01-subtotal.json"
    id: str  # e.g. "subtotal-single-item"
    description: str
    basket: Basket
    # The basket exactly as written in the JSON file - what the command line program reads.
    basket_json: dict[str, Any]
    expected: ExpectedTotals | None
    expected_error: ErrorCode | None


def load_catalog() -> Catalog:
    """Loads the standard catalog, ``spec/catalog.json``."""
    return Catalog.from_dict(_read_json(SPEC_FOLDER / "catalog.json"))


def load_spec_cases() -> list[SpecCase]:
    """Loads every case from every file in ``spec/test-cases/``, in file order."""
    cases = []
    for path in sorted((SPEC_FOLDER / "test-cases").glob("*.json")):
        for case in _read_json(path)["cases"]:
            cases.append(_spec_case_from_dict(path.name, case))
    return cases


def valid_spec_cases() -> list[SpecCase]:
    """Only the cases that expect a receipt."""
    return [case for case in load_spec_cases() if case.expected is not None]


def invalid_spec_cases() -> list[SpecCase]:
    """Only the cases that expect an error code."""
    return [case for case in load_spec_cases() if case.expected_error is not None]


def _spec_case_from_dict(file: str, data: dict[str, Any]) -> SpecCase:
    expected = data.get("expected")
    expected_error = data.get("expected_error")
    return SpecCase(
        file=file,
        id=data["id"],
        description=data["description"],
        basket=Basket.from_dict(data["basket"]),
        basket_json=data["basket"],
        expected=ExpectedTotals(**expected) if expected is not None else None,
        expected_error=ErrorCode(expected_error) if expected_error is not None else None,
    )


def _read_json(path: Path) -> Any:
    return json.loads(path.read_text(encoding="utf-8"))
