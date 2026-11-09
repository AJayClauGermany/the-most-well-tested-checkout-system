package checkout

// Receipt lines (rules §4.2): merging the basket's items and pricing a line.

// MergeItems returns one item per distinct SKU, in order of first appearance, with the
// quantities added up (rules §4.2). The basket's own items are left untouched.
//
//	[APPLE x1, PEN x2, APPLE x3]  →  [APPLE x4, PEN x2]
func MergeItems(items []BasketItem) []BasketItem {
	merged := make([]BasketItem, 0, len(items))
	position := map[string]int{} // SKU → where it is in merged
	for _, item := range items {
		if i, seen := position[item.SKU]; seen {
			merged[i].Quantity += item.Quantity
			continue
		}
		position[item.SKU] = len(merged)
		merged = append(merged, item)
	}
	return merged
}

// PriceLine builds the receipt line for quantity of product, including its promotion
// discount.
func PriceLine(product Product, quantity int) ReceiptLine {
	return ReceiptLine{
		SKU:          product.SKU,
		Name:         product.Name,
		Quantity:     quantity,
		UnitPrice:    product.UnitPrice,
		LineTotal:    product.UnitPrice * quantity,
		LineDiscount: PromotionDiscount(product.Promotion, quantity, product.UnitPrice),
	}
}
