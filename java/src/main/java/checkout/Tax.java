package checkout;

/**
 * Tax (rules section 4.6).
 */
public final class Tax {

    /** The tax rate in percent. */
    public static final int TAX_PERCENT = 8;

    private Tax() {}

    /**
     * 8% of {@code discountedSubtotal}, rounded half up. Shipping is never passed in here.
     *
     * @param discountedSubtotal the price after promotions and coupon, in cents
     * @return the tax in cents
     */
    public static int taxOn(int discountedSubtotal) {
        return Money.roundHalfUp(discountedSubtotal, TAX_PERCENT);
    }
}
