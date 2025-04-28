package models

type ProductListingRequest struct {
	ID          string  `json:"id,omitempty"`
	Photo       string  `json:"photo" binding:"required"`
	Title       string  `json:"title" binding:"required"`
	Size        string  `json:"size" binding:"required"`
	Grade       string  `json:"grade"`
	Price       float64 `json:"price" binding:"required"`
	Description string  `json:"description"`
	BundleID    string  `json:"bundle_id"`
}

type ProductResponse struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Photo       string  `json:"photo"`
	Grade       string  `json:"grade"`
	Size        string  `json:"size"`
	Status      string  `json:"status"`
	SellerID    string  `json:"seller_id"`
	Rating      float64 `json:"rating"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	BundleID    string  `json:"bundle_id"`
}
