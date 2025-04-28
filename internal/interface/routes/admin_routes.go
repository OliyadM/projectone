package routes

import (
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/auth"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/controllers"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(r *gin.Engine, ctrl *controllers.AdminController, jwtSvc auth.JWTService) {
	adminGroup := r.Group("/admin")
	adminGroup.Use(middlewares.AuthMiddleware(jwtSvc), middlewares.AuthorizeRoles("admin"))

	// GET /admin/users?role=
	adminGroup.GET("/users", ctrl.GetAllUsers)
	adminGroup.DELETE("/users/:userId", ctrl.DeleteUserIfBlacklisted)
	adminGroup.GET("/users/trust-scores", ctrl.GetTrustScores)
	adminGroup.GET("/blacklisted-users", ctrl.GetBlacklistedUsers)
	adminGroup.GET("/dashboard", ctrl.GetDashboardMetrics)

	// More admin routes can be added here (e.g. transactions, reviews, dashboards, etc.)
	// adminGroup.GET("/dashboard", ctrl.GetDashboardMetrics)
	// adminGroup.GET("/transactions", ctrl.GetAllTransactions)
}
