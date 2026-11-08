"""Shared pytest fixtures.

pytest finds this file automatically. A *fixture* is a function that prepares
something a test needs. A test asks for it simply by naming it as a parameter:

    def test_something(catalog):
        ...

Fixtures defined here can be used by every test below the tests/ folder.
"""

import pytest
from hypothesis import HealthCheck, settings
from spec_loader import SpecCase, load_catalog, load_spec_cases

from checkout import Catalog

# Hypothesis settings (tests/property/). A "deadline" fails a test whose single example is
# slow; on a busy CI machine that is noise, not a bug, so it is switched off. The profile is
# registered here because conftest.py is loaded before any test module.
settings.register_profile(
    "checkout",
    deadline=None,
    suppress_health_check=[HealthCheck.too_slow],
)
settings.load_profile("checkout")


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
