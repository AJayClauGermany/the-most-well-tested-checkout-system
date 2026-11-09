package checkout.property;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

import checkout.Basket;
import checkout.BasketItem;
import checkout.Catalog;
import checkout.Checkout;
import checkout.CheckoutException;
import checkout.ErrorCode;
import checkout.Money;
import checkout.Receipt;
import checkout.ReceiptLine;
import checkout.orders.InMemoryOrderRepository;
import checkout.orders.Order;
import checkout.orders.Orders;
import checkout.orders.Timestamps;
import checkout.spec.Schemas;
import checkout.spec.SpecFiles;
import checkout.spec.Totals;
import java.time.Clock;
import java.time.Instant;
import java.time.ZoneOffset;
import java.util.ArrayList;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Locale;
import java.util.Optional;
import java.util.Set;
import net.jqwik.api.Arbitraries;
import net.jqwik.api.Arbitrary;
import net.jqwik.api.Combinators;
import net.jqwik.api.ForAll;
import net.jqwik.api.From;
import net.jqwik.api.Property;
import net.jqwik.api.Provide;
import net.jqwik.api.constraints.IntRange;
import net.jqwik.api.constraints.Size;

/**
 * PROPERTY-BASED tests with jqwik (https://jqwik.net/).
 *
 * <p>An example test says "THIS basket gives THAT total". A property says something that must
 * be true for EVERY basket, and jqwik generates hundreds of random baskets trying to break it -
 * then shrinks any failure to the smallest example. The properties are listed once for all six
 * languages in {@code spec/properties.md}; each method is named after its id (P1 .. P14).
 *
 * <p>jqwik is its own JUnit Platform engine: {@code @Property} replaces {@code @Test}, and
 * {@code @ForAll} parameters are filled with generated values ("arbitraries").
 */
class CheckoutPropertiesTest {

    private static final Catalog CATALOG = SpecFiles.loadCatalog();
    private static final List<String> SKUS = CATALOG.products().stream().map(p -> p.sku()).toList();
    private static final Clock FIXED_CLOCK = Clock.fixed(Instant.parse("2026-01-15T10:30:00Z"), ZoneOffset.UTC);

    // -----------------------------------------------------------------------------------------
    // Generators, following the "Generating baskets" section of spec/properties.md
    // -----------------------------------------------------------------------------------------

    /** A sensible item: a real SKU, a quantity from 1 to 99. */
    @Provide
    Arbitrary<BasketItem> items() {
        return Combinators.combine(Arbitraries.of(SKUS), Arbitraries.integers().between(1, 99))
                .as(BasketItem::new);
    }

    /** SAVE10 or WELCOME5 written any old way: mixed case, spaces around it. */
    @Provide
    Arbitrary<String> sloppyCoupons() {
        return Combinators.combine(
                        Arbitraries.of("SAVE10", "WELCOME5"),
                        Arbitraries.of(true, false).list().ofSize(8),
                        Arbitraries.integers().between(0, 3),
                        Arbitraries.integers().between(0, 3))
                .as(CheckoutPropertiesTest::sloppy);
    }

    private static String sloppy(String code, List<Boolean> upper, int before, int after) {
        StringBuilder letters = new StringBuilder();
        for (int i = 0; i < code.length(); i++) {
            String letter = code.substring(i, i + 1);
            letters.append(upper.get(i) ? letter.toUpperCase(Locale.ROOT) : letter.toLowerCase(Locale.ROOT));
        }
        return " ".repeat(before) + letters + " ".repeat(after);
    }

    /** A whole basket: 1 to 8 items (the same SKU may repeat), maybe a coupon, maybe a shipping method. */
    @Provide
    Arbitrary<Basket> baskets() {
        Arbitrary<String> coupon = Arbitraries.oneOf(
                Arbitraries.just(null),
                Arbitraries.of("SAVE10", "WELCOME5"),
                sloppyCoupons(),
                Arbitraries.strings().ofMaxLength(10));
        Arbitrary<String> shipping = Arbitraries.oneOf(
                Arbitraries.just(null),
                Arbitraries.of("standard", "express", "overnight", "Express", ""));
        return Combinators.combine(items().list().ofMinSize(1).ofMaxSize(8), coupon, shipping).as(Basket::new);
    }

    /** Only baskets that check out successfully - what P1 to P9 are about. */
    @Provide
    Arbitrary<Basket> validBaskets() {
        return baskets().filter(basket -> tryCheckout(basket).isPresent());
    }

    /** Anything at all that still fits the Basket type: nonsense SKUs, silly quantities, random strings. */
    @Provide
    Arbitrary<Basket> anyBaskets() {
        Arbitrary<BasketItem> anyItem = Combinators.combine(
                        Arbitraries.oneOf(
                                Arbitraries.of(SKUS), Arbitraries.strings().ofMaxLength(8), Arbitraries.just(null)),
                        Arbitraries.oneOf(
                                Arbitraries.integers().between(-5, 150),
                                Arbitraries.integers()))
                .as(BasketItem::new);
        Arbitrary<String> coupon = Arbitraries.oneOf(
                Arbitraries.just(null),
                Arbitraries.of("SAVE10", "WELCOME5", "save10", " welcome5 "),
                Arbitraries.strings().ofMaxLength(10));
        Arbitrary<String> shipping = Arbitraries.oneOf(
                Arbitraries.just(null),
                Arbitraries.of("standard", "express"),
                Arbitraries.strings().ofMaxLength(10));
        return Combinators.combine(anyItem.list().ofMaxSize(6), coupon, shipping).as(Basket::new);
    }

    /** Only baskets that fail validation - what P14 is about. */
    @Provide
    Arbitrary<Basket> invalidBaskets() {
        return anyBaskets().filter(basket -> tryCheckout(basket).isEmpty());
    }

    /** Runs checkout; empty means the basket was invalid (a CheckoutException was thrown). */
    private static Optional<Receipt> tryCheckout(Basket basket) {
        try {
            return Optional.of(Checkout.checkout(basket, CATALOG));
        } catch (CheckoutException e) {
            return Optional.empty();
        }
    }

    // -----------------------------------------------------------------------------------------
    // Properties of a successful checkout
    // -----------------------------------------------------------------------------------------

    @Property
    void p1_theTotalsAddUp(@ForAll("validBaskets") Basket basket) {
        Receipt r = Checkout.checkout(basket, CATALOG);

        assertThat(r.total())
                .isEqualTo(r.subtotal() - r.promotionDiscount() - r.couponDiscount() + r.shipping() + r.tax());
    }

    @Property
    void p2_nothingIsNegativeAndADiscountNeverExceedsWhatItDiscounts(@ForAll("validBaskets") Basket basket) {
        Receipt r = Checkout.checkout(basket, CATALOG);

        assertThat(List.of(r.subtotal(), r.promotionDiscount(), r.couponDiscount(), r.shipping(), r.tax(), r.total()))
                .allSatisfy(amount -> assertThat(amount).isNotNegative());
        for (ReceiptLine line : r.lines()) {
            assertThat(line.unitPrice()).isNotNegative();
            assertThat(line.lineTotal()).isNotNegative();
            assertThat(line.lineDiscount()).isBetween(0, line.lineTotal());
        }
        assertThat(r.promotionDiscount()).isLessThanOrEqualTo(r.subtotal());
        assertThat(r.couponDiscount()).isLessThanOrEqualTo(r.subtotal() - r.promotionDiscount());
    }

    @Property
    void p3_theReceiptLinesExplainTheTotals(@ForAll("validBaskets") Basket basket) {
        Receipt r = Checkout.checkout(basket, CATALOG);

        assertThat(r.subtotal()).isEqualTo(r.lines().stream().mapToInt(ReceiptLine::lineTotal).sum());
        assertThat(r.promotionDiscount()).isEqualTo(r.lines().stream().mapToInt(ReceiptLine::lineDiscount).sum());
        for (ReceiptLine line : r.lines()) {
            assertThat(line.lineTotal()).isEqualTo(line.unitPrice() * line.quantity());
        }
    }

    @Property
    void p4_shufflingTheItemsGivesTheSameTotals(
            @ForAll("validBaskets") Basket basket,
            @ForAll @Size(min = 16, max = 16) List<@IntRange(min = 0) Integer> randoms) {
        // A permutation of the items driven by jqwik's random numbers, so that it shrinks too.
        List<BasketItem> shuffled = new ArrayList<>(basket.items());
        for (int i = shuffled.size() - 1; i > 0; i--) {
            int j = randoms.get(i % randoms.size()) % (i + 1);
            BasketItem swapped = shuffled.get(i);
            shuffled.set(i, shuffled.get(j));
            shuffled.set(j, swapped);
        }
        Basket shuffledBasket = new Basket(shuffled, basket.coupon(), basket.shipping());

        assertThat(Totals.of(Checkout.checkout(shuffledBasket, CATALOG)))
                .isEqualTo(Totals.of(Checkout.checkout(basket, CATALOG)));
    }

    @Property
    void p5_splittingOneLineIntoTwoLinesOfTheSameSkuGivesTheSameReceipt(
            @ForAll("validBaskets") Basket basket,
            @ForAll @IntRange(min = 0) int pick,
            @ForAll @IntRange(min = 0) int cut) {
        List<Integer> splittable = new ArrayList<>();
        for (int i = 0; i < basket.items().size(); i++) {
            if (basket.items().get(i).quantity() >= 2) {
                splittable.add(i);
            }
        }
        if (splittable.isEmpty()) {
            return; // nothing to split: the property holds trivially
        }
        int index = splittable.get(pick % splittable.size());
        BasketItem line = basket.items().get(index);
        int first = 1 + cut % (line.quantity() - 1); // 1 .. quantity-1
        List<BasketItem> split = new ArrayList<>(basket.items());
        split.set(index, new BasketItem(line.sku(), first));
        split.add(new BasketItem(line.sku(), line.quantity() - first));
        Basket splitBasket = new Basket(split, basket.coupon(), basket.shipping());

        assertThat(Checkout.checkout(splitBasket, CATALOG)).isEqualTo(Checkout.checkout(basket, CATALOG));
    }

    @Property
    void p6_theShippingTableIsAppliedExactly(@ForAll("validBaskets") Basket basket) {
        Receipt r = Checkout.checkout(basket, CATALOG);
        int discounted = r.subtotal() - r.promotionDiscount() - r.couponDiscount();

        assertThat(r.shipping()).isIn(0, 499, 999);
        if ("express".equals(basket.shipping())) {
            assertThat(r.shipping()).isEqualTo(999);
        } else {
            assertThat(r.shipping()).isEqualTo(discounted >= 5000 ? 0 : 499);
        }
    }

    @Property
    void p7_taxIsEightPercentRoundedHalfUpWithIntegerDivision(@ForAll("validBaskets") Basket basket) {
        Receipt r = Checkout.checkout(basket, CATALOG);
        int discounted = r.subtotal() - r.promotionDiscount() - r.couponDiscount();

        assertThat(r.tax()).isEqualTo((discounted * 8 + 50) / 100);
    }

    @Property
    void p8_writingTheCouponInAnotherCaseOrWithSpacesGivesTheSameReceipt(
            @ForAll("validBaskets") Basket basket, @ForAll("sloppyCoupons") String sloppy) {
        if (basket.coupon() == null || !List.of("SAVE10", "WELCOME5").contains(basket.coupon())) {
            return; // only baskets written with a real coupon in its canonical form
        }
        // Make the sloppy version spell the SAME code as the basket's.
        String sameCode = sloppy.strip().equalsIgnoreCase(basket.coupon())
                ? sloppy
                : sloppy.replace(sloppy.strip(), basket.coupon().toLowerCase(Locale.ROOT));
        Basket sloppyBasket = new Basket(basket.items(), sameCode, basket.shipping());

        assertThat(Checkout.checkout(sloppyBasket, CATALOG)).isEqualTo(Checkout.checkout(basket, CATALOG));
    }

    @Property
    void p9_oneLinePerDistinctSkuInOrderOfFirstAppearanceEachWithQuantity1To99(
            @ForAll("validBaskets") Basket basket) {
        Receipt r = Checkout.checkout(basket, CATALOG);
        Set<String> distinctInOrder = new LinkedHashSet<>(basket.items().stream().map(BasketItem::sku).toList());

        assertThat(r.lines()).extracting(ReceiptLine::sku).containsExactlyElementsOf(distinctInOrder);
        for (ReceiptLine line : r.lines()) {
            assertThat(line.quantity()).isBetween(1, 99);
            int expectedQuantity = basket.items().stream()
                    .filter(item -> item.sku().equals(line.sku()))
                    .mapToInt(BasketItem::quantity)
                    .sum();
            assertThat(line.quantity()).isEqualTo(expectedQuantity);
        }
    }

    // -----------------------------------------------------------------------------------------
    // Properties of any input at all
    // -----------------------------------------------------------------------------------------

    @Property(tries = 500)
    void p10_forAnyBasketCheckoutGivesASchemaValidReceiptOrOneOfTheSixErrorCodesNeverACrash(
            @ForAll("anyBaskets") Basket basket) {
        try {
            Receipt receipt = Checkout.checkout(basket, CATALOG);
            assertThat(Schemas.problems("receipt.schema.json", receipt)).isEmpty();
        } catch (CheckoutException e) {
            assertThat(e.getCode()).isIn((Object[]) ErrorCode.values());
        }
    }

    @Property
    void p11_roundHalfUpIsWithinHalfACentOfTheExactValueAndExactWhenTheValueIsWhole(
            @ForAll @IntRange(min = 0, max = 10_000_000) int amount,
            @ForAll @IntRange(min = 0, max = 100) int percent) {
        int rounded = Money.roundHalfUp(amount, percent);
        long exactTimes100 = (long) amount * percent; // the exact value is this / 100

        assertThat(Math.abs(rounded * 100L - exactTimes100)).isLessThanOrEqualTo(50);
        if (exactTimes100 % 100 == 0) {
            assertThat(rounded * 100L).isEqualTo(exactTimes100);
        }
        if (exactTimes100 % 100 == 50) {
            assertThat(rounded * 100L).as("halves go UP").isEqualTo(exactTimes100 + 50);
        }
    }

    // -----------------------------------------------------------------------------------------
    // Properties of saving orders
    // -----------------------------------------------------------------------------------------

    @Property
    void p12_afterSavingNOrdersTheIdsAreExactlyOrder1ToOrderNInThatOrder(
            @ForAll @Size(max = 20) List<@From("validBaskets") Basket> baskets) {
        InMemoryOrderRepository repository = new InMemoryOrderRepository(FIXED_CLOCK);

        List<String> ids = baskets.stream().map(b -> Orders.placeOrder(b, CATALOG, repository).orderId()).toList();

        List<String> expected = new ArrayList<>();
        for (int i = 1; i <= baskets.size(); i++) {
            expected.add("ORDER-" + i);
        }
        assertThat(ids).isEqualTo(expected);
    }

    @Property
    void p13_findOfSaveGivesBackTheSameReceiptAndTheClocksTime(
            @ForAll("validBaskets") Basket basket,
            @ForAll @net.jqwik.api.constraints.LongRange(min = 946_684_800L, max = 4_102_444_799L) long epochSecond) {
        // Any real moment this century, to the second (946684800 is 2000-01-01, 4102444799 is 2099-12-31).
        Instant when = Instant.ofEpochSecond(epochSecond);
        InMemoryOrderRepository repository = new InMemoryOrderRepository(Clock.fixed(when, ZoneOffset.UTC));
        Receipt receipt = Checkout.checkout(basket, CATALOG);

        String id = repository.save(receipt);
        Optional<Order> order = repository.find(id);

        assertThat(order).isPresent();
        assertThat(order.get().receipt()).isEqualTo(receipt);
        assertThat(order.get().createdAt()).isEqualTo(Timestamps.format(when)).isEqualTo(when.toString());
    }

    @Property
    void p14_aFailedCheckoutSavesNothingAndFindOfAnUnknownIdIsNotFound(
            @ForAll("invalidBaskets") Basket basket, @ForAll String randomId) {
        InMemoryOrderRepository repository = new InMemoryOrderRepository(FIXED_CLOCK);

        assertThatThrownBy(() -> Orders.placeOrder(basket, CATALOG, repository)).isInstanceOf(CheckoutException.class);

        assertThat(repository.count()).isZero();
        assertThat(repository.find(randomId)).isEmpty();
        assertThat(repository.find("ORDER-1")).isEmpty();
    }
}
