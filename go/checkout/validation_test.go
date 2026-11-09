package checkout_test

// UNIT TESTS for basket validation (rules §4.1, rows 1 to 4) and the order errors are
// reported in. The check order matters: a basket can be wrong in several ways at once.
// Boundaries are tested on both sides (0/1, 99/100) because `>` versus `>=` is the
// classic bug – and the classic surviving mutant in mutation testing.

import (
	"testing"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/fixtures"
)

func TestQuantityLimitsAreTheOnesFromTheRules(t *testing.T) {
	t.Parallel()
	if checkout.MinQuantity != 1 {
		t.Errorf("MinQuantity = %d, want 1", checkout.MinQuantity)
	}
	if checkout.MaxQuantity != 99 {
		t.Errorf("MaxQuantity = %d, want 99", checkout.MaxQuantity)
	}
}

func TestIsValidQuantity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		quantity int
		want     bool
	}{
		{quantity: 1, want: true}, // the lowest valid quantity
		{quantity: 2, want: true},
		{quantity: 50, want: true},
		{quantity: 98, want: true},
		{quantity: 99, want: true}, // the highest valid quantity
		{quantity: 0, want: false}, // one below the lowest
		{quantity: -1, want: false},
		{quantity: 100, want: false}, // one above the highest
		{quantity: 101, want: false},
		{quantity: 1000, want: false},
	}

	for _, tt := range tests {
		t.Run(itoa(tt.quantity), func(t *testing.T) {
			t.Parallel()
			if got := checkout.IsValidQuantity(tt.quantity); got != tt.want {
				t.Errorf("IsValidQuantity(%d) = %v, want %v", tt.quantity, got, tt.want)
			}
		})
	}
}

// items is a short way to write a basket with these items and nothing else.
func items(list ...checkout.BasketItem) checkout.Basket {
	return checkout.Basket{Items: list}
}

// item is a short way to write one basket item.
func item(sku string, quantity int) checkout.BasketItem {
	return checkout.BasketItem{SKU: sku, Quantity: quantity}
}

func TestValidateItemsAcceptsANormalBasket(t *testing.T) {
	t.Parallel()
	basket := items(item("APPLE", 1), item("COFFEE", 99))

	err := checkout.ValidateItems(basket, fixtures.TinyCatalog)

	if err != nil {
		t.Errorf("ValidateItems() = %v, want nil", err)
	}
}

func TestValidateItemsReportsTheRowsOfRules41(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		basket      checkout.Basket
		wantCode    checkout.ErrorCode
		wantMessage string
	}{
		{
			name:        "row 1: an empty basket",
			basket:      items(),
			wantCode:    checkout.CodeEmptyBasket,
			wantMessage: "The basket has no items",
		},
		{
			name:        "row 2: a sku that is not in the catalog",
			basket:      items(item("BANANA", 1)),
			wantCode:    checkout.CodeUnknownProduct,
			wantMessage: "Unknown product: BANANA",
		},
		{
			name:        "row 2: skus are case sensitive",
			basket:      items(item("apple", 1)),
			wantCode:    checkout.CodeUnknownProduct,
			wantMessage: "Unknown product: apple",
		},
		{
			name:        "row 2: an unknown product on any line, not only the first",
			basket:      items(item("APPLE", 1), item("PEN", 1), item("BANANA", 1)),
			wantCode:    checkout.CodeUnknownProduct,
			wantMessage: "Unknown product: BANANA",
		},
		{
			name:        "row 3: a quantity of 0",
			basket:      items(item("APPLE", 0)),
			wantCode:    checkout.CodeInvalidQuantity,
			wantMessage: "Invalid quantity for APPLE: 0",
		},
		{
			name:        "row 3: a quantity of -1",
			basket:      items(item("APPLE", -1)),
			wantCode:    checkout.CodeInvalidQuantity,
			wantMessage: "Invalid quantity for APPLE: -1",
		},
		{
			name:        "row 3: a quantity of 100",
			basket:      items(item("APPLE", 100)),
			wantCode:    checkout.CodeInvalidQuantity,
			wantMessage: "Invalid quantity for APPLE: 100",
		},
		{
			name:        "row 3: a bad quantity on a later line",
			basket:      items(item("APPLE", 1), item("PEN", 0)),
			wantCode:    checkout.CodeInvalidQuantity,
			wantMessage: "Invalid quantity for PEN: 0",
		},
		{
			name:        "row 4: lines of one sku that add up to more than 99",
			basket:      items(item("APPLE", 60), item("APPLE", 40)),
			wantCode:    checkout.CodeInvalidQuantity,
			wantMessage: "Too many APPLE in total: 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := checkout.ValidateItems(tt.basket, fixtures.TinyCatalog)

			checkoutErr := assertErrorCode(t, err, tt.wantCode)
			assertMessage(t, checkoutErr, tt.wantMessage)
		})
	}
}

func TestValidateItemsRow4BoundaryOnBothSides(t *testing.T) {
	t.Parallel()
	t.Run("lines of one sku that add up to exactly 99 are fine", func(t *testing.T) {
		t.Parallel()
		basket := items(item("APPLE", 50), item("PEN", 99), item("APPLE", 49))

		if err := checkout.ValidateItems(basket, fixtures.TinyCatalog); err != nil {
			t.Errorf("ValidateItems() = %v, want nil", err)
		}
	})

	t.Run("99 of one sku and 99 of another is fine (the limit is per sku)", func(t *testing.T) {
		t.Parallel()
		basket := items(item("APPLE", 99), item("PEN", 99))

		if err := checkout.ValidateItems(basket, fixtures.TinyCatalog); err != nil {
			t.Errorf("ValidateItems() = %v, want nil", err)
		}
	})
}

func TestValidateItemsCheckOrder(t *testing.T) {
	t.Parallel()
	t.Run("an unknown product beats a bad quantity, even on a later line", func(t *testing.T) {
		t.Parallel()
		err := checkout.ValidateItems(items(item("APPLE", 0), item("BANANA", 1)), fixtures.TinyCatalog)

		assertErrorCode(t, err, checkout.CodeUnknownProduct)
	})

	t.Run("a bad single quantity beats a bad merged quantity (both are INVALID_QUANTITY)", func(t *testing.T) {
		t.Parallel()
		err := checkout.ValidateItems(items(item("APPLE", 100), item("APPLE", 1)), fixtures.TinyCatalog)

		checkoutErr := assertErrorCode(t, err, checkout.CodeInvalidQuantity)
		assertMessage(t, checkoutErr, "Invalid quantity for APPLE: 100")
	})

	t.Run("ignores the coupon and shipping method: those are checked later, by Checkout", func(t *testing.T) {
		t.Parallel()
		basket := checkout.Basket{Items: []checkout.BasketItem{item("APPLE", 1)}, Coupon: "NOPE", Shipping: "overnight"}

		if err := checkout.ValidateItems(basket, fixtures.TinyCatalog); err != nil {
			t.Errorf("ValidateItems() = %v, want nil", err)
		}
	})
}
