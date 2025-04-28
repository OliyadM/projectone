package order

import (
	"context"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/admin"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/payment"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/warehouse"
)

type Usecase interface {
	PurchaseBundle(ctx context.Context, bundleID, resellerID string) (*Order, *payment.Payment, *warehouse.WarehouseItem, error)
	GetDashboardMetrics(ctx context.Context, supplierID string) (*DashboardMetrics, error)
	GetOrderByID(ctx context.Context, orderID string) (*Order, error)
	GetSoldBundleHistory(ctx context.Context, supplierID string) ([]*Order, map[string]string, error)
	GetResellerMetrics(ctx context.Context, resellerID string) (*ResellerMetrics, error)
	GetAdminDashboardMetrics(ctx context.Context) (*admin.Metrics, error)
}
