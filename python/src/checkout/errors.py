"""The ways a checkout can fail (see spec/checkout-rules.md, section 9)."""

from enum import StrEnum


class ErrorCode(StrEnum):
    """The six error codes from the spec.

    A ``StrEnum`` member *is* a string, so ``ErrorCode.EMPTY_BASKET == "EMPTY_BASKET"`` is
    ``True``, and ``str(ErrorCode.EMPTY_BASKET)`` prints just ``EMPTY_BASKET``.
    """

    EMPTY_BASKET = "EMPTY_BASKET"
    UNKNOWN_PRODUCT = "UNKNOWN_PRODUCT"
    INVALID_QUANTITY = "INVALID_QUANTITY"
    UNKNOWN_SHIPPING_METHOD = "UNKNOWN_SHIPPING_METHOD"
    UNKNOWN_COUPON = "UNKNOWN_COUPON"
    COUPON_MINIMUM_NOT_MET = "COUPON_MINIMUM_NOT_MET"


class CheckoutError(Exception):
    """Raised by ``checkout()`` when a basket cannot be checked out.

    The reason is in ``error.code``; ``str(error)`` is a message for humans::

        try:
            checkout(basket, catalog)
        except CheckoutError as error:
            print(error.code)  # e.g. ErrorCode.EMPTY_BASKET
            print(error)       # e.g. "The basket has no items"
    """

    def __init__(self, code: ErrorCode, message: str | None = None) -> None:
        # Without a message, the code itself is the message (as the baseline did).
        super().__init__(message if message is not None else code.value)
        self.code = code
