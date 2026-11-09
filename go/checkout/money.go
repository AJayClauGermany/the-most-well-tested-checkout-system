package checkout

// Money helpers (rules §1 and §7). All amounts are whole cents.

// RoundHalfUp returns percent% of amount, rounded half up to a whole cent (rules §7).
//
// It uses whole numbers only, as the hint in rules §7 suggests: amount × percent is the
// exact answer in hundredths of a cent; adding 50 hundredths and dividing by 100 (Go's
// integer division throws the remainder away, so for non-negative numbers it rounds
// down) is the same as rounding half up. No floating point, so no "0.1 + 0.2" surprises.
//
//	RoundHalfUp(225, 10)  →  (2250 + 50) / 100  →  23      (22.5 rounds UP)
//	RoundHalfUp(224, 10)  →  (2240 + 50) / 100  →  22.9 → 22
func RoundHalfUp(amount, percent int) int {
	return (amount*percent + 50) / 100
}
