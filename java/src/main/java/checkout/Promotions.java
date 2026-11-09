package checkout;

/**
 * The two kinds of promotion (rules section 4.2).
 */
public final class Promotions {

    private Promotions() {}

    /**
     * The discount a promotion gives one receipt line, in cents. {@code 0} without a promotion.
     *
     * @param promotion the product's promotion, or {@code null} when it has none
     * @param quantity  how many of the product are on the line
     * @param unitPrice the price of one item, in cents
     * @return the line discount in cents
     */
    public static int discount(Promotion promotion, int quantity, int unitPrice) {
        // A "switch with patterns": because Promotion is sealed, the compiler knows these are
        // the only two kinds, and would refuse to compile if a third kind were added and
        // forgotten here. "case null" handles products without a promotion.
        return switch (promotion) {
            case null -> 0;
            case Promotion.BuyXGetYFree buyXGetYFree -> {
                // For every complete group of `buy + free` items, `free` of them cost nothing.
                int groupSize = buyXGetYFree.buy() + buyXGetYFree.free();
                int freeItems = quantity / groupSize * buyXGetYFree.free();
                yield freeItems * unitPrice;
            }
            case Promotion.BulkDiscount bulk ->
                // `percent`% off the whole line, but only from `minQuantity` items.
                quantity < bulk.minQuantity() ? 0 : Money.roundHalfUp(unitPrice * quantity, bulk.percent());
        };
    }
}
