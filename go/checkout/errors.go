package checkout

import "fmt"

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

// ErrorCodes lists every code, in the order rules §9 lists them. Handy for loops and tests.
var ErrorCodes = []ErrorCode{
	CodeEmptyBasket,
	CodeUnknownProduct,
	CodeInvalidQuantity,
	CodeUnknownShippingMethod,
	CodeUnknownCoupon,
	CodeCouponMinimumNotMet,
}

// CheckoutError is returned by Checkout when the basket is not valid.
//
// In Go, errors are ordinary return values (there are no exceptions).
// Callers get the code back out with errors.As:
//
//	var checkoutErr *checkout.CheckoutError
//	if errors.As(err, &checkoutErr) {
//		fmt.Println(checkoutErr.Code) // e.g. EMPTY_BASKET
//	}
//
// Code is what programs look at; Message is for the human reading a log
// ("Unknown product: BANANA"). It may be empty.
type CheckoutError struct {
	Code    ErrorCode
	Message string
}

// Error makes *CheckoutError satisfy Go's built-in error interface.
func (e *CheckoutError) Error() string {
	if e.Message == "" {
		return "checkout failed: " + string(e.Code)
	}
	return "checkout failed: " + string(e.Code) + " (" + e.Message + ")"
}

// newError builds a *CheckoutError with a formatted message, e.g.
//
//	newError(CodeUnknownProduct, "Unknown product: %s", sku)
func newError(code ErrorCode, format string, args ...any) *CheckoutError {
	return &CheckoutError{Code: code, Message: fmt.Sprintf(format, args...)}
}
