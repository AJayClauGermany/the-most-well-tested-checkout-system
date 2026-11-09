package checkout.orders;

import checkout.Receipt;

/**
 * What {@link Orders#placeOrder} gives back: the receipt and the id it was saved under.
 *
 * @param receipt the receipt of the successful checkout
 * @param orderId the id the repository gave the order
 */
public record PlacedOrder(Receipt receipt, String orderId) {}
