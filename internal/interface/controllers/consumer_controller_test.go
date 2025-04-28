package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/order"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) CreateOrder(ctx context.Context, o *order.Order) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

func (m *MockOrderRepository) GetOrdersByConsumer(ctx context.Context, consumerID string) ([]*order.Order, error) {
	args := m.Called(ctx, consumerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*order.Order), args.Error(1)
}

func (m *MockOrderRepository) GetOrdersByReseller(ctx context.Context, resellerID string) ([]*order.Order, error) {
	args := m.Called(ctx, resellerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*order.Order), args.Error(1)
}

type ConsumerControllerTestSuite struct {
	suite.Suite
	mockRepo    *MockOrderRepository
	controller  *ConsumerController
	testContext *gin.Context
	recorder    *httptest.ResponseRecorder
}

func (suite *ConsumerControllerTestSuite) SetupTest() {
	suite.mockRepo = new(MockOrderRepository)
	suite.controller = NewConsumerController(suite.mockRepo)

	// Set up a Gin context and response recorder for testing
	gin.SetMode(gin.TestMode)
	suite.recorder = httptest.NewRecorder()
	suite.testContext, _ = gin.CreateTestContext(suite.recorder)
}

func (suite *ConsumerControllerTestSuite) TestGetOrderHistory_Success() {
	orders := []*order.Order{
		{
			ID:         "order1",
			ConsumerID: "consumer1",
			ProductIDs: []string{"product1"},
			TotalPrice: 100.0,
			Status:     order.Pending,
			CreatedAt:  time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
		},
	}

	suite.mockRepo.On("GetOrdersByConsumer", mock.Anything, "consumer1").Return(orders, nil)

	suite.testContext.Set("userID", "consumer1")
	suite.controller.GetOrderHistory(suite.testContext)

	assert.Equal(suite.T(), http.StatusOK, suite.recorder.Code)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ConsumerControllerTestSuite) TestGetOrderHistory_NoOrders() {
	suite.mockRepo.On("GetOrdersByConsumer", mock.Anything, "consumer1").Return([]*order.Order{}, nil)

	suite.testContext.Set("userID", "consumer1")
	suite.controller.GetOrderHistory(suite.testContext)

	assert.Equal(suite.T(), http.StatusOK, suite.recorder.Code)
	assert.Contains(suite.T(), suite.recorder.Body.String(), "No orders yet.")
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ConsumerControllerTestSuite) TestGetOrderHistory_ErrorFetchingOrders() {
	suite.mockRepo.On("GetOrdersByConsumer", mock.Anything, "consumer1").Return(nil, errors.New("database error"))

	suite.testContext.Set("userID", "consumer1")
	suite.controller.GetOrderHistory(suite.testContext)

	assert.Equal(suite.T(), http.StatusInternalServerError, suite.recorder.Code)
	assert.Contains(suite.T(), suite.recorder.Body.String(), "Failed to fetch orders")
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ConsumerControllerTestSuite) TestGetOrderHistory_DeliveryFailure() {
	orders := []*order.Order{
		{
			ID:         "order1",
			ConsumerID: "consumer1",
			ProductIDs: []string{"product1"},
			TotalPrice: 100.0,
			Status:     order.Pending,
			CreatedAt:  time.Now().Add(-11 * time.Minute).Format(time.RFC3339), // Simulate long delay
		},
	}

	suite.mockRepo.On("GetOrdersByConsumer", mock.Anything, "consumer1").Return(orders, nil)

	suite.testContext.Set("userID", "consumer1")
	suite.controller.GetOrderHistory(suite.testContext)

	assert.Equal(suite.T(), http.StatusOK, suite.recorder.Code)
	assert.Contains(suite.T(), suite.recorder.Body.String(), "failed") // Check for failure state
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ConsumerControllerTestSuite) TestGetOrderHistory_Pagination() {
	orders := []*order.Order{}
	for i := 1; i <= 15; i++ {
		orders = append(orders, &order.Order{
			ID:         fmt.Sprintf("order%d", i),
			ConsumerID: "consumer1",
			ProductIDs: []string{fmt.Sprintf("product%d", i)},
			TotalPrice: float64(i * 10),
			Status:     order.Pending,
			CreatedAt:  time.Now().Add(-time.Duration(i) * time.Minute).Format(time.RFC3339),
		})
	}

	suite.mockRepo.On("GetOrdersByConsumer", mock.Anything, "consumer1").Return(orders, nil)

	suite.testContext.Set("userID", "consumer1")
	suite.testContext.Request = httptest.NewRequest(http.MethodGet, "/orders/history?page=2&limit=5", nil)
	suite.controller.GetOrderHistory(suite.testContext)

	assert.Equal(suite.T(), http.StatusOK, suite.recorder.Code)
	assert.Contains(suite.T(), suite.recorder.Body.String(), "order6")
	assert.Contains(suite.T(), suite.recorder.Body.String(), "order10")
	assert.NotContains(suite.T(), suite.recorder.Body.String(), "order5")
	assert.NotContains(suite.T(), suite.recorder.Body.String(), "order11")
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ConsumerControllerTestSuite) TestGetOrderHistory_FilterByStatus() {
	orders := []*order.Order{
		{ID: "order1", ConsumerID: "consumer1", Status: order.Pending}, // Status matches filter
		{ID: "order2", ConsumerID: "consumer1", Status: order.Failed},  // Different status
		{ID: "order3", ConsumerID: "consumer1", Status: order.Pending}, // Status matches filter
	}

	// Simulate delivery logic only for orders that should match the filter
	for i := range orders {
		if orders[i].ID == "order1" && orders[i].Status == order.Pending {
			orders[i].Status = order.Delivered // Simulate status transition for order1 only
		}
	}

	suite.mockRepo.On("GetOrdersByConsumer", mock.Anything, "consumer1").Return(orders, nil)

	suite.testContext.Set("userID", "consumer1")
	suite.testContext.Request = httptest.NewRequest(http.MethodGet, "/orders/history?status=delivered", nil)
	suite.controller.GetOrderHistory(suite.testContext)

	assert.Equal(suite.T(), http.StatusOK, suite.recorder.Code)
	assert.Contains(suite.T(), suite.recorder.Body.String(), "order1") // Ensure "order1" is in the response
	assert.NotContains(suite.T(), suite.recorder.Body.String(), "order2")
	assert.NotContains(suite.T(), suite.recorder.Body.String(), "order3")
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ConsumerControllerTestSuite) TestGetOrderHistory_SortingByDate() {
	orders := []*order.Order{
		{ID: "order1", ConsumerID: "consumer1", CreatedAt: time.Now().Add(-10 * time.Minute).Format(time.RFC3339)},
		{ID: "order2", ConsumerID: "consumer1", CreatedAt: time.Now().Add(-5 * time.Minute).Format(time.RFC3339)},
		{ID: "order3", ConsumerID: "consumer1", CreatedAt: time.Now().Add(-15 * time.Minute).Format(time.RFC3339)},
	}

	suite.mockRepo.On("GetOrdersByConsumer", mock.Anything, "consumer1").Return(orders, nil)

	suite.testContext.Set("userID", "consumer1")
	suite.controller.GetOrderHistory(suite.testContext)

	assert.Equal(suite.T(), http.StatusOK, suite.recorder.Code)
	responseBody := suite.recorder.Body.String()
	assert.True(suite.T(), strings.Index(responseBody, "order2") < strings.Index(responseBody, "order1"))
	assert.True(suite.T(), strings.Index(responseBody, "order1") < strings.Index(responseBody, "order3"))
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ConsumerControllerTestSuite) TestGetOrderHistory_DefaultQueryParameters() {
	orders := []*order.Order{}
	for i := 1; i <= 15; i++ {
		orders = append(orders, &order.Order{
			ID:         fmt.Sprintf("order%d", i),
			ConsumerID: "consumer1",
			ProductIDs: []string{fmt.Sprintf("product%d", i)},
			TotalPrice: float64(i * 10),
			Status:     order.Pending,
			CreatedAt:  time.Now().Add(-time.Duration(i) * time.Minute).Format(time.RFC3339),
		})
	}

	suite.mockRepo.On("GetOrdersByConsumer", mock.Anything, "consumer1").Return(orders, nil)

	suite.testContext.Set("userID", "consumer1")
	suite.testContext.Request = httptest.NewRequest(http.MethodGet, "/orders/history", nil)
	suite.controller.GetOrderHistory(suite.testContext)

	assert.Equal(suite.T(), http.StatusOK, suite.recorder.Code)
	assert.Contains(suite.T(), suite.recorder.Body.String(), "order1")
	assert.Contains(suite.T(), suite.recorder.Body.String(), "order10")
	assert.NotContains(suite.T(), suite.recorder.Body.String(), "order11")
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestConsumerControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ConsumerControllerTestSuite))
}
