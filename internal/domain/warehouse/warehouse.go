package warehouse

type WarehouseItem struct {
	ID                 string `bson:"_id"`
	ResellerID         string `bson:"reseller_id"`
	BundleID           string `bson:"bundle_id"`
	ProductID          string `bson:"product_id"`
	Status             string `bson:"status"` // listed, skipped, pending
	CreatedAt          string `bson:"created_at"`
	DeclaredRating     int    `bson:"declared_rating"`
	RemainingItemCount int    `bson:"remaining_item_count"`
	Grade              string `bson:"grade"`
	Type               string `bson:"type"`
	Quantity           int    `bson:"quantity"`
	SortingLevel       string `bson:"sorting_level"`
	SampleImage        string `bson:"sample_image"`
}
