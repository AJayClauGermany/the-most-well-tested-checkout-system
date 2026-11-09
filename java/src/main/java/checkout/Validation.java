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
     * @param quantity the quantity of one basket line
     * @return {@code true} when it is allowed
     */
    public static boolean isValidQuantity(int quantity) {
        throw new UnsupportedOperationException("not implemented yet");
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
        throw new UnsupportedOperationException("not implemented yet");
    }
}
