package models

type OrderRequest struct {
	ProductID string `json:"product_id"`
}

type OrderResponse struct {
	ID                    string  `json:"id"`
	BuyerID               string  `json:"buyer_id"`
	ProductTitle          string  `json:"product_title"`
	Status                string  `json:"status"`
	TotalPrice            float64 `json:"total_price"`
	CreatedAt             string  `json:"created_at"`
	ImageURL              string  `json:"image_url"`
	EstimatedDeliveryTime string  `json:"estimated_delivery_time"`
}
