package checkout_test

// PROPERTY-BASED TESTS with rapid (https://pkg.go.dev/pgregory.net/rapid).
//
// An example test says "THIS basket gives THAT total". A property says something that must
// be true for EVERY basket, and rapid generates hundreds of random baskets trying to break
// it – then shrinks any failure to the smallest example. The properties are listed once for
// all six languages in spec/properties.md; each test is named after its id (P1 … P14).
//
// Run more cases than the default 100 with:  go test ./checkout -run TestP -rapid.checks=1000

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/google/go-cmp/cmp"
	"pgregory.net/rapid"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/fixtures"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/spec"
)

// ---------------------------------------------------------------------------
// Generators, following the "Generating baskets" section of spec/properties.md
// ---------------------------------------------------------------------------

// baskets knows how to generate baskets for one catalog.
type baskets struct {
	catalog checkout.Catalog
	skus    []string
}

func newBaskets(t *testing.T) baskets {
	t.Helper()
	catalog := loadCatalog(t)
	var skus []string
	for _, product := range catalog.Products {
		skus = append(skus, product.SKU)
	}
	return baskets{catalog: catalog, skus: skus}
}

// item generates a sensible item: a real SKU, a quantity from 1 to 99.
func (g baskets) item() *rapid.Generator[checkout.BasketItem] {
	return rapid.Custom(func(t *rapid.T) checkout.BasketItem {
		return checkout.BasketItem{
			SKU:      rapid.SampledFrom(g.skus).Draw(t, "sku"),
			Quantity: rapid.IntRange(1, 99).Draw(t, "quantity"),
		}
	})
}

// sloppy writes a coupon code any old way: mixed case, spaces around it.
func sloppy(t *rapid.T, code string) string {
	letters := []rune(code)
	for i, letter := range letters {
		if rapid.Bool().Draw(t, "upper") {
			letters[i] = unicode.ToUpper(letter)
		} else {
			letters[i] = unicode.ToLower(letter)
		}
	}
	before := strings.Repeat(" ", rapid.IntRange(0, 3).Draw(t, "spacesBefore"))
	after := strings.Repeat(" ", rapid.IntRange(0, 3).Draw(t, "spacesAfter"))
	return before + string(letters) + after
}

var sloppyCoupon = rapid.Custom(func(t *rapid.T) string {
	return sloppy(t, rapid.SampledFrom([]string{"SAVE10", "WELCOME5"}).Draw(t, "code"))
})

// coupon generates a coupon field: missing (""), a real code, a sloppy one, or nonsense.
// (In Go a missing coupon and a JSON null both arrive as "", so one value covers both.)
var coupon = rapid.OneOf(
	rapid.Just(""), rapid.Just("SAVE10"), rapid.Just("WELCOME5"), sloppyCoupon, rapid.String(),
)

// shipping generates a shipping field: missing, a real method, or a wrong one.
var shipping = rapid.SampledFrom([]string{"", "standard", "express", "overnight", "Express"})

// basket generates a whole basket: 1 to 8 sensible items (the same SKU may repeat), maybe a
// coupon, maybe a shipping method. Not every basket is valid – that is on purpose.
func (g baskets) basket() *rapid.Generator[checkout.Basket] {
	return rapid.Custom(func(t *rapid.T) checkout.Basket {
		return checkout.Basket{
			Items:    rapid.SliceOfN(g.item(), 1, 8).Draw(t, "items"),
			Coupon:   coupon.Draw(t, "coupon"),
			Shipping: shipping.Draw(t, "shipping"),
		}
	})
}

// isValid reports whether a basket checks out successfully.
func (g baskets) isValid(b checkout.Basket) bool {
	_, err := checkout.Checkout(b, g.catalog)
	return err == nil
}

// valid generates only baskets that check out successfully – what P1 to P9 are about.
func (g baskets) valid() *rapid.Generator[checkout.Basket] {
	return g.basket().Filter(g.isValid)
}

// any generates anything at all that still fits the Basket type: nonsense SKUs, silly
// quantities, random strings everywhere. P10 and P14 are about these.
func (g baskets) any() *rapid.Generator[checkout.Basket] {
	anyItem := rapid.Custom(func(t *rapid.T) checkout.BasketItem {
		return checkout.BasketItem{
			SKU:      rapid.OneOf(rapid.SampledFrom(g.skus), rapid.String()).Draw(t, "sku"),
			Quantity: rapid.IntRange(-5, 150).Draw(t, "quantity"),
		}
	})
	return rapid.Custom(func(t *rapid.T) checkout.Basket {
		return checkout.Basket{
			Items:    rapid.SliceOfN(anyItem, 0, 6).Draw(t, "items"),
			Coupon:   rapid.OneOf(rapid.SampledFrom([]string{"SAVE10", "WELCOME5", "save10", " welcome5 "}), rapid.String()).Draw(t, "coupon"),
			Shipping: rapid.OneOf(rapid.SampledFrom([]string{"", "standard", "express"}), rapid.String()).Draw(t, "shipping"),
		}
	})
}

// receiptOf is Checkout for a basket known to be valid.
func (g baskets) receiptOf(t *rapid.T, b checkout.Basket) checkout.Receipt {
	t.Helper()
	receipt, err := checkout.Checkout(b, g.catalog)
	if err != nil {
		t.Fatalf("Checkout(%+v) returned an error: %v", b, err)
	}
	return receipt
}

// ---------------------------------------------------------------------------
// Properties of a successful checkout
// ---------------------------------------------------------------------------

func TestP1_TheTotalsAddUp(t *testing.T) {
	t.Parallel()
	g := newBaskets(t)
	rapid.Check(t, func(t *rapid.T) {
		r := g.receiptOf(t, g.valid().Draw(t, "basket"))

		want := r.Subtotal - r.PromotionDiscount - r.CouponDiscount + r.Shipping + r.Tax
		if r.Total != want {
			t.Errorf("total = %d, want %d", r.Total, want)
		}
	})
}

func TestP2_NothingIsNegativeAndADiscountNeverExceedsWhatItDiscounts(t *testing.T) {
	t.Parallel()
	g := newBaskets(t)
	rapid.Check(t, func(t *rapid.T) {
		r := g.receiptOf(t, g.valid().Draw(t, "basket"))

		for name, amount := range map[string]int{
			"subtotal": r.Subtotal, "promotion_discount": r.PromotionDiscount, "coupon_discount": r.CouponDiscount,
			"shipping": r.Shipping, "tax": r.Tax, "total": r.Total,
		} {
			if amount < 0 {
				t.Errorf("%s = %d, want >= 0", name, amount)
			}
		}
		for _, line := range r.Lines {
			if line.UnitPrice < 0 || line.LineTotal < 0 || line.LineDiscount < 0 || line.LineDiscount > line.LineTotal {
				t.Errorf("line %s has odd amounts: %+v", line.SKU, line)
			}
		}
		if r.PromotionDiscount > r.Subtotal {
			t.Errorf("promotion_discount %d > subtotal %d", r.PromotionDiscount, r.Subtotal)
		}
		if r.CouponDiscount > r.Subtotal-r.PromotionDiscount {
			t.Errorf("coupon_discount %d > after promotions %d", r.CouponDiscount, r.Subtotal-r.PromotionDiscount)
		}
	})
}

func TestP3_TheReceiptLinesExplainTheTotals(t *testing.T) {
	t.Parallel()
	g := newBaskets(t)
	rapid.Check(t, func(t *rapid.T) {
		r := g.receiptOf(t, g.valid().Draw(t, "basket"))

		sumTotals, sumDiscounts := 0, 0
		for _, line := range r.Lines {
			if line.LineTotal != line.UnitPrice*line.Quantity {
				t.Errorf("line %s: line_total %d != %d × %d", line.SKU, line.LineTotal, line.UnitPrice, line.Quantity)
			}
			sumTotals += line.LineTotal
			sumDiscounts += line.LineDiscount
		}
		if r.Subtotal != sumTotals {
			t.Errorf("subtotal = %d, want the sum of line totals %d", r.Subtotal, sumTotals)
		}
		if r.PromotionDiscount != sumDiscounts {
			t.Errorf("promotion_discount = %d, want the sum of line discounts %d", r.PromotionDiscount, sumDiscounts)
		}
	})
}

func TestP4_ShufflingTheItemsGivesTheSameTotals(t *testing.T) {
	t.Parallel()
	g := newBaskets(t)
	rapid.Check(t, func(t *rapid.T) {
		b := g.valid().Draw(t, "basket")
		shuffled := b
		shuffled.Items = rapid.Permutation(b.Items).Draw(t, "shuffled")

		want := spec.TotalsOf(g.receiptOf(t, b))
		got := spec.TotalsOf(g.receiptOf(t, shuffled))
		if got != want {
			t.Errorf("totals of the shuffled basket = %+v, want %+v", got, want)
		}
	})
}

func TestP5_SplittingOneLineIntoTwoLinesOfTheSameSKUGivesTheSameReceipt(t *testing.T) {
	t.Parallel()
	g := newBaskets(t)
	hasSplittableLine := func(b checkout.Basket) bool {
		return slices.ContainsFunc(b.Items, func(i checkout.BasketItem) bool { return i.Quantity >= 2 })
	}
	rapid.Check(t, func(t *rapid.T) {
		b := g.valid().Filter(hasSplittableLine).Draw(t, "basket")
		var candidates []int
		for i, it := range b.Items {
			if it.Quantity >= 2 {
				candidates = append(candidates, i)
			}
		}
		index := rapid.SampledFrom(candidates).Draw(t, "line")
		line := b.Items[index]
		first := rapid.IntRange(1, line.Quantity-1).Draw(t, "first")

		// The chosen line keeps `first` items; the rest of them go on a new line at the end.
		split := b
		split.Items = make([]checkout.BasketItem, 0, len(b.Items)+1)
		split.Items = append(split.Items, b.Items[:index]...)
		split.Items = append(split.Items, checkout.BasketItem{SKU: line.SKU, Quantity: first})
		split.Items = append(split.Items, b.Items[index+1:]...)
		split.Items = append(split.Items, checkout.BasketItem{SKU: line.SKU, Quantity: line.Quantity - first})

		if diff := cmp.Diff(g.receiptOf(t, b), g.receiptOf(t, split)); diff != "" {
			t.Errorf("the split basket gives a different receipt (-original +split):\n%s", diff)
		}
	})
}

func TestP6_TheShippingTableIsAppliedExactly(t *testing.T) {
	t.Parallel()
	g := newBaskets(t)
	rapid.Check(t, func(t *rapid.T) {
		b := g.valid().Draw(t, "basket")
		r := g.receiptOf(t, b)

		if r.Shipping != 0 && r.Shipping != 499 && r.Shipping != 999 {
			t.Errorf("shipping = %d, want 0, 499 or 999", r.Shipping)
		}
		discounted := r.Subtotal - r.PromotionDiscount - r.CouponDiscount
		want := 499
		switch {
		case b.Shipping == "express":
			want = 999
		case discounted >= 5000:
			want = 0
		}
		if r.Shipping != want {
			t.Errorf("shipping = %d for %q and %d discounted, want %d", r.Shipping, b.Shipping, discounted, want)
		}
	})
}

func TestP7_TaxIsEightPercentRoundedHalfUpAndShippingIsNeverTaxed(t *testing.T) {
	t.Parallel()
	g := newBaskets(t)
	rapid.Check(t, func(t *rapid.T) {
		r := g.receiptOf(t, g.valid().Draw(t, "basket"))

		discounted := r.Subtotal - r.PromotionDiscount - r.CouponDiscount
		if want := (discounted*8 + 50) / 100; r.Tax != want {
			t.Errorf("tax = %d, want %d", r.Tax, want)
		}
	})
}

func TestP8_WritingTheCouponInAnotherCaseOrWithSpacesGivesTheSameReceipt(t *testing.T) {
	t.Parallel()
	g := newBaskets(t)
	// Baskets with a real coupon: like basket(), but the coupon is always SAVE10 or WELCOME5.
	withRealCoupon := rapid.Custom(func(t *rapid.T) checkout.Basket {
		return checkout.Basket{
			Items:    rapid.SliceOfN(g.item(), 1, 8).Draw(t, "items"),
			Coupon:   rapid.SampledFrom([]string{"SAVE10", "WELCOME5"}).Draw(t, "coupon"),
			Shipping: shipping.Draw(t, "shipping"),
		}
	}).Filter(g.isValid)
	rapid.Check(t, func(t *rapid.T) {
		b := withRealCoupon.Draw(t, "basket")
		sloppyBasket := b
		sloppyBasket.Coupon = sloppy(t, b.Coupon)

		if diff := cmp.Diff(g.receiptOf(t, b), g.receiptOf(t, sloppyBasket)); diff != "" {
			t.Errorf("coupon %q gives a different receipt than %q (-clean +sloppy):\n%s", sloppyBasket.Coupon, b.Coupon, diff)
		}
	})
}

func TestP9_OneLinePerDistinctSKUInOrderOfFirstAppearanceEachWithQuantity1To99(t *testing.T) {
	t.Parallel()
	g := newBaskets(t)
	rapid.Check(t, func(t *rapid.T) {
		b := g.valid().Draw(t, "basket")
		r := g.receiptOf(t, b)

		var distinctInOrder []string
		quantities := map[string]int{}
		for _, it := range b.Items {
			if _, seen := quantities[it.SKU]; !seen {
				distinctInOrder = append(distinctInOrder, it.SKU)
			}
			quantities[it.SKU] += it.Quantity
		}
		var lineSKUs []string
		for _, line := range r.Lines {
			lineSKUs = append(lineSKUs, line.SKU)
			if line.Quantity < 1 || line.Quantity > 99 {
				t.Errorf("line %s has quantity %d, want 1 to 99", line.SKU, line.Quantity)
			}
			if line.Quantity != quantities[line.SKU] {
				t.Errorf("line %s has quantity %d, want %d", line.SKU, line.Quantity, quantities[line.SKU])
			}
		}
		if diff := cmp.Diff(distinctInOrder, lineSKUs); diff != "" {
			t.Errorf("line order mismatch (-want +got):\n%s", diff)
		}
	})
}

// ---------------------------------------------------------------------------
// Properties of any input at all
// ---------------------------------------------------------------------------

func TestP10_ForAnyBasketCheckoutGivesASchemaValidReceiptOrOneOfTheSixErrorCodes(t *testing.T) {
	t.Parallel()
	g := newBaskets(t)
	rapid.Check(t, func(t *rapid.T) {
		b := g.any().Draw(t, "basket")

		receipt, err := checkout.Checkout(b, g.catalog) // must never panic
		if err != nil {
			var checkoutErr *checkout.CheckoutError
			if !errors.As(err, &checkoutErr) {
				t.Fatalf("err = %v, want a *CheckoutError", err)
			}
			if !slices.Contains(checkout.ErrorCodes, checkoutErr.Code) {
				t.Errorf("code %q is not one of the six error codes", checkoutErr.Code)
			}
			return
		}
		problems, schemaErr := spec.SchemaErrors("receipt.schema.json", receipt)
		if schemaErr != nil {
			t.Fatalf("loading the schema: %v", schemaErr)
		}
		if len(problems) != 0 {
			t.Errorf("the receipt does not follow receipt.schema.json: %v", problems)
		}
	})
}

func TestP11_RoundHalfUpIsWithinHalfACentAndExactWhenTheValueIsWhole(t *testing.T) {
	t.Parallel()
	rapid.Check(t, func(t *rapid.T) {
		amount := rapid.IntRange(0, 10_000_000).Draw(t, "amount")
		percent := rapid.IntRange(0, 100).Draw(t, "percent")

		rounded := checkout.RoundHalfUp(amount, percent)
		exactTimes100 := amount * percent // the exact value is this ÷ 100

		if diff := rounded*100 - exactTimes100; diff > 50 || diff < -50 {
			t.Errorf("RoundHalfUp(%d, %d) = %d is more than half a cent from %d/100", amount, percent, rounded, exactTimes100)
		}
		if exactTimes100%100 == 0 && rounded*100 != exactTimes100 {
			t.Errorf("RoundHalfUp(%d, %d) = %d, want exactly %d", amount, percent, rounded, exactTimes100/100)
		}
		if exactTimes100%100 == 50 && rounded*100 != exactTimes100+50 {
			t.Errorf("RoundHalfUp(%d, %d) = %d, want the half rounded UP", amount, percent, rounded)
		}
	})
}

// ---------------------------------------------------------------------------
// Properties of saving orders
// ---------------------------------------------------------------------------

func TestP12_AfterSavingNOrdersTheIDsAreExactlyORDER1ToORDERn(t *testing.T) {
	t.Parallel()
	g := newBaskets(t)
	rapid.Check(t, func(t *rapid.T) {
		bs := rapid.SliceOfN(g.valid(), 0, 20).Draw(t, "baskets")
		repository := checkout.NewInMemoryOrderRepository(fixtures.FixedClock())

		var ids, want []string
		for i, b := range bs {
			placed, err := checkout.PlaceOrder(b, g.catalog, repository)
			if err != nil {
				t.Fatalf("PlaceOrder() returned an error: %v", err)
			}
			ids = append(ids, placed.OrderID)
			want = append(want, checkout.OrderID(i+1))
		}
		if diff := cmp.Diff(want, ids); diff != "" {
			t.Errorf("ids mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestP13_FindOfSaveGivesBackTheSameReceiptAndTheClocksTime(t *testing.T) {
	t.Parallel()
	g := newBaskets(t)
	// Any moment this century, with any fraction of a second.
	start := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	end := time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC).Unix()
	rapid.Check(t, func(t *rapid.T) {
		receipt := g.receiptOf(t, g.valid().Draw(t, "basket"))
		seconds := rapid.Int64Range(start, end).Draw(t, "seconds")
		nanos := rapid.Int64Range(0, 999_999_999).Draw(t, "nanos")
		when := time.Unix(seconds, nanos).UTC()
		repository := checkout.NewInMemoryOrderRepository(checkout.FixedClock(when))

		id, err := repository.Save(receipt)
		if err != nil {
			t.Fatalf("Save() returned an error: %v", err)
		}
		order, err := repository.Find(id)
		if err != nil || order == nil {
			t.Fatalf("Find(%q) = %v, %v; want the order", id, order, err)
		}

		if diff := cmp.Diff(receipt, order.Receipt); diff != "" {
			t.Errorf("receipt mismatch (-saved +found):\n%s", diff)
		}
		// time.RFC3339 is an independent way to write the same thing (and prints Z for UTC).
		if want := when.Truncate(time.Second).Format(time.RFC3339); order.CreatedAt != want {
			t.Errorf("created_at = %q, want %q", order.CreatedAt, want)
		}
	})
}

func TestP14_AFailedCheckoutSavesNothingAndFindOfAnUnknownIDIsNotFound(t *testing.T) {
	t.Parallel()
	g := newBaskets(t)
	invalid := g.any().Filter(func(b checkout.Basket) bool { return !g.isValid(b) })
	rapid.Check(t, func(t *rapid.T) {
		b := invalid.Draw(t, "basket")
		randomID := rapid.String().Draw(t, "id")
		repository := checkout.NewInMemoryOrderRepository(fixtures.FixedClock())

		_, err := checkout.PlaceOrder(b, g.catalog, repository)

		var checkoutErr *checkout.CheckoutError
		if !errors.As(err, &checkoutErr) {
			t.Fatalf("err = %v, want a *CheckoutError", err)
		}
		if repository.Count() != 0 {
			t.Errorf("Count() = %d after a failed checkout, want 0", repository.Count())
		}
		for _, id := range []string{randomID, "ORDER-1"} {
			if order, err := repository.Find(id); order != nil || err != nil {
				t.Errorf("Find(%q) = %v, %v; want nil, nil", id, order, err)
			}
		}
	})
}
