package review

import "context"

type Repository interface {
	CreateReview(ctx context.Context, r *Review) error
	GetReviewByUserAndProduct(ctx context.Context, userID, productID string) (*Review, error)
	GetReviewsByReseller(ctx context.Context, resellerID string) ([]*Review, error)
}
