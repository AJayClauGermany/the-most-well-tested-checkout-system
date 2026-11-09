package checkout.spec;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Test helper: turns {@code "$12.39"} into {@code 1239} cents - without ever touching a
 * floating point number ({@code spec/features/README.md} insists, and rules section 1
 * explains why: {@code 0.1 + 0.2 != 0.3}).
 */
public final class Dollars {

    /** Optional "$", some digits, optionally a dot and one or two more digits. */
    private static final Pattern AMOUNT = Pattern.compile("^\\$?(\\d+)(?:\\.(\\d{1,2}))?$");

    private Dollars() {}

    /**
     * Converts a dollar amount to whole cents.
     *
     * @param text {@code "$12.39"}, {@code "12.39"}, {@code "$5"} or {@code "$5.5"}
     * @return 1239, 1239, 500 or 550
     * @throws IllegalArgumentException when the text is not a dollar amount
     */
    public static int toCents(String text) {
        Matcher matcher = AMOUNT.matcher(text.strip());
        if (!matcher.matches()) {
            throw new IllegalArgumentException("not a dollar amount: \"" + text + "\"");
        }
        int dollars = Integer.parseInt(matcher.group(1));
        String centsText = matcher.group(2) == null ? "0" : matcher.group(2);
        // "5" means 50 cents, not 5: pad to two digits first.
        int cents = Integer.parseInt(centsText.length() == 1 ? centsText + "0" : centsText);
        return dollars * 100 + cents;
    }
}
