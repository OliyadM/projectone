package controllers

import (
	"fmt"
	"net/http"

	OrderUsecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/order"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/models/common"
	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderUseCase OrderUsecase.OrderUseCase
}

func NewOrderController(orderUseCase OrderUsecase.OrderUseCase) *OrderController {
	return &OrderController{orderUseCase: orderUseCase}
}

func (c *OrderController) PurchaseBundle(ctx *gin.Context) {
	type Request struct {
		BundleID string `json:"bundle_id"`
	}

	var req Request
	if err := ctx.ShouldBindJSON(&req); err != nil || req.BundleID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload; bundleId is required"})
		return
	}

	resellerID, _ := ctx.Get("userID")

	resellerIDStr, ok := resellerID.(string)
	if !ok || resellerIDStr == "" {
		ctx.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Message: "invalid or empty user ID in context",
		})
		return
	}

	order, payment, warehouseItem, err := c.orderUseCase.PurchaseBundle(ctx, req.BundleID, resellerIDStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"order":         order,
		"payment":       payment,
		"warehouseItem": warehouseItem,
	})
}

func (c *OrderController) GetOrderByID(ctx *gin.Context) {
	orderID := ctx.Param("id")
	if orderID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "order id is required"})
		return
	}

	order, err := c.orderUseCase.GetOrderByID(ctx, orderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, order)
}

func (c *OrderController) GetSoldBundleHistory(ctx *gin.Context) {
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
			Message: "invalid user ID format",
		})
		return
	}

	orders, userNames, err := c.orderUseCase.GetSoldBundleHistory(ctx, supplierIDStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Message: fmt.Sprintf("failed to get sold bundle history: %v", err),
		})
		return
	}

	var formattedOrders []map[string]interface{}
	for _, order := range orders {
		formattedOrders = append(formattedOrders, map[string]interface{}{
			"order": order,
			"resellerUsername": userNames[order.ResellerID],
		})
	}

	ctx.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Data: gin.H{
			"orders": formattedOrders,
		},
	})
}

func (c *OrderController) GetOrdersByReseller(ctx *gin.Context) {
	resellerID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Message: "user ID not found in context",
		})
		return
	}

	resellerIDStr, ok := resellerID.(string)
	if !ok || resellerIDStr == "" {
		ctx.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Message: "invalid user ID format",
		})
		return
	}

	orders, userNames, err := c.orderUseCase.GetOrdersByReseller(ctx, resellerIDStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Message: fmt.Sprintf("failed to get reseller orders: %v", err),
		})
		return
	}

	// Separate orders into sold and bought
	var soldOrders []map[string]interface{}
	var boughtOrders []map[string]interface{}

	for _, order := range orders {
		if len(order.ProductIDs) > 0 { // Sold order
			soldOrders = append(soldOrders, map[string]interface{}{
				"order":        order,
				"consumerName": userNames[order.ConsumerID],
			})
		} else if order.BundleID != "" { // Bought order
			boughtOrders = append(boughtOrders, map[string]interface{}{
				"order":      order,
				"supplierName": userNames[order.SupplierID],
			})
		}
	}

	ctx.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Data: gin.H{
			"sold":   soldOrders,
			"bought": boughtOrders,
		},
	})
}

func (c *OrderController) GetOrderHistory(ctx *gin.Context) {
	consumerID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Message: "user ID not found in context",
		})
		return
	}

	consumerIDStr, ok := consumerID.(string)
	if !ok || consumerIDStr == "" {
		ctx.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Message: "invalid user ID format",
		})
		return
	}

	orders, userNames, productNames, err := c.orderUseCase.GetOrdersByConsumer(ctx, consumerIDStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Message: fmt.Sprintf("failed to get order history: %v", err),
		})
		return
	}

	var formattedOrders []map[string]interface{}
	for _, order := range orders {
		// Get product details for this order
		var products []map[string]interface{}
		for _, productID := range order.ProductIDs {
			if name, exists := productNames[productID]; exists {
				products = append(products, map[string]interface{}{
					"id":    productID,
					"title": name,
				})
			}
		}

		formattedOrders = append(formattedOrders, map[string]interface{}{
			"order": order,
			"resellerUsername": userNames[order.ResellerID],
			"products": products,
		})
	}

	ctx.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Data: gin.H{
			"orders": formattedOrders,
		},
	})
}