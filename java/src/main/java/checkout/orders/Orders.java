package checkout.orders;

import checkout.Basket;
import checkout.Catalog;
import checkout.CheckoutException;

/**
 * Checkout plus saving (rules section 6).
 */
public final class Orders {

    private Orders() {}

    /**
     * Checks out, then saves the receipt as an order. A failed checkout saves nothing:
     * the {@link CheckoutException} is thrown before the repository is touched.
     *
     * @param basket     what the shopper wants to buy
     * @param catalog    the products and their prices
     * @param repository where to save the order
     * @return the receipt and the new order's id
     * @throws CheckoutException when the basket is invalid - and then nothing was saved
     */
    public static PlacedOrder placeOrder(Basket basket, Catalog catalog, OrderRepository repository) {
        throw new UnsupportedOperationException("not implemented yet");
    }
}
