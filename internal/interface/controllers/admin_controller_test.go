package controllers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/admin"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/order"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/payment"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/warehouse"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// -------------------- Mock User Usecase --------------------

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) ListByRole(ctx context.Context, role user.Role) ([]*user.User, error) {
	args := m.Called(ctx, role)
	return args.Get(0).([]*user.User), args.Error(1)
}

func (m *MockUserUsecase) GetByID(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserUsecase) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockUserUsecase) GetBlacklistedUsers(ctx context.Context) ([]*user.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*user.User), args.Error(1)
}

func (m *MockUserUsecase) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserUsecase) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*user.User), args.Error(1)
}

// -------------------- Mock Order Usecase --------------------

// Update this type to implement the full interface
type AdminMockOrderUsecase struct {
	mock.Mock
}

func (m *AdminMockOrderUsecase) PurchaseBundle(ctx context.Context, bundleID, resellerID string) (*order.Order, *payment.Payment, *warehouse.WarehouseItem, error) {
	args := m.Called(ctx, bundleID, resellerID)
	return nil, nil, nil, args.Error(3)
}

func (m *AdminMockOrderUsecase) GetDashboardMetrics(ctx context.Context, supplierID string) (*order.DashboardMetrics, error) {
	args := m.Called(ctx, supplierID)
	return nil, args.Error(1)
}

func (m *AdminMockOrderUsecase) GetOrderByID(ctx context.Context, orderID string) (*order.Order, error) {
	args := m.Called(ctx, orderID)
	return nil, args.Error(1)
}

func (m *AdminMockOrderUsecase) GetSoldBundleHistory(ctx context.Context, supplierID string) ([]*order.Order, map[string]string, error) {
	args := m.Called(ctx, supplierID)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]*order.Order), args.Get(1).(map[string]string), args.Error(2)
}

func (m *AdminMockOrderUsecase) GetAdminDashboardMetrics(ctx context.Context) (*admin.Metrics, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*admin.Metrics), args.Error(1)
}
func (m *AdminMockOrderUsecase) GetResellerMetrics(ctx context.Context, resellerID string) (*order.ResellerMetrics, error) {
	args := m.Called(ctx, resellerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order.ResellerMetrics), args.Error(1)
}

// -------------------- Test Suite --------------------

type AdminControllerTestSuite struct {
	suite.Suite
	controller  *AdminController
	mockUC      *MockUserUsecase
	mockOrderUC *AdminMockOrderUsecase
	router      *gin.Engine
}

func (suite *AdminControllerTestSuite) SetupTest() {
	suite.mockUC = new(MockUserUsecase)
	suite.mockOrderUC = new(AdminMockOrderUsecase)
	suite.controller = NewAdminController(suite.mockUC, suite.mockOrderUC)
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()
}

func (suite *AdminControllerTestSuite) TestGetAllUsers_NoRoleParam() {
	users := []*user.User{
		{ID: "1", Name: "User1", Role: string(user.RoleSupplier)},
		{ID: "2", Name: "User2", Role: string(user.RoleConsumer)},
	}

	suite.mockUC.On("ListByRole", mock.Anything, user.RoleSupplier).Return(users[:1], nil)
	suite.mockUC.On("ListByRole", mock.Anything, user.RoleReseller).Return([]*user.User{}, nil)
	suite.mockUC.On("ListByRole", mock.Anything, user.RoleConsumer).Return(users[1:], nil)
	suite.mockUC.On("ListByRole", mock.Anything, user.RoleAdmin).Return([]*user.User{}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/users", nil)
	suite.router.GET("/api/admin/users", suite.controller.GetAllUsers)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockUC.AssertExpectations(suite.T())
}

func (suite *AdminControllerTestSuite) TestGetAllUsers_WithRoleParam() {
	users := []*user.User{
		{ID: "1", Name: "User1", Role: string(user.RoleSupplier)},
	}
	suite.mockUC.On("ListByRole", mock.Anything, user.RoleSupplier).Return(users, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/users?role=supplier", nil)
	suite.router.GET("/api/admin/users", suite.controller.GetAllUsers)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockUC.AssertExpectations(suite.T())
}

func (suite *AdminControllerTestSuite) TestDeleteUserIfBlacklisted_Success() {
	userID := "1"
	userData := &user.User{
		ID:         userID,
		Role:       string(user.RoleSupplier),
		TrustScore: 50,
	}
	suite.mockUC.On("GetByID", mock.Anything, userID).Return(userData, nil)
	suite.mockUC.On("Update", mock.Anything, userID, map[string]interface{}{"is_deleted": true}).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/admin/users/"+userID, nil)
	suite.router.DELETE("/api/admin/users/:userId", suite.controller.DeleteUserIfBlacklisted)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockUC.AssertExpectations(suite.T())
}

func (suite *AdminControllerTestSuite) TestDeleteUserIfBlacklisted_InvalidRole() {
	userID := "1"
	userData := &user.User{
		ID:         userID,
		Role:       string(user.RoleConsumer),
		TrustScore: 50,
	}
	suite.mockUC.On("GetByID", mock.Anything, userID).Return(userData, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/admin/users/"+userID, nil)
	suite.router.DELETE("/api/admin/users/:userId", suite.controller.DeleteUserIfBlacklisted)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
	suite.mockUC.AssertExpectations(suite.T())
}

func (suite *AdminControllerTestSuite) TestGetTrustScores() {
	users := []*user.User{
		{ID: "1", Name: "User1", Role: string(user.RoleSupplier), TrustScore: 70},
		{ID: "2", Name: "User2", Role: string(user.RoleReseller), TrustScore: 40},
	}
	suite.mockUC.On("ListByRole", mock.Anything, user.RoleSupplier).Return(users[:1], nil)
	suite.mockUC.On("ListByRole", mock.Anything, user.RoleReseller).Return(users[1:], nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/users/trust-scores", nil)
	suite.router.GET("/api/admin/users/trust-scores", suite.controller.GetTrustScores)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *AdminControllerTestSuite) TestGetBlacklistedUsers() {
	users := []*user.User{
		{ID: "1", Name: "User1", Role: string(user.RoleSupplier), TrustScore: 40, IsBlacklisted: true},
	}
	suite.mockUC.On("GetBlacklistedUsers", mock.Anything).Return(users, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/blacklisted-users", nil)
	suite.router.GET("/api/admin/blacklisted-users", suite.controller.GetBlacklistedUsers)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *AdminControllerTestSuite) TestGetDashboardMetrics() {
	metrics := &admin.Metrics{
		TotalBundles:    10,
		TotalUsers:      20,
		TotalSales:      1500,
		RevenueFromFees: 30,
		SkippedClothes:  0,
	}
	suite.mockOrderUC.On("GetAdminDashboardMetrics", mock.Anything).Return(metrics, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/dashboard", nil)
	suite.router.GET("/api/admin/dashboard", suite.controller.GetDashboardMetrics)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockOrderUC.AssertExpectations(suite.T())
}

func TestAdminControllerSuite(t *testing.T) {
	suite.Run(t, new(AdminControllerTestSuite))
}
