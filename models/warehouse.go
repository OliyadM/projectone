package models

type WarehouseItemResponse struct {
	ID             string  `json:"id"`
	ResellerID     string  `json:"reseller_id"`
	BundleID       string  `json:"bundle_id"`
	Status         string  `json:"status"`
	CreatedAt      string  `json:"created_at"`
	Title          string  `json:"title"`
	SampleImage    string  `json:"sample_image"`
	DeclaredRating float64 `json:"declared_rating"`
	RemainingItems int     `json:"remaining_items"`
}
