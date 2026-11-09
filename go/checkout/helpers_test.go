package checkout_test

// Small helpers shared by the black-box tests in this folder. Helpers call t.Helper(), so
// a failure is reported at the line in the test that called them, not in here.

import (
	"errors"
	"testing"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/internal/spec"
)

// assertErrorCode fails the test unless err is a *CheckoutError carrying want. It returns
// the error so the caller can also look at its message.
func assertErrorCode(t *testing.T, err error, want checkout.ErrorCode) *checkout.CheckoutError {
	t.Helper()

	var checkoutErr *checkout.CheckoutError
	if !errors.As(err, &checkoutErr) {
		t.Fatalf("got err = %v, want a *CheckoutError with code %s", err, want)
	}
	if checkoutErr.Code != want {
		t.Errorf("error code = %s, want %s", checkoutErr.Code, want)
	}
	return checkoutErr
}

// assertMessage fails the test unless the error's message is exactly want.
func assertMessage(t *testing.T, err *checkout.CheckoutError, want string) {
	t.Helper()
	if err.Message != want {
		t.Errorf("error message = %q, want %q", err.Message, want)
	}
}

// loadCatalog reads the standard catalog from spec/catalog.json, or stops the test.
func loadCatalog(t testing.TB) checkout.Catalog {
	t.Helper()
	catalog, err := spec.LoadCatalog()
	if err != nil {
		t.Fatalf("loading the catalog: %v", err)
	}
	return catalog
}

// loadCases reads the 60 shared cases from spec/test-cases/, or stops the test.
func loadCases(t testing.TB) []spec.Case {
	t.Helper()
	cases, err := spec.LoadCases()
	if err != nil {
		t.Fatalf("loading the spec cases: %v", err)
	}
	return cases
}
