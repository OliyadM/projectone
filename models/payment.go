package models

type PaymentRequest struct {
	BuyerID  string  `json:"buyer_id"`
	SellerID string  `json:"seller_id"`
	ItemID   string  `json:"item_id,omitempty"`
	BundleID string  `json:"bundle_id,omitempty"`
	Amount   float64 `json:"amount"`
}

type PaymentResponse struct {
	ID             string  `json:"id"`
	Amount         float64 `json:"amount"`
	SellerEarnings float64 `json:"seller_earnings"`
	PlatformFee    float64 `json:"platform_fee"`
	Status         string  `json:"status"`
}
