package review

import "context"


type Usecase interface {
	SubmitReview(ctx context.Context, r *Review) error
	GetResellerReviews(ctx context.Context, resellerID string) ([]*Review, error)
}
