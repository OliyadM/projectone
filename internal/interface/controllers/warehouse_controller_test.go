package controllers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockWarehouseUsecase struct {
	mock.Mock
}

func (m *MockWarehouseUsecase) GetWarehouseItems(ctx context.Context, resellerID string) ([]*models.WarehouseItemResponse, error) {
	args := m.Called(ctx, resellerID)
	return args.Get(0).([]*models.WarehouseItemResponse), args.Error(1)
}

type WarehouseControllerTestSuite struct {
	suite.Suite
	usecase    *MockWarehouseUsecase
	controller *WarehouseController
	router     *gin.Engine
}

func (suite *WarehouseControllerTestSuite) SetupTest() {
	suite.usecase = new(MockWarehouseUsecase)
	suite.controller = NewWarehouseController(suite.usecase)
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()
}

func TestWarehouseControllerTestSuite(t *testing.T) {
	suite.Run(t, new(WarehouseControllerTestSuite))
}

func (suite *WarehouseControllerTestSuite) TestGetWarehouseItems_Success() {
	// Setup
	expectedItems := []*models.WarehouseItemResponse{
		{
			ID:             "item1",
			ResellerID:     "reseller123",
			BundleID:       "bundle1",
			Status:         "listed",
			CreatedAt:      "2024-01-01",
			Title:          "Winter Bundle",
			SampleImage:    "https://example.com/image.jpg",
			DeclaredRating: 85,
			RemainingItems: 5,
		},
		{
			ID:             "item2",
			ResellerID:     "reseller123",
			BundleID:       "bundle2",
			Status:         "pending",
			CreatedAt:      "2024-01-02",
			Title:          "Summer Mix",
			SampleImage:    "https://example.com/summer.jpg",
			DeclaredRating: 90,
			RemainingItems: 3,
		},
	}

	suite.usecase.On("GetWarehouseItems", mock.Anything, "reseller123").
		Return(expectedItems, nil)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "reseller123")
	c.Request = httptest.NewRequest("GET", "/warehouse/items", nil)

	// Execute
	suite.controller.GetWarehouseItems(c)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.usecase.AssertExpectations(suite.T())
}

func (suite *WarehouseControllerTestSuite) TestGetWarehouseItems_Unauthorized() {
	// Setup
	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Note: No userID set
	c.Request = httptest.NewRequest("GET", "/warehouse/items", nil)

	// Execute
	suite.controller.GetWarehouseItems(c)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	suite.usecase.AssertNotCalled(suite.T(), "GetWarehouseItems")
}

func (suite *WarehouseControllerTestSuite) TestGetWarehouseItems_InvalidUserID() {
	// Setup
	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", 123) // Invalid type
	c.Request = httptest.NewRequest("GET", "/warehouse/items", nil)

	// Execute
	suite.controller.GetWarehouseItems(c)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	suite.usecase.AssertNotCalled(suite.T(), "GetWarehouseItems")
}

func (suite *WarehouseControllerTestSuite) TestGetWarehouseItems_EmptyUserID() {
	// Setup
	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "") // Empty string
	c.Request = httptest.NewRequest("GET", "/warehouse/items", nil)

	// Execute
	suite.controller.GetWarehouseItems(c)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	suite.usecase.AssertNotCalled(suite.T(), "GetWarehouseItems")
}

func (suite *WarehouseControllerTestSuite) TestGetWarehouseItems_UseCaseError() {
	// Setup
	suite.usecase.On("GetWarehouseItems", mock.Anything, "reseller123").
		Return([]*models.WarehouseItemResponse{}, assert.AnError) // Return an empty slice instead of nil

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "reseller123")
	c.Request = httptest.NewRequest("GET", "/warehouse/items", nil)

	// Execute
	suite.controller.GetWarehouseItems(c)

	// Assert
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	suite.usecase.AssertExpectations(suite.T())
}
