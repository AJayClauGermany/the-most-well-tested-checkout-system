package checkout;

import static org.assertj.core.api.Assertions.assertThat;

import org.junit.jupiter.api.Test;

/**
 * UNIT tests for the error codes and the {@link CheckoutException} class.
 */
class CheckoutExceptionTest {

    @Test
    void errorCodesAreExactlyTheSixFromRulesSection9SpelledTheSameWayInOrder() {
        // ErrorCode.name() is what ends up in the JSON output, so the spelling matters.
        assertThat(ErrorCode.values())
                .extracting(ErrorCode::name)
                .containsExactly(
                        "EMPTY_BASKET",
                        "UNKNOWN_PRODUCT",
                        "INVALID_QUANTITY",
                        "UNKNOWN_SHIPPING_METHOD",
                        "UNKNOWN_COUPON",
                        "COUPON_MINIMUM_NOT_MET");
    }

    @Test
    void carriesTheErrorCode() {
        CheckoutException error = new CheckoutException(ErrorCode.UNKNOWN_COUPON);

        assertThat(error.getCode()).isEqualTo(ErrorCode.UNKNOWN_COUPON);
    }

    @Test
    void usesTheCodeAsTheMessageWhenNoMessageIsGiven() {
        assertThat(new CheckoutException(ErrorCode.UNKNOWN_COUPON)).hasMessage("UNKNOWN_COUPON");
    }

    @Test
    void keepsTheCodeAndTheMessageApart() {
        CheckoutException error = new CheckoutException(ErrorCode.EMPTY_BASKET, "The basket has no items");

        assertThat(error.getCode()).isEqualTo(ErrorCode.EMPTY_BASKET);
        assertThat(error).hasMessage("The basket has no items");
    }

    @Test
    void isAnUncheckedExceptionSoCallersAreNotForcedToCatchIt() {
        assertThat(new CheckoutException(ErrorCode.EMPTY_BASKET)).isInstanceOf(RuntimeException.class);
    }
}
