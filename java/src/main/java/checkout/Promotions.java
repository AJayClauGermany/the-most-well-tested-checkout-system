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
        throw new UnsupportedOperationException("not implemented yet");
    }
}
