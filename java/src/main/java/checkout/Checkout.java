package checkout;

import java.util.List;

/**
 * The checkout: turns a {@link Basket} into a {@link Receipt}, following
 * {@code spec/checkout-rules.md}.
 *
 * <p>It reads like the hint in rules section 4.7: validate, merge lines, price, coupon,
 * shipping, tax, total. Each step is a small method in its own class ({@link Validation},
 * {@link Lines}, {@link Promotions}, {@link Coupons}, {@link ShippingMethod}, {@link Tax},
 * {@link Money}) with its own unit tests.
 */
public final class Checkout {

    // Nobody needs a Checkout object: everything is a static method.
    private Checkout() {}

    /**
     * Calculates the receipt for a basket (rules section 4).
     *
     * @param basket  what the shopper wants to buy
     * @param catalog the products and their prices
     * @return the receipt with every line and all totals, in cents
     * @throws CheckoutException when the basket is invalid (see {@link ErrorCode})
     */
    public static Receipt checkout(Basket basket, Catalog catalog) {
        Validation.validateItems(basket, catalog); // rules section 4.1 rows 1-4
        ShippingMethod method = ShippingMethod.parse(basket.shipping()); // row 5

        // After validation every merged quantity is 1..99, so it fits the int of a receipt line.
        // Math.toIntExact would throw (instead of cutting the number off) if it ever did not.
        List<ReceiptLine> lines = Lines.mergeItems(basket.items()).stream()
                .map(item -> Lines.priceLine(catalog.requireBySku(item.sku()), Math.toIntExact(item.quantity())))
                .toList();
        int subtotal = lines.stream().mapToInt(ReceiptLine::lineTotal).sum();
        int promotionDiscount = lines.stream().mapToInt(ReceiptLine::lineDiscount).sum();
        int afterPromotions = subtotal - promotionDiscount;

        int couponDiscount = Coupons.discount(basket.coupon(), afterPromotions); // rows 6-7
        int discountedSubtotal = afterPromotions - couponDiscount;

        int shipping = method.cost(discountedSubtotal);
        int tax = Tax.taxOn(discountedSubtotal);

        return new Receipt(
                lines,
                subtotal,
                promotionDiscount,
                couponDiscount,
                shipping,
                tax,
                discountedSubtotal + shipping + tax);
    }
}
