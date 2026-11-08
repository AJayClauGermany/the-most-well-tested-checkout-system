package checkout

// Unit tests for the error types. They make sure our Go constants are spelled
// exactly like the codes in spec/checkout-rules.md §9 – a typo here would make
// every shared spec case about errors fail in a confusing way.

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorCodesMatchTheSpecSpelling(t *testing.T) {
	tests := []struct {
		code ErrorCode
		want string
	}{
		{CodeEmptyBasket, "EMPTY_BASKET"},
		{CodeUnknownProduct, "UNKNOWN_PRODUCT"},
		{CodeInvalidQuantity, "INVALID_QUANTITY"},
		{CodeUnknownShippingMethod, "UNKNOWN_SHIPPING_METHOD"},
		{CodeUnknownCoupon, "UNKNOWN_COUPON"},
		{CodeCouponMinimumNotMet, "COUPON_MINIMUM_NOT_MET"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := string(tt.code); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCheckoutErrorCodeCanBeFoundWithErrorsAs(t *testing.T) {
	// Arrange: wrap the error, like a caller adding context might do.
	err := fmt.Errorf("pricing basket: %w", &CheckoutError{Code: CodeUnknownCoupon})

	// Act: errors.As looks through the wrapping for a *CheckoutError.
	var checkoutErr *CheckoutError
	found := errors.As(err, &checkoutErr)

	// Assert
	if !found {
		t.Fatalf("errors.As did not find a *CheckoutError in %v", err)
	}
	if checkoutErr.Code != CodeUnknownCoupon {
		t.Errorf("code = %s, want %s", checkoutErr.Code, CodeUnknownCoupon)
	}
}
