package OrderUsecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/order"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/payment"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/product"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/warehouse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock Repositories
type MockBundleRepo struct {
	mock.Mock
}

func (m *MockBundleRepo) CreateBundle(ctx context.Context, b *bundle.Bundle) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockBundleRepo) GetBundleByID(ctx context.Context, id string) (*bundle.Bundle, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bundle.Bundle), args.Error(1)
}

func (m *MockBundleRepo) ListBundles(ctx context.Context, supplierID string) ([]*bundle.Bundle, error) {
	args := m.Called(ctx, supplierID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bundle.Bundle), args.Error(1)
}

func (m *MockBundleRepo) ListAvailableBundles(ctx context.Context) ([]*bundle.Bundle, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bundle.Bundle), args.Error(1)
}

func (m *MockBundleRepo) ListPurchasedByReseller(ctx context.Context, resellerID string) ([]*bundle.Bundle, error) {
	args := m.Called(ctx, resellerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bundle.Bundle), args.Error(1)
}

func (m *MockBundleRepo) UpdateBundleStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockBundleRepo) MarkAsPurchased(ctx context.Context, bundleID string, resellerID string) error {
	args := m.Called(ctx, bundleID, resellerID)
	return args.Error(0)
}

func (m *MockBundleRepo) DeleteBundle(ctx context.Context, bundleID string) error {
	args := m.Called(ctx, bundleID)
	return args.Error(0)
}

func (m *MockBundleRepo) UpdateBundle(ctx context.Context, id string, updatedData map[string]interface{}) error {
	args := m.Called(ctx, id, updatedData)
	return args.Error(0)
}

func (m *MockBundleRepo) DecreaseBundleQuantity(ctx context.Context, bundleID string) error {
	args := m.Called(ctx, bundleID)
	return args.Error(0)
}

type MockOrderRepo struct {
	mock.Mock
}

func (m *MockOrderRepo) CreateOrder(ctx context.Context, o *order.Order) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

func (m *MockOrderRepo) GetOrdersByConsumer(ctx context.Context, consumerID string) ([]*order.Order, error) {
	args := m.Called(ctx, consumerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*order.Order), args.Error(1)
}

func (m *MockOrderRepo) GetOrderByID(ctx context.Context, orderID string) (*order.Order, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order.Order), args.Error(1)
}

func (m *MockOrderRepo) UpdateOrderStatus(ctx context.Context, orderID string, status order.OrderStatus) error {
	args := m.Called(ctx, orderID, status)
	return args.Error(0)
}

func (m *MockOrderRepo) DeleteOrder(ctx context.Context, orderID string) error {
	args := m.Called(ctx, orderID)
	return args.Error(0)
}

func (m *MockOrderRepo) GetOrdersBySupplier(ctx context.Context, supplierID string) ([]*order.Order, error) {
	args := m.Called(ctx, supplierID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*order.Order), args.Error(1)
}

func (m *MockOrderRepo) GetOrdersByReseller(ctx context.Context, resellerID string) ([]*order.Order, error) {
	args := m.Called(ctx, resellerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*order.Order), args.Error(1)
}

type MockWarehouseRepo struct {
	mock.Mock
}

func (m *MockWarehouseRepo) AddItem(ctx context.Context, item *warehouse.WarehouseItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockWarehouseRepo) GetItemsByReseller(ctx context.Context, resellerID string) ([]*warehouse.WarehouseItem, error) {
	args := m.Called(ctx, resellerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*warehouse.WarehouseItem), args.Error(1)
}

func (m *MockWarehouseRepo) GetItemsByBundle(ctx context.Context, bundleID string) ([]*warehouse.WarehouseItem, error) {
	args := m.Called(ctx, bundleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*warehouse.WarehouseItem), args.Error(1)
}

func (m *MockWarehouseRepo) MarkItemAsListed(ctx context.Context, itemID string) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}

func (m *MockWarehouseRepo) MarkItemAsSkipped(ctx context.Context, itemID string) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}
func (m *MockWarehouseRepo) CountByStatus(ctx context.Context, status string) (int, error) {
	args := m.Called(ctx, status)
	return args.Int(0), args.Error(1)
}

func (m *MockWarehouseRepo) DeleteItem(ctx context.Context, itemID string) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}

func (m *MockWarehouseRepo) HasResellerReceivedBundle(ctx context.Context, resellerID string, bundleID string) (bool, error) {
	args := m.Called(ctx, resellerID, bundleID)
	return args.Bool(0), args.Error(1)
}

type MockPaymentRepo struct {
	mock.Mock
}

func (m *MockPaymentRepo) RecordPayment(ctx context.Context, p *payment.Payment) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockPaymentRepo) GetPaymentsByUser(ctx context.Context, userID string) ([]*payment.Payment, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*payment.Payment), args.Error(1)
}
func (m *MockPaymentRepo) GetAllPlatformFees(ctx context.Context) (float64, float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Get(1).(float64), args.Error(2)
}

func (m *MockPaymentRepo) GetPaymentsByType(ctx context.Context, userID string, pType payment.PaymentType) ([]*payment.Payment, error) {
	args := m.Called(ctx, userID, pType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*payment.Payment), args.Error(1)
}

func (m *MockPaymentRepo) CreatePayment(ctx context.Context, p *payment.Payment) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepo) CreateUser(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}
func (m *MockUserRepo) CountActiveUsers(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepo) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepo) ListUsersByRole(ctx context.Context, role user.Role) ([]*user.User, error) {
	args := m.Called(ctx, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.User), args.Error(1)
}

func (m *MockUserRepo) UpdateUser(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockUserRepo) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepo) FindUserByUsername(ctx context.Context, username string) (*user.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepo) UpdateTrustData(ctx context.Context, user *user.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) GetBlacklistedUsers(ctx context.Context) ([]*user.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.User), args.Error(1)
}
func (m *MockBundleRepo) CountBundles(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

type MockProductRepo struct {
	mock.Mock
}

func (m *MockProductRepo) GetProductByTitle(ctx context.Context, title string) (*product.Product, error) {
	args := m.Called(ctx, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Product), args.Error(1)
}
func (m *MockProductRepo) AddProduct(ctx context.Context, p *product.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockProductRepo) GetProductByID(ctx context.Context, id string) (*product.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Product), args.Error(1)
}

func (m *MockProductRepo) ListProducts(ctx context.Context) ([]*product.Product, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*product.Product), args.Error(1)
}

func (m *MockProductRepo) UpdateProduct(ctx context.Context, id string, updatedData map[string]interface{}) error {
	args := m.Called(ctx, id, updatedData)
	return args.Error(0)
}

func (m *MockProductRepo) DeleteProduct(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepo) GetProductsByBundleID(ctx context.Context, bundleID string) ([]*product.Product, error) {
	args := m.Called(ctx, bundleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*product.Product), args.Error(1)
}

func (m *MockProductRepo) GetSoldProductsByReseller(ctx context.Context, resellerID string) ([]*product.Product, error) {
	args := m.Called(ctx, resellerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*product.Product), args.Error(1)
}

func (m *MockProductRepo) ListAvailableProducts(ctx context.Context, page, limit int) ([]*product.Product, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*product.Product), args.Error(1)
}

func (m *MockProductRepo) ListProductsByReseller(ctx context.Context, resellerID string, page, limit int) ([]*product.Product, error) {
	args := m.Called(ctx, resellerID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*product.Product), args.Error(1)
}
func (m *MockBundleRepo) GetBundleByTitle(ctx context.Context, title string) (*bundle.Bundle, error) {
	args := m.Called(ctx, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bundle.Bundle), args.Error(1)
}

// OrderUsecaseTestSuite is the test suite for order usecase
type OrderUsecaseTestSuite struct {
	suite.Suite
	ctx           context.Context
	bundleRepo    *MockBundleRepo
	orderRepo     *MockOrderRepo
	warehouseRepo *MockWarehouseRepo
	paymentRepo   *MockPaymentRepo
	userRepo      *MockUserRepo
	productRepo   *MockProductRepo
	useCase       *orderUseCaseImpl
}

// SetupTest runs before each test
func (suite *OrderUsecaseTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.bundleRepo = new(MockBundleRepo)
	suite.orderRepo = new(MockOrderRepo)
	suite.warehouseRepo = new(MockWarehouseRepo)
	suite.paymentRepo = new(MockPaymentRepo)
	suite.userRepo = new(MockUserRepo)
	suite.productRepo = new(MockProductRepo)
	suite.useCase = NewOrderUsecase(
		suite.bundleRepo,
		suite.orderRepo,
		suite.warehouseRepo,
		suite.paymentRepo,
		suite.userRepo,
		suite.productRepo,
	)
}

// TearDownTest runs after each test
func (suite *OrderUsecaseTestSuite) TearDownTest() {
	suite.bundleRepo.AssertExpectations(suite.T())
	suite.orderRepo.AssertExpectations(suite.T())
	suite.warehouseRepo.AssertExpectations(suite.T())
	suite.paymentRepo.AssertExpectations(suite.T())
	suite.userRepo.AssertExpectations(suite.T())
	suite.productRepo.AssertExpectations(suite.T())
}

// TestOrderUsecaseTestSuite runs all the tests in the suite
func TestOrderUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(OrderUsecaseTestSuite))
}

// TestNewOrderUsecase tests the constructor
func (suite *OrderUsecaseTestSuite) TestNewOrderUsecase() {
	// Act
	uc := NewOrderUsecase(
		suite.bundleRepo,
		suite.orderRepo,
		suite.warehouseRepo,
		suite.paymentRepo,
		suite.userRepo,
		suite.productRepo,
	)

	// Assert
	assert.NotNil(suite.T(), uc)
	assert.Equal(suite.T(), suite.orderRepo, uc.orderRepo)
	assert.Equal(suite.T(), suite.bundleRepo, uc.bundleRepo)
	assert.Equal(suite.T(), suite.warehouseRepo, uc.warehouseRepo)
	assert.Equal(suite.T(), suite.paymentRepo, uc.paymentRepo)
	assert.Equal(suite.T(), suite.userRepo, uc.userRepo)
	assert.Equal(suite.T(), suite.productRepo, uc.prodRepo)
}

// TestPurchaseBundle tests the PurchaseBundle method
func (suite *OrderUsecaseTestSuite) TestPurchaseBundle() {
	tests := []struct {
		name         string
		bundleID     string
		resellerID   string
		mockBundle   *bundle.Bundle
		mockError    error
		expectError  bool
		errorMessage string
	}{
		{
			name:       "Success - Valid purchase",
			bundleID:   "bundle1",
			resellerID: "reseller1",
			mockBundle: &bundle.Bundle{
				ID:         "bundle1",
				SupplierID: "supplier1",
				Price:      100.0,
				Status:     "available",
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:         "Error - Bundle not found",
			bundleID:     "nonexistent",
			resellerID:   "reseller1",
			mockBundle:   nil,
			mockError:    errors.New("bundle not found"),
			expectError:  true,
			errorMessage: "bundle not found",
		},
		{
			name:       "Error - Self purchase attempt",
			bundleID:   "bundle1",
			resellerID: "supplier1",
			mockBundle: &bundle.Bundle{
				ID:         "bundle1",
				SupplierID: "supplier1",
				Price:      100.0,
				Status:     "available",
			},
			mockError:    nil,
			expectError:  true,
			errorMessage: "reseller cannot purchase their own bundle",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.bundleRepo.On("GetBundleByID", suite.ctx, tt.bundleID).Return(tt.mockBundle, tt.mockError)
			if tt.mockBundle != nil {
				suite.bundleRepo.On("ListAvailableBundles", suite.ctx).Return([]*bundle.Bundle{tt.mockBundle}, nil)
			}

			if !tt.expectError {
				suite.orderRepo.On("CreateOrder", suite.ctx, mock.AnythingOfType("*order.Order")).Return(nil)
				suite.paymentRepo.On("RecordPayment", suite.ctx, mock.AnythingOfType("*payment.Payment")).Return(nil)
				suite.bundleRepo.On("MarkAsPurchased", suite.ctx, tt.bundleID, tt.resellerID).Return(nil)
				suite.warehouseRepo.On("AddItem", suite.ctx, mock.AnythingOfType("*warehouse.WarehouseItem")).Return(nil)
			}

			// Act
			order, payment, warehouseItem, err := suite.useCase.PurchaseBundle(suite.ctx, tt.bundleID, tt.resellerID)

			// Assert
			if tt.expectError {
				assert.Error(suite.T(), err)
				assert.Contains(suite.T(), err.Error(), tt.errorMessage)
				assert.Nil(suite.T(), order)
				assert.Nil(suite.T(), payment)
				assert.Nil(suite.T(), warehouseItem)
			} else {
				assert.NoError(suite.T(), err)
				assert.NotNil(suite.T(), order)
				assert.NotNil(suite.T(), payment)
				assert.NotNil(suite.T(), warehouseItem)
			}
		})
	}
}

// TestGetDashboardMetrics tests the GetDashboardMetrics method
func (suite *OrderUsecaseTestSuite) TestGetDashboardMetrics() {
	tests := []struct {
		name           string
		supplierID     string
		mockBundles    []*bundle.Bundle
		mockUser       *user.User
		mockError      error
		expectError    bool
		expectedSales  float64
		expectedCounts order.PerformanceMetrics
		expectedRating int
		expectedBest   float64
	}{
		{
			name:       "Success - With bundles",
			supplierID: "supplier1",
			mockBundles: []*bundle.Bundle{
				{
					ID:         "bundle1",
					Status:     "purchased",
					Price:      100.0,
					DateListed: time.Now(),
				},
				{
					ID:         "bundle2",
					Status:     "available",
					Price:      200.0,
					DateListed: time.Now(),
				},
			},
			mockUser: &user.User{
				ID:         "supplier1",
				TrustScore: 85,
			},
			mockError:     nil,
			expectError:   false,
			expectedSales: 100.0,
			expectedCounts: order.PerformanceMetrics{
				TotalBundlesListed: 2,
				ActiveCount:        1,
				SoldCount:          1,
			},
			expectedRating: 85,
			expectedBest:   100.0,
		},
		{
			name:        "Success - No bundles",
			supplierID:  "supplier2",
			mockBundles: []*bundle.Bundle{},
			mockUser: &user.User{
				ID:         "supplier2",
				TrustScore: 90,
			},
			mockError:     nil,
			expectError:   false,
			expectedSales: 0.0,
			expectedCounts: order.PerformanceMetrics{
				TotalBundlesListed: 0,
				ActiveCount:        0,
				SoldCount:          0,
			},
			expectedRating: 90,
			expectedBest:   0.0,
		},
		{
			name:           "Error - Repository error",
			supplierID:     "supplier3",
			mockBundles:    nil,
			mockUser:       nil,
			mockError:      errors.New("database error"),
			expectError:    true,
			expectedSales:  0.0,
			expectedCounts: order.PerformanceMetrics{},
			expectedRating: 0,
			expectedBest:   0.0,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.bundleRepo.On("ListBundles", suite.ctx, tt.supplierID).Return(tt.mockBundles, tt.mockError)
			if tt.mockUser != nil {
				suite.userRepo.On("GetByID", suite.ctx, tt.supplierID).Return(tt.mockUser, nil)
			}

			// Act
			metrics, err := suite.useCase.GetDashboardMetrics(suite.ctx, tt.supplierID)

			// Assert
			if tt.expectError {
				assert.Error(suite.T(), err)
				assert.Nil(suite.T(), metrics)
			} else {
				assert.NoError(suite.T(), err)
				assert.NotNil(suite.T(), metrics)
				assert.Equal(suite.T(), tt.expectedSales, metrics.TotalSales)
				assert.Equal(suite.T(), tt.expectedCounts, metrics.PerformanceMetrics)
				assert.Equal(suite.T(), tt.expectedRating, metrics.Rating)
				assert.Equal(suite.T(), tt.expectedBest, metrics.BestSelling)
			}
		})
	}
}

// TestGetOrderByID tests the GetOrderByID method
func (suite *OrderUsecaseTestSuite) TestGetOrderByID() {
	tests := []struct {
		name        string
		orderID     string
		mockOrder   *order.Order
		mockError   error
		expectError bool
	}{
		{
			name:    "Success - Order found",
			orderID: "order1",
			mockOrder: &order.Order{
				ID:         "order1",
				ResellerID: "reseller1",
				Status:     order.OrderStatusCompleted,
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "Error - Order not found",
			orderID:     "nonexistent",
			mockOrder:   nil,
			mockError:   errors.New("order not found"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.orderRepo.On("GetOrderByID", suite.ctx, tt.orderID).Return(tt.mockOrder, tt.mockError)

			// Act
			order, err := suite.useCase.GetOrderByID(suite.ctx, tt.orderID)

			// Assert
			if tt.expectError {
				assert.Error(suite.T(), err)
				assert.Nil(suite.T(), order)
			} else {
				assert.NoError(suite.T(), err)
				assert.NotNil(suite.T(), order)
				assert.Equal(suite.T(), tt.mockOrder.ID, order.ID)
			}
		})
	}
}

// TestGetSoldBundleHistory tests the GetSoldBundleHistory method
func (suite *OrderUsecaseTestSuite) TestGetSoldBundleHistory() {
	tests := []struct {
		name          string
		supplierID    string
		mockOrders    []*order.Order
		mockError     error
		expectError   bool
		expectedCount int
	}{
		{
			name:       "Success - With sold bundles",
			supplierID: "supplier1",
			mockOrders: []*order.Order{
				{
					ID:         "order1",
					BundleID:   "bundle1",
					ProductIDs: []string{},
				},
				{
					ID:         "order2",
					BundleID:   "bundle2",
					ProductIDs: []string{},
				},
			},
			mockError:     nil,
			expectError:   false,
			expectedCount: 2,
		},
		{
			name:          "Success - No sold bundles",
			supplierID:    "supplier2",
			mockOrders:    []*order.Order{},
			mockError:     nil,
			expectError:   false,
			expectedCount: 0,
		},
		{
			name:          "Error - Repository error",
			supplierID:    "supplier3",
			mockOrders:    nil,
			mockError:     errors.New("database error"),
			expectError:   true,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.orderRepo.On("GetOrdersBySupplier", suite.ctx, tt.supplierID).Return(tt.mockOrders, tt.mockError)

			orders, userNames, err := suite.useCase.GetSoldBundleHistory(suite.ctx, tt.supplierID)

			// Assert
			if tt.expectError {
				assert.Error(suite.T(), err)
				assert.Nil(suite.T(), orders)
				assert.Nil(suite.T(), userNames)
			} else {
				assert.NoError(suite.T(), err)
				if tt.expectedCount == 0 {
					assert.Empty(suite.T(), orders)
					assert.Empty(suite.T(), userNames)
				} else {
					assert.NotNil(suite.T(), orders)
					assert.NotNil(suite.T(), userNames)
					assert.Len(suite.T(), orders, tt.expectedCount)
				}
			}
		})
	}
}

// TestPurchaseProduct tests the PurchaseProduct method
func (suite *OrderUsecaseTestSuite) TestPurchaseProduct() {
	productID := "test-product-id"
	userID := "test-user-id"
	price := 100.0

	resellerObjID := primitive.NewObjectID()

	// Mock product
	mockProduct := &product.Product{
		ID:         productID,
		ResellerID: resellerObjID,
		Price:      price,
	}
	suite.productRepo.On("GetProductByID", suite.ctx, productID).Return(mockProduct, nil)

	// Mock order creation
	suite.orderRepo.On("CreateOrder", suite.ctx, mock.AnythingOfType("*order.Order")).Return(nil)

	// Mock payment recording
	suite.paymentRepo.On("RecordPayment", suite.ctx, mock.AnythingOfType("*payment.Payment")).Return(nil)

	// Act
	order, payment, err := suite.useCase.PurchaseProduct(suite.ctx, productID, userID, price)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), order)
	assert.NotNil(suite.T(), payment)
}

// TestPurchaseProduct_ProductNotFound tests the PurchaseProduct method when product is not found
func (suite *OrderUsecaseTestSuite) TestPurchaseProduct_ProductNotFound() {
	productID := "test-product-id"
	userID := "test-user-id"
	price := 100.0

	// Mock product not found
	suite.productRepo.On("GetProductByID", suite.ctx, productID).Return(nil, errors.New("product not found"))

	// No need to mock other calls since the test should fail before reaching them

	// Act
	order, payment, err := suite.useCase.PurchaseProduct(suite.ctx, productID, userID, price)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), order)
	assert.Nil(suite.T(), payment)
	assert.Contains(suite.T(), err.Error(), "product not found")
}
