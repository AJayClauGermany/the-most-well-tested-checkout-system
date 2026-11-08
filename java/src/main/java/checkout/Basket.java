package checkout;

import java.util.List;

/**
 * The checkout input (rules section 3).
 *
 * <p>JSON: {@code {"items": [...], "coupon": "SAVE10", "shipping": "express"}}
 *
 * @param items    the basket lines, in the order the shopper added them
 * @param coupon   a coupon code, or {@code null} for no coupon
 * @param shipping {@code "standard"} or {@code "express"}; {@code null} means {@code "standard"}
 */
public record Basket(List<BasketItem> items, String coupon, String shipping) {

    /** The shipping method used when the basket does not say which one it wants. */
    public static final String DEFAULT_SHIPPING = "standard";

    /**
     * This is a "compact constructor": it runs every time a {@code Basket} is created
     * (also when Jackson reads one from JSON) and fills in the default shipping method.
     */
    public Basket {
        if (shipping == null) {
            shipping = DEFAULT_SHIPPING;
        }
    }

    /** Convenience constructor: no coupon and standard shipping. */
    public Basket(List<BasketItem> items) {
        this(items, null, null);
    }
}
