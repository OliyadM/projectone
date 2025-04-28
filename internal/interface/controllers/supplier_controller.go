package controllers

import (
    "net/http"
    "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/order"
    "github.com/Zeamanuel-Admasu/afro-vintage-backend/models/common"
    "github.com/gin-gonic/gin"
    "fmt"
)

type SupplierController struct {
    orderUseCase order.Usecase
}

func NewSupplierController(orderUseCase order.Usecase) *SupplierController {
    return &SupplierController{orderUseCase: orderUseCase}
}

func (c *SupplierController) GetDashboardMetrics(ctx *gin.Context) {
    supplierID, exists := ctx.Get("userID")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, common.APIResponse{
            Success: false,
            Message: "Unauthorized",
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

    metrics, err := c.orderUseCase.GetDashboardMetrics(ctx, supplierIDStr)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, metrics)
}

func (c *SupplierController) GetResellerMetrics(ctx *gin.Context) {
	fmt.Println("🔍 Starting GetResellerMetrics request")
	
	resellerID, exists := ctx.Get("userID")
	if !exists {
		fmt.Println("❌ No user ID found in context")
		ctx.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Message: "Unauthorized: user ID not found in context",
		})
		return
	}

	resellerIDStr, ok := resellerID.(string)
	if !ok || resellerIDStr == "" {
		fmt.Println("❌ Invalid user ID in context")
		ctx.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Message: "Unauthorized: invalid or empty user ID in context",
		})
		return
	}

	fmt.Printf("👤 Processing metrics for reseller: %s\n", resellerIDStr)
	metrics, err := c.orderUseCase.GetResellerMetrics(ctx, resellerIDStr)
	if err != nil {
		fmt.Printf("❌ Error getting metrics: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Message: "Failed to get reseller metrics",
			
		})
		return
	}

	fmt.Println("✅ Successfully retrieved metrics")
	ctx.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Message: "Reseller metrics retrieved successfully",
		Data:    metrics,
	})
}