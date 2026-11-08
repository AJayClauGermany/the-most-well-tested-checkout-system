package checkout;

/**
 * One line of a basket: which product, and how many.
 *
 * <p>JSON: {@code {"sku": "COFFEE", "quantity": 3}}
 *
 * @param sku      the product code (may be unknown - the checkout must report that)
 * @param quantity how many items (may be invalid, e.g. 0 or 100 - the checkout must report that)
 */
public record BasketItem(String sku, int quantity) {}
