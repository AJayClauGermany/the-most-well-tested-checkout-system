package checkout;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * One product in the catalog.
 *
 * <p>A {@code record} is a short way to write a class that only holds data: Java
 * creates the constructor, the getters ({@code sku()}, {@code name()}, ...),
 * {@code equals}, {@code hashCode} and {@code toString} for us.
 *
 * <p>{@code @JsonProperty("unit_price")} tells Jackson that the JSON field
 * {@code unit_price} (snake_case) goes into the Java field {@code unitPrice} (camelCase).
 *
 * @param sku       unique product code, e.g. {@code "COFFEE"} (case sensitive)
 * @param name      human readable name, e.g. {@code "Bag of Coffee Beans"}
 * @param unitPrice price of one item in whole cents ({@code 899} means $8.99)
 * @param promotion the product's promotion, or {@code null} when it has none
 */
public record Product(
        String sku,
        String name,
        @JsonProperty("unit_price") int unitPrice,
        Promotion promotion) {

    /** Convenience constructor for a product without a promotion. */
    public Product(String sku, String name, int unitPrice) {
        this(sku, name, unitPrice, null);
    }
}
