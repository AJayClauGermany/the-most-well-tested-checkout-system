package checkout;

/**
 * One line of a basket: which product, and how many.
 *
 * <p>JSON: {@code {"sku": "COFFEE", "quantity": 3}}
 *
 * <p>The quantity is a {@code long} (64 bits), not an {@code int} (32 bits): the basket schema
 * allows any whole number, and 10000000000 must be reported as INVALID_QUANTITY, not rejected
 * because it does not fit. Go and Rust read quantities as 64-bit numbers too.
 *
 * @param sku      the product code (may be unknown - the checkout must report that)
 * @param quantity how many items (may be invalid, e.g. 0, 100 or 10000000000 - the checkout must report that)
 */
public record BasketItem(String sku, long quantity) {}
