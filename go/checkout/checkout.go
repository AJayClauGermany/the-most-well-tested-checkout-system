package checkout

// Checkout prices a basket using the products in catalog (see spec/checkout-rules.md §4).
//
// It returns either a Receipt and a nil error, or an empty Receipt and an error.
// For an invalid basket the error is a *CheckoutError carrying one of the six codes;
// Checkout returns no other kind of error.
//
// It reads like the hint in rules §4.7: validate, merge lines, price, coupon, shipping,
// tax, total. Each step is a small function in its own file with its own unit tests.
func Checkout(basket Basket, catalog Catalog) (Receipt, error) {
	if err := ValidateItems(basket, catalog); err != nil { // rules §4.1 rows 1-4
		return Receipt{}, err
	}
	method, err := ParseShippingMethod(basket.Shipping) // row 5
	if err != nil {
		return Receipt{}, err
	}

	merged := MergeItems(basket.Items)
	lines := make([]ReceiptLine, 0, len(merged))
	subtotal, promotionDiscount := 0, 0
	for _, item := range merged {
		product, _ := catalog.Find(item.SKU) // ValidateItems made sure every SKU exists
		line := PriceLine(product, item.Quantity)
		lines = append(lines, line)
		subtotal += line.LineTotal
		promotionDiscount += line.LineDiscount
	}
	afterPromotions := subtotal - promotionDiscount

	couponDiscount, err := CouponDiscount(basket.Coupon, afterPromotions) // rows 6-7
	if err != nil {
		return Receipt{}, err
	}
	discountedSubtotal := afterPromotions - couponDiscount

	shipping := ShippingCost(method, discountedSubtotal)
	tax := TaxOn(discountedSubtotal)

	return Receipt{
		Lines:             lines,
		Subtotal:          subtotal,
		PromotionDiscount: promotionDiscount,
		CouponDiscount:    couponDiscount,
		Shipping:          shipping,
		Tax:               tax,
		Total:             discountedSubtotal + shipping + tax,
	}, nil
}
