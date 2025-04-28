package routes

import (
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/auth"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/controllers"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterWarehouseRoutes(r *gin.Engine, warehouse_ctrl *controllers.WarehouseController, jwtSvc auth.JWTService) {
	warehouseGroup := r.Group("/warehouse")
	warehouseGroup.Use(middlewares.AuthMiddleware(jwtSvc))

	warehouseGroup.GET("", middlewares.AuthorizeRoles("reseller"), warehouse_ctrl.GetWarehouseItems)
}
