package controllers

import (
	"net/http"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/warehouse"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/models/common"
	"github.com/gin-gonic/gin"
)

type WarehouseController struct {
	warehouseUsecase warehouse.WarehouseUseCase
}

func NewWarehouseController(warehouseUsecase warehouse.WarehouseUseCase) *WarehouseController {
	return &WarehouseController{warehouseUsecase: warehouseUsecase}
}

func (c *WarehouseController) GetWarehouseItems(ctx *gin.Context) {

	resellerID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}
	resellerIDStr, ok := resellerID.(string)
	if !ok || resellerIDStr == "" {
		ctx.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Message: "invalid or empty user ID in context",
		})
		return
	}

	items, err := c.warehouseUsecase.GetWarehouseItems(ctx, resellerIDStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, items)
}
