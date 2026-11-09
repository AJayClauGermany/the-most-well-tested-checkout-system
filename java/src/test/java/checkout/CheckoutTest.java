package checkout;

import static checkout.spec.Fixtures.TINY_CATALOG;
import static checkout.spec.Fixtures.basketOf;
import static checkout.spec.Fixtures.item;
import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.assertj.core.groups.Tuple.tuple;

import java.util.List;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;

/**
 * UNIT tests of {@link Checkout#checkout} as a whole, with the tiny made-up catalog from
 * {@code checkout.spec.Fixtures}. The 60 shared spec cases ({@link SharedSpecCasesTest})
 * cover the standard catalog; these tests check that the pieces fit together and the things
 * the six totals alone cannot show, such as the receipt lines.
 */
class CheckoutTest {

    private static void assertFailsWith(Basket basket, ErrorCode code) {
        assertThatThrownBy(() -> Checkout.checkout(basket, TINY_CATALOG))
                .isInstanceOf(CheckoutException.class)
                .extracting("code")
                .isEqualTo(code);
    }

    @Test
    void pricesOneApple() {
        // 50 + $4.99 shipping + 4 tax = $5.53
        Receipt receipt = Checkout.checkout(basketOf("APPLE", 1), TINY_CATALOG);

        assertThat(receipt).isEqualTo(new Receipt(
                List.of(new ReceiptLine("APPLE", "Apple", 1, 50, 50, 0)),
                50, 0, 0, 499, 4, 553));
    }

    @Test
    void writesOneReceiptLinePerDistinctSkuInOrderOfFirstAppearanceWithMergedQuantities() {
        Receipt receipt = Checkout.checkout(
                basketOf(item("PEN", 1), item("APPLE", 2), item("PEN", 3)), TINY_CATALOG);

        assertThat(receipt.lines())
                .extracting(ReceiptLine::sku, ReceiptLine::quantity)
                .containsExactly(tuple("PEN", 4), tuple("APPLE", 2));
    }

    @Test
    void worksTheTotalsOutFromTheLines() {
        // 3 coffees ($26.97, one free: -$8.99) and 10 notebooks ($34.90, 10% off: -$3.49)
        Receipt receipt = Checkout.checkout(basketOf(item("COFFEE", 3), item("NOTEBOOK", 10)), TINY_CATALOG);

        assertThat(receipt.subtotal()).isEqualTo(2697 + 3490);
        assertThat(receipt.promotionDiscount()).isEqualTo(899 + 349);
        assertThat(receipt.couponDiscount()).isZero();
        assertThat(receipt.shipping()).as("$49.39 after promotions: not free").isEqualTo(499);
        assertThat(receipt.tax()).as("8% of 4939 = 395.12").isEqualTo(395);
        assertThat(receipt.total()).isEqualTo(4939 + 499 + 395);
    }

    @Test
    void takesTheCouponOffAfterPromotionsAndBeforeShippingAndTax() {
        // TV $100.00 with SAVE10: $90.00 -> free shipping, tax 8% of $90.00 = $7.20
        Basket basket = new Basket(List.of(item("TV", 1)), "SAVE10", null);

        Receipt receipt = Checkout.checkout(basket, TINY_CATALOG);

        assertThat(receipt)
                .returns(10000, Receipt::subtotal)
                .returns(0, Receipt::promotionDiscount)
                .returns(1000, Receipt::couponDiscount)
                .returns(0, Receipt::shipping)
                .returns(720, Receipt::tax)
                .returns(9720, Receipt::total);
    }

    @Test
    void chargesExpressShippingOnTopUntaxed() {
        Basket basket = new Basket(List.of(item("TV", 1)), null, "express");

        Receipt receipt = Checkout.checkout(basket, TINY_CATALOG);

        assertThat(receipt)
                .returns(999, Receipt::shipping)
                .returns(800, Receipt::tax)
                .returns(10000 + 999 + 800, Receipt::total);
    }

    @Test
    void doesNotChangeTheBasketItIsGiven() {
        Basket basket = new Basket(List.of(item("APPLE", 1), item("APPLE", 1)), " save10 ", null);

        Checkout.checkout(basket, TINY_CATALOG);

        assertThat(basket).isEqualTo(new Basket(List.of(item("APPLE", 1), item("APPLE", 1)), " save10 ", "standard"));
    }

    /** Rules section 4.1: when several rows apply, the first one in the table wins. */
    @Nested
    class ReportsTheFirstFailingRow {

        @Test
        void emptyBasketBeforeAnythingElse() {
            assertFailsWith(new Basket(List.of(), "NOPE", "overnight"), ErrorCode.EMPTY_BASKET);
        }

        @Test
        void unknownProductBeforeBadQuantityShippingAndCoupon() {
            assertFailsWith(
                    new Basket(List.of(item("APPLE", 0), item("BANANA", 1)), "NOPE", "x"),
                    ErrorCode.UNKNOWN_PRODUCT);
        }

        @Test
        void badQuantityBeforeShippingAndCoupon() {
            assertFailsWith(new Basket(List.of(item("APPLE", 0)), "NOPE", "x"), ErrorCode.INVALID_QUANTITY);
        }

        @Test
        void mergedQuantityOverNinetyNineBeforeShippingAndCoupon() {
            assertFailsWith(
                    new Basket(List.of(item("APPLE", 60), item("APPLE", 40)), "NOPE", "x"),
                    ErrorCode.INVALID_QUANTITY);
        }

        @Test
        void badShippingBeforeBadCoupon() {
            assertFailsWith(new Basket(List.of(item("APPLE", 1)), "NOPE", "x"), ErrorCode.UNKNOWN_SHIPPING_METHOD);
        }

        @Test
        void badShippingBeforeAnUnmetCouponMinimum() {
            assertFailsWith(
                    new Basket(List.of(item("APPLE", 1)), "WELCOME5", "x"), ErrorCode.UNKNOWN_SHIPPING_METHOD);
        }

        @Test
        void unknownCoupon() {
            assertFailsWith(new Basket(List.of(item("APPLE", 1)), "NOPE", null), ErrorCode.UNKNOWN_COUPON);
        }

        @Test
        void couponMinimumNotMetJudgedOnThePriceAfterPromotions() {
            // 3 coffees are $26.97 but $17.98 after the free one: below WELCOME5's $25.00.
            assertFailsWith(
                    new Basket(List.of(item("COFFEE", 3)), "WELCOME5", null), ErrorCode.COUPON_MINIMUM_NOT_MET);
        }
    }
}
