package checkout;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.List;

/**
 * The checkout output (rules section 5). All amounts are whole cents.
 *
 * <p>We use {@code int} for money: a whole number of cents can never have the
 * rounding problems of {@code double} ({@code 0.1 + 0.2 != 0.3}!).
 *
 * @param lines             one line per product, in order of first appearance in the basket
 * @param subtotal          sum of all line totals
 * @param promotionDiscount sum of all line discounts
 * @param couponDiscount    the coupon discount ({@code 0} without a coupon)
 * @param shipping          the shipping cost ({@code 0} when shipping is free)
 * @param tax               8% tax on the discounted subtotal
 * @param total             what the shopper pays
 */
public record Receipt(
        List<ReceiptLine> lines,
        int subtotal,
        @JsonProperty("promotion_discount") int promotionDiscount,
        @JsonProperty("coupon_discount") int couponDiscount,
        int shipping,
        int tax,
        int total) {}
