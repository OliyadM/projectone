package routes

import (
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/auth"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/product"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/trust"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/controllers"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(
	r *gin.Engine,
	productCtrl *controllers.ProductController,
	jwtSvc auth.JWTService,
	reviewCtrl *controllers.ReviewController,
	trustUC trust.Usecase,
	productUC product.Usecase,
) {
	products := r.Group("/products")
	products.Use(middlewares.AuthMiddleware(jwtSvc))

	{
		products.POST("", middlewares.AuthorizeRoles("reseller"), productCtrl.Create)
		products.GET("", productCtrl.ListAvailable)
		products.GET("/title/:title", productCtrl.GetByTitle) 
		products.GET("/:id", productCtrl.GetByID)
		products.GET("/reseller/:id", productCtrl.ListByReseller)
		products.PUT("/:id", middlewares.AuthorizeRoles("reseller"), productCtrl.Update)
		products.DELETE("/:id", middlewares.AuthorizeRoles("reseller"), productCtrl.Delete)
		products.POST("/:id/reviews", middlewares.AuthorizeRoles("consumer"), reviewCtrl.SubmitReview)
	}

	// Separate reviews group
	reviews := r.Group("/reviews")
	reviews.Use(middlewares.AuthMiddleware(jwtSvc))
	{
		reviews.GET("/reseller/:id", reviewCtrl.GetResellerReviews)
	}
}
