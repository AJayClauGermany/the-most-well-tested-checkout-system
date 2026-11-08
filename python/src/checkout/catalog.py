"""Looking products up in the catalog (rules §2 and §4.1 row 2)."""

from checkout.models import Catalog, Product


def require_product(catalog: Catalog, sku: str) -> Product:
    """Like ``Catalog.find``, but a missing product is an error (UNKNOWN_PRODUCT).

    Not implemented yet - see tests/unit/test_catalog.py.
    """
    raise NotImplementedError(f"require_product({catalog}, {sku!r}) is not implemented yet")
