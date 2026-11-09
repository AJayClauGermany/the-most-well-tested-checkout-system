package checkout;

/**
 * Money helpers (rules sections 1 and 7). All amounts are whole cents.
 */
public final class Money {

    private Money() {}

    /**
     * {@code percent}% of {@code amount}, rounded half up to a whole cent (rules section 7).
     *
     * @param amount  an amount in cents
     * @param percent a percentage, e.g. {@code 10} for 10%
     * @return the rounded result in cents
     */
    public static int roundHalfUp(int amount, int percent) {
        throw new UnsupportedOperationException("not implemented yet");
    }
}
