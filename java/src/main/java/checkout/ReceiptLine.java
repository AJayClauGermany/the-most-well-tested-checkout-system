package checkout;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * One line of the receipt, after duplicate basket lines were merged (rules section 5).
 *
 * @param sku          the product code
 * @param name         the product name
 * @param quantity     how many items (the sum of all basket lines with this SKU)
 * @param unitPrice    price of one item, in cents
 * @param lineTotal    {@code unitPrice * quantity}, in cents
 * @param lineDiscount the promotion discount for this line, in cents ({@code 0} if none)
 */
public record ReceiptLine(
        String sku,
        String name,
        int quantity,
        @JsonProperty("unit_price") int unitPrice,
        @JsonProperty("line_total") int lineTotal,
        @JsonProperty("line_discount") int lineDiscount) {}
