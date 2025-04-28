package rating

type Rating struct {
	ID         string
	ResellerID string
	SupplierID string
	Score      int // e.g., 1–5
	Comment    string
	SkipRate   float64
	CreatedAt  string
}
