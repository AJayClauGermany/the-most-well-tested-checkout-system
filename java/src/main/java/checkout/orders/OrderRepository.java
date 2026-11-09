package checkout.orders;

import checkout.Receipt;
import java.util.Optional;

/**
 * Where orders are kept (rules section 6).
 *
 * <p>An interface with two methods is all the checkout needs to know about storage. That
 * is what makes it easy to test: the real program can use a folder of JSON files or a
 * database, while a unit test passes a tiny in-memory version or a Mockito mock.
 * Implemented by {@link InMemoryOrderRepository}, {@code checkout.store.JsonFileOrderRepository}
 * and {@code checkout.store.JdbcOrderRepository}.
 */
public interface OrderRepository {

    /**
     * Saves a receipt as a new order.
     *
     * @param receipt the receipt of a successful checkout
     * @return the new order's id: {@code "ORDER-1"} for the first order saved here, then {@code "ORDER-2"}, ...
     */
    String save(Receipt receipt);

    /**
     * Looks an order up by its id.
     *
     * @param id an order id
     * @return the order, or empty when no order with this id was ever saved here
     *         ("not found" is a normal answer, not an error)
     */
    Optional<Order> find(String id);
}
