package models

type RatingRequest struct {
	RateeID  string `json:"ratee_id"`
	Score    int    `json:"score" binding:"required"`
	Comment  string `json:"comment"`
	BundleID string `json:"bundle_id,omitempty"`
}

type RatingResponse struct {
	ID        string `json:"id"`
	RaterID   string `json:"rater_id"`
	RateeID   string `json:"ratee_id"`
	Score     int    `json:"score"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
}
