package checkout;

import java.util.Locale;
import java.util.Map;
import java.util.Optional;
import java.util.function.IntUnaryOperator;

/**
 * Coupons (rules section 4.4): {@code SAVE10} and {@code WELCOME5}.
 */
public final class Coupons {

    /**
     * What the shop knows about one coupon code.
     *
     * @param discount     the discount in cents, given the price after promotions
     * @param minimumSpend the least the customer must spend (after promotions) to use it, in cents
     */
    private record Coupon(IntUnaryOperator discount, int minimumSpend) {}

    /** The coupons that exist, by their normalized code. */
    private static final Map<String, Coupon> COUPONS = Map.of(
            "SAVE10", new Coupon(afterPromotions -> Money.roundHalfUp(afterPromotions, 10), 0),
            "WELCOME5", new Coupon(afterPromotions -> 500, 2500));

    private Coupons() {}

    /**
     * The coupon code as the shop knows it: trimmed and upper case.
     * {@code null}, empty or blank all mean "no coupon" and give an empty Optional.
     *
     * @param coupon the code as the shopper typed it, or {@code null}
     * @return the normalized code, or empty for "no coupon"
     */
    public static Optional<String> normalizeCode(String coupon) {
        // Locale.ROOT: upper-casing must not depend on the computer's language settings
        // (in Turkish, "i".toUpperCase() is "İ", which would break "welcome5").
        return Optional.ofNullable(coupon)
                .map(code -> code.strip().toUpperCase(Locale.ROOT))
                .filter(code -> !code.isEmpty());
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
        Optional<String> normalized = normalizeCode(coupon);
        if (normalized.isEmpty()) {
            return 0;
        }
        String code = normalized.get();

        Coupon known = COUPONS.get(code);
        if (known == null) {
            throw new CheckoutException(ErrorCode.UNKNOWN_COUPON, "Unknown coupon: " + code);
        }
        if (afterPromotions < known.minimumSpend()) {
            throw new CheckoutException(
                    ErrorCode.COUPON_MINIMUM_NOT_MET,
                    code + " needs a spend of at least " + known.minimumSpend() + " cents, got " + afterPromotions);
        }
        return known.discount().applyAsInt(afterPromotions);
    }
}
