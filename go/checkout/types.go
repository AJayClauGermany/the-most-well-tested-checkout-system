// Package checkout calculates the price of a basket of products.
//
// The rules live in spec/checkout-rules.md at the root of the repository.
// This file only describes the *shape* of the data (the "domain types").
// The `json:"..."` tags tell Go's encoding/json package which snake_case
// JSON field belongs to which Go field, so the shared spec files can be
// loaded straight into these types.
//
// All money is a whole number of cents, stored in an int (199 means $1.99).
// We never use float64 for money.
package checkout

// PromotionType says which kind of promotion a product has.
type PromotionType string

// The two kinds of promotion from the rules (§4.2).
const (
	BuyXGetYFree PromotionType = "buy_x_get_y_free"
	BulkDiscount PromotionType = "bulk_discount"
)

// Promotion is a discount attached to a product.
//
// Go has no "one of several shapes" type (like a union or sealed class), so we
// keep it simple: one struct with a Type field plus the fields of every kind.
// Only the fields that belong to Type are used; the others stay 0.
//
//	{"type": "buy_x_get_y_free", "buy": 2, "free": 1}           → uses Buy and Free
//	{"type": "bulk_discount", "min_quantity": 10, "percent": 10} → uses MinQuantity and Percent
type Promotion struct {
	Type PromotionType `json:"type"`

	// Used when Type is BuyXGetYFree: for every Buy+Free items, Free are free.
	Buy  int `json:"buy,omitempty"`
	Free int `json:"free,omitempty"`

	// Used when Type is BulkDiscount: Percent % off once you buy MinQuantity or more.
	MinQuantity int `json:"min_quantity,omitempty"`
	Percent     int `json:"percent,omitempty"`
}

// Product is one entry in the catalog.
type Product struct {
	SKU       string `json:"sku"`
	Name      string `json:"name"`
	UnitPrice int    `json:"unit_price"` // in cents

	// Promotion is nil when the product has no promotion (JSON null).
	Promotion *Promotion `json:"promotion"`
}

// Catalog is the list of products the shop sells.
type Catalog struct {
	Products []Product `json:"products"`
}

// Find looks up a product by its SKU. SKUs are case sensitive.
// The second result is false when no product has that SKU
// (this "value, ok" pattern is the usual Go way to say "maybe not found").
func (c Catalog) Find(sku string) (Product, bool) {
	for _, product := range c.Products {
		if product.SKU == sku {
			return product, true
		}
	}
	return Product{}, false
}

// Shipping methods (rules §5). An empty Basket.Shipping means ShippingStandard.
const (
	ShippingStandard = "standard"
	ShippingExpress  = "express"
)

// BasketItem is one line of the basket: which product and how many.
type BasketItem struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"` // may be invalid (0, -1, 100) – checkout must reject it
}

// Basket is the input to Checkout.
type Basket struct {
	Items []BasketItem `json:"items"`

	// Coupon is optional. A missing coupon, JSON null and "" all end up as ""
	// in Go, which is handy because the rules say they all mean "no coupon".
	Coupon string `json:"coupon,omitempty"`

	// Shipping is optional. "" (missing) means ShippingStandard.
	Shipping string `json:"shipping,omitempty"`
}

// ReceiptLine is one priced line of the receipt (after merging duplicate SKUs).
type ReceiptLine struct {
	SKU          string `json:"sku"`
	Name         string `json:"name"`
	Quantity     int    `json:"quantity"`
	UnitPrice    int    `json:"unit_price"`
	LineTotal    int    `json:"line_total"`
	LineDiscount int    `json:"line_discount"`
}

// Receipt is the result of a successful checkout. All amounts are in cents.
type Receipt struct {
	Lines             []ReceiptLine `json:"lines"`
	Subtotal          int           `json:"subtotal"`
	PromotionDiscount int           `json:"promotion_discount"`
	CouponDiscount    int           `json:"coupon_discount"`
	Shipping          int           `json:"shipping"`
	Tax               int           `json:"tax"`
	Total             int           `json:"total"`
}
