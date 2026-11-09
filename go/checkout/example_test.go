package checkout_test

// EXAMPLE TESTS: documentation that is also a test.
//
// A function named ExampleXxx with an "// Output:" comment at the end is run by `go test`,
// and its printed output must match the comment exactly. `go doc` and pkg.go.dev show
// these functions as examples next to the code they document – so the examples in the
// documentation can never go stale.

import (
	"errors"
	"fmt"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/fixtures"
)

func ExampleCheckout() {
	catalog := checkout.Catalog{Products: []checkout.Product{
		{SKU: "APPLE", Name: "Apple", UnitPrice: 50},
		{SKU: "COFFEE", Name: "Bag of Coffee Beans", UnitPrice: 899,
			Promotion: &checkout.Promotion{Type: checkout.BuyXGetYFree, Buy: 2, Free: 1}},
	}}
	basket := checkout.Basket{
		Items:  []checkout.BasketItem{{SKU: "COFFEE", Quantity: 3}, {SKU: "APPLE", Quantity: 2}},
		Coupon: "SAVE10",
	}

	receipt, err := checkout.Checkout(basket, catalog)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	for _, line := range receipt.Lines {
		fmt.Printf("%-6s x%d  %5d cents, %d off\n", line.SKU, line.Quantity, line.LineTotal, line.LineDiscount)
	}
	fmt.Println("subtotal:", receipt.Subtotal)
	fmt.Println("promotions:", receipt.PromotionDiscount)
	fmt.Println("coupon:", receipt.CouponDiscount)
	fmt.Println("shipping:", receipt.Shipping)
	fmt.Println("tax:", receipt.Tax)
	fmt.Println("total:", receipt.Total)
	// Output:
	// COFFEE x3   2697 cents, 899 off
	// APPLE  x2    100 cents, 0 off
	// subtotal: 2797
	// promotions: 899
	// coupon: 190
	// shipping: 499
	// tax: 137
	// total: 2344
}

func ExampleCheckout_invalidBasket() {
	basket := checkout.Basket{} // no items

	_, err := checkout.Checkout(basket, fixtures.TinyCatalog)

	var checkoutErr *checkout.CheckoutError
	if errors.As(err, &checkoutErr) {
		fmt.Println("code:", checkoutErr.Code)
		fmt.Println("error:", err)
	}
	// Output:
	// code: EMPTY_BASKET
	// error: checkout failed: EMPTY_BASKET (The basket has no items)
}

func ExampleRoundHalfUp() {
	fmt.Println(checkout.RoundHalfUp(224, 10)) // 22.4
	fmt.Println(checkout.RoundHalfUp(225, 10)) // 22.5 – an exact half goes up
	fmt.Println(checkout.RoundHalfUp(226, 10)) // 22.6
	// Output:
	// 22
	// 23
	// 23
}

func ExampleMergeItems() {
	merged := checkout.MergeItems([]checkout.BasketItem{
		{SKU: "APPLE", Quantity: 1},
		{SKU: "PEN", Quantity: 2},
		{SKU: "APPLE", Quantity: 3},
	})
	for _, item := range merged {
		fmt.Printf("%s x%d\n", item.SKU, item.Quantity)
	}
	// Output:
	// APPLE x4
	// PEN x2
}

func ExamplePlaceOrder() {
	clock := checkout.FixedClock(fixtures.FixedTime) // a real program would use checkout.SystemClock()
	repository := checkout.NewInMemoryOrderRepository(clock)
	basket := checkout.Basket{Items: []checkout.BasketItem{{SKU: "APPLE", Quantity: 1}}}

	placed, err := checkout.PlaceOrder(basket, fixtures.TinyCatalog, repository)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	order, _ := repository.Find(placed.OrderID)
	fmt.Println(order.ID, order.CreatedAt, order.Receipt.Total)
	// Output:
	// ORDER-1 2026-01-15T10:30:00Z 553
}
