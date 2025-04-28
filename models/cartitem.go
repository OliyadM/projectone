package models

type CreateCartItemRequest struct {
	ListingID string `json:"listing_id" binding:"required"`
}

// CartItemResponse remains the same.
type CartItemResponse struct {
	ID        string  `json:"id"`
	ListingID string  `json:"listing_id"`
	Title     string  `json:"title"`
	Price     float64 `json:"price"`
	ImageURL  string  `json:"image_url"`
	Grade     string  `json:"grade"`
	CreatedAt string  `json:"created_at"`
	Rating    float64 `json:"rating"`
}

type CheckoutItemResponse struct {
	ListingID string  `json:"listingId"`
	Title     string  `json:"title"`
	Price     float64 `json:"price"`
	SellerID  string  `json:"sellerId"`
	Status    string  `json:"status"` // "available", "sold"
}

type CheckoutResponse struct {
	TotalAmount float64                `json:"totalAmount"`
	Items       []CheckoutItemResponse `json:"items"`
	PlatformFee float64                `json:"platformFee"` // 2%
	NetPayable  float64                `json:"netPayable"`  // Total - fee
}

type PaymentRecord struct {
	ID             string  `json:"id"`
	BuyerID        string  `json:"buyerId"`
	SellerID       string  `json:"sellerId"`
	ProductID      string  `json:"productId"`
	Amount         float64 `json:"amount"`
	PlatformFee    float64 `json:"platformFee"`    // 2% fee for Admin
	SellerEarnings float64 `json:"sellerEarnings"` // Amount - PlatformFee
	Status         string  `json:"status"`         // "paid"
}
