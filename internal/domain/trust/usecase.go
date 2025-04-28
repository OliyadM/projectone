package trust

import "context"

type Usecase interface {
	UpdateSupplierTrustScoreOnNewRating(ctx context.Context, supplierID string, declaredRating float64, productRating float64) error
	UpdateResellerTrustScoreOnNewRating(ctx context.Context, resellerID string, declaredRating float64, productRating float64) error
}
