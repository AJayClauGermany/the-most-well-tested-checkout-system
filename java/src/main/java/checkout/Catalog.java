package checkout;

import java.util.List;
import java.util.Optional;

/**
 * The list of products the shop sells (rules section 2).
 *
 * <p>The checkout receives the catalog as an input, so tests can use a tiny,
 * made-up catalog instead of the standard one.
 *
 * @param products every product, in catalog order
 */
public record Catalog(List<Product> products) {

    /**
     * Finds the product with exactly this SKU.
     *
     * <p>SKUs are case sensitive: {@code "apple"} does not find {@code "APPLE"}.
     *
     * @param sku the product code to look for
     * @return the product, or {@link Optional#empty()} when there is no such product
     */
    public Optional<Product> findBySku(String sku) {
        return products.stream()
                .filter(product -> product.sku().equals(sku))
                .findFirst();
    }
}
