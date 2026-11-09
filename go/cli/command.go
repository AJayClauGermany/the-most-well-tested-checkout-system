// Package cli holds the logic of the command line program (rules §8), kept free of
// os.Stdin and os.Stdout so it can be unit tested in-process. cmd/checkout/main.go is the
// thin wrapper that connects it to the terminal.
package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/AJayClauGermany/the-most-well-tested-checkout-system/go/checkout"
)

// Result is what the program prints (one line, without the newline) and the exit code it
// ends with: 0 with a receipt, 1 with {"error": "<CODE>"}.
type Result struct {
	Output   string
	ExitCode int
}

// Run runs one checkout from the JSON text of a basket.
//
//   - a valid basket gives the receipt as JSON and exit code 0,
//   - an invalid basket gives {"error": "<CODE>"} and exit code 1,
//   - text that is not a basket at all gives an error (main turns that into exit code 2).
func Run(input string, catalog checkout.Catalog) (Result, error) {
	basket, err := parseBasket(input)
	if err != nil {
		return Result{}, err
	}

	// Checkout's only error is a *CheckoutError (its documentation says so, and property
	// P10 and the fuzz test check it), so there is no third case to handle here.
	receipt, err := checkout.Checkout(basket, catalog)
	var checkoutErr *checkout.CheckoutError
	if errors.As(err, &checkoutErr) {
		return Result{Output: compactJSON(cliError{Error: checkoutErr.Code}), ExitCode: 1}, nil
	}
	return Result{Output: compactJSON(receipt), ExitCode: 0}, nil
}

// cliError is what the program prints on failure (spec/schema/cli-error.schema.json).
type cliError struct {
	Error checkout.ErrorCode `json:"error"`
}

// rawBasket is the basket as it arrives, with pointers so that a missing field can be told
// apart from an empty one.
type rawBasket struct {
	Items *[]struct {
		SKU      *string `json:"sku"`
		Quantity *int    `json:"quantity"`
	} `json:"items"`
	Coupon   *string `json:"coupon"`
	Shipping *string `json:"shipping"`
}

// parseBasket parses the JSON text of a basket. Only the shape that matters is checked –
// an object with an `items` list of objects that each have a string `sku` and a
// whole-number `quantity` – so that checkout.Checkout can do the real validation.
func parseBasket(input string) (checkout.Basket, error) {
	if !json.Valid([]byte(input)) {
		return checkout.Basket{}, errors.New("the input is not JSON")
	}
	if !strings.HasPrefix(strings.TrimSpace(input), "{") {
		return checkout.Basket{}, errors.New("the input is not a basket: expected a JSON object")
	}

	var raw rawBasket
	if err := json.Unmarshal([]byte(input), &raw); err != nil {
		return checkout.Basket{}, fmt.Errorf("the input is not a basket: %w", err)
	}
	if raw.Items == nil {
		return checkout.Basket{}, errors.New("the input is not a basket: 'items' must be a list")
	}

	basket := checkout.Basket{Items: make([]checkout.BasketItem, 0, len(*raw.Items))}
	for _, item := range *raw.Items {
		if item.SKU == nil {
			return checkout.Basket{}, errors.New("the input is not a basket: every item needs a string 'sku'")
		}
		if item.Quantity == nil {
			return checkout.Basket{}, errors.New("the input is not a basket: every item needs a whole-number 'quantity'")
		}
		basket.Items = append(basket.Items, checkout.BasketItem{SKU: *item.SKU, Quantity: *item.Quantity})
	}
	if raw.Coupon != nil {
		basket.Coupon = *raw.Coupon
	}
	if raw.Shipping != nil {
		basket.Shipping = *raw.Shipping
	}
	return basket, nil
}

// compactJSON encodes a value as compact, single-line JSON. json.Marshal only fails for
// things like channels, functions or NaN; a Receipt and a cliError hold nothing but
// strings and ints, so the error can safely be ignored.
func compactJSON(value any) string {
	data, _ := json.Marshal(value)
	return string(data)
}
