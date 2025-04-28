package models

type CreateBundleRequest struct {
	Title              string         `json:"title" binding:"required"`
	SampleImage        string         `json:"sample_image"`
	NumberOfItems      int            `json:"number_of_items" binding:"required"`
	Grade              string         `json:"grade" binding:"required"`
	Price              float64        `json:"price" binding:"required"`
	Description        string         `json:"description"`
	SizeRange          string         `json:"size_range"`
	Type               string         `json:"type" binding:"required"`
	EstimatedBreakdown map[string]int `json:"estimated_breakdown"`
	DeclaredRating     int            `json:"declared_rating" binding:"required"`
}

type BundleResponse struct {
	ID                 string         `json:"id"`
	Title              string         `json:"title"`
	SampleImage        string         `json:"sample_image"`
	Quantity           int            `json:"quantity"`
	Grade              string         `json:"grade"`
	Description        string         `json:"description"`
	SizeRange          string         `json:"size_range"`
	Type               string         `json:"type"`
	Price              float64        `json:"price"`
	Status             string         `json:"status"`
	EstimatedBreakdown map[string]int `json:"estimated_breakdown"`
	DeclaredRating     int            `json:"declared_rating"`
	SortingLevel       string         `json:"sorting_level,omitempty"` // optional
	CreatedAt          string         `json:"created_at,omitempty"`
}

type BundleDetailResponse struct {
	Bundle struct {
		ID                 string         `json:"id"`
		Title              string         `json:"title"`
		Description        string         `json:"description"`
		SampleImage        string         `json:"sample_image"`
		Quantity           int            `json:"quantity"`
		Grade              string         `json:"grade"`
		SortingLevel       string         `json:"sorting_level"`
		EstimatedBreakdown map[string]int `json:"estimated_breakdown"`
		Type               string         `json:"type"`
		Price              float64        `json:"price"`
		Status             string         `json:"status"`
		DeclaredRating     int            `json:"declared_rating"`
		RemainingItemCount int            `json:"remaining_item_count"`
	} `json:"bundle"`
	Supplier struct {
		ID     string  `json:"id"`
		Name   string  `json:"name"`
		Rating float64 `json:"rating"`
	} `json:"supplier"`
}
