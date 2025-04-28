package review

import (
	"errors"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID         string    `bson:"_id" json:"id"`
	OrderID    string    `bson:"order_id" json:"order_id"`
	ProductID  string    `bson:"product_id" json:"product_id"`
	UserID     string    `bson:"user_id" json:"user_id"`
	ResellerID string    `bson:"reseller_id" json:"reseller_id"`
	Rating     int       `bson:"rating" json:"rating"`
	Comment    string    `bson:"comment" json:"comment"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
}

func (r *Review) Validate() error {
	if r.OrderID == "" || r.ProductID == "" || r.UserID == "" {
		return errors.New("order_id, product_id, and user_id are required")
	}
	if r.Rating < 0 || r.Rating > 5 {
		return errors.New("rating must be between 0 and 5")
	}
	return nil
}

func (r *Review) CalculateTrustImpact(declaredRating float64) float64 {
	return math.Abs(float64(r.Rating) - declaredRating)
}

func (r *Review) UpdateRating(newRating int) error {
	if newRating < 0 || newRating > 5 {
		return errors.New("rating must be between 0 and 5")
	}
	r.Rating = newRating
	return nil
}

func NewReview(orderID, productID, userID, resellerID string, rating int, comment string) *Review {
	return &Review{
		ID:         primitive.NewObjectID().Hex(),
		OrderID:    orderID,
		ProductID:  productID,
		UserID:     userID,
		ResellerID: resellerID,
		Rating:     rating,
		Comment:    comment,
		CreatedAt:  time.Now(),
	}
}
