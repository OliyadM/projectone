package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/product"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/warehouse"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type ProductControllerTestSuite struct {
	suite.Suite
	productUseCase *MockProductUseCase
	trustUseCase   *MockTrustUseCase
	bundleUseCase  *MockBundleUseCase
	warehouseRepo  *MockWarehouseRepo
	controller     *ProductController
	router         *gin.Engine
}

type MockProductUseCase struct {
	mock.Mock
}

// GetProductByTitle implements product.Usecase.
func (m *MockProductUseCase) GetProductByTitle(ctx context.Context, title string) (*product.Product, error) {
	args := m.Called(ctx, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Product), args.Error(1)
}

func (m *MockProductUseCase) AddProduct(ctx context.Context, p *product.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockProductUseCase) GetProductByID(ctx context.Context, id string) (*product.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Product), args.Error(1)
}

func (m *MockProductUseCase) ListAvailableProducts(ctx context.Context, page, limit int) ([]*product.Product, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*product.Product), args.Error(1)
}

func (m *MockProductUseCase) ListProductsByReseller(ctx context.Context, resellerID string, page, limit int) ([]*product.Product, error) {
	args := m.Called(ctx, resellerID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*product.Product), args.Error(1)
}

func (m *MockProductUseCase) UpdateProduct(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockProductUseCase) DeleteProduct(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockTrustUseCase struct {
	mock.Mock
}

func (m *MockTrustUseCase) UpdateSupplierTrustScoreOnNewRating(ctx context.Context, supplierID string, declaredRating, actualRating float64) error {
	args := m.Called(ctx, supplierID, declaredRating, actualRating)
	return args.Error(0)
}

func (m *MockTrustUseCase) UpdateResellerTrustScoreOnNewRating(ctx context.Context, resellerID string, declaredRating, actualRating float64) error {
	args := m.Called(ctx, resellerID, declaredRating, actualRating)
	return args.Error(0)
}

type MockBundleUseCase struct {
	mock.Mock
}

func (m *MockBundleUseCase) CreateBundle(ctx context.Context, supplierID string, bundle *bundle.Bundle) error {
	args := m.Called(ctx, supplierID, bundle)
	return args.Error(0)
}

func (m *MockBundleUseCase) DeleteBundle(ctx context.Context, supplierID, bundleID string) error {
	args := m.Called(ctx, supplierID, bundleID)
	return args.Error(0)
}

func (m *MockBundleUseCase) GetBundleByID(ctx context.Context, supplierID, bundleID string) (*bundle.Bundle, error) {
	args := m.Called(ctx, supplierID, bundleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bundle.Bundle), args.Error(1)
}

func (m *MockBundleUseCase) GetBundlePublicByID(ctx context.Context, id string) (*bundle.Bundle, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bundle.Bundle), args.Error(1)
}

func (m *MockBundleUseCase) ListAvailableBundles(ctx context.Context) ([]*bundle.Bundle, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bundle.Bundle), args.Error(1)
}

func (m *MockBundleUseCase) ListBundles(ctx context.Context, supplierID string) ([]*bundle.Bundle, error) {
	args := m.Called(ctx, supplierID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bundle.Bundle), args.Error(1)
}

func (m *MockBundleUseCase) UpdateBundle(ctx context.Context, supplierID string, bundleID string, updates map[string]interface{}) error {
	args := m.Called(ctx, supplierID, bundleID, updates)
	return args.Error(0)
}

func (m *MockBundleUseCase) DecreaseRemainingItemCount(ctx context.Context, bundleID string) error {
	args := m.Called(ctx, bundleID)
	return args.Error(0)
}

func (m *MockBundleUseCase) GetBundleByTitle(ctx context.Context, title string) (*bundle.Bundle, error) {
	args := m.Called(ctx, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bundle.Bundle), args.Error(1)
}

type MockWarehouseRepo struct {
	mock.Mock
}

func (m *MockWarehouseRepo) HasResellerReceivedBundle(ctx context.Context, resellerID, bundleID string) (bool, error) {
	args := m.Called(ctx, resellerID, bundleID)
	return args.Bool(0), args.Error(1)
}

func (m *MockWarehouseRepo) AddItem(ctx context.Context, item *warehouse.WarehouseItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockWarehouseRepo) DeleteItem(ctx context.Context, itemID string) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}

func (m *MockWarehouseRepo) GetItemsByBundle(ctx context.Context, bundleID string) ([]*warehouse.WarehouseItem, error) {
	args := m.Called(ctx, bundleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*warehouse.WarehouseItem), args.Error(1)
}

func (m *MockWarehouseRepo) GetItemsByReseller(ctx context.Context, resellerID string) ([]*warehouse.WarehouseItem, error) {
	args := m.Called(ctx, resellerID)
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

func (suite *ProductControllerTestSuite) SetupTest() {
	suite.productUseCase = new(MockProductUseCase)
	suite.trustUseCase = new(MockTrustUseCase)
	suite.bundleUseCase = new(MockBundleUseCase)
	suite.warehouseRepo = new(MockWarehouseRepo)
	suite.controller = NewProductController(
		suite.productUseCase,
		suite.trustUseCase,
		suite.bundleUseCase,
		suite.warehouseRepo,
	)
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()
}

func TestProductControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ProductControllerTestSuite))
}

func (suite *ProductControllerTestSuite) TestCreate_Success() {
	// Setup
	userID := primitive.NewObjectID()
	product := &product.Product{
		BundleID: "bundle123",
		Rating:   4.5,
	}
	bundle := &bundle.Bundle{
		SupplierID:         "supplier123",
		DeclaredRating:     4.0,
		RemainingItemCount: 5,
	}

	suite.warehouseRepo.On("HasResellerReceivedBundle", mock.Anything, userID.Hex(), product.BundleID).
		Return(true, nil)
	suite.bundleUseCase.On("GetBundlePublicByID", mock.Anything, product.BundleID).
		Return(bundle, nil)
	suite.productUseCase.On("AddProduct", mock.Anything, mock.Anything).
		Return(nil)
	suite.bundleUseCase.On("DecreaseRemainingItemCount", mock.Anything, product.BundleID).
		Return(nil)
	suite.trustUseCase.On("UpdateSupplierTrustScoreOnNewRating", mock.Anything, bundle.SupplierID, float64(bundle.DeclaredRating), product.Rating).
		Return(nil).Run(func(args mock.Arguments) {
		// Add a small delay to allow the goroutine to complete
		time.Sleep(100 * time.Millisecond)
	})

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", userID.Hex())

	body, _ := json.Marshal(product)
	req := httptest.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	// Execute
	suite.controller.Create(c)

	// Add a small delay to allow the goroutine to complete
	time.Sleep(200 * time.Millisecond)

	// Assert
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	suite.warehouseRepo.AssertExpectations(suite.T())
	suite.bundleUseCase.AssertExpectations(suite.T())
	suite.productUseCase.AssertExpectations(suite.T())
	suite.trustUseCase.AssertExpectations(suite.T())
}

func (suite *ProductControllerTestSuite) TestCreate_InvalidPayload() {
	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal("invalid")
	c.Request = httptest.NewRequest("POST", "/products", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	suite.controller.Create(c)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *ProductControllerTestSuite) TestCreate_Unauthorized() {
	// Setup
	product := &product.Product{
		BundleID: "bundle123",
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(product)
	c.Request = httptest.NewRequest("POST", "/products", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	suite.controller.Create(c)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *ProductControllerTestSuite) TestGetByID_Success() {
	// Setup
	expectedProduct := &product.Product{ID: "product123"}
	suite.productUseCase.On("GetProductByID", mock.Anything, "product123").
		Return(expectedProduct, nil)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "product123"}}
	c.Request = httptest.NewRequest("GET", "/products/product123", nil)

	// Execute
	suite.controller.GetByID(c)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.productUseCase.AssertExpectations(suite.T())
}

func (suite *ProductControllerTestSuite) TestGetByID_NotFound() {
	// Setup
	suite.productUseCase.On("GetProductByID", mock.Anything, "product123").
		Return(nil, ErrProductNotFound)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "product123"}}
	c.Request = httptest.NewRequest("GET", "/products/product123", nil)

	// Execute
	suite.controller.GetByID(c)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
	suite.productUseCase.AssertExpectations(suite.T())
}

func (suite *ProductControllerTestSuite) TestListAvailable_Success() {
	// Setup
	expectedProducts := []*product.Product{
		{ID: "product1"},
		{ID: "product2"},
	}
	suite.productUseCase.On("ListAvailableProducts", mock.Anything, 1, 10).
		Return(expectedProducts, nil)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/products?page=1&limit=10", nil)

	// Execute
	suite.controller.ListAvailable(c)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.productUseCase.AssertExpectations(suite.T())
}

func (suite *ProductControllerTestSuite) TestListByReseller_Success() {
	// Setup
	expectedProducts := []*product.Product{
		{ID: "product1"},
		{ID: "product2"},
	}
	suite.productUseCase.On("ListProductsByReseller", mock.Anything, "reseller123", 1, 10).
		Return(expectedProducts, nil)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "reseller123"}}
	c.Request = httptest.NewRequest("GET", "/resellers/reseller123/products?page=1&limit=10", nil)

	// Execute
	suite.controller.ListByReseller(c)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.productUseCase.AssertExpectations(suite.T())
}

func (suite *ProductControllerTestSuite) TestUpdate_Success() {
	// Setup
	updates := map[string]interface{}{
		"name": "Updated Product",
	}
	suite.productUseCase.On("UpdateProduct", mock.Anything, "product123", updates).
		Return(nil)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "product123"}}

	body, _ := json.Marshal(updates)
	c.Request = httptest.NewRequest("PUT", "/products/product123", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	suite.controller.Update(c)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.productUseCase.AssertExpectations(suite.T())
}

func (suite *ProductControllerTestSuite) TestDelete_Success() {
	// Setup
	suite.productUseCase.On("DeleteProduct", mock.Anything, "product123").
		Return(nil)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "product123"}}
	c.Request = httptest.NewRequest("DELETE", "/products/product123", nil)

	// Execute
	suite.controller.Delete(c)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.productUseCase.AssertExpectations(suite.T())
}

func (suite *ProductControllerTestSuite) TestGetByTitle_Success() {
	// Setup
	expectedProduct := &product.Product{Title: "Test Product"}
	suite.productUseCase.On("GetProductByTitle", mock.Anything, "Test Product").
		Return(expectedProduct, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "title", Value: "Test Product"}}
	c.Request = httptest.NewRequest("GET", "/products/title/Test%20Product", nil) // Properly encode the title

	suite.controller.GetByTitle(c)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.productUseCase.AssertExpectations(suite.T())
}

func (suite *ProductControllerTestSuite) TestGetByTitle_NotFound() {
	suite.productUseCase.On("GetProductByTitle", mock.Anything, "Nonexistent Product").
		Return(nil, mongo.ErrNoDocuments)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "title", Value: "Nonexistent Product"}}
	c.Request = httptest.NewRequest("GET", "/products/title/Nonexistent%20Product", nil)

	suite.controller.GetByTitle(c)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
	suite.productUseCase.AssertExpectations(suite.T())
}

func (suite *ProductControllerTestSuite) TestGetByTitle_EmptyTitle() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "title", Value: ""}}
	c.Request = httptest.NewRequest("GET", "/products/title/", nil)

	suite.controller.GetByTitle(c)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}
