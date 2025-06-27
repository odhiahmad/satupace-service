package helper

func DeterminePromoType(amount float64) string {
	if amount <= 1.0 {
		return "percent"
	}
	return "fixed"
}
