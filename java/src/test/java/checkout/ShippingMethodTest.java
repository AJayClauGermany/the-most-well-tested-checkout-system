package checkout;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;

/**
 * UNIT tests for shipping (rules section 4.5). The free-shipping threshold is tested on both
 * sides (4999 / 5000).
 */
class ShippingMethodTest {

    @Nested
    class Parse {

        @Test
        void defaultsToStandardWhenTheBasketSaysNothing() {
            assertThat(ShippingMethod.parse(null)).isEqualTo(ShippingMethod.STANDARD);
        }

        @Test
        void acceptsTheTwoKnownMethods() {
            assertThat(ShippingMethod.parse("standard")).isEqualTo(ShippingMethod.STANDARD);
            assertThat(ShippingMethod.parse("express")).isEqualTo(ShippingMethod.EXPRESS);
        }

        @ParameterizedTest(name = "rejects ''{0}''")
        @ValueSource(strings = {"overnight", "Express", "STANDARD", " standard", "express ", ""})
        void rejectsEverythingElseExactMatchOnly(String method) {
            assertThatThrownBy(() -> ShippingMethod.parse(method))
                    .isInstanceOf(CheckoutException.class)
                    .extracting("code")
                    .isEqualTo(ErrorCode.UNKNOWN_SHIPPING_METHOD);
        }

        @Test
        void namesTheMethodItDidNotRecogniseInTheErrorMessage() {
            assertThatThrownBy(() -> ShippingMethod.parse("overnight"))
                    .hasMessage("Unknown shipping method: overnight");
        }
    }

    @Nested
    class Cost {

        @Test
        void standardCostsFourNinetyNineBelowFiftyDollars() {
            assertThat(ShippingMethod.STANDARD.cost(0)).isEqualTo(499);
            assertThat(ShippingMethod.STANDARD.cost(4999)).isEqualTo(499);
        }

        @Test
        void standardIsFreeFromExactlyFiftyDollars() {
            assertThat(ShippingMethod.STANDARD.cost(5000)).isZero();
            assertThat(ShippingMethod.STANDARD.cost(5001)).isZero();
            assertThat(ShippingMethod.STANDARD.cost(100000)).isZero();
        }

        @Test
        void expressAlwaysCostsNineNinetyNine() {
            assertThat(ShippingMethod.EXPRESS.cost(0)).isEqualTo(999);
            assertThat(ShippingMethod.EXPRESS.cost(4999)).isEqualTo(999);
            assertThat(ShippingMethod.EXPRESS.cost(5000)).isEqualTo(999);
            assertThat(ShippingMethod.EXPRESS.cost(100000)).isEqualTo(999);
        }

        @Test
        void theBaseCostsAndThresholdAreTheOnesFromTheRules() {
            assertThat(ShippingMethod.STANDARD.baseCost()).isEqualTo(499);
            assertThat(ShippingMethod.EXPRESS.baseCost()).isEqualTo(999);
            assertThat(ShippingMethod.FREE_STANDARD_SHIPPING_FROM).isEqualTo(5000);
        }
    }
}
