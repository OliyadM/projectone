package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/admin"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/order"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/payment"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/warehouse"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type OrderControllerTestSuite struct {
	suite.Suite
	orderUseCase *MockOrderUseCase
	controller   *OrderController
	router       *gin.Engine
}

type MockOrderUseCase struct {
	mock.Mock
}

func (m *MockOrderUseCase) PurchaseBundle(ctx context.Context, bundleID, resellerID string) (*order.Order, *payment.Payment, *warehouse.WarehouseItem, error) {
	args := m.Called(ctx, bundleID, resellerID)
	if args.Get(0) == nil {
		return nil, nil, nil, args.Error(3)
	}
	return args.Get(0).(*order.Order), args.Get(1).(*payment.Payment), args.Get(2).(*warehouse.WarehouseItem), args.Error(3)
}

func (m *MockOrderUseCase) GetOrderByID(ctx context.Context, orderID string) (*order.Order, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order.Order), args.Error(1)
}

func (m *MockOrderUseCase) GetDashboardMetrics(ctx context.Context, supplierID string) (*order.DashboardMetrics, error) {
	args := m.Called(ctx, supplierID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order.DashboardMetrics), args.Error(1)
}

func (m *MockOrderUseCase) GetResellerMetrics(ctx context.Context, resellerID string) (*order.ResellerMetrics, error) {
	args := m.Called(ctx, resellerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order.ResellerMetrics), args.Error(1)
}

func (m *MockOrderUseCase) GetSoldBundleHistory(ctx context.Context, supplierID string) ([]*order.Order, map[string]string, error) {
	args := m.Called(ctx, supplierID)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]*order.Order), args.Get(1).(map[string]string), args.Error(2)
}

func (m *MockOrderUseCase) GetOrdersByReseller(ctx context.Context, resellerID string) ([]*order.Order, map[string]string, error) {
	args := m.Called(ctx, resellerID)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]*order.Order), args.Get(1).(map[string]string), args.Error(2)
}

func (m *MockOrderUseCase) GetAdminDashboardMetrics(ctx context.Context) (*admin.Metrics, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*admin.Metrics), args.Error(1)
}

func (m *MockOrderUseCase) PurchaseProduct(ctx context.Context, productID, consumerID string, totalPrice float64) (*order.Order, *payment.Payment, error) {
	args := m.Called(ctx, productID, consumerID, totalPrice)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*order.Order), args.Get(1).(*payment.Payment), args.Error(2)
}

func (m *MockOrderUseCase) GetOrdersByConsumer(ctx context.Context, consumerID string) ([]*order.Order, map[string]string, map[string]string, error) {
	args := m.Called(ctx, consumerID)
	if args.Get(0) == nil {
		return nil, nil, nil, args.Error(3)
	}
	return args.Get(0).([]*order.Order), args.Get(1).(map[string]string), args.Get(2).(map[string]string), args.Error(3)
}

func (suite *OrderControllerTestSuite) SetupTest() {
	suite.orderUseCase = new(MockOrderUseCase)
	suite.controller = NewOrderController(suite.orderUseCase)
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()
}

func TestOrderControllerTestSuite(t *testing.T) {
	suite.Run(t, new(OrderControllerTestSuite))
}

func (suite *OrderControllerTestSuite) TestPurchaseBundle_Success() {
	// Setup
	expectedOrder := &order.Order{ID: "order123"}
	expectedPayment := &payment.Payment{ID: "payment123"}
	expectedWarehouseItem := &warehouse.WarehouseItem{ID: "warehouse123"}

	suite.orderUseCase.On("PurchaseBundle", mock.Anything, "bundle123", "reseller123").
		Return(expectedOrder, expectedPayment, expectedWarehouseItem, nil)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "reseller123")

	body, _ := json.Marshal(gin.H{"bundle_id": "bundle123"})
	c.Request = httptest.NewRequest("POST", "/purchase", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	suite.controller.PurchaseBundle(c)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.orderUseCase.AssertExpectations(suite.T())
}

func (suite *OrderControllerTestSuite) TestPurchaseBundle_InvalidPayload() {
	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(gin.H{})
	c.Request = httptest.NewRequest("POST", "/purchase", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	suite.controller.PurchaseBundle(c)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *OrderControllerTestSuite) TestPurchaseBundle_InvalidUserID() {
	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(gin.H{"bundle_id": "bundle123"})
	c.Request = httptest.NewRequest("POST", "/purchase", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	suite.controller.PurchaseBundle(c)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *OrderControllerTestSuite) TestPurchaseBundle_UseCaseError() {
	// Setup
	suite.orderUseCase.On("PurchaseBundle", mock.Anything, "bundle123", "reseller123").
		Return(nil, nil, nil, errors.New("use case error"))

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "reseller123")

	body, _ := json.Marshal(gin.H{"bundle_id": "bundle123"})
	c.Request = httptest.NewRequest("POST", "/purchase", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	suite.controller.PurchaseBundle(c)

	// Assert
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	suite.orderUseCase.AssertExpectations(suite.T())
}

func (suite *OrderControllerTestSuite) TestGetOrderByID_Success() {
	// Setup
	expectedOrder := &order.Order{ID: "order123"}
	suite.orderUseCase.On("GetOrderByID", mock.Anything, "order123").
		Return(expectedOrder, nil)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "order123"}}

	// Execute
	suite.controller.GetOrderByID(c)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.orderUseCase.AssertExpectations(suite.T())
}

func (suite *OrderControllerTestSuite) TestGetOrderByID_MissingID() {
	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Execute
	suite.controller.GetOrderByID(c)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *OrderControllerTestSuite) TestGetOrderByID_UseCaseError() {
	// Setup
	suite.orderUseCase.On("GetOrderByID", mock.Anything, "order123").
		Return(nil, errors.New("use case error"))

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "order123"}}

	// Execute
	suite.controller.GetOrderByID(c)

	// Assert
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	suite.orderUseCase.AssertExpectations(suite.T())
}
