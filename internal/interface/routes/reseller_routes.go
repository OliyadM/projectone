package routes

import (
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/auth"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/controllers"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterResellerRoutes(r *gin.Engine, ctrl *controllers.SupplierController, jwtSvc auth.JWTService) {
	resellerGroup := r.Group("/reseller")
	resellerGroup.Use(middlewares.AuthMiddleware(jwtSvc))

	resellerGroup.GET("/metrics", middlewares.AuthorizeRoles("reseller"), ctrl.GetResellerMetrics)
} 