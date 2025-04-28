package productusecase

import (
	"context"
	"errors"
	"testing"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/product"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ---------------- Mock Implementations ----------------

type MockRepository struct {
	mock.Mock
}

// GetProductByTitle implements product.Repository.
func (m *MockRepository) GetProductByTitle(ctx context.Context, title string) (*product.Product, error) {
	args := m.Called(ctx, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Product), args.Error(1)
}

func (m *MockRepository) AddProduct(ctx context.Context, p *product.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockRepository) GetProductByID(ctx context.Context, id string) (*product.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Product), args.Error(1)
}

func (m *MockRepository) ListProductsByReseller(ctx context.Context, resellerID string, page, limit int) ([]*product.Product, error) {
	args := m.Called(ctx, resellerID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*product.Product), args.Error(1)
}

func (m *MockRepository) ListAvailableProducts(ctx context.Context, page, limit int) ([]*product.Product, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*product.Product), args.Error(1)
}

func (m *MockRepository) GetProductsByBundleID(ctx context.Context, bundleID string) ([]*product.Product, error) {
	args := m.Called(ctx, bundleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*product.Product), args.Error(1)
}

func (m *MockRepository) DeleteProduct(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) UpdateProduct(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockRepository) GetSoldProductsByReseller(ctx context.Context, resellerID string) ([]*product.Product, error) {
	args := m.Called(ctx, resellerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*product.Product), args.Error(1)
}

type MockBundleRepository struct {
	mock.Mock
}

func (m *MockBundleRepository) CreateBundle(ctx context.Context, b *bundle.Bundle) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockBundleRepository) GetBundleByID(ctx context.Context, id string) (*bundle.Bundle, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bundle.Bundle), args.Error(1)
}

func (m *MockBundleRepository) GetBundleByTitle(ctx context.Context, title string) (*bundle.Bundle, error) {
	args := m.Called(ctx, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bundle.Bundle), args.Error(1)
}

func (m *MockBundleRepository) ListBundles(ctx context.Context, supplierID string) ([]*bundle.Bundle, error) {
	args := m.Called(ctx, supplierID)
	return args.Get(0).([]*bundle.Bundle), args.Error(1)
}

func (m *MockBundleRepository) ListAvailableBundles(ctx context.Context) ([]*bundle.Bundle, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*bundle.Bundle), args.Error(1)
}

func (m *MockBundleRepository) ListPurchasedByReseller(ctx context.Context, resellerID string) ([]*bundle.Bundle, error) {
	args := m.Called(ctx, resellerID)
	return args.Get(0).([]*bundle.Bundle), args.Error(1)
}

func (m *MockBundleRepository) UpdateBundleStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockBundleRepository) MarkAsPurchased(ctx context.Context, bundleID string, resellerID string) error {
	args := m.Called(ctx, bundleID, resellerID)
	return args.Error(0)
}

func (m *MockBundleRepository) DeleteBundle(ctx context.Context, bundleID string) error {
	args := m.Called(ctx, bundleID)
	return args.Error(0)
}

func (m *MockBundleRepository) UpdateBundle(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockBundleRepository) DecreaseBundleQuantity(ctx context.Context, bundleID string) error {
	args := m.Called(ctx, bundleID)
	return args.Error(0)
}

// ---------------- Test Suite ----------------

type ProductUsecaseTestSuite struct {
	suite.Suite
	mockRepo       *MockRepository
	mockBundleRepo *MockBundleRepository
	usecase        product.Usecase
}

func (suite *ProductUsecaseTestSuite) SetupTest() {
	suite.mockRepo = new(MockRepository)
	suite.mockBundleRepo = new(MockBundleRepository)
	suite.usecase = NewProductUsecase(suite.mockRepo, suite.mockBundleRepo)
}

func (suite *ProductUsecaseTestSuite) TestAddProduct_Success() {
	ctx := context.Background()
	resellerID := primitive.NewObjectID()
	product := &product.Product{
		ID:         "test-id",
		Title:      "Test Product",
		ResellerID: resellerID,
	}

	suite.mockRepo.On("AddProduct", ctx, product).Return(nil)
	err := suite.usecase.AddProduct(ctx, product)
	suite.NoError(err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductUsecaseTestSuite) TestAddProduct_WithBundle_Success() {
	ctx := context.Background()
	resellerID := primitive.NewObjectID()
	product := &product.Product{
		ID:         "test-id",
		Title:      "Test Product",
		ResellerID: resellerID,
		BundleID:   "bundle-1",
	}
	bundle := &bundle.Bundle{
		ID:       "bundle-1",
		Quantity: 5,
	}

	suite.mockBundleRepo.On("GetBundleByID", ctx, "bundle-1").Return(bundle, nil)
	suite.mockRepo.On("AddProduct", ctx, product).Return(nil)
	suite.mockBundleRepo.On("UpdateBundle", ctx, "bundle-1", map[string]interface{}{"quantity": 4}).Return(nil)

	err := suite.usecase.AddProduct(ctx, product)
	suite.NoError(err)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockBundleRepo.AssertExpectations(suite.T())
}

func (suite *ProductUsecaseTestSuite) TestAddProduct_WithBundle_OutOfStock() {
	ctx := context.Background()
	resellerID := primitive.NewObjectID()
	product := &product.Product{
		ID:         "test-id",
		Title:      "Test Product",
		ResellerID: resellerID,
		BundleID:   "bundle-1",
	}
	bundle := &bundle.Bundle{
		ID:       "bundle-1",
		Quantity: 0,
	}

	suite.mockBundleRepo.On("GetBundleByID", ctx, "bundle-1").Return(bundle, nil)

	err := suite.usecase.AddProduct(ctx, product)
	suite.Error(err)
	suite.Equal("bundle is out of stock", err.Error())
	suite.mockBundleRepo.AssertExpectations(suite.T())
}

func (suite *ProductUsecaseTestSuite) TestGetProductByID_Success() {
	ctx := context.Background()
	resellerID := primitive.NewObjectID()
	expectedProduct := &product.Product{
		ID:         "test-id",
		Title:      "Test Product",
		ResellerID: resellerID,
	}

	suite.mockRepo.On("GetProductByID", ctx, "test-id").Return(expectedProduct, nil)
	product, err := suite.usecase.GetProductByID(ctx, "test-id")
	suite.NoError(err)
	suite.Equal(expectedProduct, product)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductUsecaseTestSuite) TestGetProductByID_NotFound() {
	ctx := context.Background()

	suite.mockRepo.On("GetProductByID", ctx, "non-existent-id").Return(nil, errors.New("product not found"))
	product, err := suite.usecase.GetProductByID(ctx, "non-existent-id")
	suite.Error(err)
	suite.Nil(product)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductUsecaseTestSuite) TestGetProductByTitle_Success() {
	ctx := context.Background()
	resellerID := primitive.NewObjectID()
	expectedProduct := &product.Product{
		ID:         "test-id",
		Title:      "Test Product",
		ResellerID: resellerID,
	}

	suite.mockRepo.On("GetProductByTitle", ctx, "Test Product").Return(expectedProduct, nil)
	product, err := suite.usecase.GetProductByTitle(ctx, "Test Product")
	suite.NoError(err)
	suite.Equal(expectedProduct, product)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductUsecaseTestSuite) TestGetProductByTitle_EmptyTitle() {
	ctx := context.Background()

	product, err := suite.usecase.GetProductByTitle(ctx, "")
	suite.Error(err)
	suite.Equal("title cannot be empty", err.Error())
	suite.Nil(product)
}

func (suite *ProductUsecaseTestSuite) TestListProductsByReseller_Success() {
	ctx := context.Background()
	resellerID := primitive.NewObjectID()
	expectedProducts := []*product.Product{
		{
			ID:         "test-id-1",
			Title:      "Test Product 1",
			ResellerID: resellerID,
		},
		{
			ID:         "test-id-2",
			Title:      "Test Product 2",
			ResellerID: resellerID,
		},
	}

	suite.mockRepo.On("ListProductsByReseller", ctx, resellerID.Hex(), 1, 10).Return(expectedProducts, nil)
	products, err := suite.usecase.ListProductsByReseller(ctx, resellerID.Hex(), 1, 10)
	suite.NoError(err)
	suite.Equal(expectedProducts, products)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductUsecaseTestSuite) TestListAvailableProducts_Success() {
	ctx := context.Background()
	resellerID1 := primitive.NewObjectID()
	resellerID2 := primitive.NewObjectID()
	expectedProducts := []*product.Product{
		{
			ID:         "test-id-1",
			Title:      "Test Product 1",
			ResellerID: resellerID1,
		},
		{
			ID:         "test-id-2",
			Title:      "Test Product 2",
			ResellerID: resellerID2,
		},
	}

	suite.mockRepo.On("ListAvailableProducts", ctx, 1, 10).Return(expectedProducts, nil)
	products, err := suite.usecase.ListAvailableProducts(ctx, 1, 10)
	suite.NoError(err)
	suite.Equal(expectedProducts, products)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductUsecaseTestSuite) TestDeleteProduct_Success() {
	ctx := context.Background()

	suite.mockRepo.On("DeleteProduct", ctx, "test-id").Return(nil)
	err := suite.usecase.DeleteProduct(ctx, "test-id")
	suite.NoError(err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductUsecaseTestSuite) TestUpdateProduct_Success() {
	ctx := context.Background()
	updates := map[string]interface{}{
		"title": "Updated Title",
		"price": 99.99,
	}

	suite.mockRepo.On("UpdateProduct", ctx, "test-id", updates).Return(nil)
	err := suite.usecase.UpdateProduct(ctx, "test-id", updates)
	suite.NoError(err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestProductUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(ProductUsecaseTestSuite))
}

func (m *MockBundleRepository) CountBundles(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}
