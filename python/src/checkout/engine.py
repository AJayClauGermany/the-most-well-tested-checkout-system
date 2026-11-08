"""The checkout engine: turns a basket into a receipt."""

from checkout.models import Basket, Catalog, Receipt


def checkout(basket: Basket, catalog: Catalog) -> Receipt:
    """Prices ``basket`` using the products in ``catalog``.

    Follows the steps in spec/checkout-rules.md, section 4. Returns a Receipt,
    or raises ``CheckoutError`` with one of the error codes from section 9.
    """
    # This project is built test-first: the tests in tests/ already describe what
    # this function must do. Write the code to make them pass!
    raise NotImplementedError("checkout() is not implemented yet")
