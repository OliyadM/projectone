package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/models"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/models/common"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockBundleUsecase struct {
	mock.Mock
}

func (m *MockBundleUsecase) CreateBundle(ctx context.Context, supplierID string, b *bundle.Bundle) error {
	args := m.Called(ctx, supplierID, b)
	return args.Error(0)
}

func (m *MockBundleUsecase) ListBundles(ctx context.Context, supplierID string) ([]*bundle.Bundle, error) {
	args := m.Called(ctx, supplierID)
	return args.Get(0).([]*bundle.Bundle), args.Error(1)
}

func (m *MockBundleUsecase) DeleteBundle(ctx context.Context, supplierID, bundleID string) error {
	args := m.Called(ctx, supplierID, bundleID)
	return args.Error(0)
}

func (m *MockBundleUsecase) UpdateBundle(ctx context.Context, supplierID, bundleID string, updates map[string]interface{}) error {
	args := m.Called(ctx, supplierID, bundleID, updates)
	return args.Error(0)
}

func (m *MockBundleUsecase) GetBundleByID(ctx context.Context, supplierID, bundleID string) (*bundle.Bundle, error) {
	args := m.Called(ctx, supplierID, bundleID)
	return args.Get(0).(*bundle.Bundle), args.Error(1)
}

func (m *MockBundleUsecase) GetBundlePublicByID(ctx context.Context, bundleID string) (*bundle.Bundle, error) {
	args := m.Called(ctx, bundleID)
	return args.Get(0).(*bundle.Bundle), args.Error(1)
}

func (m *MockBundleUsecase) ListAvailableBundles(ctx context.Context) ([]*bundle.Bundle, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*bundle.Bundle), args.Error(1)
}

func (m *MockBundleUsecase) DecreaseRemainingItemCount(ctx context.Context, bundleID string) error {
	args := m.Called(ctx, bundleID)
	return args.Error(0)
}

func (m *MockBundleUsecase) GetBundleByTitle(ctx context.Context, title string) (*bundle.Bundle, error) {
	args := m.Called(ctx, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bundle.Bundle), args.Error(1)
}

type BundleControllerTestSuite struct {
	suite.Suite
	controller    *BundleController
	mockBundleUC  *MockBundleUsecase
	mockUserUC    *MockUserUsecase
	router        *gin.Engine
	supplierID    string
	supplierToken string
}

func (suite *BundleControllerTestSuite) SetupTest() {
	suite.mockBundleUC = new(MockBundleUsecase)
	suite.mockUserUC = new(MockUserUsecase)
	suite.controller = NewBundleController(suite.mockBundleUC, suite.mockUserUC)
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()
	suite.supplierID = "supplier123"
	suite.supplierToken = "supplier-token"

	// Add auth middleware
	suite.router.Use(func(c *gin.Context) {
		c.Set("userID", suite.supplierID)
		c.Set("role", "supplier")
	})
}

func (suite *BundleControllerTestSuite) TestCreateBundle_Success() {
	// Setup
	req := models.CreateBundleRequest{
		Title:              "Test Bundle",
		Description:        "Test Description",
		SampleImage:        "test.jpg",
		NumberOfItems:      10,
		Grade:              "A",
		Type:               "basic",
		EstimatedBreakdown: map[string]int{"shirts": 5, "pants": 5},
		Price:              100.0,
		DeclaredRating:     4,
	}

	user := &user.User{
		ID:            suite.supplierID,
		IsBlacklisted: false,
	}

	suite.mockUserUC.On("GetByID", mock.Anything, suite.supplierID).Return(user, nil)
	suite.mockBundleUC.On("CreateBundle", mock.Anything, suite.supplierID, mock.Anything).Return(nil)

	// Execute
	jsonData, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/bundles", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+suite.supplierToken)
	suite.router.POST("/bundles", suite.controller.CreateBundle)
	suite.router.ServeHTTP(w, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var response common.APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(suite.T(), response.Success)
	suite.mockUserUC.AssertExpectations(suite.T())
	suite.mockBundleUC.AssertExpectations(suite.T())
}

func (suite *BundleControllerTestSuite) TestCreateBundle_BlacklistedUser() {
	// Setup
	req := models.CreateBundleRequest{
		Title:              "Test Bundle",
		Description:        "Test Description",
		SampleImage:        "test.jpg",
		NumberOfItems:      10,
		Grade:              "A",
		Type:               "basic",
		EstimatedBreakdown: map[string]int{"shirts": 5, "pants": 5},
		Price:              100.0,
		DeclaredRating:     4,
	}

	user := &user.User{
		ID:            suite.supplierID,
		IsBlacklisted: true,
	}

	suite.mockUserUC.On("GetByID", mock.Anything, suite.supplierID).Return(user, nil)

	// Execute
	jsonData, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/bundles", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+suite.supplierToken)
	suite.router.POST("/bundles", suite.controller.CreateBundle)
	suite.router.ServeHTTP(w, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
	var response common.APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(suite.T(), response.Success)
	suite.mockUserUC.AssertExpectations(suite.T())
}

func (suite *BundleControllerTestSuite) TestListBundles_Success() {
	// Setup
	bundles := []*bundle.Bundle{
		{
			ID:           "bundle1",
			Title:        "Test Bundle 1",
			Grade:        "A",
			Price:        100.0,
			SortingLevel: "basic",
			Status:       "available",
		},
	}

	suite.mockBundleUC.On("ListBundles", mock.Anything, suite.supplierID).Return(bundles, nil)

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/bundles", nil)
	req.Header.Set("Authorization", "Bearer "+suite.supplierToken)
	suite.router.GET("/bundles", suite.controller.ListBundles)
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var response common.APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(suite.T(), response.Success)
	suite.mockBundleUC.AssertExpectations(suite.T())
}

func (suite *BundleControllerTestSuite) TestDeleteBundle_Success() {
	// Setup
	bundleID := "bundle123"
	suite.mockBundleUC.On("DeleteBundle", mock.Anything, suite.supplierID, bundleID).Return(nil)

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/bundles/"+bundleID, nil)
	req.Header.Set("Authorization", "Bearer "+suite.supplierToken)
	suite.router.DELETE("/bundles/:id", suite.controller.DeleteBundle)
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var response common.APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(suite.T(), response.Success)
	suite.mockBundleUC.AssertExpectations(suite.T())
}

func (suite *BundleControllerTestSuite) TestUpdateBundle_Success() {
	// Setup
	bundleID := "bundle123"
	updates := map[string]interface{}{
		"title": "Updated Title",
		"price": 150.0,
	}

	updatedBundle := &bundle.Bundle{
		ID:           bundleID,
		Title:        "Updated Title",
		Grade:        "A",
		Price:        150.0,
		SortingLevel: "basic",
		Status:       "available",
	}

	suite.mockBundleUC.On("UpdateBundle", mock.Anything, suite.supplierID, bundleID, updates).Return(nil)
	suite.mockBundleUC.On("GetBundleByID", mock.Anything, suite.supplierID, bundleID).Return(updatedBundle, nil)

	// Execute
	jsonData, _ := json.Marshal(updates)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/bundles/"+bundleID, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.supplierToken)
	suite.router.PUT("/bundles/:id", suite.controller.UpdateBundle)
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var response common.APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(suite.T(), response.Success)
	suite.mockBundleUC.AssertExpectations(suite.T())
}

func (suite *BundleControllerTestSuite) TestGetBundle_Success() {
	// Setup
	bundleID := "bundle123"
	b := &bundle.Bundle{
		ID:           bundleID,
		Title:        "Test Bundle",
		Grade:        "A",
		Price:        100.0,
		SortingLevel: "basic",
		Status:       "available",
	}

	suite.mockBundleUC.On("GetBundleByID", mock.Anything, suite.supplierID, bundleID).Return(b, nil)

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/bundles/"+bundleID, nil)
	req.Header.Set("Authorization", "Bearer "+suite.supplierToken)
	suite.router.GET("/bundles/:id", suite.controller.GetBundle)
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var response common.APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(suite.T(), response.Success)
	suite.mockBundleUC.AssertExpectations(suite.T())
}

func (suite *BundleControllerTestSuite) TestListAvailableBundles_Success() {
	// Setup
	bundles := []*bundle.Bundle{
		{
			ID:           "bundle1",
			Title:        "Test Bundle 1",
			Grade:        "A",
			Price:        100.0,
			SortingLevel: "basic",
			Status:       "available",
		},
	}

	suite.mockBundleUC.On("ListAvailableBundles", mock.Anything).Return(bundles, nil)

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/bundles/available", nil)
	suite.router.GET("/bundles/available", suite.controller.ListAvailableBundles)
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockBundleUC.AssertExpectations(suite.T())
}

func TestBundleControllerSuite(t *testing.T) {
	suite.Run(t, new(BundleControllerTestSuite))
}
