package checkout

// Unit tests for the error types. They make sure our Go constants are spelled
// exactly like the codes in spec/checkout-rules.md §9 – a typo here would make
// every shared spec case about errors fail in a confusing way.

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
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

func TestErrorCodesListsTheSixCodesInTheOrderOfTheRules(t *testing.T) {
	want := []ErrorCode{
		"EMPTY_BASKET", "UNKNOWN_PRODUCT", "INVALID_QUANTITY",
		"UNKNOWN_SHIPPING_METHOD", "UNKNOWN_COUPON", "COUPON_MINIMUM_NOT_MET",
	}

	if diff := cmp.Diff(want, ErrorCodes); diff != "" {
		t.Errorf("ErrorCodes mismatch (-want +got):\n%s", diff)
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

func TestCheckoutErrorReadsWellWhenPrinted(t *testing.T) {
	tests := []struct {
		name string
		err  *CheckoutError
		want string
	}{
		{
			name: "code only",
			err:  &CheckoutError{Code: CodeEmptyBasket},
			want: "checkout failed: EMPTY_BASKET",
		},
		{
			name: "code and message",
			err:  &CheckoutError{Code: CodeEmptyBasket, Message: "The basket has no items"},
			want: "checkout failed: EMPTY_BASKET (The basket has no items)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

// A white-box test of the unexported helper every step uses to build its error.
func TestNewErrorFormatsTheMessage(t *testing.T) {
	err := newError(CodeInvalidQuantity, "Invalid quantity for %s: %d", "APPLE", 0)

	if err.Code != CodeInvalidQuantity {
		t.Errorf("code = %s, want %s", err.Code, CodeInvalidQuantity)
	}
	if err.Message != "Invalid quantity for APPLE: 0" {
		t.Errorf("message = %q, want %q", err.Message, "Invalid quantity for APPLE: 0")
	}
}
