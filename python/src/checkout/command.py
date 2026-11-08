"""The command line program's logic (rules §8), kept free of the terminal so it can be unit
tested in-process. ``src/checkout/__main__.py`` is the thin wrapper that calls ``main()``.
"""

from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path
from typing import Literal

from checkout.models import Basket, Catalog
from checkout.spec_folder import find_spec_folder


def standard_catalog_path() -> Path:
    """Where the standard catalog lives: spec/catalog.json at the root of the repository."""
    return find_spec_folder() / "catalog.json"


def load_standard_catalog() -> Catalog:
    """Reads spec/catalog.json.

    Not implemented yet - see tests/unit/test_command.py.
    """
    raise NotImplementedError("load_standard_catalog() is not implemented yet")


class NotABasketError(ValueError):
    """The input was valid JSON, but not the shape of a basket (rules §3)."""


def parse_basket(text: str) -> Basket:
    """Parses the JSON text of a basket.

    Not implemented yet - see tests/unit/test_command.py.
    """
    raise NotImplementedError(f"parse_basket({text!r}) is not implemented yet")


@dataclass(frozen=True)
class CommandResult:
    """What the program prints and the exit code it ends with."""

    output: str
    exit_code: Literal[0, 1]


def run_checkout_command(text: str, catalog: Catalog) -> CommandResult:
    """Runs one checkout from the JSON text of a basket.

    Not implemented yet - see tests/unit/test_command.py.
    """
    raise NotImplementedError("run_checkout_command() is not implemented yet")


def main(catalog: Catalog | None = None) -> int:
    """The whole program: stdin -> run_checkout_command -> stdout, returning the exit code.

    Not implemented yet - see tests/unit/test_command.py.
    """
    raise NotImplementedError("main() is not implemented yet")
