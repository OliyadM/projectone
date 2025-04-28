package order

import "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"

type OrderStatus string

const (
	Pending   OrderStatus = "pending"
	Shipped   OrderStatus = "shipped"
	Delivered OrderStatus = "delivered"
	Failed    OrderStatus = "failed"

	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCanceled   OrderStatus = "canceled"
)

type Order struct {
	ID          string      `json:"id"`
	ResellerID  string      `json:"reseller_id"`
	SupplierID  string      `json:"supplier_id"`
	BundleID    string      `json:"bundle_id"`
	PlatformFee float64     `json:"platform_fee"`
	ConsumerID  string      `json:"consumer_id"`
	ProductIDs  []string    `json:"product_ids"`
	TotalPrice  float64     `json:"total_price"`
	Status      OrderStatus `json:"status"`
	CreatedAt   string      `json:"created_at"`
}

type PerformanceMetrics struct {
	TotalBundlesListed int `json:"totalBundlesListed"`
	ActiveCount        int `json:"activeCount"`
	SoldCount          int `json:"soldCount"`
}

type DashboardMetrics struct {
	TotalSales         float64            `json:"totalSales"`
	ActiveBundles      []*bundle.Bundle   `json:"activeBundles"`
	PerformanceMetrics PerformanceMetrics `json:"performanceMetrics"`
	Rating             int                `json:"rating"`
	BestSelling        float64            `json:"bestSelling"`
}

type ResellerMetrics struct {
	TotalBoughtBundles int               `json:"totalBoughtBundles"`
	TotalItemsSold     int               `json:"totalItemsSold"`
	Rating             int               `json:"rating"`
	BestSelling        float64           `json:"bestSelling"`
	BoughtBundles      []*bundle.Bundle  `json:"boughtBundles"`
}
