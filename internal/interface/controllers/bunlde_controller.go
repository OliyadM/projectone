package controllers

import (
	"net/http"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/models"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/models/common"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BundleController struct {
	bundleUsecase bundle.Usecase
	userUsecase   user.Usecase
}

func NewBundleController(bundleUsecase bundle.Usecase, userUsecase user.Usecase) *BundleController {
	return &BundleController{
		bundleUsecase: bundleUsecase,
		userUsecase:   userUsecase,
	}
}

func (c *BundleController) CreateBundle(ctx *gin.Context) {
	supplierID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Message: "user ID not found in context",
		})
		return
	}

	supplierIDStr, ok := supplierID.(string)
	if !ok || supplierIDStr == "" {
		ctx.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Message: "invalid or empty user ID in context",
		})
		return
	}

	role, exists := ctx.Get("role")
	if !exists || role != "supplier" {
		ctx.JSON(http.StatusForbidden, common.APIResponse{
			Success: false,
			Message: "only suppliers can create bundles",
		})
		return
	}
	user, err := c.userUsecase.GetByID(ctx, supplierIDStr)
	if err != nil || user.IsBlacklisted {
		ctx.JSON(http.StatusForbidden, common.APIResponse{
			Success: false,
			Message: "you are blacklisted and cannot create bundles",
		})
		return
	}

	var req models.CreateBundleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Message: "invalid request: " + err.Error(),
		})
		return
	}

	b := &bundle.Bundle{
		ID:                 "bundle_" + primitive.NewObjectID().Hex(),
		SupplierID:         supplierIDStr,
		Title:              req.Title,
		Description:        req.Description,
		SampleImage:        req.SampleImage,
		Quantity:           req.NumberOfItems,
		Grade:              req.Grade,
		SortingLevel:       bundle.SortingLevel(req.Type),
		EstimatedBreakdown: req.EstimatedBreakdown,
		Type:               req.Type,
		Price:              req.Price,
		Status:             "available",
		CreatedAt:          time.Now().Format(time.RFC3339),
		DeclaredRating:     req.DeclaredRating, // ✅ included here
		RemainingItemCount: req.NumberOfItems,
	}

	if err := c.bundleUsecase.CreateBundle(ctx, supplierIDStr, b); err != nil {
		ctx.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	resp := models.BundleResponse{
		ID:                 b.ID,
		Title:              b.Title,
		SampleImage:        b.SampleImage,
		Quantity:           b.Quantity, // from NumberOfItems
		Grade:              b.Grade,
		Description:        b.Description,
		SizeRange:          b.SizeRange,
		Type:               b.Type,
		Price:              b.Price,
		Status:             b.Status,
		EstimatedBreakdown: b.EstimatedBreakdown,
		DeclaredRating:     b.DeclaredRating,
		SortingLevel:       string(b.SortingLevel),
		CreatedAt:          b.CreatedAt,
	}

	ctx.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Message: "Bundle successfully created and listed!",
		Data:    resp,
	})
}

func (c *BundleController) ListBundles(ctx *gin.Context) {
	supplierID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Message: "user ID not found in context",
		})
		return
	}

	supplierIDStr, ok := supplierID.(string)
	if !ok || supplierIDStr == "" {
		ctx.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Message: "invalid or empty user ID in context",
		})
		return
	}

	role, exists := ctx.Get("role")
	if !exists || role != "supplier" {
		ctx.JSON(http.StatusForbidden, common.APIResponse{
			Success: false,
			Message: "only suppliers can list bundles",
		})
		return
	}

	bundles, err := c.bundleUsecase.ListBundles(ctx, supplierIDStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	var resp []models.BundleResponse
	for _, b := range bundles {
		resp = append(resp, models.BundleResponse{
			ID:                 b.ID,
			Title:              b.Title,
			SampleImage:        b.SampleImage,
			Quantity:           b.Quantity,
			Grade:              b.Grade,
			Description:        b.Description,
			SizeRange:          b.SizeRange,
			Type:               b.Type,
			Price:              b.Price,
			Status:             b.Status,
			EstimatedBreakdown: b.EstimatedBreakdown,
			DeclaredRating:     b.DeclaredRating,
			SortingLevel:       string(b.SortingLevel),
			CreatedAt:          b.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Message: "Bundles retrieved successfully",
		Data:    resp,
	})
}

func (c *BundleController) DeleteBundle(ctx *gin.Context) {
	supplierID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Message: "user ID not found in context",
		})
		return
	}

	supplierIDStr, ok := supplierID.(string)
	if !ok || supplierIDStr == "" {
		ctx.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Message: "invalid or empty user ID in context",
		})
		return
	}

	role, exists := ctx.Get("role")
	if !exists || role != "supplier" {
		ctx.JSON(http.StatusForbidden, common.APIResponse{
			Success: false,
			Message: "only suppliers can delete bundles",
		})
		return
	}

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Message: "bundle ID is required",
		})
		return
	}

	err := c.bundleUsecase.DeleteBundle(ctx, supplierIDStr, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Message: "Bundle successfully deactivated",
		Data:    nil,
	})
}

func (c *BundleController) UpdateBundle(ctx *gin.Context) {
	// Extract Supplier ID from JWT
	supplierID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Message: "user ID not found in context",
		})
		return
	}

	supplierIDStr, ok := supplierID.(string)
	if !ok || supplierIDStr == "" {
		ctx.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Message: "invalid or empty user ID in context",
		})
		return
	}

	// Validate role
	role, exists := ctx.Get("role")
	if !exists || role != "supplier" {
		ctx.JSON(http.StatusForbidden, common.APIResponse{
			Success: false,
			Message: "only suppliers can update bundles",
		})
		return
	}

	// Extract bundle ID from URL
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Message: "bundle ID is required",
		})
		return
	}

	// Parse request body
	var updatedData map[string]interface{}
	if err := ctx.ShouldBindJSON(&updatedData); err != nil {
		ctx.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Message: "invalid request: " + err.Error(),
		})
		return
	}

	// Call the use case to update the bundle
	err := c.bundleUsecase.UpdateBundle(ctx, supplierIDStr, id, updatedData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Fetch the updated bundle to return in the response
	updatedBundle, err := c.bundleUsecase.GetBundleByID(ctx, supplierIDStr, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Map to response DTO
	resp := models.BundleResponse{
		ID:     updatedBundle.ID,
		Title:  updatedBundle.Title,
		Grade:  updatedBundle.Grade,
		Price:  updatedBundle.Price,
		Type:   string(updatedBundle.SortingLevel),
		Status: updatedBundle.Status,
	}

	ctx.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Message: "Bundle successfully updated",
		Data:    resp,
	})
}

func (c *BundleController) GetBundle(ctx *gin.Context) { // Added
	// Extract Supplier ID from JWT
	supplierID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Message: "user ID not found in context",
		})
		return
	}

	supplierIDStr, ok := supplierID.(string)
	if !ok || supplierIDStr == "" {
		ctx.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Message: "invalid or empty user ID in context",
		})
		return
	}

	// Validate role
	role, exists := ctx.Get("role")
	if !exists || role != "supplier" {
		ctx.JSON(http.StatusForbidden, common.APIResponse{
			Success: false,
			Message: "only suppliers can view bundles",
		})
		return
	}

	// Extract bundle ID from URL
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Message: "bundle ID is required",
		})
		return
	}

	// Fetch the bundle using the use case
	b, err := c.bundleUsecase.GetBundleByID(ctx, supplierIDStr, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Map to response DTO
	resp := models.BundleResponse{
		ID:     b.ID,
		Title:  b.Title,
		Grade:  b.Grade,
		Price:  b.Price,
		Type:   string(b.SortingLevel),
		Status: b.Status,
	}

	ctx.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Message: "Bundle retrieved successfully",
		Data:    resp,
	})
}

func (c *BundleController) ListAvailableBundles(ctx *gin.Context) {
	bundles, err := c.bundleUsecase.ListAvailableBundles(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, bundles)
}

func (c *BundleController) GetBundleDetail(ctx *gin.Context) {
	// Check authentication
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Message: "user ID not found in context",
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		ctx.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Message: "invalid or empty user ID in context",
		})
		return
	}

	// Validate role
	role, exists := ctx.Get("role")
	if !exists || (role != "supplier" && role != "reseller") {
		ctx.JSON(http.StatusForbidden, common.APIResponse{
			Success: false,
			Message: "only suppliers and resellers can view bundle details",
		})
		return
	}

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Message: "bundle ID is required",
		})
		return
	}

	// Get bundle details
	bundle, err := c.bundleUsecase.GetBundlePublicByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, common.APIResponse{
			Success: false,
			Message: "bundle not found",
		})
		return
	}

	// Get supplier details
	supplier, err := c.userUsecase.GetByID(ctx, bundle.SupplierID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Message: "error fetching supplier details",
		})
		return
	}

	// Calculate supplier rating from trust score
	supplierRating := float64(supplier.TrustScore) / 100.0

	// Construct response
	response := models.BundleDetailResponse{}

	// Fill bundle details
	response.Bundle.ID = bundle.ID
	response.Bundle.Title = bundle.Title
	response.Bundle.Description = bundle.Description
	response.Bundle.SampleImage = bundle.SampleImage
	response.Bundle.Quantity = bundle.Quantity
	response.Bundle.Grade = bundle.Grade
	response.Bundle.SortingLevel = string(bundle.SortingLevel)
	response.Bundle.EstimatedBreakdown = bundle.EstimatedBreakdown
	response.Bundle.Type = bundle.Type
	response.Bundle.Price = bundle.Price
	response.Bundle.Status = bundle.Status
	response.Bundle.DeclaredRating = bundle.DeclaredRating
	response.Bundle.RemainingItemCount = bundle.RemainingItemCount

	// Fill supplier details
	response.Supplier.ID = supplier.ID
	response.Supplier.Name = supplier.Name
	response.Supplier.Rating = supplierRating

	ctx.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Message: "bundle detail fetched",
		Data:    response,
	})
}

func (c *BundleController) GetBundleByTitle(ctx *gin.Context) {
	title := ctx.Param("title")
	if title == "" {
		ctx.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Message: "title parameter is required",
		})
		return
	}

	bundle, err := c.bundleUsecase.GetBundleByTitle(ctx, title)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	resp := models.BundleResponse{
		ID:                 bundle.ID,
		Title:              bundle.Title,
		SampleImage:        bundle.SampleImage,
		Quantity:           bundle.Quantity,
		Grade:              bundle.Grade,
		Description:        bundle.Description,
		SizeRange:          bundle.SizeRange,
		Type:               bundle.Type,
		Price:              bundle.Price,
		Status:             bundle.Status,
		EstimatedBreakdown: bundle.EstimatedBreakdown,
		DeclaredRating:     bundle.DeclaredRating,
		SortingLevel:       string(bundle.SortingLevel),
		CreatedAt:          bundle.CreatedAt,
	}

	ctx.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Message: "Bundle retrieved successfully",
		Data:    resp,
	})
}
