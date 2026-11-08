package checkout;

/**
 * The checkout: turns a {@link Basket} into a {@link Receipt}.
 *
 * <p><strong>Not implemented yet!</strong> This project is built test-first: the
 * shared spec cases in {@code spec/test-cases/} are already loaded by
 * {@code SharedSpecCasesTest}, waiting (disabled) for an implementation.
 * See {@code spec/checkout-rules.md} for the rules to implement.
 */
public final class Checkout {

    // Nobody needs a Checkout object: everything is a static method.
    private Checkout() {}

    /**
     * Calculates the receipt for a basket (rules section 4, steps 1 to 7).
     *
     * @param basket  what the shopper wants to buy
     * @param catalog the products and their prices
     * @return the receipt with every line and all totals, in cents
     * @throws CheckoutException when the basket is invalid (see {@link ErrorCode})
     */
    public static Receipt checkout(Basket basket, Catalog catalog) {
        throw new UnsupportedOperationException("not implemented yet");
    }
}
