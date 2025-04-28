package admin

import "context"

type Repository interface {
	FetchPlatformMetrics(ctx context.Context) (*Metrics, error)
	GetActiveUsersCount(ctx context.Context) (int, error)
	GetRevenueReport(ctx context.Context) (float64, error)
}
