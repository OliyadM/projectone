package controllers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/admin"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/order"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/payment"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/warehouse"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockOrderUsecase struct {
	mock.Mock
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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order.DashboardMetrics), args.Error(1)
}

func (m *MockOrderUsecase) GetResellerMetrics(ctx context.Context, resellerID string) (*order.ResellerMetrics, error) {
	args := m.Called(ctx, resellerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order.ResellerMetrics), args.Error(1)
}

func (m *MockOrderUsecase) GetOrderByID(ctx context.Context, orderID string) (*order.Order, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order.Order), args.Error(1)
}

func (m *MockOrderUsecase) GetSoldBundleHistory(ctx context.Context, supplierID string) ([]*order.Order, map[string]string, error) {
	args := m.Called(ctx, supplierID)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]*order.Order), args.Get(1).(map[string]string), args.Error(2)
}

func (m *MockOrderUsecase) GetAdminDashboardMetrics(ctx context.Context) (*admin.Metrics, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*admin.Metrics), args.Error(1)
}

type SupplierControllerTestSuite struct {
	suite.Suite
	usecase    *MockOrderUsecase
	controller *SupplierController
	router     *gin.Engine
}

func (suite *SupplierControllerTestSuite) SetupTest() {
	suite.usecase = new(MockOrderUsecase)
	suite.controller = NewSupplierController(suite.usecase)
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()
}

func TestSupplierControllerTestSuite(t *testing.T) {
	suite.Run(t, new(SupplierControllerTestSuite))
}

func (suite *SupplierControllerTestSuite) TestGetDashboardMetrics_Success() {
	// Setup
	expectedMetrics := &order.DashboardMetrics{
		TotalSales: 1000.0,
		ActiveBundles: []*bundle.Bundle{
			{
				ID: "bundle1",
			},
		},
		PerformanceMetrics: order.PerformanceMetrics{
			TotalBundlesListed: 10,
			ActiveCount:        5,
			SoldCount:          5,
		},
	}

	suite.usecase.On("GetDashboardMetrics", mock.Anything, "supplier123").
		Return(expectedMetrics, nil)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "supplier123")
	c.Request = httptest.NewRequest("GET", "/supplier/dashboard", nil)

	// Execute
	suite.controller.GetDashboardMetrics(c)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.usecase.AssertExpectations(suite.T())
}

func (suite *SupplierControllerTestSuite) TestGetDashboardMetrics_Unauthorized() {
	// Setup
	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Note: No userID set
	c.Request = httptest.NewRequest("GET", "/supplier/dashboard", nil)

	// Execute
	suite.controller.GetDashboardMetrics(c)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	suite.usecase.AssertNotCalled(suite.T(), "GetDashboardMetrics")
}

func (suite *SupplierControllerTestSuite) TestGetDashboardMetrics_InvalidUserID() {
	// Setup
	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", 123) // Invalid type
	c.Request = httptest.NewRequest("GET", "/supplier/dashboard", nil)

	// Execute
	suite.controller.GetDashboardMetrics(c)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	suite.usecase.AssertNotCalled(suite.T(), "GetDashboardMetrics")
}

func (suite *SupplierControllerTestSuite) TestGetDashboardMetrics_UseCaseError() {
	// Setup
	suite.usecase.On("GetDashboardMetrics", mock.Anything, "supplier123").
		Return(nil, assert.AnError)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "supplier123")
	c.Request = httptest.NewRequest("GET", "/supplier/dashboard", nil)

	// Execute
	suite.controller.GetDashboardMetrics(c)

	// Assert
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	suite.usecase.AssertExpectations(suite.T())
}
