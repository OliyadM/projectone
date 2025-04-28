package OrderUsecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/admin"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/order"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/payment"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/warehouse"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/product"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



func NewOrderUsecase(bRepo bundle.Repository, oRepo order.Repository, wRepo warehouse.Repository, pRepo payment.Repository, uRepo user.Repository, prodRepo product.Repository) *orderUseCaseImpl {
	return &orderUseCaseImpl{
		bundleRepo:    bRepo,
		orderRepo:     oRepo,
		warehouseRepo: wRepo,
		paymentRepo:   pRepo,
		userRepo:      uRepo,
		prodRepo:      prodRepo,
	}
}

func simulateStripePayment(_ float64) (string, error) {
	time.Sleep(500 * time.Millisecond)
	return fmt.Sprintf("ch_%d", rand.Intn(1000000)), nil
}

func processPayment(total float64) (fee float64, net float64, err error) {
	fee = total * 0.02
	net = total - fee
	_, err = simulateStripePayment(total)
	return
}

func (uc *orderUseCaseImpl) PurchaseBundle(ctx context.Context, bundleID, resellerID string) (*order.Order, *payment.Payment, *warehouse.WarehouseItem, error) {
	b, err := uc.bundleRepo.GetBundleByID(ctx, bundleID)
	if err != nil {
		return nil, nil, nil, err
	}

	availables, err := uc.bundleRepo.ListAvailableBundles(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	found := false
	for _, available := range availables {
		if available.ID == b.ID {
			found = true
			break
		}
	}
	if !found {
		return nil, nil, nil, errors.New("bundle not available")
	}

	if b.SupplierID == resellerID {
		return nil, nil, nil, errors.New("reseller cannot purchase their own bundle")
	}

	fee, net, err := processPayment(b.Price)
	if err != nil {
		return nil, nil, nil, err
	}

	order := &order.Order{
		ID:          primitive.NewObjectID().Hex(),
		BundleID:    b.ID,
		ResellerID:  resellerID,
		SupplierID:  b.SupplierID,
		TotalPrice:  b.Price,
		PlatformFee: fee,
		Status:      order.OrderStatusCompleted,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}
	if err := uc.orderRepo.CreateOrder(ctx, order); err != nil {
		return nil, nil, nil, err
	}

	payment := &payment.Payment{
		FromUserID:    resellerID,
		ToUserID:      b.SupplierID,
		Amount:        b.Price,
		PlatformFee:   fee,
		SellerEarning: net,
		Status:        "Paid",
		ReferenceID:   b.ID,
		Type:          payment.B2B,
		CreatedAt:     time.Now().Format(time.RFC3339),
	}
	if err := uc.paymentRepo.RecordPayment(ctx, payment); err != nil {
		return nil, nil, nil, err
	}

	if err := uc.bundleRepo.MarkAsPurchased(ctx, b.ID, resellerID); err != nil {
		return nil, nil, nil, err
	}

	warehouseItem := &warehouse.WarehouseItem{
		ID:                 primitive.NewObjectID().Hex(),
		BundleID:           b.ID,
		ResellerID:         resellerID,
		Status:             "pending",
		DeclaredRating:     b.DeclaredRating,
		RemainingItemCount: b.RemainingItemCount,
		Grade:              b.Grade,
		Type:               b.Type,
		Quantity:           b.Quantity,
		SortingLevel:       string(b.SortingLevel),
		SampleImage:        b.SampleImage,
		CreatedAt:          time.Now().Format(time.RFC3339),
	}
	if err := uc.warehouseRepo.AddItem(ctx, warehouseItem); err != nil {
		return nil, nil, nil, err
	}

	go func(itemID string) {
		time.Sleep(3 * time.Minute)
		if err := uc.warehouseRepo.MarkItemAsListed(context.Background(), itemID); err != nil {
			fmt.Println("Failed to mark warehouse item as listed:", err)
		}
	}(warehouseItem.ID)

	return order, payment, warehouseItem, nil
}

func (uc *orderUseCaseImpl) GetDashboardMetrics(ctx context.Context, supplierID string) (*order.DashboardMetrics, error) {
	bundles, err := uc.bundleRepo.ListBundles(ctx, supplierID)
	if err != nil {
		return nil, err
	}

	totalSales := 0.0
	activeCount := 0
	soldCount := 0
	bestSelling := 0.0
	var activeBundles []*bundle.Bundle

	userData, err := uc.userRepo.GetByID(ctx, supplierID)
	if err != nil {
		return nil, err
	}

	for _, b := range bundles {
		if b.Status == "purchased" {
			totalSales += b.Price
			soldCount++
			if b.Price > bestSelling {
				bestSelling = b.Price
			}
		} else if b.Status == "available" {
			activeCount++
			activeBundles = append(activeBundles, b)
		}
	}

	sort.Slice(activeBundles, func(i, j int) bool {
		return activeBundles[i].DateListed.After(activeBundles[j].DateListed)
	})

	return &order.DashboardMetrics{
		TotalSales:         totalSales,
		ActiveBundles:      activeBundles,
		PerformanceMetrics: order.PerformanceMetrics{TotalBundlesListed: len(bundles), ActiveCount: activeCount, SoldCount: soldCount},
		Rating:             userData.TrustScore,
		BestSelling:        bestSelling,
	}, nil
}

func (uc *orderUseCaseImpl) GetResellerMetrics(ctx context.Context, resellerID string) (*order.ResellerMetrics, error) {
	fmt.Printf("\n🔍 Starting GetResellerMetrics for reseller: %s\n", resellerID)

	// Get purchased bundles
	bundles, err := uc.bundleRepo.ListPurchasedByReseller(ctx, resellerID)
	if err != nil {
		fmt.Printf("❌ Error getting bundles: %v\n", err)
		return nil, err
	}
	fmt.Printf("📦 Found %d purchased bundles\n", len(bundles))

	// Get reseller info
	reseller, err := uc.userRepo.GetByID(ctx, resellerID)
	if err != nil {
		fmt.Printf("❌ Error getting reseller info: %v\n", err)
		return nil, err
	}
	fmt.Printf("👤 Reseller found: %s\n", reseller.Username)
	fmt.Printf("📊 Reseller Trust Data:\n")
	fmt.Printf("  - Trust Score: %d\n", reseller.TrustScore)
	fmt.Printf("  - Trust Rated Count: %d\n", reseller.TrustRatedCount)
	fmt.Printf("  - Trust Total Error: %f\n", reseller.TrustTotalError)
	fmt.Printf("  - Is Blacklisted: %v\n", reseller.IsBlacklisted)

	// Get sold products directly from product collection
	soldProducts, err := uc.prodRepo.GetSoldProductsByReseller(ctx, resellerID)
	if err != nil {
		fmt.Printf("❌ Error getting sold products: %v\n", err)
		return nil, err
	}
	fmt.Printf("🛍️ Found %d sold products\n", len(soldProducts))

	// Find best selling item
	bestSelling := 0.0
	for _, product := range soldProducts {
		if product.Price > bestSelling {
			bestSelling = product.Price
			fmt.Printf("💰 New best selling item found: %s with price %.2f\n", product.Title, product.Price)
		}
	}
	fmt.Printf("💰 Best selling item price: %.2f\n", bestSelling)

	metrics := &order.ResellerMetrics{
		TotalBoughtBundles: len(bundles),
		TotalItemsSold:     len(soldProducts),
		Rating:             reseller.TrustScore,
		BestSelling:        bestSelling,
		BoughtBundles:      bundles,
	}

	fmt.Printf("📊 Final Metrics:\n")
	fmt.Printf("  - Total Bought Bundles: %d\n", metrics.TotalBoughtBundles)
	fmt.Printf("  - Total Items Sold: %d\n", metrics.TotalItemsSold)
	fmt.Printf("  - Rating (Trust Score): %d\n", metrics.Rating)
	fmt.Printf("  - Best Selling: %.2f\n", metrics.BestSelling)
	fmt.Printf("✅ GetResellerMetrics completed\n\n")

	return metrics, nil
}

func (uc *orderUseCaseImpl) GetSoldBundleHistory(ctx context.Context, supplierID string) ([]*order.Order, map[string]string, error) {
	log.Printf("Getting sold bundle history for supplier: %s", supplierID)
	orders, err := uc.orderRepo.GetOrdersBySupplier(ctx, supplierID)
	if err != nil {
		log.Printf("Error getting orders: %v", err)
		return nil, nil, err
	}

	userNames := make(map[string]string)
	for _, order := range orders {
		if order.ResellerID != "" {
			user, err := uc.userRepo.GetByID(ctx, order.ResellerID)
			if err != nil {
				log.Printf("Error getting user name for ID %s: %v", order.ResellerID, err)
				continue
			}
			userNames[order.ResellerID] = user.Username
		}
	}

	log.Printf("Found %d orders for supplier %s", len(orders), supplierID)
	return orders, userNames, nil
}

func (uc *orderUseCaseImpl) GetAdminDashboardMetrics(ctx context.Context) (*admin.Metrics, error) {
	totalBundles, err := uc.bundleRepo.CountBundles(ctx)
	if err != nil {
		return nil, err
	}

	totalUsers, err := uc.userRepo.CountActiveUsers(ctx)
	if err != nil {
		return nil, err
	}

	totalSales, platformFees, err := uc.paymentRepo.GetAllPlatformFees(ctx)
	if err != nil {
		return nil, err
	}

	skippedClothes, err := uc.warehouseRepo.CountByStatus(ctx, "skipped")
	if err != nil {
		return nil, err
	}

	return &admin.Metrics{
		TotalBundles:    totalBundles,
		TotalUsers:      totalUsers,
		TotalSales:      totalSales,
		RevenueFromFees: platformFees,
		SkippedClothes:  skippedClothes,
	}, nil
}

func (uc *orderUseCaseImpl) GetOrderByID(ctx context.Context, orderID string) (*order.Order, error) {
	return uc.orderRepo.GetOrderByID(ctx, orderID)
}

func (uc *orderUseCaseImpl) PurchaseProduct(ctx context.Context, productID, consumerID string, totalPrice float64) (*order.Order, *payment.Payment, error) {
	// Get product details to get reseller ID
	prod, err := uc.prodRepo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get product details: %w", err)
	}

	// Calculate platform fee and net amount
	platformFee := totalPrice * 0.02
	netAmount := totalPrice - platformFee

	// Create order
	order := &order.Order{
		ID:          primitive.NewObjectID().Hex(),
		ResellerID:  prod.ResellerID.Hex(),
		SupplierID:  "", // Not needed for product orders
		BundleID:    "", // Not needed for product orders
		ConsumerID:  consumerID,
		ProductIDs:  []string{productID},
		TotalPrice:  totalPrice,
		PlatformFee: platformFee,
		Status:      order.OrderStatusCompleted,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	if err := uc.orderRepo.CreateOrder(ctx, order); err != nil {
		return nil, nil, err
	}

	// Create payment record
	payment := &payment.Payment{
		FromUserID:    consumerID,
		ToUserID:      prod.ResellerID.Hex(),
		Amount:        totalPrice,
		PlatformFee:   platformFee,
		SellerEarning: netAmount,
		Status:        "Paid",
		ReferenceID:   productID,
		Type:          payment.B2C,
		CreatedAt:     time.Now().Format(time.RFC3339),
	}

	if err := uc.paymentRepo.RecordPayment(ctx, payment); err != nil {
		return nil, nil, err
	}

	return order, payment, nil
}
func (uc *orderUseCaseImpl) GetOrdersByReseller(ctx context.Context, resellerID string) ([]*order.Order, map[string]string, error) {
	orders, err := uc.orderRepo.GetOrdersByReseller(ctx, resellerID)
	if err != nil {
		return nil, nil, err
	}

	userNames := make(map[string]string)
	for _, order := range orders {
		if len(order.ProductIDs) > 0 { // Sold order
			if order.ConsumerID != "" {
				user, err := uc.userRepo.GetByID(ctx, order.ConsumerID)
				if err == nil && user != nil {
					userNames[order.ConsumerID] = user.Username
				}
			}
		} else if order.BundleID != "" { // Bought order
			if order.SupplierID != "" {
				user, err := uc.userRepo.GetByID(ctx, order.SupplierID)
				if err == nil && user != nil {
					userNames[order.SupplierID] = user.Username
				}
			}
		}
	}

	return orders, userNames, nil
}

func (uc *orderUseCaseImpl) GetOrdersByConsumer(ctx context.Context, consumerID string) ([]*order.Order, map[string]string, map[string]string, error) {
	fmt.Printf("🔍 Getting orders for consumer: %s\n", consumerID)
	
	orders, err := uc.orderRepo.GetOrdersByConsumer(ctx, consumerID)
	if err != nil {
		fmt.Printf("❌ Error getting consumer orders: %v\n", err)
		return nil, nil, nil, fmt.Errorf("failed to get consumer orders: %w", err)
	}

	userNames := make(map[string]string)
	productNames := make(map[string]string)

	// Get unique product IDs
	productIDs := make(map[string]bool)
	for _, order := range orders {
		if order.ResellerID != "" {
			user, err := uc.userRepo.GetByID(ctx, order.ResellerID)
			if err == nil && user != nil {
				userNames[order.ResellerID] = user.Username
			}
		}
		for _, productID := range order.ProductIDs {
			productIDs[productID] = true
		}
	}

	// Get product names
	for productID := range productIDs {
		product, err := uc.prodRepo.GetProductByID(ctx, productID)
		if err == nil && product != nil {
			productNames[productID] = product.Title
		}
	}

	fmt.Printf("✅ Found %d orders for consumer\n", len(orders))
	return orders, userNames, productNames, nil
}

