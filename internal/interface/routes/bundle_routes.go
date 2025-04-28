package routes

import (
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/auth"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/controllers"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterBundleRoutes(r *gin.Engine, ctrl *controllers.BundleController, jwtSvc auth.JWTService) {
	bundleGroup := r.Group("/bundles")
	bundleGroup.Use(middlewares.AuthMiddleware(jwtSvc)) // All routes require valid token

	bundleGroup.POST("", middlewares.AuthorizeRoles("supplier"), ctrl.CreateBundle)
	bundleGroup.GET("", middlewares.AuthorizeRoles("supplier"), ctrl.ListBundles)
	bundleGroup.GET("/:id", middlewares.AuthorizeRoles("supplier", "reseller"), ctrl.GetBundle)
	bundleGroup.DELETE("/:id", middlewares.AuthorizeRoles("supplier"), ctrl.DeleteBundle)
	bundleGroup.PUT("/:id", middlewares.AuthorizeRoles("supplier"), ctrl.UpdateBundle)
	bundleGroup.GET("/available", middlewares.AuthorizeRoles("reseller", "supplier"), ctrl.ListAvailableBundles)
	bundleGroup.GET("/detail/:id", middlewares.AuthorizeRoles("reseller", "supplier"), ctrl.GetBundleDetail)
	bundleGroup.GET("/title/:title", middlewares.AuthorizeRoles("reseller", "supplier"), ctrl.GetBundleByTitle)
}
