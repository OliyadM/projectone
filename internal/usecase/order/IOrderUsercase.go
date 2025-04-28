package OrderUsecase

import (
	"context"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/order"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/payment"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/warehouse"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/product"
)

type OrderUseCase interface {
	PurchaseBundle(ctx context.Context, bundleID, resellerID string) (*order.Order, *payment.Payment, *warehouse.WarehouseItem, error)
	GetDashboardMetrics(ctx context.Context, supplierID string) (*order.DashboardMetrics, error)
	GetOrderByID(ctx context.Context, orderID string) (*order.Order, error)
	GetResellerMetrics(ctx context.Context, resellerID string) (*order.ResellerMetrics, error)
	GetSoldBundleHistory(ctx context.Context, supplierID string) ([]*order.Order, map[string]string, error)
	GetOrdersByReseller(ctx context.Context, resellerID string) ([]*order.Order, map[string]string, error)
	GetOrdersByConsumer(ctx context.Context, consumerID string) ([]*order.Order, map[string]string, map[string]string, error)
	PurchaseProduct(ctx context.Context, productID, consumerID string, totalPrice float64) (*order.Order, *payment.Payment, error)
}

type orderUseCaseImpl struct {
	bundleRepo    bundle.Repository
	orderRepo     order.Repository
	warehouseRepo warehouse.Repository
	paymentRepo   payment.Repository
	userRepo      user.Repository
	prodRepo      product.Repository
}

