"""Looking products up in the catalog (rules §2 and §4.1 row 2)."""

from checkout.errors import CheckoutError, ErrorCode
from checkout.models import Catalog, Product


def require_product(catalog: Catalog, sku: str) -> Product:
    """Like ``Catalog.find``, but a missing product is an error.

    SKUs are case sensitive (rules §4.1): "apple" does not find "APPLE".
    Raises ``CheckoutError(UNKNOWN_PRODUCT)`` when there is no such product.
    """
    product = catalog.find(sku)
    if product is None:
        raise CheckoutError(ErrorCode.UNKNOWN_PRODUCT, f"Unknown product: {sku}")
    return product
