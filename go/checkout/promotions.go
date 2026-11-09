package checkout

// Promotions (rules §4.2).

// PromotionDiscount returns the discount a promotion gives one receipt line, in cents.
// It is 0 when promotion is nil (the product has no promotion).
func PromotionDiscount(promotion *Promotion, quantity, unitPrice int) int {
	if promotion == nil {
		return 0
	}

	switch promotion.Type {
	case BuyXGetYFree:
		// For every complete group of Buy+Free items, Free of them cost nothing.
		groupSize := promotion.Buy + promotion.Free
		if groupSize <= 0 {
			return 0 // a broken promotion gives no discount rather than dividing by zero
		}
		completeGroups := quantity / groupSize // integer division rounds down
		freeItems := completeGroups * promotion.Free
		return freeItems * unitPrice

	case BulkDiscount:
		// Percent% off the whole line, but only from MinQuantity items.
		if quantity < promotion.MinQuantity {
			return 0
		}
		return RoundHalfUp(unitPrice*quantity, promotion.Percent)
	}

	return 0 // a promotion type we do not know gives no discount
}
