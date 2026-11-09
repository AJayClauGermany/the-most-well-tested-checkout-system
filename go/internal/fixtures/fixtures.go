// Package fixtures holds small, hand-made test data shared by the tests of several
// packages.
//
// Unit tests prefer tiny made-up data over the real catalog: when a test fails you can
// work out the right answer in your head. The shared spec cases use the real catalog.
//
// It is test-support code under internal/, like package spec.
package fixtures

import (
	"time"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
)

// The products of the tiny catalog. Apple, Pen and TV have no promotion.
var (
	Apple  = checkout.Product{SKU: "APPLE", Name: "Apple", UnitPrice: 50}
	Pen    = checkout.Product{SKU: "PEN", Name: "Pen", UnitPrice: 120}
	Coffee = checkout.Product{SKU: "COFFEE", Name: "Bag of Coffee Beans", UnitPrice: 899,
		Promotion: &checkout.Promotion{Type: checkout.BuyXGetYFree, Buy: 2, Free: 1}}
	Notebook = checkout.Product{SKU: "NOTEBOOK", Name: "Notebook", UnitPrice: 349,
		Promotion: &checkout.Promotion{Type: checkout.BulkDiscount, MinQuantity: 10, Percent: 10}}
	// TV is a $100.00 item: handy for reaching free shipping and coupon minimums quickly.
	TV = checkout.Product{SKU: "TV", Name: "Television", UnitPrice: 10000}
)

// TinyCatalog has five products, two of them with a promotion.
var TinyCatalog = checkout.Catalog{Products: []checkout.Product{Apple, Pen, Coffee, Notebook, TV}}

// OneAppleReceipt returns a valid receipt for one apple, used wherever "some receipt" is
// needed. It is a function so that every caller gets its own copy to play with.
func OneAppleReceipt() checkout.Receipt {
	return checkout.Receipt{
		Lines: []checkout.ReceiptLine{
			{SKU: "APPLE", Name: "Apple", Quantity: 1, UnitPrice: 50, LineTotal: 50, LineDiscount: 0},
		},
		Subtotal:          50,
		PromotionDiscount: 0,
		CouponDiscount:    0,
		Shipping:          499,
		Tax:               4,
		Total:             553,
	}
}

// FixedTime is 10:30 on 15 January 2026, UTC – the timestamp used in spec/examples/order.json.
var FixedTime = time.Date(2026, time.January, 15, 10, 30, 0, 0, time.UTC)

// FixedTimestamp is FixedTime written the way created_at must look.
const FixedTimestamp = "2026-01-15T10:30:00Z"

// FixedClock returns a stub clock that always answers FixedTime.
func FixedClock() checkout.Clock {
	return checkout.FixedClock(FixedTime)
}

// ExampleBasket is the basket behind spec/examples/receipt.json (rules §5): 3 bags of
// coffee and 2 apples with the coupon SAVE10.
var ExampleBasket = checkout.Basket{
	Items: []checkout.BasketItem{
		{SKU: "COFFEE", Quantity: 3},
		{SKU: "APPLE", Quantity: 2},
	},
	Coupon:   "SAVE10",
	Shipping: "standard",
}

// EverythingBasket is the biggest shared case ("full-everything"): every product of the
// standard catalog, both promotions, WELCOME5 and express shipping.
var EverythingBasket = checkout.Basket{
	Items: []checkout.BasketItem{
		{SKU: "HEADPHONES", Quantity: 1},
		{SKU: "COFFEE", Quantity: 6},
		{SKU: "NOTEBOOK", Quantity: 12},
		{SKU: "PEN", Quantity: 5},
		{SKU: "BREAD", Quantity: 1},
		{SKU: "APPLE", Quantity: 3},
	},
	Coupon:   "WELCOME5",
	Shipping: "express",
}
