package trustusecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"
)

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *mockUserRepo) UpdateTrustData(ctx context.Context, user *user.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepo) CountActiveUsers(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Get(0).(int), args.Error(1)
}

func (m *mockUserRepo) CreateUser(ctx context.Context, user *user.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepo) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockUserRepo) FindUserByUsername(ctx context.Context, username string) (*user.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *mockUserRepo) GetBlacklistedUsers(ctx context.Context) ([]*user.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.User), args.Error(1)
}

func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *mockUserRepo) ListUsersByRole(ctx context.Context, role user.Role) ([]*user.User, error) {
	args := m.Called(ctx, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.User), args.Error(1)
}

func (m *mockUserRepo) UpdateUser(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func TestTrustUsecase_UpdateSupplierTrustScoreOnNewRating(t *testing.T) {
	tests := []struct {
		name           string
		supplierID     string
		declaredRating float64
		productRating  float64
		initialScore   int
		initialError   float64
		initialCount   int
		expectedScore  int
		expectedError  float64
		expectedCount  int
	}{
		{
			name:           "perfect match",
			supplierID:     primitive.NewObjectID().Hex(),
			declaredRating: 4.5,
			productRating:  4.5,
			initialScore:   100,
			initialError:   0,
			initialCount:   0,
			expectedScore:  100,
			expectedError:  0,
			expectedCount:  1,
		},
		{
			name:           "small difference",
			supplierID:     primitive.NewObjectID().Hex(),
			declaredRating: 4.0,
			productRating:  4.5,
			initialScore:   100,
			initialError:   0,
			initialCount:   0,
			expectedScore:  95,
			expectedError:  0.5,
			expectedCount:  1,
		},
		{
			name:           "large difference",
			supplierID:     primitive.NewObjectID().Hex(),
			declaredRating: 2.0,
			productRating:  4.5,
			initialScore:   100,
			initialError:   0,
			initialCount:   0,
			expectedScore:  75,
			expectedError:  2.5,
			expectedCount:  1,
		},
		{
			name:           "multiple ratings",
			supplierID:     primitive.NewObjectID().Hex(),
			declaredRating: 4.0,
			productRating:  4.5,
			initialScore:   90,
			initialError:   0.5,
			initialCount:   1,
			expectedScore:  96,
			expectedError:  1.0,
			expectedCount:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockRepo := new(mockUserRepo)
			uc := NewTrustUsecase(nil, nil, mockRepo)

			// Mock user
			supplier := &user.User{
				ID:              tt.supplierID,
				TrustScore:      tt.initialScore,
				TrustTotalError: tt.initialError,
				TrustRatedCount: tt.initialCount,
			}

			mockRepo.On("GetByID", mock.Anything, tt.supplierID).Return(supplier, nil)
			mockRepo.On("UpdateTrustData", mock.Anything, mock.Anything).Return(nil)

			// Execute
			err := uc.UpdateSupplierTrustScoreOnNewRating(
				context.Background(),
				tt.supplierID,
				tt.declaredRating,
				tt.productRating,
			)

			// Assert
			assert.NoError(t, err)
			mockRepo.AssertExpectations(t)

			// Verify the final state
			assert.Equal(t, tt.expectedScore, supplier.TrustScore)
			assert.Equal(t, tt.expectedError, supplier.TrustTotalError)
			assert.Equal(t, tt.expectedCount, supplier.TrustRatedCount)
		})
	}
}

func TestTrustUsecase_UpdateResellerTrustScoreOnNewRating(t *testing.T) {
	tests := []struct {
		name           string
		resellerID     string
		declaredRating float64
		productRating  float64
		initialScore   int
		initialError   float64
		initialCount   int
		expectedScore  int
		expectedError  float64
		expectedCount  int
	}{
		{
			name:           "perfect match",
			resellerID:     primitive.NewObjectID().Hex(),
			declaredRating: 4.5,
			productRating:  4.5,
			initialScore:   100,
			initialError:   0,
			initialCount:   0,
			expectedScore:  100,
			expectedError:  0,
			expectedCount:  1,
		},
		{
			name:           "small difference",
			resellerID:     primitive.NewObjectID().Hex(),
			declaredRating: 4.0,
			productRating:  4.5,
			initialScore:   100,
			initialError:   0,
			initialCount:   0,
			expectedScore:  95,
			expectedError:  0.5,
			expectedCount:  1,
		},
		{
			name:           "large difference",
			resellerID:     primitive.NewObjectID().Hex(),
			declaredRating: 2.0,
			productRating:  4.5,
			initialScore:   100,
			initialError:   0,
			initialCount:   0,
			expectedScore:  75,
			expectedError:  2.5,
			expectedCount:  1,
		},
		{
			name:           "multiple ratings",
			resellerID:     primitive.NewObjectID().Hex(),
			declaredRating: 4.0,
			productRating:  4.5,
			initialScore:   90,
			initialError:   0.5,
			initialCount:   1,
			expectedScore:  96,
			expectedError:  1.0,
			expectedCount:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockRepo := new(mockUserRepo)
			uc := NewTrustUsecase(nil, nil, mockRepo)

			// Mock user
			reseller := &user.User{
				ID:              tt.resellerID,
				TrustScore:      tt.initialScore,
				TrustTotalError: tt.initialError,
				TrustRatedCount: tt.initialCount,
			}

			mockRepo.On("GetByID", mock.Anything, tt.resellerID).Return(reseller, nil)
			mockRepo.On("UpdateTrustData", mock.Anything, mock.Anything).Return(nil)

			// Execute
			err := uc.UpdateResellerTrustScoreOnNewRating(
				context.Background(),
				tt.resellerID,
				tt.declaredRating,
				tt.productRating,
			)

			// Assert
			assert.NoError(t, err)
			mockRepo.AssertExpectations(t)

			// Verify the final state
			assert.Equal(t, tt.expectedScore, reseller.TrustScore)
			assert.Equal(t, tt.expectedError, reseller.TrustTotalError)
			assert.Equal(t, tt.expectedCount, reseller.TrustRatedCount)
		})
	}
}