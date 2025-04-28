package warehouse_usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/warehouse"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockRepository is a mock implementation of the warehouse.Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) AddItem(ctx context.Context, item *warehouse.WarehouseItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockRepository) GetItemsByReseller(ctx context.Context, resellerID string) ([]*warehouse.WarehouseItem, error) {
	args := m.Called(ctx, resellerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*warehouse.WarehouseItem), args.Error(1)
}

func (m *MockRepository) GetItemsByBundle(ctx context.Context, bundleID string) ([]*warehouse.WarehouseItem, error) {
	args := m.Called(ctx, bundleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*warehouse.WarehouseItem), args.Error(1)
}

func (m *MockRepository) MarkItemAsListed(ctx context.Context, itemID string) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}

func (m *MockRepository) MarkItemAsSkipped(ctx context.Context, itemID string) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}

func (m *MockRepository) DeleteItem(ctx context.Context, itemID string) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}

func (m *MockRepository) HasResellerReceivedBundle(ctx context.Context, resellerID string, bundleID string) (bool, error) {
	args := m.Called(ctx, resellerID, bundleID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) CountByStatus(ctx context.Context, status string) (int, error) {
	args := m.Called(ctx, status)
	return args.Int(0), args.Error(1)
}

// MockBundleRepository is a mock implementation of the bundle.Repository interface
type MockBundleRepository struct {
	mock.Mock
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

func (m *MockBundleRepository) CountBundles(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockBundleRepository) CreateBundle(ctx context.Context, b *bundle.Bundle) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockBundleRepository) DecreaseBundleQuantity(ctx context.Context, bundleID string) error {
	args := m.Called(ctx, bundleID)
	return args.Error(0)
}

func (m *MockBundleRepository) DeleteBundle(ctx context.Context, bundleID string) error {
	args := m.Called(ctx, bundleID)
	return args.Error(0)
}

func (m *MockBundleRepository) ListAvailableBundles(ctx context.Context) ([]*bundle.Bundle, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*bundle.Bundle), args.Error(1)
}

func (m *MockBundleRepository) ListBundles(ctx context.Context, supplierID string) ([]*bundle.Bundle, error) {
	args := m.Called(ctx, supplierID)
	return args.Get(0).([]*bundle.Bundle), args.Error(1)
}

func (m *MockBundleRepository) ListPurchasedByReseller(ctx context.Context, resellerID string) ([]*bundle.Bundle, error) {
	args := m.Called(ctx, resellerID)
	return args.Get(0).([]*bundle.Bundle), args.Error(1)
}

func (m *MockBundleRepository) MarkAsPurchased(ctx context.Context, bundleID, resellerID string) error {
	args := m.Called(ctx, bundleID, resellerID)
	return args.Error(0)
}

func (m *MockBundleRepository) UpdateBundle(ctx context.Context, bundleID string, updates map[string]interface{}) error {
	args := m.Called(ctx, bundleID, updates)
	return args.Error(0)
}

func (m *MockBundleRepository) UpdateBundleStatus(ctx context.Context, bundleID string, status string) error {
	args := m.Called(ctx, bundleID, status)
	return args.Error(0)
}

// WarehouseUsecaseTestSuite is the test suite for warehouse usecase
type WarehouseUsecaseTestSuite struct {
	suite.Suite
	mockRepo       *MockRepository
	mockBundleRepo *MockBundleRepository
	usecase        warehouse.WarehouseUseCase
	ctx            context.Context
}

// SetupTest runs before each test
func (suite *WarehouseUsecaseTestSuite) SetupTest() {
	suite.mockRepo = new(MockRepository)
	suite.mockBundleRepo = new(MockBundleRepository)
	suite.usecase = NewWarehouseUseCase(suite.mockRepo, suite.mockBundleRepo)
	suite.ctx = context.Background()
}

// TestNewWarehouseUseCase tests the constructor
func (suite *WarehouseUsecaseTestSuite) TestNewWarehouseUseCase() {
	useCase := NewWarehouseUseCase(suite.mockRepo, suite.mockBundleRepo)
	suite.NotNil(useCase)
}

// TestGetWarehouseItems_Success tests successful retrieval of warehouse items with bundle info
func (suite *WarehouseUsecaseTestSuite) TestGetWarehouseItems_Success() {
	// Arrange
	resellerID := "reseller1"
	items := []*warehouse.WarehouseItem{
		{
			ID:         "item1",
			ResellerID: resellerID,
			BundleID:   "bundle1",
			ProductID:  "product1",
			Status:     "pending",
			CreatedAt:  time.Now().Format(time.RFC3339),
		},
	}
	bundle := &bundle.Bundle{
		ID:                 "bundle1",
		Title:              "Test Bundle",
		SampleImage:        "image.jpg",
		DeclaredRating:     4,
		RemainingItemCount: 10,
	}

	suite.mockRepo.On("GetItemsByReseller", suite.ctx, resellerID).Return(items, nil)
	suite.mockBundleRepo.On("GetBundleByID", suite.ctx, "bundle1").Return(bundle, nil)

	// Act
	responses, err := suite.usecase.GetWarehouseItems(suite.ctx, resellerID)

	// Assert
	suite.NoError(err)
	suite.NotNil(responses)
	suite.Len(responses, 1)
	suite.Equal("item1", responses[0].ID)
	suite.Equal("Test Bundle", responses[0].Title)
	suite.Equal("image.jpg", responses[0].SampleImage)
	suite.Equal(float64(4), responses[0].DeclaredRating)
	suite.Equal(10, responses[0].RemainingItems)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockBundleRepo.AssertExpectations(suite.T())
}

// TestGetWarehouseItems_WithMissingBundle tests handling of items with missing bundle info
func (suite *WarehouseUsecaseTestSuite) TestGetWarehouseItems_WithMissingBundle() {
	// Arrange
	resellerID := "reseller1"
	items := []*warehouse.WarehouseItem{
		{
			ID:         "item1",
			ResellerID: resellerID,
			BundleID:   "bundle1",
			ProductID:  "product1",
			Status:     "pending",
			CreatedAt:  time.Now().Format(time.RFC3339),
		},
	}

	suite.mockRepo.On("GetItemsByReseller", suite.ctx, resellerID).Return(items, nil)
	suite.mockBundleRepo.On("GetBundleByID", suite.ctx, "bundle1").Return(nil, errors.New("bundle not found"))

	// Act
	responses, err := suite.usecase.GetWarehouseItems(suite.ctx, resellerID)

	// Assert
	suite.NoError(err)
	suite.NotNil(responses)
	suite.Empty(responses)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockBundleRepo.AssertExpectations(suite.T())
}

// TestGetWarehouseItems_EmptyList tests handling of empty list
func (suite *WarehouseUsecaseTestSuite) TestGetWarehouseItems_EmptyList() {
	// Arrange
	resellerID := "reseller1"
	items := []*warehouse.WarehouseItem{}

	suite.mockRepo.On("GetItemsByReseller", suite.ctx, resellerID).Return(items, nil)

	// Act
	responses, err := suite.usecase.GetWarehouseItems(suite.ctx, resellerID)

	// Assert
	suite.NoError(err)
	suite.NotNil(responses)
	suite.Empty(responses)
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestGetWarehouseItems_RepositoryError tests handling of repository errors
func (suite *WarehouseUsecaseTestSuite) TestGetWarehouseItems_RepositoryError() {
	// Arrange
	resellerID := "reseller1"
	expectedErr := errors.New("database error")

	suite.mockRepo.On("GetItemsByReseller", suite.ctx, resellerID).Return(nil, expectedErr)

	// Act
	responses, err := suite.usecase.GetWarehouseItems(suite.ctx, resellerID)

	// Assert
	suite.Error(err)
	suite.Equal(expectedErr, err)
	suite.Nil(responses)
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestGetWarehouseItems_MultipleItems tests handling of multiple items
func (suite *WarehouseUsecaseTestSuite) TestGetWarehouseItems_MultipleItems() {
	// Arrange
	resellerID := "reseller1"
	items := []*warehouse.WarehouseItem{
		{
			ID:         "item1",
			ResellerID: resellerID,
			BundleID:   "bundle1",
			ProductID:  "product1",
			Status:     "pending",
			CreatedAt:  time.Now().Format(time.RFC3339),
		},
		{
			ID:         "item2",
			ResellerID: resellerID,
			BundleID:   "bundle2",
			ProductID:  "product2",
			Status:     "pending",
			CreatedAt:  time.Now().Format(time.RFC3339),
		},
	}
	bundle1 := &bundle.Bundle{
		ID:                 "bundle1",
		Title:              "Test Bundle 1",
		SampleImage:        "image1.jpg",
		DeclaredRating:     4,
		RemainingItemCount: 10,
	}
	bundle2 := &bundle.Bundle{
		ID:                 "bundle2",
		Title:              "Test Bundle 2",
		SampleImage:        "image2.jpg",
		DeclaredRating:     4,
		RemainingItemCount: 5,
	}

	suite.mockRepo.On("GetItemsByReseller", suite.ctx, resellerID).Return(items, nil)
	suite.mockBundleRepo.On("GetBundleByID", suite.ctx, "bundle1").Return(bundle1, nil)
	suite.mockBundleRepo.On("GetBundleByID", suite.ctx, "bundle2").Return(bundle2, nil)

	// Act
	responses, err := suite.usecase.GetWarehouseItems(suite.ctx, resellerID)

	// Assert
	suite.NoError(err)
	suite.NotNil(responses)
	suite.Len(responses, 2)
	suite.Equal("item1", responses[0].ID)
	suite.Equal("Test Bundle 1", responses[0].Title)
	suite.Equal("item2", responses[1].ID)
	suite.Equal("Test Bundle 2", responses[1].Title)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockBundleRepo.AssertExpectations(suite.T())
}

// TestWarehouseUsecase runs the test suite
func TestWarehouseUsecase(t *testing.T) {
	suite.Run(t, new(WarehouseUsecaseTestSuite))
}
