package cartitem

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/cartitem"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/order"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/payment"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/product"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/warehouse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- Mocks ---

type MockCartItemRepository struct {
	mock.Mock
}

func (m *MockCartItemRepository) CreateCartItem(ctx context.Context, item *cartitem.CartItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockCartItemRepository) GetCartItems(ctx context.Context, userID string) ([]*cartitem.CartItem, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*cartitem.CartItem), args.Error(1)
}

func (m *MockCartItemRepository) DeleteCartItem(ctx context.Context, userID, listingID string) error {
	args := m.Called(ctx, userID, listingID)
	return args.Error(0)
}

func (m *MockCartItemRepository) ClearCart(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

type MockProductRepository struct {
	mock.Mock
}

// GetProductByTitle implements product.Repository.
func (m *MockProductRepository) GetProductByTitle(ctx context.Context, title string) (*product.Product, error) {
	panic("unimplemented")
}

func (m *MockProductRepository) GetProductByID(ctx context.Context, id string) (*product.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Product), args.Error(1)
}

// Added dummy implementation so that it satisfies product.Repository.
func (m *MockProductRepository) AddProduct(ctx context.Context, prod *product.Product) error {
	return nil
}

// Added dummy DeleteProduct method as required by product.Repository.
func (m *MockProductRepository) DeleteProduct(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepository) ListProductsByReseller(ctx context.Context, resellerID string, page, limit int) ([]*product.Product, error) {
	args := m.Called(ctx, resellerID, page, limit)
	return args.Get(0).([]*product.Product), args.Error(1)
}

func (m *MockProductRepository) ListAvailableProducts(ctx context.Context, page, limit int) ([]*product.Product, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]*product.Product), args.Error(1)
}

func (m *MockProductRepository) UpdateProduct(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockProductRepository) GetProductsByBundleID(ctx context.Context, bundleID string) ([]*product.Product, error) {
	args := m.Called(ctx, bundleID)
	return args.Get(0).([]*product.Product), args.Error(1)
}

func (m *MockProductRepository) GetSoldProductsByReseller(ctx context.Context, resellerID string) ([]*product.Product, error) {
	args := m.Called(ctx, resellerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*product.Product), args.Error(1)
}

type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) RecordPayment(ctx context.Context, p *payment.Payment) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockPaymentRepository) GetAllPlatformFees(ctx context.Context) (float64, float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Get(1).(float64), args.Error(2)
}

func (m *MockPaymentRepository) GetPaymentsByType(ctx context.Context, userID string, pType payment.PaymentType) ([]*payment.Payment, error) {
	args := m.Called(ctx, userID, pType)
	return args.Get(0).([]*payment.Payment), args.Error(1)
}

func (m *MockPaymentRepository) GetPaymentsByUser(ctx context.Context, userID string) ([]*payment.Payment, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*payment.Payment), args.Error(1)
}

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) CreateOrder(ctx context.Context, o *order.Order) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

func (m *MockOrderRepository) GetOrdersByConsumer(ctx context.Context, consumerID string) ([]*order.Order, error) {
	args := m.Called(ctx, consumerID)
	return args.Get(0).([]*order.Order), args.Error(1)
}

func (m *MockOrderRepository) GetOrderByID(ctx context.Context, orderID string) (*order.Order, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order.Order), args.Error(1)
}

func (m *MockOrderRepository) UpdateOrderStatus(ctx context.Context, orderID string, status order.OrderStatus) error {
	args := m.Called(ctx, orderID, status)
	return args.Error(0)
}

func (m *MockOrderRepository) DeleteOrder(ctx context.Context, orderID string) error {
	args := m.Called(ctx, orderID)
	return args.Error(0)
}

func (m *MockOrderRepository) GetOrdersBySupplier(ctx context.Context, supplierID string) ([]*order.Order, error) {
	args := m.Called(ctx, supplierID)
	return args.Get(0).([]*order.Order), args.Error(1)
}

func (m *MockOrderRepository) GetOrdersByReseller(ctx context.Context, resellerID string) ([]*order.Order, error) {
	args := m.Called(ctx, resellerID)
	return args.Get(0).([]*order.Order), args.Error(1)
}

type MockOrderUsecase struct {
	mock.Mock
}

func (m *MockOrderUsecase) PurchaseProduct(ctx context.Context, productID, consumerID string, totalPrice float64) (*order.Order, *payment.Payment, error) {
	args := m.Called(ctx, productID, consumerID, totalPrice)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*order.Order), args.Get(1).(*payment.Payment), args.Error(2)
}

func (m *MockOrderUsecase) PurchaseBundle(ctx context.Context, bundleID, resellerID string) (*order.Order, *payment.Payment, *warehouse.WarehouseItem, error) {
	args := m.Called(ctx, bundleID, resellerID)
	if args.Get(0) == nil {
		return nil, nil, nil, args.Error(3)
	}
	return args.Get(0).(*order.Order), args.Get(1).(*payment.Payment), args.Get(2).(*warehouse.WarehouseItem), args.Error(3)
}

func (m *MockOrderUsecase) GetDashboardMetrics(ctx context.Context, supplierID string) (*order.DashboardMetrics, error) {
	args := m.Called(ctx, supplierID)
	return args.Get(0).(*order.DashboardMetrics), args.Error(1)
}

func (m *MockOrderUsecase) GetOrderByID(ctx context.Context, orderID string) (*order.Order, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).(*order.Order), args.Error(1)
}

func (m *MockOrderUsecase) GetOrdersByReseller(ctx context.Context, resellerID string) ([]*order.Order, map[string]string, error) {
	args := m.Called(ctx, resellerID)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]*order.Order), args.Get(1).(map[string]string), args.Error(2)
}

func (m *MockOrderUsecase) GetSoldBundleHistory(ctx context.Context, supplierID string) ([]*order.Order, map[string]string, error) {
	args := m.Called(ctx, supplierID)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]*order.Order), args.Get(1).(map[string]string), args.Error(2)
}

func (m *MockOrderUsecase) GetResellerMetrics(ctx context.Context, resellerID string) (*order.ResellerMetrics, error) {
	args := m.Called(ctx, resellerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order.ResellerMetrics), args.Error(1)
}

func (m *MockOrderUsecase) GetOrdersByConsumer(ctx context.Context, consumerID string) ([]*order.Order, map[string]string, map[string]string, error) {
	args := m.Called(ctx, consumerID)
	if args.Get(0) == nil {
		return nil, nil, nil, args.Error(3)
	}
	return args.Get(0).([]*order.Order), args.Get(1).(map[string]string), args.Get(2).(map[string]string), args.Error(3)
}

// --- Test Suite ---

type CartItemUsecaseTestSuite struct {
	suite.Suite
	ctx             context.Context
	usecase         cartitem.Usecase
	mockCartRepo    *MockCartItemRepository
	mockProductRepo *MockProductRepository
	mockPaymentRepo *MockPaymentRepository
	mockOrderRepo   *MockOrderRepository
	mockOrderUC     *MockOrderUsecase
	userID          string
}

func (suite *CartItemUsecaseTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.mockCartRepo = new(MockCartItemRepository)
	suite.mockProductRepo = new(MockProductRepository)
	suite.mockPaymentRepo = new(MockPaymentRepository)
	suite.mockOrderRepo = new(MockOrderRepository)
	suite.mockOrderUC = new(MockOrderUsecase)
	suite.usecase = NewCartItemUsecase(suite.mockCartRepo, suite.mockProductRepo, suite.mockPaymentRepo, suite.mockOrderUC, suite.mockOrderRepo)
	suite.userID = "user123"
}

// --- Helper: create a dummy product ---
func createTestProduct(id string, price float64, status string, title string) *product.Product {
	// Use primitive.NewObjectID() to assign the ResellerID.
	return &product.Product{
		ID:         id,
		Title:      title,
		Price:      price,
		ImageURL:   "image.jpg",
		Grade:      "A",
		Status:     status,
		ResellerID: primitive.NewObjectID(),
	}
}

// --- Tests for AddCartItem ---

func (suite *CartItemUsecaseTestSuite) TestAddCartItem_Success() {
	testListingID := "prod123"
	prod := createTestProduct(testListingID, 100.0, "available", "Test Product")
	// Expect product lookup.
	suite.mockProductRepo.On("GetProductByID", suite.ctx, testListingID).Return(prod, nil).Once()
	// Expect repository CreateCartItem call.
	suite.mockCartRepo.On("CreateCartItem", suite.ctx, mock.MatchedBy(func(item *cartitem.CartItem) bool {
		return item.UserID == suite.userID && item.ListingID == prod.ID && item.Title == prod.Title
	})).Return(nil).Once()

	err := suite.usecase.AddCartItem(suite.ctx, suite.userID, testListingID)
	assert.NoError(suite.T(), err)
	suite.mockProductRepo.AssertExpectations(suite.T())
	suite.mockCartRepo.AssertExpectations(suite.T())
}

func (suite *CartItemUsecaseTestSuite) TestAddCartItem_ProductNotFound() {
	testListingID := "prodNotFound"
	suite.mockProductRepo.On("GetProductByID", suite.ctx, testListingID).Return(nil, nil).Once()

	err := suite.usecase.AddCartItem(suite.ctx, suite.userID, testListingID)
	assert.EqualError(suite.T(), err, "product not found")
	suite.mockProductRepo.AssertExpectations(suite.T())
}

func (suite *CartItemUsecaseTestSuite) TestAddCartItem_ProductNotAvailable() {
	testListingID := "prod123"
	prod := createTestProduct(testListingID, 100.0, "sold", "Test Product")
	suite.mockProductRepo.On("GetProductByID", suite.ctx, testListingID).Return(prod, nil).Once()

	err := suite.usecase.AddCartItem(suite.ctx, suite.userID, testListingID)
	expectedErr := fmt.Sprintf("product %s is not available", testListingID)
	assert.EqualError(suite.T(), err, expectedErr)
	suite.mockProductRepo.AssertExpectations(suite.T())
}

// --- Tests for GetCartItems ---

func (suite *CartItemUsecaseTestSuite) TestGetCartItems_Success() {
	cartItems := []*cartitem.CartItem{
		{
			ID:        "item1",
			UserID:    suite.userID,
			ListingID: "prod123",
			Title:     "Test Product",
			Price:     100.0,
			ImageURL:  "img.jpg",
			Grade:     "A",
			CreatedAt: time.Now(),
		},
	}
	suite.mockCartRepo.On("GetCartItems", suite.ctx, suite.userID).Return(cartItems, nil).Once()

	items, err := suite.usecase.GetCartItems(suite.ctx, suite.userID)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), items, 1)
	suite.mockCartRepo.AssertExpectations(suite.T())
}

// --- Tests for RemoveCartItem ---

func (suite *CartItemUsecaseTestSuite) TestRemoveCartItem_Success() {
	listingID := "prod123"
	suite.mockCartRepo.On("DeleteCartItem", suite.ctx, suite.userID, listingID).Return(nil).Once()

	err := suite.usecase.RemoveCartItem(suite.ctx, suite.userID, listingID)
	assert.NoError(suite.T(), err)
	suite.mockCartRepo.AssertExpectations(suite.T())
}

// --- Tests for CheckoutCart ---

func (suite *CartItemUsecaseTestSuite) TestCheckoutCart_Success() {
	now := time.Now()
	cartItems := []*cartitem.CartItem{
		{
			ID:        "item1",
			UserID:    suite.userID,
			ListingID: "prod1",
			Title:     "Test Product 1",
			Price:     100.0,
			ImageURL:  "img1.jpg",
			Grade:     "A",
			CreatedAt: now,
		},
		{
			ID:        "item2",
			UserID:    suite.userID,
			ListingID: "prod2",
			Title:     "Test Product 2",
			Price:     200.0,
			ImageURL:  "img2.jpg",
			Grade:     "B",
			CreatedAt: now,
		},
	}
	suite.mockCartRepo.On("GetCartItems", suite.ctx, suite.userID).Return(cartItems, nil).Once()

	prod1 := createTestProduct("prod1", 100.0, "available", "Test Product 1")
	prod2 := createTestProduct("prod2", 200.0, "available", "Test Product 2")
	suite.mockProductRepo.On("GetProductByID", suite.ctx, "prod1").Return(prod1, nil).Once()
	suite.mockProductRepo.On("GetProductByID", suite.ctx, "prod2").Return(prod2, nil).Once()

	// Mock order and payment creation
	order1 := &order.Order{ID: "order1"}
	payment1 := &payment.Payment{ID: "payment1"}
	order2 := &order.Order{ID: "order2"}
	payment2 := &payment.Payment{ID: "payment2"}

	suite.mockOrderUC.On("PurchaseProduct", suite.ctx, "prod1", suite.userID, 100.0).Return(order1, payment1, nil).Once()
	suite.mockOrderUC.On("PurchaseProduct", suite.ctx, "prod2", suite.userID, 200.0).Return(order2, payment2, nil).Once()

	// Mock payment updates
	suite.mockPaymentRepo.On("RecordPayment", suite.ctx, mock.Anything).Return(nil).Twice()

	// Mock product status updates
	suite.mockProductRepo.On("UpdateProduct", suite.ctx, "prod1", mock.Anything).Return(nil).Once()
	suite.mockProductRepo.On("UpdateProduct", suite.ctx, "prod2", mock.Anything).Return(nil).Once()

	// Mock cart clearing
	suite.mockCartRepo.On("ClearCart", suite.ctx, suite.userID).Return(nil).Once()

	resp, err := suite.usecase.CheckoutCart(suite.ctx, suite.userID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 300.0, resp.TotalAmount)
	assert.Equal(suite.T(), 6.0, resp.PlatformFee)
	assert.Equal(suite.T(), 294.0, resp.NetPayable)
	assert.Len(suite.T(), resp.Items, 2)

	suite.mockCartRepo.AssertExpectations(suite.T())
	suite.mockProductRepo.AssertExpectations(suite.T())
	suite.mockPaymentRepo.AssertExpectations(suite.T())
	suite.mockOrderUC.AssertExpectations(suite.T())
}

func (suite *CartItemUsecaseTestSuite) TestCheckoutCart_EmptyCart() {
	suite.mockCartRepo.On("GetCartItems", suite.ctx, suite.userID).Return([]*cartitem.CartItem{}, nil).Once()

	resp, err := suite.usecase.CheckoutCart(suite.ctx, suite.userID)
	assert.Nil(suite.T(), resp)
	assert.EqualError(suite.T(), err, "cart is empty")
	suite.mockCartRepo.AssertExpectations(suite.T())
}

// --- Tests for CheckoutSingleItem ---

func (suite *CartItemUsecaseTestSuite) TestCheckoutSingleItem_Success() {
	now := time.Now()
	cartItems := []*cartitem.CartItem{
		{
			ID:        "item1",
			UserID:    suite.userID,
			ListingID: "prod1",
			Title:     "Test Product 1",
			Price:     100.0,
			ImageURL:  "img1.jpg",
			Grade:     "A",
			CreatedAt: now,
		},
	}
	suite.mockCartRepo.On("GetCartItems", suite.ctx, suite.userID).Return(cartItems, nil).Once()

	prod1 := createTestProduct("prod1", 100.0, "available", "Test Product 1")
	suite.mockProductRepo.On("GetProductByID", suite.ctx, "prod1").Return(prod1, nil).Once()

	// Mock order and payment creation
	order := &order.Order{ID: "order1"}
	payment := &payment.Payment{ID: "payment1"}
	suite.mockOrderUC.On("PurchaseProduct", suite.ctx, "prod1", suite.userID, 100.0).Return(order, payment, nil).Once()

	// Mock payment update
	suite.mockPaymentRepo.On("RecordPayment", suite.ctx, mock.Anything).Return(nil).Once()

	// Mock product status update
	suite.mockProductRepo.On("UpdateProduct", suite.ctx, "prod1", mock.Anything).Return(nil).Once()

	// Mock cart item deletion
	suite.mockCartRepo.On("DeleteCartItem", suite.ctx, suite.userID, "prod1").Return(nil).Once()

	resp, err := suite.usecase.CheckoutSingleItem(suite.ctx, suite.userID, "prod1")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 100.0, resp.TotalAmount)
	assert.Equal(suite.T(), 2.0, resp.PlatformFee)
	assert.Equal(suite.T(), 98.0, resp.NetPayable)
	assert.Len(suite.T(), resp.Items, 1)

	suite.mockCartRepo.AssertExpectations(suite.T())
	suite.mockProductRepo.AssertExpectations(suite.T())
	suite.mockPaymentRepo.AssertExpectations(suite.T())
	suite.mockOrderUC.AssertExpectations(suite.T())
}

func (suite *CartItemUsecaseTestSuite) TestCheckoutSingleItem_ItemNotFoundInCart() {
	// Empty cart scenario.
	suite.mockCartRepo.On("GetCartItems", suite.ctx, suite.userID).Return([]*cartitem.CartItem{}, nil).Once()

	resp, err := suite.usecase.CheckoutSingleItem(suite.ctx, suite.userID, "prod1")
	assert.Nil(suite.T(), resp)
	assert.EqualError(suite.T(), err, "item not found in cart")
	suite.mockCartRepo.AssertExpectations(suite.T())
}

func (suite *CartItemUsecaseTestSuite) TestCheckoutSingleItem_ProductNotAvailable() {
	// Set up a cart with one item.
	now := time.Now()
	cartItems := []*cartitem.CartItem{
		{
			ID:        "item1",
			UserID:    suite.userID,
			ListingID: "prod1",
			Title:     "Test Product 1",
			Price:     100.0,
			ImageURL:  "img1.jpg",
			Grade:     "A",
			CreatedAt: now,
		},
	}
	suite.mockCartRepo.On("GetCartItems", suite.ctx, suite.userID).Return(cartItems, nil).Once()
	prod1 := createTestProduct("prod1", 100.0, "sold", "Test Product 1")
	suite.mockProductRepo.On("GetProductByID", suite.ctx, "prod1").Return(prod1, nil).Once()

	resp, err := suite.usecase.CheckoutSingleItem(suite.ctx, suite.userID, "prod1")
	assert.Nil(suite.T(), resp)
	expectedErr := fmt.Sprintf("item %q is no longer available", prod1.Title)
	assert.EqualError(suite.T(), err, expectedErr)
	suite.mockCartRepo.AssertExpectations(suite.T())
	suite.mockProductRepo.AssertExpectations(suite.T())
}

func TestCartItemUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(CartItemUsecaseTestSuite))
}
