package bdd

// Step definitions: the glue between the Given/When/Then sentences in
// spec/features/*.feature and the checkout code. Every sentence listed in
// spec/features/README.md appears exactly once here.
//
// godog matches each sentence with a regular expression; the parts in (brackets) are
// handed to the step function as arguments. A step returns nil when it passes and an
// error when it fails.

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/cucumber/godog"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/spec"
)

// world holds the state of ONE scenario as it goes from step to step (the basket being
// built, the receipt or error that came out, …). Every scenario gets a fresh one, so
// scenarios can never leak into each other.
type world struct {
	catalog checkout.Catalog
	basket  checkout.Basket
	receipt *checkout.Receipt
	err     *checkout.CheckoutError

	// now is what the (replaceable) clock answers. "Given the clock is fixed at …" changes it.
	now         time.Time
	repository  *checkout.InMemoryOrderRepository
	lastOrderID string
}

// InitializeScenario registers every step for one scenario, bound to a fresh world.
func InitializeScenario(sc *godog.ScenarioContext) {
	w := &world{now: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}

	// "$12.39" in a sentence: the dollar amount is captured as text and turned into cents
	// with whole-number maths (spec.DollarsToCents) – never via a float.
	const money = `\$(\d+(?:\.\d{1,2})?)`

	// --- Given --------------------------------------------------------------
	sc.Given(`^the standard catalog$`, w.theStandardCatalog)
	sc.Given(`^an empty basket$`, w.anEmptyBasket)
	sc.Given(`^a basket containing:$`, w.aBasketContaining)
	sc.Given(`^the customer uses the coupon "([^"]*)"$`, w.theCustomerUsesTheCoupon)
	sc.Given(`^the customer chooses "([^"]*)" shipping$`, w.theCustomerChoosesShipping)
	sc.Given(`^an empty order repository$`, w.anEmptyOrderRepository)
	sc.Given(`^the clock is fixed at "([^"]*)"$`, w.theClockIsFixedAt)

	// --- When ---------------------------------------------------------------
	sc.When(`^the customer checks out$`, w.theCustomerChecksOut)
	sc.When(`^the customer checks out and saves the order$`, w.theCustomerChecksOutAndSavesTheOrder)

	// --- Then: the receipt --------------------------------------------------
	sc.Then(`^the subtotal is `+money+`$`, w.amountIs("subtotal", func(r *checkout.Receipt) int { return r.Subtotal }))
	sc.Then(`^the promotion discount is `+money+`$`, w.amountIs("promotion discount", func(r *checkout.Receipt) int { return r.PromotionDiscount }))
	sc.Then(`^the coupon discount is `+money+`$`, w.amountIs("coupon discount", func(r *checkout.Receipt) int { return r.CouponDiscount }))
	sc.Then(`^the shipping cost is `+money+`$`, w.amountIs("shipping cost", func(r *checkout.Receipt) int { return r.Shipping }))
	sc.Then(`^the tax is `+money+`$`, w.amountIs("tax", func(r *checkout.Receipt) int { return r.Tax }))
	sc.Then(`^the total is `+money+`$`, w.amountIs("total", func(r *checkout.Receipt) int { return r.Total }))
	sc.Then(`^the receipt has (\d+) lines$`, w.theReceiptHasLines)
	sc.Then(`^the receipt line for "([^"]*)" has a quantity of (\d+)$`, w.theReceiptLineHasAQuantityOf)
	sc.Then(`^checkout fails with the error "([^"]*)"$`, w.checkoutFailsWithTheError)

	// --- Then: saved orders -------------------------------------------------
	sc.Then(`^the order is saved with the id "([^"]*)"$`, w.theOrderIsSavedWithTheID)
	sc.Then(`^the saved order "([^"]*)" has a total of `+money+`$`, w.theSavedOrderHasATotalOf)
	sc.Then(`^the order repository is empty$`, w.theOrderRepositoryIsEmpty)
	sc.Then(`^the order "([^"]*)" cannot be found$`, w.theOrderCannotBeFound)
	sc.Then(`^the saved order "([^"]*)" was created at "([^"]*)"$`, w.theSavedOrderWasCreatedAt)
}

// --- Given -----------------------------------------------------------------

func (w *world) theStandardCatalog() error {
	catalog, err := spec.LoadCatalog()
	w.catalog = catalog
	return err
}

func (w *world) anEmptyBasket() error {
	w.basket = checkout.Basket{Items: []checkout.BasketItem{}}
	return nil
}

func (w *world) aBasketContaining(table *godog.Table) error {
	// The first row is the header: | sku | quantity |
	for _, row := range table.Rows[1:] {
		quantity, err := strconv.Atoi(row.Cells[1].Value)
		if err != nil {
			return fmt.Errorf("quantity %q is not a whole number", row.Cells[1].Value)
		}
		w.basket.Items = append(w.basket.Items, checkout.BasketItem{SKU: row.Cells[0].Value, Quantity: quantity})
	}
	return nil
}

func (w *world) theCustomerUsesTheCoupon(code string) error {
	w.basket.Coupon = code
	return nil
}

func (w *world) theCustomerChoosesShipping(method string) error {
	w.basket.Shipping = method
	return nil
}

func (w *world) anEmptyOrderRepository() error {
	// The repository's clock reads w.now, so a later "the clock is fixed at" step still applies.
	w.repository = checkout.NewInMemoryOrderRepository(checkout.ClockFunc(func() time.Time { return w.now }))
	return nil
}

func (w *world) theClockIsFixedAt(timestamp string) error {
	now, err := time.Parse(time.RFC3339, timestamp)
	w.now = now
	return err
}

// --- When ------------------------------------------------------------------

// remember stores the receipt or the checkout error. Any other error is a real failure.
func (w *world) remember(receipt checkout.Receipt, err error) error {
	if err == nil {
		w.receipt = &receipt
		return nil
	}
	if !errors.As(err, &w.err) {
		return err
	}
	return nil
}

func (w *world) theCustomerChecksOut() error {
	return w.remember(checkout.Checkout(w.basket, w.catalog))
}

func (w *world) theCustomerChecksOutAndSavesTheOrder() error {
	if w.repository == nil {
		return errors.New("no order repository – add 'Given an empty order repository'")
	}
	placed, err := checkout.PlaceOrder(w.basket, w.catalog, w.repository)
	if err == nil {
		w.lastOrderID = placed.OrderID
	}
	return w.remember(placed.Receipt, err)
}

// --- Then: the receipt -------------------------------------------------------

// theReceipt returns the receipt, or an error explaining why there is none.
func (w *world) theReceipt() (*checkout.Receipt, error) {
	if w.err != nil {
		return nil, fmt.Errorf("checkout failed with %s", w.err.Code)
	}
	if w.receipt == nil {
		return nil, errors.New("the customer has not checked out yet")
	}
	return w.receipt, nil
}

// amountIs builds the step "the <name> is $<amount>" for one field of the receipt. All six
// amount steps work the same way, so one function makes all of them.
func (w *world) amountIs(name string, field func(*checkout.Receipt) int) func(string) error {
	return func(dollars string) error {
		want, err := spec.DollarsToCents(dollars)
		if err != nil {
			return err
		}
		receipt, err := w.theReceipt()
		if err != nil {
			return err
		}
		if got := field(receipt); got != want {
			return fmt.Errorf("the %s is %d cents, want %d", name, got, want)
		}
		return nil
	}
}

func (w *world) theReceiptHasLines(count int) error {
	receipt, err := w.theReceipt()
	if err != nil {
		return err
	}
	if len(receipt.Lines) != count {
		return fmt.Errorf("the receipt has %d lines, want %d", len(receipt.Lines), count)
	}
	return nil
}

func (w *world) theReceiptLineHasAQuantityOf(sku string, quantity int) error {
	receipt, err := w.theReceipt()
	if err != nil {
		return err
	}
	for _, line := range receipt.Lines {
		if line.SKU == sku {
			if line.Quantity != quantity {
				return fmt.Errorf("the line for %s has quantity %d, want %d", sku, line.Quantity, quantity)
			}
			return nil
		}
	}
	return fmt.Errorf("no receipt line for %s", sku)
}

func (w *world) checkoutFailsWithTheError(code string) error {
	if w.receipt != nil {
		return errors.New("checkout succeeded but was expected to fail")
	}
	if w.err == nil || string(w.err.Code) != code {
		return fmt.Errorf("checkout failed with %v, want %s", w.err, code)
	}
	return nil
}

// --- Then: saved orders ------------------------------------------------------

// theOrder looks a saved order up, or explains why it cannot.
func (w *world) theOrder(id string) (*checkout.Order, error) {
	if w.repository == nil {
		return nil, errors.New("no order repository – add 'Given an empty order repository'")
	}
	order, err := w.repository.Find(id)
	if err == nil && order == nil {
		err = fmt.Errorf("order %s not found", id)
	}
	return order, err
}

func (w *world) theOrderIsSavedWithTheID(id string) error {
	if w.lastOrderID != id {
		return fmt.Errorf("the last order was saved as %q, want %q", w.lastOrderID, id)
	}
	return nil
}

func (w *world) theSavedOrderHasATotalOf(id, dollars string) error {
	want, err := spec.DollarsToCents(dollars)
	if err != nil {
		return err
	}
	order, err := w.theOrder(id)
	if err != nil {
		return err
	}
	if order.Receipt.Total != want {
		return fmt.Errorf("order %s has a total of %d cents, want %d", id, order.Receipt.Total, want)
	}
	return nil
}

func (w *world) theOrderRepositoryIsEmpty() error {
	if w.repository == nil {
		return errors.New("no order repository – add 'Given an empty order repository'")
	}
	if n := w.repository.Count(); n != 0 {
		return fmt.Errorf("the repository holds %d orders, want none", n)
	}
	return nil
}

func (w *world) theOrderCannotBeFound(id string) error {
	if w.repository == nil {
		return errors.New("no order repository – add 'Given an empty order repository'")
	}
	order, err := w.repository.Find(id)
	if err != nil {
		return err
	}
	if order != nil {
		return fmt.Errorf("order %s was found, want not found", id)
	}
	return nil
}

func (w *world) theSavedOrderWasCreatedAt(id, timestamp string) error {
	order, err := w.theOrder(id)
	if err != nil {
		return err
	}
	if order.CreatedAt != timestamp {
		return fmt.Errorf("order %s was created at %q, want %q", id, order.CreatedAt, timestamp)
	}
	return nil
}
