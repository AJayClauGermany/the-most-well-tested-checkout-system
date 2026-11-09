package checkout;

import static checkout.spec.Fixtures.APPLE;
import static checkout.spec.Fixtures.COFFEE;
import static checkout.spec.Fixtures.NOTEBOOK;
import static checkout.spec.Fixtures.item;
import static org.assertj.core.api.Assertions.assertThat;

import java.util.ArrayList;
import java.util.List;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;

/**
 * UNIT tests for merging basket items into receipt lines and pricing them (rules section 4.2).
 */
class LinesTest {

    @Nested
    class MergeItems {

        @Test
        void keepsItemsWithDifferentSkusApartInBasketOrder() {
            List<BasketItem> merged = Lines.mergeItems(List.of(item("PEN", 2), item("APPLE", 1)));

            assertThat(merged).containsExactly(item("PEN", 2), item("APPLE", 1));
        }

        @Test
        void addsUpTheQuantitiesOfTheSameSkuKeepingThePlaceOfItsFirstAppearance() {
            // The example from rules section 4.2: [APPLE x1, PEN x2, APPLE x3] gives [APPLE x4, PEN x2].
            List<BasketItem> merged = Lines.mergeItems(List.of(item("APPLE", 1), item("PEN", 2), item("APPLE", 3)));

            assertThat(merged).containsExactly(item("APPLE", 4), item("PEN", 2));
        }

        @Test
        void mergesThreeLinesOfOneSkuIntoOne() {
            List<BasketItem> merged =
                    Lines.mergeItems(List.of(item("COFFEE", 1), item("COFFEE", 1), item("COFFEE", 1)));

            assertThat(merged).containsExactly(item("COFFEE", 3));
        }

        @Test
        void givesAnEmptyListForAnEmptyList() {
            assertThat(Lines.mergeItems(List.of())).isEmpty();
        }

        @Test
        void doesNotChangeTheBasketsOwnItems() {
            List<BasketItem> items = new ArrayList<>(List.of(item("APPLE", 1), item("APPLE", 1)));

            Lines.mergeItems(items);

            assertThat(items).containsExactly(item("APPLE", 1), item("APPLE", 1));
        }
    }

    @Nested
    class PriceLine {

        @Test
        void copiesSkuNameAndUnitPriceFromTheProductAndMultipliesForTheLineTotal() {
            ReceiptLine line = Lines.priceLine(APPLE, 3);

            assertThat(line).isEqualTo(new ReceiptLine("APPLE", "Apple", 3, 50, 150, 0));
        }

        @Test
        void includesThePromotionDiscount() {
            assertThat(Lines.priceLine(COFFEE, 3))
                    .returns(2697, ReceiptLine::lineTotal)
                    .returns(899, ReceiptLine::lineDiscount);
            assertThat(Lines.priceLine(NOTEBOOK, 10))
                    .returns(3490, ReceiptLine::lineTotal)
                    .returns(349, ReceiptLine::lineDiscount);
        }

        @Test
        void givesNoDiscountBelowThePromotionsThreshold() {
            assertThat(Lines.priceLine(COFFEE, 2).lineDiscount()).isZero();
            assertThat(Lines.priceLine(NOTEBOOK, 9).lineDiscount()).isZero();
        }
    }
}
