package checkout;

/**
 * Why a checkout failed (rules section 9).
 *
 * <p>The names are spelled exactly like in the spec, so {@code code.name()} gives
 * the text used in the JSON test cases and in the command line output.
 */
public enum ErrorCode {
    /** The basket has no items. */
    EMPTY_BASKET,
    /** An item's SKU is not in the catalog. */
    UNKNOWN_PRODUCT,
    /** A quantity is not a whole number from 1 to 99. */
    INVALID_QUANTITY,
    /** Shipping is not {@code standard} or {@code express}. */
    UNKNOWN_SHIPPING_METHOD,
    /** The coupon code does not exist. */
    UNKNOWN_COUPON,
    /** The coupon needs a bigger spend. */
    COUPON_MINIMUM_NOT_MET
}
