package checkout.orders;

import checkout.Receipt;
import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * A saved order (rules section 6), exactly as {@code spec/schema/order.schema.json} describes it.
 *
 * @param id        {@code "ORDER-1"}, {@code "ORDER-2"}, ... - see {@link OrderIds}
 * @param createdAt when it was saved: UTC, ISO 8601, whole seconds, trailing Z
 *                  ({@code "2026-01-15T10:30:00Z"}) - see {@link Timestamps}
 * @param receipt   the receipt exactly as it was saved
 */
public record Order(
        String id,
        @JsonProperty("created_at") String createdAt,
        Receipt receipt) {}
