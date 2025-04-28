package models

type Review struct {
	ID         string `json:"id" bson:"_id"`
	OrderID    string `json:"order_id" bson:"order_id"`
	ProductID  string `json:"product_id" bson:"product_id"`
	UserID     string `json:"user_id" bson:"user_id"`
	ResellerID string `json:"reseller_id" bson:"reseller_id"`
	Rating     int    `json:"rating" bson:"rating"` // 1-100
	Comment    string `json:"comment" bson:"comment"`
	CreatedAt  string `json:"created_at" bson:"created_at"`
}

type CreateReviewRequest struct {
	OrderID   string `json:"order_id" binding:"required"`
	ProductID string `json:"product_id" binding:"required"`
	Rating    int    `json:"rating" binding:"required,min=1,max=100"`
	Comment   string `json:"comment"`
}
