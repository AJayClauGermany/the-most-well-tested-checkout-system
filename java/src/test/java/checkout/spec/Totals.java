package checkout.spec;

import checkout.Receipt;
import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * The six totals a spec case expects (everything on the receipt except the lines).
 *
 * <p>Because this is a record, two {@code Totals} are {@code equals} when all six numbers
 * are the same. So a test can compare all six at once, and when they differ JUnit
 * prints both records so you can see which number is wrong.
 *
 * @param subtotal          sum of all line totals, in cents
 * @param promotionDiscount sum of all line discounts, in cents
 * @param couponDiscount    the coupon discount, in cents
 * @param shipping          the shipping cost, in cents
 * @param tax               the tax, in cents
 * @param total             what the shopper pays, in cents
 */
public record Totals(
        int subtotal,
        @JsonProperty("promotion_discount") int promotionDiscount,
        @JsonProperty("coupon_discount") int couponDiscount,
        int shipping,
        int tax,
        int total) {

    /**
     * Picks the six totals out of a receipt.
     *
     * @param receipt a receipt returned by the checkout
     * @return its totals
     */
    public static Totals of(Receipt receipt) {
        return new Totals(
                receipt.subtotal(),
                receipt.promotionDiscount(),
                receipt.couponDiscount(),
                receipt.shipping(),
                receipt.tax(),
                receipt.total());
    }
}
