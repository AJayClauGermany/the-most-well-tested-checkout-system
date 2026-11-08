package checkout

import "errors"

// ErrorCode is one of the six error codes from the rules (§9).
// It is a string type, so the value is exactly the spelling used in the spec.
type ErrorCode string

// The error codes, in the order the rules check for them.
const (
	CodeEmptyBasket           ErrorCode = "EMPTY_BASKET"
	CodeUnknownProduct        ErrorCode = "UNKNOWN_PRODUCT"
	CodeInvalidQuantity       ErrorCode = "INVALID_QUANTITY"
	CodeUnknownShippingMethod ErrorCode = "UNKNOWN_SHIPPING_METHOD"
	CodeUnknownCoupon         ErrorCode = "UNKNOWN_COUPON"
	CodeCouponMinimumNotMet   ErrorCode = "COUPON_MINIMUM_NOT_MET"
)

// CheckoutError is returned by Checkout when the basket is not valid.
//
// In Go, errors are ordinary return values (there are no exceptions).
// Callers get the code back out with errors.As:
//
//	var checkoutErr *checkout.CheckoutError
//	if errors.As(err, &checkoutErr) {
//		fmt.Println(checkoutErr.Code) // e.g. EMPTY_BASKET
//	}
type CheckoutError struct {
	Code ErrorCode
}

// Error makes *CheckoutError satisfy Go's built-in error interface.
func (e *CheckoutError) Error() string {
	return "checkout failed: " + string(e.Code)
}

// ErrNotImplemented is returned by Checkout until someone implements it.
// A "sentinel error" like this is compared with errors.Is(err, ErrNotImplemented).
var ErrNotImplemented = errors.New("checkout: not implemented yet")
