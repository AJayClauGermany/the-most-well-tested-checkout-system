"""Shared pytest fixtures.

pytest finds this file automatically. A *fixture* is a function that prepares
something a test needs. A test asks for it simply by naming it as a parameter:

    def test_something(catalog):
        ...

Fixtures defined here can be used by every test below the tests/ folder.
"""

import pytest
from spec_loader import SpecCase, load_catalog, load_spec_cases

from checkout import Catalog


@pytest.fixture(scope="session")
def catalog() -> Catalog:
    """The standard catalog from spec/catalog.json.

    scope="session" means the file is read once and shared by all tests.
    That is fine because tests only read the catalog, they never change it.
    """
    return load_catalog()


@pytest.fixture(scope="session")
def spec_cases() -> list[SpecCase]:
    """Every case from spec/test-cases/*.json."""
    return load_spec_cases()
