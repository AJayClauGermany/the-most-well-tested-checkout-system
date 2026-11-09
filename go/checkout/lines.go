package checkout

// Receipt lines (rules §4.2): merging the basket's items and pricing a line.

// MergeItems returns one item per distinct SKU, in order of first appearance, with the
// quantities added up (rules §4.2).
//
// TODO: not implemented yet.
func MergeItems(items []BasketItem) []BasketItem {
	return nil
}

// PriceLine builds the receipt line for quantity of product, including its promotion
// discount.
//
// TODO: not implemented yet.
func PriceLine(product Product, quantity int) ReceiptLine {
	return ReceiptLine{}
}
