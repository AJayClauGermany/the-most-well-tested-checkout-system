package checkout.bdd;

import static org.assertj.core.api.Assertions.assertThat;

import checkout.Basket;
import checkout.BasketItem;
import checkout.Catalog;
import checkout.Checkout;
import checkout.CheckoutException;
import checkout.Receipt;
import checkout.ReceiptLine;
import checkout.orders.InMemoryOrderRepository;
import checkout.orders.Order;
import checkout.orders.Orders;
import checkout.orders.PlacedOrder;
import checkout.spec.Dollars;
import checkout.spec.SpecFiles;
import io.cucumber.datatable.DataTable;
import io.cucumber.java.ParameterType;
import io.cucumber.java.en.Given;
import io.cucumber.java.en.Then;
import io.cucumber.java.en.When;
import java.time.Clock;
import java.time.Instant;
import java.time.ZoneId;
import java.time.ZoneOffset;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Optional;

/**
 * Step definitions: the glue between the Given/When/Then sentences in {@code spec/features/*.feature}
 * and the checkout code. Every sentence used by the features is listed in
 * {@code spec/features/README.md}; each one appears exactly once here.
 *
 * <p>Cucumber creates a new instance of this class for every scenario, so the fields below are
 * the scenario's state: the basket being built, then the receipt or the error that came out.
 */
public class CheckoutSteps {

    private Catalog catalog;
    private final List<BasketItem> items = new ArrayList<>();
    private String coupon;
    private String shipping;

    private Receipt receipt;
    private CheckoutException error;

    /** What the (replaceable) clock answers. "Given the clock is fixed at ..." changes it. */
    private Instant now = Instant.parse("2026-01-01T00:00:00Z");
    private InMemoryOrderRepository repository;
    private String lastOrderId;

    // A custom parameter type: "{money}" in a step matches "$12.39" and hands the step 1239 cents.
    @ParameterType("\\$\\d+(?:\\.\\d{1,2})?")
    public int money(String text) {
        return Dollars.toCents(text);
    }

    // --- Given --------------------------------------------------------------------------------

    @Given("the standard catalog")
    public void theStandardCatalog() {
        catalog = SpecFiles.loadCatalog();
    }

    @Given("an empty basket")
    public void anEmptyBasket() {
        items.clear();
    }

    @Given("a basket containing:")
    public void aBasketContaining(DataTable table) {
        for (Map<String, String> row : table.asMaps()) {
            items.add(new BasketItem(row.get("sku"), Integer.parseInt(row.get("quantity"))));
        }
    }

    @Given("the customer uses the coupon {string}")
    public void theCustomerUsesTheCoupon(String code) {
        coupon = code;
    }

    @Given("the customer chooses {string} shipping")
    public void theCustomerChoosesShipping(String method) {
        shipping = method;
    }

    @Given("an empty order repository")
    public void anEmptyOrderRepository() {
        // The repository's clock reads `now`, so a later "the clock is fixed at" step still applies.
        repository = new InMemoryOrderRepository(new ScenarioClock());
    }

    @Given("the clock is fixed at {string}")
    public void theClockIsFixedAt(String timestamp) {
        now = Instant.parse(timestamp);
    }

    // --- When ---------------------------------------------------------------------------------

    @When("the customer checks out")
    public void theCustomerChecksOut() {
        try {
            receipt = Checkout.checkout(basket(), catalog);
        } catch (CheckoutException e) {
            error = e;
        }
    }

    @When("the customer checks out and saves the order")
    public void theCustomerChecksOutAndSavesTheOrder() {
        try {
            PlacedOrder placed = Orders.placeOrder(basket(), catalog, repositoryOrFail());
            receipt = placed.receipt();
            lastOrderId = placed.orderId();
        } catch (CheckoutException e) {
            error = e;
        }
    }

    // --- Then: the receipt --------------------------------------------------------------------

    @Then("the subtotal is {money}")
    public void theSubtotalIs(int cents) {
        assertThat(receiptOrFail().subtotal()).isEqualTo(cents);
    }

    @Then("the promotion discount is {money}")
    public void thePromotionDiscountIs(int cents) {
        assertThat(receiptOrFail().promotionDiscount()).isEqualTo(cents);
    }

    @Then("the coupon discount is {money}")
    public void theCouponDiscountIs(int cents) {
        assertThat(receiptOrFail().couponDiscount()).isEqualTo(cents);
    }

    @Then("the shipping cost is {money}")
    public void theShippingCostIs(int cents) {
        assertThat(receiptOrFail().shipping()).isEqualTo(cents);
    }

    @Then("the tax is {money}")
    public void theTaxIs(int cents) {
        assertThat(receiptOrFail().tax()).isEqualTo(cents);
    }

    @Then("the total is {money}")
    public void theTotalIs(int cents) {
        assertThat(receiptOrFail().total()).isEqualTo(cents);
    }

    @Then("the receipt has {int} lines")
    public void theReceiptHasLines(int count) {
        assertThat(receiptOrFail().lines()).hasSize(count);
    }

    @Then("the receipt line for {string} has a quantity of {int}")
    public void theReceiptLineForHasAQuantityOf(String sku, int quantity) {
        Optional<ReceiptLine> line = receiptOrFail().lines().stream().filter(l -> l.sku().equals(sku)).findFirst();

        assertThat(line).as("receipt line for %s", sku).isPresent();
        assertThat(line.get().quantity()).isEqualTo(quantity);
    }

    @Then("checkout fails with the error {string}")
    public void checkoutFailsWithTheError(String code) {
        assertThat(receipt).as("checkout succeeded but was expected to fail").isNull();
        assertThat(error).as("checkout has not run or did not fail").isNotNull();
        assertThat(error.getCode().name()).isEqualTo(code);
    }

    // --- Then: saved orders -------------------------------------------------------------------

    @Then("the order is saved with the id {string}")
    public void theOrderIsSavedWithTheId(String id) {
        assertThat(lastOrderId).isEqualTo(id);
    }

    @Then("the saved order {string} has a total of {money}")
    public void theSavedOrderHasATotalOf(String id, int cents) {
        Optional<Order> order = repositoryOrFail().find(id);

        assertThat(order).as("order %s", id).isPresent();
        assertThat(order.get().receipt().total()).isEqualTo(cents);
    }

    @Then("the order repository is empty")
    public void theOrderRepositoryIsEmpty() {
        assertThat(repositoryOrFail().count()).isZero();
    }

    @Then("the order {string} cannot be found")
    public void theOrderCannotBeFound(String id) {
        assertThat(repositoryOrFail().find(id)).isEmpty();
    }

    @Then("the saved order {string} was created at {string}")
    public void theSavedOrderWasCreatedAt(String id, String timestamp) {
        Optional<Order> order = repositoryOrFail().find(id);

        assertThat(order).as("order %s", id).isPresent();
        assertThat(order.get().createdAt()).isEqualTo(timestamp);
    }

    // --- helpers ------------------------------------------------------------------------------

    private Basket basket() {
        return new Basket(List.copyOf(items), coupon, shipping);
    }

    private Receipt receiptOrFail() {
        assertThat(error).as("checkout failed with %s", error == null ? "" : error.getCode()).isNull();
        assertThat(receipt).as("the customer has not checked out yet").isNotNull();
        return receipt;
    }

    private InMemoryOrderRepository repositoryOrFail() {
        assertThat(repository).as("no order repository - add 'Given an empty order repository'").isNotNull();
        return repository;
    }

    /** A clock that answers whatever `now` currently is - a stub the scenario can adjust. */
    private final class ScenarioClock extends Clock {
        @Override
        public ZoneId getZone() {
            return ZoneOffset.UTC;
        }

        @Override
        public Clock withZone(ZoneId zone) {
            return this;
        }

        @Override
        public Instant instant() {
            return now;
        }
    }
}
