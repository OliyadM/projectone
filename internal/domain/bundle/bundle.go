package bundle

import "time"

type SortingLevel string

const (
	Sorted     SortingLevel = "sorted"
	SemiSorted SortingLevel = "semi_sorted"
	Unsorted   SortingLevel = "unsorted"
)

type Bundle struct {
	ID                 string         `bson:"_id"`
	SupplierID         string         `bson:"supplierid"`
	Title              string         `bson:"title"`
	Description        string         `bson:"description"`
	SampleImage        string         `bson:"sampleimage"`
	Quantity           int            `bson:"quantity"`
	Grade              string         `bson:"grade"`
	SortingLevel       SortingLevel   `bson:"sortinglevel"`
	EstimatedBreakdown map[string]int `bson:"estimatedBreakdown,omitempty"`
	Type               string         `bson:"type,omitempty"`
	Price              float64        `bson:"price"`
	Status             string         `bson:"status"`
	CreatedAt          string         `bson:"createdat"`
	DateListed         time.Time      `json:"dateListed" bson:"datelisted"`
	DeclaredRating     int            `bson:"declared_rating"`
	EstimatedItemCount int            `bson:"estimated_item_count"`
	RemainingItemCount int            `bson:"remaining_item_count"`
	SizeRange          string         `bson:"size_range" json:"size_range"`
}
