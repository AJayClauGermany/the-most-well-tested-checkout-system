package checkout.spec;

import checkout.Basket;
import checkout.BasketItem;
import checkout.Catalog;
import checkout.Product;
import checkout.Promotion;
import checkout.Receipt;
import checkout.ReceiptLine;
import java.util.List;

/**
 * Test helper: small, hand-made test data for unit tests.
 *
 * <p>Unit tests prefer tiny made-up data over the real catalog: when a test fails you can
 * work out the right answer in your head. The shared spec cases use the real catalog.
 */
public final class Fixtures {

    public static final Product APPLE = new Product("APPLE", "Apple", 50);
    public static final Product PEN = new Product("PEN", "Pen", 120);
    public static final Product COFFEE =
            new Product("COFFEE", "Bag of Coffee Beans", 899, new Promotion.BuyXGetYFree(2, 1));
    public static final Product NOTEBOOK =
            new Product("NOTEBOOK", "Notebook", 349, new Promotion.BulkDiscount(10, 10));
    /** A $100.00 item: handy for reaching free shipping and coupon minimums quickly. */
    public static final Product TV = new Product("TV", "Television", 10000);

    /** Five products, two of them with a promotion. */
    public static final Catalog TINY_CATALOG = new Catalog(List.of(APPLE, PEN, COFFEE, NOTEBOOK, TV));

    /** A valid receipt for one apple, used wherever "some receipt" is needed. */
    public static final Receipt ONE_APPLE_RECEIPT = new Receipt(
            List.of(new ReceiptLine("APPLE", "Apple", 1, 50, 50, 0)),
            50, 0, 0, 499, 4, 553);

    private Fixtures() {}

    /**
     * A basket line.
     *
     * @param sku      the product code
     * @param quantity how many (a long, like {@link BasketItem#quantity()})
     * @return the item
     */
    public static BasketItem item(String sku, long quantity) {
        return new BasketItem(sku, quantity);
    }

    /**
     * A basket with these items, no coupon and standard shipping.
     *
     * @param items the basket lines, in order
     * @return the basket
     */
    public static Basket basketOf(BasketItem... items) {
        return new Basket(List.of(items));
    }

    /**
     * A basket of one line.
     *
     * @param sku      the product code
     * @param quantity how many
     * @return the basket
     */
    public static Basket basketOf(String sku, long quantity) {
        return basketOf(item(sku, quantity));
    }
}
