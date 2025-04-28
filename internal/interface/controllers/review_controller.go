package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/product"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/review"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/trust"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/models"
	"github.com/gin-gonic/gin"
)

type ReviewController struct {
	usecase        review.Usecase
	trustUsecase   trust.Usecase
	productUsecase product.Usecase
}

func NewReviewController(usecase review.Usecase, trustUsecase trust.Usecase, productUsecase product.Usecase) *ReviewController {
	return &ReviewController{
		usecase:        usecase,
		trustUsecase:   trustUsecase,
		productUsecase: productUsecase,
	}
}

func (ctrl *ReviewController) SubmitReview(c *gin.Context) {
	fmt.Printf("🔍 Starting review submission\n")
	
	var req models.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("❌ Invalid request body: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}
	fmt.Printf("📝 Review request: %+v\n", req)

	userID := c.GetString("userID")
	if userID == "" {
		fmt.Printf("❌ No user ID found in context\n")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	fmt.Printf("👤 User ID from context: %s\n", userID)

	// Get the product to get the reseller ID and current rating
	fmt.Printf("🔍 Fetching product details for ID: %s\n", req.ProductID)
	product, err := ctrl.productUsecase.GetProductByID(c.Request.Context(), req.ProductID)
	if err != nil {
		fmt.Printf("❌ Error fetching product: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "product not found: " + err.Error()})
		return
	}
	fmt.Printf("✅ Found product: %+v\n", product)

	r := review.NewReview(
		req.OrderID,
		req.ProductID,
		userID,
		product.ResellerID.Hex(),
		req.Rating,
		req.Comment,
	)
	fmt.Printf("📝 Creating review: %+v\n", r)

	if err := ctrl.usecase.SubmitReview(c.Request.Context(), r); err != nil {
		fmt.Printf("❌ Error submitting review: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update reseller trust score
	if ctrl.trustUsecase != nil {
		fmt.Printf("📊 Updating trust score for reseller: %s\n", product.ResellerID.Hex())
		go ctrl.trustUsecase.UpdateResellerTrustScoreOnNewRating(
			context.Background(),
			product.ResellerID.Hex(),
			product.Rating,
			float64(req.Rating),
		)
	}

	fmt.Printf("✅ Review submitted successfully\n")
	c.JSON(http.StatusCreated, gin.H{"message": "review submitted"})
}

func (ctrl *ReviewController) GetResellerReviews(c *gin.Context) {
	fmt.Printf("🔍 Starting to fetch reseller reviews\n")
	
	resellerID := c.Param("id")
	if resellerID == "" {
		fmt.Printf("❌ No reseller ID provided\n")
		c.JSON(http.StatusBadRequest, gin.H{"error": "reseller ID is required"})
		return
	}

	reviews, err := ctrl.usecase.GetResellerReviews(c.Request.Context(), resellerID)
	if err != nil {
		fmt.Printf("❌ Error fetching reviews: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("✅ Successfully fetched %d reviews\n", len(reviews))
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"reviews": reviews,
		},
	})
}
