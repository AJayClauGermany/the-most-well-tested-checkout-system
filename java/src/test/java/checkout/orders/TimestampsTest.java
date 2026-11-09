package checkout.orders;

import static org.assertj.core.api.Assertions.assertThat;

import java.time.Clock;
import java.time.Instant;
import java.time.ZoneId;
import java.time.ZoneOffset;
import org.junit.jupiter.api.Test;

/**
 * UNIT tests for the {@code created_at} format (rules section 6).
 */
class TimestampsTest {

    @Test
    void writesUtcWithWholeSecondsAndATrailingZ() {
        assertThat(Timestamps.format(Instant.parse("2026-01-15T10:30:00Z"))).isEqualTo("2026-01-15T10:30:00Z");
    }

    @Test
    void dropsTheMillisecondsInsteadOfRoundingThem() {
        assertThat(Timestamps.format(Instant.parse("2026-01-15T10:30:00.999Z"))).isEqualTo("2026-01-15T10:30:00Z");
    }

    @Test
    void alwaysWritesTwoDigitsAndAFullYear() {
        assertThat(Timestamps.format(Instant.parse("2030-02-03T04:05:06Z"))).isEqualTo("2030-02-03T04:05:06Z");
    }

    @Test
    void readsTheClockAndConvertsOtherTimeZonesToUtc() {
        // A STUB clock: it always answers 12:30 in Berlin (UTC+2 in summer... here written as +02:00).
        Clock berlin = Clock.fixed(Instant.parse("2026-01-15T10:30:00Z"), ZoneId.of("Europe/Berlin"));
        Clock offset = Clock.fixed(Instant.parse("2026-01-15T12:30:00+02:00"), ZoneOffset.ofHours(2));

        assertThat(Timestamps.now(berlin)).isEqualTo("2026-01-15T10:30:00Z");
        assertThat(Timestamps.now(offset)).isEqualTo("2026-01-15T10:30:00Z");
    }
}
