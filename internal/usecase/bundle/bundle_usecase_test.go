package bundle

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
)

// MockRepository is a mock implementation of the bundle.Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateBundle(ctx context.Context, b *bundle.Bundle) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockRepository) GetBundleByID(ctx context.Context, id string) (*bundle.Bundle, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bundle.Bundle), args.Error(1)
}

func (m *MockRepository) ListBundles(ctx context.Context, supplierID string) ([]*bundle.Bundle, error) {
	args := m.Called(ctx, supplierID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bundle.Bundle), args.Error(1)
}

func (m *MockRepository) ListAvailableBundles(ctx context.Context) ([]*bundle.Bundle, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bundle.Bundle), args.Error(1)
}

func (m *MockRepository) ListPurchasedByReseller(ctx context.Context, resellerID string) ([]*bundle.Bundle, error) {
	args := m.Called(ctx, resellerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bundle.Bundle), args.Error(1)
}

func (m *MockRepository) UpdateBundleStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockRepository) MarkAsPurchased(ctx context.Context, bundleID string, resellerID string) error {
	args := m.Called(ctx, bundleID, resellerID)
	return args.Error(0)
}

func (m *MockRepository) DeleteBundle(ctx context.Context, bundleID string) error {
	args := m.Called(ctx, bundleID)
	return args.Error(0)
}
func (m *MockRepository) CountBundles(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) UpdateBundle(ctx context.Context, id string, updatedData map[string]interface{}) error {
	args := m.Called(ctx, id, updatedData)
	return args.Error(0)
}

func (m *MockRepository) DecreaseBundleQuantity(ctx context.Context, bundleID string) error {
	args := m.Called(ctx, bundleID)
	return args.Error(0)
}

func (m *MockRepository) GetBundleByTitle(ctx context.Context, title string) (*bundle.Bundle, error) {
	args := m.Called(ctx, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bundle.Bundle), args.Error(1)
}

// Helper function to create a test bundle
func createTestBundle(supplierID string) *bundle.Bundle {
	return &bundle.Bundle{
		ID:                 "test-bundle-id",
		SupplierID:         supplierID,
		Title:              "Test Bundle",
		Description:        "Test Description",
		SampleImage:        "test-image.jpg",
		Quantity:           10,
		Grade:              "A",
		SortingLevel:       bundle.Sorted,
		EstimatedBreakdown: map[string]int{"shirts": 5, "pants": 5},
		Type:               "clothing",
		Price:              100.00,
		Status:             "available",
		CreatedAt:          time.Now().Format(time.RFC3339),
		DateListed:         time.Now(),
		DeclaredRating:     4,
		EstimatedItemCount: 10,
		RemainingItemCount: 10,
	}
}

// ---------------- Test Suite ----------------

type BundleUsecaseTestSuite struct {
	suite.Suite
	mockRepo *MockRepository
	usecase  bundle.Usecase
	ctx      context.Context
}

func (suite *BundleUsecaseTestSuite) SetupTest() {
	suite.mockRepo = new(MockRepository)
	suite.usecase = NewBundleUsecase(suite.mockRepo)
	suite.ctx = context.Background()
}

// ---------------- Test Cases ----------------

func (suite *BundleUsecaseTestSuite) TestCreateBundle() {
	tests := []struct {
		name        string
		supplierID  string
		bundle      *bundle.Bundle
		setupMock   func()
		expectError bool
	}{
		{
			name:       "Successful bundle creation",
			supplierID: "supplier-1",
			bundle:     createTestBundle("supplier-1"),
			setupMock: func() {
				suite.mockRepo.On("CreateBundle", suite.ctx, mock.AnythingOfType("*bundle.Bundle")).Return(nil)
			},
			expectError: false,
		},
		{
			name:        "Supplier ID mismatch",
			supplierID:  "supplier-1",
			bundle:      createTestBundle("supplier-2"),
			setupMock:   func() {},
			expectError: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockRepo.ExpectedCalls = nil // Reset mock expectations
			tt.setupMock()
			err := suite.usecase.CreateBundle(suite.ctx, tt.supplierID, tt.bundle)
			if tt.expectError {
				assert.Error(suite.T(), err)
			} else {
				assert.NoError(suite.T(), err)
				suite.mockRepo.AssertExpectations(suite.T())
			}
		})
	}
}

func (suite *BundleUsecaseTestSuite) TestListBundles() {
	tests := []struct {
		name        string
		supplierID  string
		setupMock   func()
		expectError bool
	}{
		{
			name:       "Successful bundle listing",
			supplierID: "supplier-1",
			setupMock: func() {
				suite.mockRepo.On("ListBundles", suite.ctx, "supplier-1").Return([]*bundle.Bundle{createTestBundle("supplier-1")}, nil)
			},
			expectError: false,
		},
		{
			name:       "Repository error",
			supplierID: "supplier-1",
			setupMock: func() {
				suite.mockRepo.On("ListBundles", suite.ctx, "supplier-1").Return(([]*bundle.Bundle)(nil), errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockRepo.ExpectedCalls = nil // Reset mock expectations
			tt.setupMock()
			bundles, err := suite.usecase.ListBundles(suite.ctx, tt.supplierID)
			if tt.expectError {
				assert.Error(suite.T(), err)
				assert.Nil(suite.T(), bundles)
			} else {
				assert.NoError(suite.T(), err)
				assert.NotNil(suite.T(), bundles)
			}
			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *BundleUsecaseTestSuite) TestDeleteBundle() {
	tests := []struct {
		name        string
		supplierID  string
		bundleID    string
		setupMock   func()
		expectError bool
	}{
		{
			name:       "Successful bundle deletion",
			supplierID: "supplier-1",
			bundleID:   "test-bundle-id",
			setupMock: func() {
				bundle := createTestBundle("supplier-1")
				suite.mockRepo.On("GetBundleByID", suite.ctx, "test-bundle-id").Return(bundle, nil)
				suite.mockRepo.On("DeleteBundle", suite.ctx, "test-bundle-id").Return(nil)
			},
			expectError: false,
		},
		{
			name:       "Bundle not found",
			supplierID: "supplier-1",
			bundleID:   "non-existent",
			setupMock: func() {
				suite.mockRepo.On("GetBundleByID", suite.ctx, "non-existent").Return(nil, mongo.ErrNoDocuments)
			},
			expectError: true,
		},
		{
			name:       "Unauthorized deletion",
			supplierID: "supplier-1",
			bundleID:   "test-bundle-id",
			setupMock: func() {
				bundle := &bundle.Bundle{
					ID:         "test-bundle-id",
					SupplierID: "supplier-2", // Different supplier
				}
				suite.mockRepo.On("GetBundleByID", suite.ctx, "test-bundle-id").Return(bundle, nil)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockRepo.ExpectedCalls = nil // Reset mock expectations
			tt.setupMock()
			err := suite.usecase.DeleteBundle(suite.ctx, tt.supplierID, tt.bundleID)
			if tt.expectError {
				assert.Error(suite.T(), err)
			} else {
				assert.NoError(suite.T(), err)
			}
			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *BundleUsecaseTestSuite) TestGetBundleByID() {
	tests := []struct {
		name        string
		supplierID  string
		bundleID    string
		setupMock   func()
		expectError bool
	}{
		{
			name:       "Successful bundle retrieval",
			supplierID: "supplier-1",
			bundleID:   "test-bundle-id",
			setupMock: func() {
				b := createTestBundle("supplier-1")
				suite.mockRepo.On("GetBundleByID", suite.ctx, "test-bundle-id").Return(b, nil)
			},
			expectError: false,
		},
		{
			name:       "Bundle not found",
			supplierID: "supplier-1",
			bundleID:   "non-existent",
			setupMock: func() {
				suite.mockRepo.On("GetBundleByID", suite.ctx, "non-existent").Return(nil, mongo.ErrNoDocuments)
			},
			expectError: true,
		},
		{
			name:       "Unauthorized access",
			supplierID: "supplier-1",
			bundleID:   "test-bundle-id",
			setupMock: func() {
				b := createTestBundle("supplier-2") // Different supplier
				suite.mockRepo.On("GetBundleByID", suite.ctx, "test-bundle-id").Return(b, nil)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockRepo.ExpectedCalls = nil // Reset mock expectations
			tt.setupMock()
			bundle, err := suite.usecase.GetBundleByID(suite.ctx, tt.supplierID, tt.bundleID)
			if tt.expectError {
				assert.Error(suite.T(), err)
				assert.Nil(suite.T(), bundle)
			} else {
				assert.NoError(suite.T(), err)
				assert.NotNil(suite.T(), bundle)
				assert.Equal(suite.T(), tt.supplierID, bundle.SupplierID)
			}
			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *BundleUsecaseTestSuite) TestUpdateBundle() {
	tests := []struct {
		name        string
		supplierID  string
		bundleID    string
		updateData  map[string]interface{}
		setupMock   func()
		expectError bool
	}{
		{
			name:       "Successful bundle update",
			supplierID: "supplier-1",
			bundleID:   "test-bundle-id",
			updateData: map[string]interface{}{
				"title": "Updated Title",
			},
			setupMock: func() {
				b := createTestBundle("supplier-1")
				b.Status = "available"
				suite.mockRepo.On("GetBundleByID", suite.ctx, "test-bundle-id").Return(b, nil)
				suite.mockRepo.On("UpdateBundle", suite.ctx, "test-bundle-id", mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name:       "Bundle not found",
			supplierID: "supplier-1",
			bundleID:   "non-existent",
			updateData: map[string]interface{}{
				"title": "Updated Title",
			},
			setupMock: func() {
				suite.mockRepo.On("GetBundleByID", suite.ctx, "non-existent").Return(nil, mongo.ErrNoDocuments)
			},
			expectError: true,
		},
		{
			name:       "Unauthorized update",
			supplierID: "supplier-1",
			bundleID:   "test-bundle-id",
			updateData: map[string]interface{}{
				"title": "Updated Title",
			},
			setupMock: func() {
				b := createTestBundle("supplier-2") // Different supplier
				b.Status = "available"
				suite.mockRepo.On("GetBundleByID", suite.ctx, "test-bundle-id").Return(b, nil)
			},
			expectError: true,
		},
		{
			name:       "Bundle not in available status",
			supplierID: "supplier-1",
			bundleID:   "test-bundle-id",
			updateData: map[string]interface{}{
				"title": "Updated Title",
			},
			setupMock: func() {
				b := createTestBundle("supplier-1")
				b.Status = "sold" // Not available
				suite.mockRepo.On("GetBundleByID", suite.ctx, "test-bundle-id").Return(b, nil)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockRepo.ExpectedCalls = nil // Reset mock expectations
			tt.setupMock()
			err := suite.usecase.UpdateBundle(suite.ctx, tt.supplierID, tt.bundleID, tt.updateData)
			if tt.expectError {
				assert.Error(suite.T(), err)
			} else {
				assert.NoError(suite.T(), err)
			}
			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *BundleUsecaseTestSuite) TestListAvailableBundles() {
	tests := []struct {
		name        string
		setupMock   func()
		expectError bool
	}{
		{
			name: "Successful available bundles listing",
			setupMock: func() {
				suite.mockRepo.On("ListAvailableBundles", suite.ctx).Return([]*bundle.Bundle{createTestBundle("supplier-1")}, nil)
			},
			expectError: false,
		},
		{
			name: "Repository error",
			setupMock: func() {
				suite.mockRepo.On("ListAvailableBundles", suite.ctx).Return(([]*bundle.Bundle)(nil), errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockRepo.ExpectedCalls = nil // Reset mock expectations
			tt.setupMock()
			bundles, err := suite.usecase.ListAvailableBundles(suite.ctx)
			if tt.expectError {
				assert.Error(suite.T(), err)
				assert.Nil(suite.T(), bundles)
			} else {
				assert.NoError(suite.T(), err)
				assert.NotNil(suite.T(), bundles)
			}
			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *BundleUsecaseTestSuite) TestDecreaseRemainingItemCount() {
	tests := []struct {
		name        string
		bundleID    string
		setupMock   func()
		expectError bool
	}{
		{
			name:     "Successful decrease",
			bundleID: "test-bundle-id",
			setupMock: func() {
				b := createTestBundle("supplier-1")
				b.RemainingItemCount = 5
				suite.mockRepo.On("GetBundleByID", suite.ctx, "test-bundle-id").Return(b, nil)
				suite.mockRepo.On("UpdateBundle", suite.ctx, "test-bundle-id", mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name:     "Bundle not found",
			bundleID: "non-existent",
			setupMock: func() {
				suite.mockRepo.On("GetBundleByID", suite.ctx, "non-existent").Return(nil, mongo.ErrNoDocuments)
			},
			expectError: true,
		},
		{
			name:     "Bundle fully unpacked",
			bundleID: "test-bundle-id",
			setupMock: func() {
				b := createTestBundle("supplier-1")
				b.RemainingItemCount = 0
				suite.mockRepo.On("GetBundleByID", suite.ctx, "test-bundle-id").Return(b, nil)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockRepo.ExpectedCalls = nil // Reset mock expectations
			tt.setupMock()
			err := suite.usecase.DecreaseRemainingItemCount(suite.ctx, tt.bundleID)
			if tt.expectError {
				assert.Error(suite.T(), err)
			} else {
				assert.NoError(suite.T(), err)
			}
			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *BundleUsecaseTestSuite) TestGetBundlePublicByID() {
	tests := []struct {
		name        string
		bundleID    string
		setupMock   func()
		expectError bool
	}{
		{
			name:     "Successful public bundle retrieval",
			bundleID: "test-bundle-id",
			setupMock: func() {
				suite.mockRepo.On("GetBundleByID", suite.ctx, "test-bundle-id").Return(createTestBundle("supplier-1"), nil)
			},
			expectError: false,
		},
		{
			name:     "Bundle not found",
			bundleID: "non-existent",
			setupMock: func() {
				suite.mockRepo.On("GetBundleByID", suite.ctx, "non-existent").Return(nil, mongo.ErrNoDocuments)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.mockRepo.ExpectedCalls = nil // Reset mock expectations
			tt.setupMock()
			bundle, err := suite.usecase.GetBundlePublicByID(suite.ctx, tt.bundleID)
			if tt.expectError {
				assert.Error(suite.T(), err)
				assert.Nil(suite.T(), bundle)
			} else {
				assert.NoError(suite.T(), err)
				assert.NotNil(suite.T(), bundle)
			}
			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

// Run the test suite
func TestBundleUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(BundleUsecaseTestSuite))
}
