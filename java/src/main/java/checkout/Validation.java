package checkout;

/**
 * Basket validation (rules section 4.1, rows 1 to 4). Rows 5 to 7 (shipping and coupon)
 * are checked by {@link ShippingMethod} and {@link Coupons}, because the coupon check
 * needs the price after promotions.
 */
public final class Validation {

    /** The smallest quantity a basket line may have. */
    public static final int MIN_QUANTITY = 1;
    /** The largest quantity a basket line (and all lines of one SKU together) may have. */
    public static final int MAX_QUANTITY = 99;

    private Validation() {}

    /**
     * Returns whether a quantity is a whole number from 1 to 99 (rules section 4.1 row 3).
     *
     * <p>In Java a quantity is an {@code int}, so "a whole number" is guaranteed by the type:
     * {@code 1.5} cannot even get this far (the JSON reader rejects it, see
     * {@code checkout.cli.CheckoutCommand}). Only the range is left to check.
     *
     * @param quantity the quantity of one basket line
     * @return {@code true} when it is allowed
     */
    public static boolean isValidQuantity(int quantity) {
        return quantity >= MIN_QUANTITY && quantity <= MAX_QUANTITY;
    }

    /**
     * Throws a {@link CheckoutException} for rows 1 to 4 of rules section 4.1, in that order.
     * Does nothing when the items are fine.
     *
     * @param basket  the basket to check
     * @param catalog the products that exist
     * @throws CheckoutException EMPTY_BASKET, UNKNOWN_PRODUCT or INVALID_QUANTITY
     */
    public static void validateItems(Basket basket, Catalog catalog) {
        // Row 1: no items at all.
        if (basket.items().isEmpty()) {
            throw new CheckoutException(ErrorCode.EMPTY_BASKET, "The basket has no items");
        }

        // Row 2: every SKU must exist. All SKUs are checked before any quantity, so an unknown
        // product on line 3 still wins over a bad quantity on line 1.
        for (BasketItem item : basket.items()) {
            catalog.requireBySku(item.sku());
        }

        // Row 3: every quantity must be from 1 to 99.
        for (BasketItem item : basket.items()) {
            if (!isValidQuantity(item.quantity())) {
                throw new CheckoutException(
                        ErrorCode.INVALID_QUANTITY, "Invalid quantity for " + item.sku() + ": " + item.quantity());
            }
        }

        // Row 4: the lines of one SKU, added together, may not go over 99 either.
        // (Each line is at most 99 by now, so the sum cannot overflow an int.)
        for (BasketItem line : Lines.mergeItems(basket.items())) {
            if (line.quantity() > MAX_QUANTITY) {
                throw new CheckoutException(
                        ErrorCode.INVALID_QUANTITY, "Too many " + line.sku() + " in total: " + line.quantity());
            }
        }
    }
}
