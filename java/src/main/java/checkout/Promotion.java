package checkout;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonSubTypes;
import com.fasterxml.jackson.annotation.JsonTypeInfo;

/**
 * A promotion that lowers the price of a product line (rules section 3).
 *
 * <p>There are exactly two kinds of promotion, so {@code Promotion} is a
 * <em>sealed</em> interface: the only classes allowed to implement it are the two
 * records written inside it. Later, the implementation can use a {@code switch}
 * and the compiler will check that both kinds are handled:
 *
 * <pre>{@code
 * int discount = switch (promotion) {
 *     case Promotion.BuyXGetYFree p -> ...;
 *     case Promotion.BulkDiscount p -> ...;
 * };
 * }</pre>
 *
 * <p>In JSON the kind is stored in a {@code "type"} field. The two Jackson
 * annotations below tell Jackson which record to create for which {@code "type"}.
 */
@JsonTypeInfo(use = JsonTypeInfo.Id.NAME, property = "type")
@JsonSubTypes({
    @JsonSubTypes.Type(value = Promotion.BuyXGetYFree.class, name = "buy_x_get_y_free"),
    @JsonSubTypes.Type(value = Promotion.BulkDiscount.class, name = "bulk_discount")
})
public sealed interface Promotion {

    /**
     * "Buy X, get Y free": for every complete group of {@code buy + free} items,
     * {@code free} of them cost nothing.
     *
     * <p>JSON: {@code {"type": "buy_x_get_y_free", "buy": 2, "free": 1}}
     *
     * @param buy  how many items you pay for in each group
     * @param free how many items you get for free in each group
     */
    record BuyXGetYFree(int buy, int free) implements Promotion {}

    /**
     * "Bulk discount": {@code percent}% off the line total when you buy at least
     * {@code minQuantity} items.
     *
     * <p>JSON: {@code {"type": "bulk_discount", "min_quantity": 10, "percent": 10}}
     *
     * @param minQuantity the smallest quantity that gets the discount
     * @param percent     the discount in percent (10 means 10%)
     */
    record BulkDiscount(@JsonProperty("min_quantity") int minQuantity, int percent) implements Promotion {}
}
