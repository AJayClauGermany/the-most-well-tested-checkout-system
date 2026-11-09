package checkout;

/**
 * Money helpers (rules sections 1 and 7). All amounts are whole cents.
 */
public final class Money {

    private Money() {}

    /**
     * {@code percent}% of {@code amount}, rounded half up to a whole cent (rules section 7).
     *
     * <p>Done with whole numbers only, as the hint in rules section 7 suggests:
     * {@code amount * percent} is the exact answer in hundredths of a cent; adding 50 hundredths
     * and dividing by 100 (Java's {@code /} on whole numbers rounds down for positive numbers)
     * is the same as rounding half up. No {@code double}, so no {@code 0.1 + 0.2} surprises.
     *
     * <pre>
     * roundHalfUp(225, 10)  -&gt;  (2250 + 50) / 100  -&gt;  23     (22.5 rounds UP)
     * roundHalfUp(224, 10)  -&gt;  (2240 + 50) / 100  -&gt;  22     (22.9 rounded down)
     * </pre>
     *
     * <p>The multiplication is done in {@code long}, so even 100% of the largest {@code int}
     * cannot overflow on the way.
     *
     * @param amount  an amount in cents, not negative
     * @param percent a percentage from 0 to 100, e.g. {@code 10} for 10%
     * @return the rounded result in cents
     */
    public static int roundHalfUp(int amount, int percent) {
        return (int) (((long) amount * percent + 50) / 100);
    }
}
