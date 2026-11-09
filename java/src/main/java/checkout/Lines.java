package checkout;

import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * Merging basket items into receipt lines and pricing them (rules section 4.2).
 */
public final class Lines {

    private Lines() {}

    /**
     * One item per distinct SKU, in order of first appearance, quantities added up.
     *
     * <p>{@code [APPLE x1, PEN x2, APPLE x3]} gives {@code [APPLE x4, PEN x2]}.
     *
     * @param items the basket items, in basket order
     * @return the merged items; the input list is not changed
     */
    public static List<BasketItem> mergeItems(List<BasketItem> items) {
        // A LinkedHashMap remembers the order keys were first added - exactly the order the
        // rules want. (A plain HashMap would give the SKUs back in an unpredictable order.)
        Map<String, Long> quantities = new LinkedHashMap<>();
        for (BasketItem item : items) {
            // Math::addExact throws an ArithmeticException instead of silently wrapping around
            // to a wrong (negative) number when the sum does not fit a long.
            quantities.merge(item.sku(), item.quantity(), Math::addExact);
        }
        return quantities.entrySet().stream()
                .map(entry -> new BasketItem(entry.getKey(), entry.getValue()))
                .toList();
    }

    /**
     * The receipt line for {@code quantity} of {@code product}, including its promotion discount.
     *
     * @param product  the product
     * @param quantity how many of it (already merged and validated)
     * @return the receipt line
     */
    public static ReceiptLine priceLine(Product product, int quantity) {
        return new ReceiptLine(
                product.sku(),
                product.name(),
                quantity,
                product.unitPrice(),
                product.unitPrice() * quantity,
                Promotions.discount(product.promotion(), quantity, product.unitPrice()));
    }
}
