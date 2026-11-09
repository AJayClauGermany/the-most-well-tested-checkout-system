package checkout_test

// UNIT TESTS of Checkout as a whole, with the tiny made-up catalog from internal/fixtures.
// The 60 shared spec cases (spec_cases_test.go) cover the standard catalog; these tests
// check that the pieces fit together, and the things the totals alone cannot show, such
// as the receipt lines. Whole receipts are compared with cmp.Diff.

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/fixtures"
)

// mustCheckout runs Checkout on the tiny catalog and stops the test if it fails.
func mustCheckout(t *testing.T, basket checkout.Basket) checkout.Receipt {
	t.Helper()
	receipt, err := checkout.Checkout(basket, fixtures.TinyCatalog)
	if err != nil {
		t.Fatalf("Checkout(%+v) returned an error: %v", basket, err)
	}
	return receipt
}

func TestCheckoutPricesOneApple(t *testing.T) {
	t.Parallel()
	// Arrange
	basket := items(item("APPLE", 1))

	// Act
	receipt := mustCheckout(t, basket)

	// Assert: 50 + $4.99 shipping + 4 tax = $5.53, and exactly one line.
	want := fixtures.OneAppleReceipt()
	if diff := cmp.Diff(want, receipt); diff != "" {
		t.Errorf("Checkout() mismatch (-want +got):\n%s", diff)
	}
}

func TestCheckoutWritesOneLinePerDistinctSKUInOrderOfFirstAppearance(t *testing.T) {
	t.Parallel()
	receipt := mustCheckout(t, items(item("PEN", 1), item("APPLE", 2), item("PEN", 3)))

	want := []checkout.ReceiptLine{
		{SKU: "PEN", Name: "Pen", Quantity: 4, UnitPrice: 120, LineTotal: 480, LineDiscount: 0},
		{SKU: "APPLE", Name: "Apple", Quantity: 2, UnitPrice: 50, LineTotal: 100, LineDiscount: 0},
	}
	if diff := cmp.Diff(want, receipt.Lines); diff != "" {
		t.Errorf("receipt lines mismatch (-want +got):\n%s", diff)
	}
}

func TestCheckoutWorksTheTotalsOutFromTheLines(t *testing.T) {
	t.Parallel()
	// 3 coffees ($26.97, one free: -$8.99) and 10 notebooks ($34.90, 10% off: -$3.49).
	receipt := mustCheckout(t, items(item("COFFEE", 3), item("NOTEBOOK", 10)))

	want := checkout.Receipt{
		Lines: []checkout.ReceiptLine{
			{SKU: "COFFEE", Name: "Bag of Coffee Beans", Quantity: 3, UnitPrice: 899, LineTotal: 2697, LineDiscount: 899},
			{SKU: "NOTEBOOK", Name: "Notebook", Quantity: 10, UnitPrice: 349, LineTotal: 3490, LineDiscount: 349},
		},
		Subtotal:          2697 + 3490,
		PromotionDiscount: 899 + 349,
		CouponDiscount:    0,
		Shipping:          499, // $49.39 after promotions: not free
		Tax:               395, // 8% of 4939 = 395.12
		Total:             4939 + 499 + 395,
	}
	if diff := cmp.Diff(want, receipt); diff != "" {
		t.Errorf("Checkout() mismatch (-want +got):\n%s", diff)
	}
}

func TestCheckoutTakesTheCouponOffAfterPromotionsAndBeforeShippingAndTax(t *testing.T) {
	t.Parallel()
	// TV $100.00 with SAVE10: $90.00 → free shipping, tax 8% of $90.00 = $7.20.
	receipt := mustCheckout(t, checkout.Basket{Items: []checkout.BasketItem{item("TV", 1)}, Coupon: "SAVE10"})

	got := [6]int{receipt.Subtotal, receipt.PromotionDiscount, receipt.CouponDiscount, receipt.Shipping, receipt.Tax, receipt.Total}
	want := [6]int{10000, 0, 1000, 0, 720, 9720}
	if got != want {
		t.Errorf("totals (subtotal, promotion, coupon, shipping, tax, total) = %v, want %v", got, want)
	}
}

func TestCheckoutChargesExpressShippingOnTopUntaxed(t *testing.T) {
	t.Parallel()
	receipt := mustCheckout(t, checkout.Basket{Items: []checkout.BasketItem{item("TV", 1)}, Shipping: "express"})

	if receipt.Shipping != 999 {
		t.Errorf("shipping = %d, want 999", receipt.Shipping)
	}
	if receipt.Tax != 800 {
		t.Errorf("tax = %d, want 800 (8%% of $100.00, shipping not taxed)", receipt.Tax)
	}
	if receipt.Total != 10000+999+800 {
		t.Errorf("total = %d, want %d", receipt.Total, 10000+999+800)
	}
}

func TestCheckoutDoesNotChangeTheBasketItIsGiven(t *testing.T) {
	t.Parallel()
	basket := checkout.Basket{Items: []checkout.BasketItem{item("APPLE", 1), item("APPLE", 1)}, Coupon: " save10 "}

	mustCheckout(t, basket)

	want := checkout.Basket{Items: []checkout.BasketItem{item("APPLE", 1), item("APPLE", 1)}, Coupon: " save10 "}
	if diff := cmp.Diff(want, basket); diff != "" {
		t.Errorf("the basket was changed (-want +got):\n%s", diff)
	}
}

func TestCheckoutReturnsAnEmptyReceiptWithTheError(t *testing.T) {
	t.Parallel()
	receipt, err := checkout.Checkout(items(), fixtures.TinyCatalog)

	assertErrorCode(t, err, checkout.CodeEmptyBasket)
	if diff := cmp.Diff(checkout.Receipt{}, receipt); diff != "" {
		t.Errorf("receipt should be empty on failure (-want +got):\n%s", diff)
	}
}

// Rules §4.1: when several things are wrong, the FIRST row of the table that applies wins.
func TestCheckoutReportsTheFirstFailingRowOfRules41(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		basket checkout.Basket
		want   checkout.ErrorCode
	}{
		{
			name:   "empty basket before anything else",
			basket: checkout.Basket{Items: nil, Coupon: "NOPE", Shipping: "overnight"},
			want:   checkout.CodeEmptyBasket,
		},
		{
			name:   "unknown product before bad quantity, shipping and coupon",
			basket: checkout.Basket{Items: []checkout.BasketItem{item("APPLE", 0), item("BANANA", 1)}, Coupon: "NOPE", Shipping: "x"},
			want:   checkout.CodeUnknownProduct,
		},
		{
			name:   "bad quantity before shipping and coupon",
			basket: checkout.Basket{Items: []checkout.BasketItem{item("APPLE", 0)}, Coupon: "NOPE", Shipping: "x"},
			want:   checkout.CodeInvalidQuantity,
		},
		{
			name:   "merged quantity over 99 before shipping and coupon",
			basket: checkout.Basket{Items: []checkout.BasketItem{item("APPLE", 60), item("APPLE", 40)}, Coupon: "NOPE", Shipping: "x"},
			want:   checkout.CodeInvalidQuantity,
		},
		{
			name:   "bad shipping before bad coupon",
			basket: checkout.Basket{Items: []checkout.BasketItem{item("APPLE", 1)}, Coupon: "NOPE", Shipping: "x"},
			want:   checkout.CodeUnknownShippingMethod,
		},
		{
			name:   "bad shipping before an unmet coupon minimum",
			basket: checkout.Basket{Items: []checkout.BasketItem{item("APPLE", 1)}, Coupon: "WELCOME5", Shipping: "x"},
			want:   checkout.CodeUnknownShippingMethod,
		},
		{
			name:   "unknown coupon",
			basket: checkout.Basket{Items: []checkout.BasketItem{item("APPLE", 1)}, Coupon: "NOPE"},
			want:   checkout.CodeUnknownCoupon,
		},
		{
			// 3 coffees are $26.97 but $17.98 after the free one: below WELCOME5's $25.00.
			name:   "coupon minimum not met, judged on the price after promotions",
			basket: checkout.Basket{Items: []checkout.BasketItem{item("COFFEE", 3)}, Coupon: "WELCOME5"},
			want:   checkout.CodeCouponMinimumNotMet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := checkout.Checkout(tt.basket, fixtures.TinyCatalog)

			assertErrorCode(t, err, tt.want)
		})
	}
}
