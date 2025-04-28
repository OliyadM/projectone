package routes

import (
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/auth"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/controllers"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(r *gin.Engine, order_ctrl *controllers.OrderController, consumer_ctrl *controllers.ConsumerController, jwtSvc auth.JWTService) {
	consumerGroup := r.Group("/orders")
	consumerGroup.Use(middlewares.AuthMiddleware(jwtSvc))

	consumerGroup.POST("", middlewares.AuthorizeRoles("reseller"), order_ctrl.PurchaseBundle)
	consumerGroup.POST("/:id", middlewares.AuthorizeRoles("reseller", "consumer"), order_ctrl.GetOrderByID)
	consumerGroup.GET("/history", middlewares.AuthorizeRoles("reseller", "consumer"), order_ctrl.GetOrderHistory)
	consumerGroup.GET("/supplier/history", middlewares.AuthorizeRoles("supplier"), order_ctrl.GetSoldBundleHistory)
	consumerGroup.GET("/reseller/history", middlewares.AuthorizeRoles("reseller"), order_ctrl.GetOrdersByReseller)
}
