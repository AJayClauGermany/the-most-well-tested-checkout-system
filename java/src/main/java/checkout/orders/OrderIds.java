package checkout.orders;

import java.util.OptionalInt;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * The shape of an order id: {@code "ORDER-"} followed by a number from 1 upwards (rules section 6).
 */
public final class OrderIds {

    /** Every order id starts with this. */
    public static final String PREFIX = "ORDER-";

    /** A positive whole number with no leading zero, exactly as order.schema.json spells it. */
    private static final Pattern ORDER_ID = Pattern.compile("ORDER-([1-9][0-9]*)");

    private OrderIds() {}

    /**
     * The id of the n-th order saved into a repository.
     *
     * @param n 1 for the first order, 2 for the second, ...
     * @return {@code "ORDER-1"}, {@code "ORDER-2"}, ...
     */
    public static String of(int n) {
        return PREFIX + n;
    }

    /**
     * The inverse of {@link #of}: {@code "ORDER-7"} gives 7.
     *
     * @param id any text
     * @return the number, or empty when the text is not an order id
     *         (wrong prefix, leading zero, {@code "ORDER-0"}, spaces, ...)
     */
    public static OptionalInt numberOf(String id) {
        Matcher matcher = ORDER_ID.matcher(id);
        if (!matcher.matches()) {
            return OptionalInt.empty();
        }
        try {
            return OptionalInt.of(Integer.parseInt(matcher.group(1)));
        } catch (NumberFormatException tooBig) {
            // "ORDER-99999999999" has the right shape, but no repository can ever hand it out.
            return OptionalInt.empty();
        }
    }
}
