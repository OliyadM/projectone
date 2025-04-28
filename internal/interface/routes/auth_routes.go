package routes

import (
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine, authCtrl *controllers.AuthController) {
	authGroup := r.Group("/auth")

	authGroup.POST("/register", authCtrl.Register)
	authGroup.POST("/login", authCtrl.Login)
}
