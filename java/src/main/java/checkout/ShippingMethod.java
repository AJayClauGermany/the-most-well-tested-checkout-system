package checkout;

/**
 * The two shipping methods and what they cost (rules section 4.5).
 *
 * <p>An {@code enum} is the natural Java type for "one of a fixed set of values". Once the
 * basket's free-form text has been {@link #parse}d into one of these, the rest of the code
 * cannot be handed an unknown method any more.
 */
public enum ShippingMethod {
    /** $4.99, free from $50.00 (after all discounts). */
    STANDARD(499),
    /** $9.99, never free. */
    EXPRESS(999);

    /** Standard shipping is free from this discounted subtotal upwards, in cents. */
    public static final int FREE_STANDARD_SHIPPING_FROM = 5000;

    private final int baseCost;

    ShippingMethod(int baseCost) {
        this.baseCost = baseCost;
    }

    /**
     * Turns the basket's {@code shipping} field into a known method.
     *
     * @param shipping {@code "standard"} or {@code "express"}; {@code null} means standard
     * @return the method
     * @throws CheckoutException UNKNOWN_SHIPPING_METHOD for anything else (rules section 4.1 row 5)
     */
    public static ShippingMethod parse(String shipping) {
        throw new UnsupportedOperationException("not implemented yet");
    }

    /**
     * The shipping cost in cents for this method and this discounted subtotal.
     *
     * @param discountedSubtotal the price after promotions and coupon, in cents
     * @return {@code 0}, {@code 499} or {@code 999}
     */
    public int cost(int discountedSubtotal) {
        throw new UnsupportedOperationException("not implemented yet");
    }

    /**
     * The cost before any free-shipping rule.
     *
     * @return {@code 499} for standard, {@code 999} for express
     */
    public int baseCost() {
        return baseCost;
    }
}
