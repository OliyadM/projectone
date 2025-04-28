package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"errors"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/product"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/review"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockReviewUsecase struct {
	mock.Mock
}

func (m *MockReviewUsecase) SubmitReview(ctx context.Context, r *review.Review) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockReviewUsecase) GetResellerReviews(ctx context.Context, resellerID string) ([]*review.Review, error) {
	args := m.Called(ctx, resellerID)
	return args.Get(0).([]*review.Review), args.Error(1)
}

type MockTrustUsecase struct {
	mock.Mock
}

func (m *MockTrustUsecase) UpdateSupplierTrustScoreOnNewRating(ctx context.Context, supplierID string, declaredRating, productRating float64) error {
	args := m.Called(ctx, supplierID, declaredRating, productRating)
	return args.Error(0)
}

func (m *MockTrustUsecase) UpdateResellerTrustScoreOnNewRating(ctx context.Context, resellerID string, declaredRating, productRating float64) error {
	args := m.Called(ctx, resellerID, declaredRating, productRating)
	return args.Error(0)
}

func (m *MockProductUsecase) UpdateProductRating(ctx context.Context, productID string, rating float64) error {
	args := m.Called(ctx, productID, rating)
	return args.Error(0)
}

type ReviewControllerTestSuite struct {
	suite.Suite
	reviewUsecase  *MockReviewUsecase
	trustUsecase   *MockTrustUsecase
	productUsecase *MockProductUsecase
	controller     *ReviewController
	router         *gin.Engine
}

func (suite *ReviewControllerTestSuite) SetupTest() {
	suite.reviewUsecase = new(MockReviewUsecase)
	suite.trustUsecase = new(MockTrustUsecase)
	suite.productUsecase = new(MockProductUsecase)
	suite.controller = NewReviewController(suite.reviewUsecase, suite.trustUsecase, suite.productUsecase)
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()
}

func TestReviewControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ReviewControllerTestSuite))
}

func (suite *ReviewControllerTestSuite) TestSubmitReview_Success() {
	// Setup
	req := models.CreateReviewRequest{
		OrderID:   "order123",
		ProductID: "product123",
		Rating:    4,
		Comment:   "Great product!",
	}

	product := &product.Product{
		ResellerID: primitive.NewObjectID(),
		Rating:     4.5,
	}

	suite.productUsecase.On("GetProductByID", mock.Anything, req.ProductID).
		Return(product, nil)
	suite.reviewUsecase.On("SubmitReview", mock.Anything, mock.Anything).
		Return(nil)
	suite.trustUsecase.On("UpdateResellerTrustScoreOnNewRating", mock.Anything, product.ResellerID.Hex(), product.Rating, float64(req.Rating)).
		Return(nil).Maybe()

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "user123")

	body, _ := json.Marshal(req)
	request := httptest.NewRequest("POST", "/reviews", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	c.Request = request

	// Execute
	suite.controller.SubmitReview(c)

	// Assert
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	suite.reviewUsecase.AssertExpectations(suite.T())
	suite.productUsecase.AssertExpectations(suite.T())

	// Wait a bit for the goroutine to complete
	time.Sleep(100 * time.Millisecond)
	suite.trustUsecase.AssertExpectations(suite.T())
}

func (suite *ReviewControllerTestSuite) TestSubmitReview_InvalidPayload() {
	// Setup
	invalidPayload := "invalid json"

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "user123")

	request := httptest.NewRequest("POST", "/reviews", bytes.NewBufferString(invalidPayload))
	request.Header.Set("Content-Type", "application/json")
	c.Request = request

	// Execute
	suite.controller.SubmitReview(c)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	suite.reviewUsecase.AssertNotCalled(suite.T(), "SubmitReview")
}

func (suite *ReviewControllerTestSuite) TestSubmitReview_Unauthorized() {
	// Setup
	req := models.CreateReviewRequest{
		OrderID:   "order123",
		ProductID: "product123",
		Rating:    4,
		Comment:   "Great product!",
	}

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Note: No userID set

	body, _ := json.Marshal(req)
	request := httptest.NewRequest("POST", "/reviews", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	c.Request = request

	// Execute
	suite.controller.SubmitReview(c)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	suite.reviewUsecase.AssertNotCalled(suite.T(), "SubmitReview")
}

func (suite *ReviewControllerTestSuite) TestSubmitReview_UseCaseError() {
	// Setup
	req := models.CreateReviewRequest{
		OrderID:   "order123",
		ProductID: "product123",
		Rating:    4,
		Comment:   "Great product!",
	}

	product := &product.Product{
		ResellerID: primitive.NewObjectID(),
		Rating:     4.5,
	}

	expectedError := errors.New("review already exists")
	suite.productUsecase.On("GetProductByID", mock.Anything, req.ProductID).
		Return(product, nil)
	suite.reviewUsecase.On("SubmitReview", mock.Anything, mock.Anything).
		Return(expectedError)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "user123")

	body, _ := json.Marshal(req)
	request := httptest.NewRequest("POST", "/reviews", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	c.Request = request

	// Execute
	suite.controller.SubmitReview(c)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	suite.reviewUsecase.AssertExpectations(suite.T())
	suite.productUsecase.AssertExpectations(suite.T())
}
