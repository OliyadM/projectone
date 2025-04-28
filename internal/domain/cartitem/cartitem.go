package cartitem

import "time"

type CartItem struct {
	ID        string    `bson:"_id" json:"id"`        // UUID
	UserID    string    `bson:"userid" json:"userid"` // consumer
	ListingID string    `bson:"listingid" json:"listingid"`
	Title     string    `bson:"title" json:"title"`
	Price     float64   `bson:"price" json:"price"`
	ImageURL  string    `bson:"imageurl" json:"imageurl"`
	Grade     string    `bson:"grade" json:"grade"` // Reseller's assigned rating (e.g., 93)
	CreatedAt time.Time `bson:"createdat" json:"createdat"`
}
