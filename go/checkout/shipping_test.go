package checkout_test

// UNIT TESTS for shipping (rules §4.5). The free-shipping threshold (5000) is tested on
// both sides: 4999 pays, 5000 is free.

import (
	"fmt"
	"testing"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
)

func TestParseShippingMethodDefaultsToStandardWhenTheBasketSaysNothing(t *testing.T) {
	t.Parallel()
	method, err := checkout.ParseShippingMethod("")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != checkout.ShippingStandard {
		t.Errorf("method = %q, want %q", method, checkout.ShippingStandard)
	}
}

func TestParseShippingMethodAcceptsTheTwoKnownMethods(t *testing.T) {
	t.Parallel()
	for _, want := range []string{checkout.ShippingStandard, checkout.ShippingExpress} {
		method, err := checkout.ParseShippingMethod(want)

		if err != nil {
			t.Fatalf("ParseShippingMethod(%q) returned an error: %v", want, err)
		}
		if string(method) != want {
			t.Errorf("ParseShippingMethod(%q) = %q", want, method)
		}
	}
}

func TestParseShippingMethodRejectsEverythingElse(t *testing.T) {
	t.Parallel()
	// Shipping method names are matched exactly, in lower case (rules §4.5).
	for _, input := range []string{"overnight", "Express", "STANDARD", " standard", "express "} {
		t.Run(fmt.Sprintf("%q", input), func(t *testing.T) {
			t.Parallel()
			_, err := checkout.ParseShippingMethod(input)

			checkoutErr := assertErrorCode(t, err, checkout.CodeUnknownShippingMethod)
			assertMessage(t, checkoutErr, "Unknown shipping method: "+input)
		})
	}
}

func TestShippingCostsAreTheOnesFromTheRules(t *testing.T) {
	t.Parallel()
	if checkout.StandardShippingCost != 499 || checkout.ExpressShippingCost != 999 || checkout.FreeShippingFrom != 5000 {
		t.Errorf("costs = %d / %d, free from %d; want 499 / 999, free from 5000",
			checkout.StandardShippingCost, checkout.ExpressShippingCost, checkout.FreeShippingFrom)
	}
}

func TestStandardShippingCostsFourNinetyNineBelowFiftyDollars(t *testing.T) {
	t.Parallel()
	for _, discountedSubtotal := range []int{0, 1, 4999} {
		if got := checkout.ShippingCost(checkout.ShippingStandard, discountedSubtotal); got != 499 {
			t.Errorf("standard shipping on %d = %d, want 499", discountedSubtotal, got)
		}
	}
}

func TestStandardShippingIsFreeFromExactlyFiftyDollars(t *testing.T) {
	t.Parallel()
	for _, discountedSubtotal := range []int{5000, 5001, 100000} {
		if got := checkout.ShippingCost(checkout.ShippingStandard, discountedSubtotal); got != 0 {
			t.Errorf("standard shipping on %d = %d, want 0", discountedSubtotal, got)
		}
	}
}

func TestExpressShippingAlwaysCostsNineNinetyNine(t *testing.T) {
	t.Parallel()
	for _, discountedSubtotal := range []int{0, 4999, 5000, 100000} {
		if got := checkout.ShippingCost(checkout.ShippingExpress, discountedSubtotal); got != 999 {
			t.Errorf("express shipping on %d = %d, want 999", discountedSubtotal, got)
		}
	}
}
