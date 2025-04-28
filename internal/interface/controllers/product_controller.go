package controllers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/product"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/trust"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/models"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/warehouse"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductController struct {
	Usecase       product.Usecase
	TrustUsecase  trust.Usecase
	BundleUsecase bundle.Usecase
	WarehouseRepo warehouse.Repository
}

func NewProductController(
	prodUC product.Usecase,
	trustUC trust.Usecase,
	bundleUC bundle.Usecase,
	warehouseRepo warehouse.Repository,
) *ProductController {
	return &ProductController{
		Usecase:       prodUC,
		TrustUsecase:  trustUC,
		BundleUsecase: bundleUC,
		WarehouseRepo: warehouseRepo,
	}
}

func (h *ProductController) Create(c *gin.Context) {
	var p product.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	resellerID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
		return
	}

	owns, err := h.WarehouseRepo.HasResellerReceivedBundle(c.Request.Context(), userIDStr, p.BundleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "warehouse check failed"})
		return
	}
	if !owns {
		c.JSON(http.StatusForbidden, gin.H{"error": "you have not received this bundle in your warehouse yet"})
		return
	}

	b, err := h.BundleUsecase.GetBundlePublicByID(c.Request.Context(), p.BundleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bundle ID"})
		return
	}

	if b.RemainingItemCount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bundle is fully unpacked"})
		return
	}

	p.ResellerID = resellerID
	p.SupplierID = b.SupplierID
	p.ID = p.GenerateID()
	p.Status = "available"

	if err := h.Usecase.AddProduct(c.Request.Context(), &p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.BundleUsecase.DecreaseRemainingItemCount(c.Request.Context(), p.BundleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decrease bundle quantity"})
		return
	}

	if h.TrustUsecase != nil && p.SupplierID != "" {
		go h.TrustUsecase.UpdateSupplierTrustScoreOnNewRating(
			context.Background(),
			p.SupplierID,
			float64(b.DeclaredRating),
			p.Rating,
		)
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Product created successfully",
		"data": models.ProductResponse{
			ID:          p.ID,
			Title:       p.Title,
			Price:       p.Price,
			Photo:       p.ImageURL,
			Grade:       p.Grade,
			Size:        p.Size,
			Status:      p.Status,
			SellerID:    p.ResellerID.Hex(),
			Rating:      p.Rating,
			Description: p.Description,
			Type:        p.Type,
			BundleID:    p.BundleID,
		},
	})

}

func (h *ProductController) GetByID(c *gin.Context) {
	id := c.Param("id")
	prod, err := h.Usecase.GetProductByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, prod)
}

func (h *ProductController) GetByTitle(c *gin.Context) {
	title := c.Param("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title cannot be empty"})
		return
	}

	prod, err := h.Usecase.GetProductByTitle(c.Request.Context(), title)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve product"})
		}
		return
	}

	c.JSON(http.StatusOK, prod)
}

func (h *ProductController) ListAvailable(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	products, err := h.Usecase.ListAvailableProducts(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load products", "details": err.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no products available", "products": []product.Product{}})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductController) ListByReseller(c *gin.Context) {
	resellerID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, err := h.Usecase.ListProductsByReseller(c.Request.Context(), resellerID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load reseller products", "details": err.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no products found for this reseller", "products": []product.Product{}})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductController) Update(c *gin.Context) {
	id := c.Param("id")
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid update payload"})
		return
	}
	if err := h.Usecase.UpdateProduct(c.Request.Context(), id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product updated"})
}

func (h *ProductController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.Usecase.DeleteProduct(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}
