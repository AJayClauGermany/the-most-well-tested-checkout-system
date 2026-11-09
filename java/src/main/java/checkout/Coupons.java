package checkout;

import java.util.Optional;

/**
 * Coupons (rules section 4.4): {@code SAVE10} and {@code WELCOME5}.
 */
public final class Coupons {

    private Coupons() {}

    /**
     * The coupon code as the shop knows it: trimmed and upper case.
     * {@code null}, empty or blank all mean "no coupon" and give an empty Optional.
     *
     * @param coupon the code as the shopper typed it, or {@code null}
     * @return the normalized code, or empty for "no coupon"
     */
    public static Optional<String> normalizeCode(String coupon) {
        throw new UnsupportedOperationException("not implemented yet");
    }

    /**
     * The coupon discount in cents for {@code afterPromotions} cents of shopping.
     *
     * @param coupon          the code as the shopper typed it, or {@code null}
     * @param afterPromotions the subtotal minus the promotion discount, in cents
     * @return the discount in cents ({@code 0} without a coupon)
     * @throws CheckoutException UNKNOWN_COUPON or COUPON_MINIMUM_NOT_MET (rules section 4.1 rows 6 and 7)
     */
    public static int discount(String coupon, int afterPromotions) {
        throw new UnsupportedOperationException("not implemented yet");
    }
}
