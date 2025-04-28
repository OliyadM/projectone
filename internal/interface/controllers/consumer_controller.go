package controllers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/order"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/models/common"
	"github.com/gin-gonic/gin"
)

type ConsumerController struct {
	orderRepo order.Repository
}

func NewConsumerController(orderRepo order.Repository) *ConsumerController {
	return &ConsumerController{orderRepo: orderRepo}
}

func (c *ConsumerController) GetOrderHistory(ctx *gin.Context) {
	// Get user ID from context
	userID, exists := ctx.Get("userID")
	if !exists {
		fmt.Printf("❌ User ID not found in context\n")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get query parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	status := ctx.Query("status")

	fmt.Printf("🔍 Fetching order history - UserID: %s, Page: %d, Limit: %d, Status: %s\n", 
		userID.(string), page, limit, status)

	// Get orders from use case
	orders, err := c.orderRepo.GetOrdersByConsumer(ctx, userID.(string))
	if err != nil {
		fmt.Printf("❌ Error fetching orders: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Message: "Failed to fetch orders",
		})
		return
	}

	fmt.Printf("✅ Found %d orders\n", len(orders))

	// Simulate delivery status
	for i := range orders {
		createdAt, _ := time.Parse(time.RFC3339, orders[i].CreatedAt)
		if time.Since(createdAt) > 10*time.Minute && orders[i].Status == order.Pending {
			orders[i].Status = order.Failed
		} else if time.Since(createdAt) > 3*time.Minute && orders[i].Status == order.Pending {
			orders[i].Status = order.Delivered
		}
	}

	// Filter by status if provided
	if status != "" {
		filteredOrders := []*order.Order{}
		for _, o := range orders {
			if strings.EqualFold(string(o.Status), status) {
				filteredOrders = append(filteredOrders, o)
			}
		}
		orders = filteredOrders
	}

	// Sort orders by CreatedAt in descending order
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].CreatedAt > orders[j].CreatedAt
	})

	// Paginate results
	start := (page - 1) * limit
	end := start + limit
	if start > len(orders) {
		start = len(orders)
	}
	if end > len(orders) {
		end = len(orders)
	}
	paginatedOrders := orders[start:end]

	// Map to response format
	var response []map[string]interface{}
	for _, o := range paginatedOrders {
		itemTitle := ""
		if len(o.ProductIDs) > 0 {
			itemTitle = o.ProductIDs[0]
		}

		response = append(response, map[string]interface{}{
			"orderId":               o.ID,
			"itemTitle":             itemTitle,
			"price":                 o.TotalPrice,
			"imageUrl":              "https://example.com/image.jpg",
			"status":                o.Status,
			"purchaseDate":          o.CreatedAt,
			"estimatedDeliveryTime": "3 minutes",
		})
	}

	if len(response) == 0 {
		ctx.JSON(http.StatusOK, common.APIResponse{
			Success: true,
			Message: "No orders yet.",
			Data:    []map[string]interface{}{},
		})
		return
	}

	ctx.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Message: "Orders retrieved successfully",
		Data:    response,
	})
}
