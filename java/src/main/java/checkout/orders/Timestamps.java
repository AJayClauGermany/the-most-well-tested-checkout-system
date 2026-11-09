package checkout.orders;

import java.time.Clock;
import java.time.Instant;
import java.time.temporal.ChronoUnit;

/**
 * The {@code created_at} format of a saved order (rules section 6):
 * UTC, ISO 8601, whole seconds, trailing {@code Z}.
 *
 * <p>The current time comes from a {@link Clock}. Production code uses
 * {@link Clock#systemUTC()}; tests use {@link Clock#fixed} so that {@code created_at}
 * is predictable. That is the textbook reason for a test double.
 */
public final class Timestamps {

    private Timestamps() {}

    /**
     * Formats an instant the way {@code created_at} must look.
     *
     * @param instant a moment in time
     * @return e.g. {@code "2026-01-15T10:30:00Z"} (milliseconds are cut off, not rounded)
     */
    public static String format(Instant instant) {
        // Instant.toString() is already ISO 8601 in UTC with a Z; it only adds fractions of a
        // second when there are some, so cut those off first.
        return instant.truncatedTo(ChronoUnit.SECONDS).toString();
    }

    /**
     * What the clock says right now, formatted with {@link #format}.
     *
     * @param clock the time source
     * @return the current time as {@code created_at} text
     */
    public static String now(Clock clock) {
        return format(clock.instant());
    }
}
