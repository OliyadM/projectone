package routes

import (
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/auth"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/controllers"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterCartItemRoutes(r *gin.Engine, ctrl *controllers.CartItemController, jwtSvc auth.JWTService) {
	// Cart group for cart item related routes.
	cartGroup := r.Group("/api/cart")
	cartGroup.Use(middlewares.AuthMiddleware(jwtSvc))

	// Route to add an item to a cart => POST /api/cart/items
	cartGroup.POST("/items", middlewares.AuthorizeRoles("consumer"), ctrl.AddCartItem)

	// Route to retrieve all cart items for a user => GET /api/cart
	cartGroup.GET("", middlewares.AuthorizeRoles("consumer"), ctrl.GetCartItems)

	// Route to remove a cart item => DELETE /api/cart/items/:listingID
	cartGroup.DELETE("/items/:listingID", middlewares.AuthorizeRoles("consumer"), ctrl.RemoveCartItem)

	// Checkout route. Although related to the cart, it is defined separately.
	checkoutGroup := r.Group("/api/checkout")
	checkoutGroup.Use(middlewares.AuthMiddleware(jwtSvc))
	checkoutGroup.POST("", middlewares.AuthorizeRoles("consumer"), ctrl.CheckoutCart)
	checkoutGroup.POST("/:listingId", middlewares.AuthorizeRoles("consumer"), ctrl.CheckoutSingleItem)
}
