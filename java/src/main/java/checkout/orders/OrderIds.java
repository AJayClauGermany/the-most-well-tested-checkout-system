package checkout.orders;

import java.util.OptionalInt;

/**
 * The shape of an order id: {@code "ORDER-"} followed by a number from 1 upwards (rules section 6).
 */
public final class OrderIds {

    /** Every order id starts with this. */
    public static final String PREFIX = "ORDER-";

    private OrderIds() {}

    /**
     * The id of the n-th order saved into a repository.
     *
     * @param n 1 for the first order, 2 for the second, ...
     * @return {@code "ORDER-1"}, {@code "ORDER-2"}, ...
     */
    public static String of(int n) {
        throw new UnsupportedOperationException("not implemented yet");
    }

    /**
     * The inverse of {@link #of}: {@code "ORDER-7"} gives 7.
     *
     * @param id any text
     * @return the number, or empty when the text is not an order id
     *         (wrong prefix, leading zero, {@code "ORDER-0"}, spaces, ...)
     */
    public static OptionalInt numberOf(String id) {
        throw new UnsupportedOperationException("not implemented yet");
    }
}
