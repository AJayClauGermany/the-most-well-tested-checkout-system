package checkout;

import java.util.List;

/**
 * Merging basket items into receipt lines and pricing them (rules section 4.2).
 */
public final class Lines {

    private Lines() {}

    /**
     * One item per distinct SKU, in order of first appearance, quantities added up.
     *
     * @param items the basket items, in basket order
     * @return the merged items; the input list is not changed
     */
    public static List<BasketItem> mergeItems(List<BasketItem> items) {
        throw new UnsupportedOperationException("not implemented yet");
    }

    /**
     * The receipt line for {@code quantity} of {@code product}, including its promotion discount.
     *
     * @param product  the product
     * @param quantity how many of it (already merged and validated)
     * @return the receipt line
     */
    public static ReceiptLine priceLine(Product product, int quantity) {
        throw new UnsupportedOperationException("not implemented yet");
    }
}
