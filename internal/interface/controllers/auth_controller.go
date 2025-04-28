package controllers

import (
	"net/http"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/auth"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUC auth.AuthUsecase
}

func NewAuthController(authUC auth.AuthUsecase) *AuthController {
	return &AuthController{authUC: authUC}
}

// POST /auth/register
func (a *AuthController) Register(c *gin.Context) {
	var newUser user.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	result, err := a.authUC.Register(c.Request.Context(), newUser)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": result.Token,
		"user": gin.H{
			"id":       result.ID,
			"username": result.Username,
			"role":     result.Role,
		},
	})
}

// POST /auth/login
func (a *AuthController) Login(c *gin.Context) {
	var creds auth.LoginCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	result, err := a.authUC.Login(c.Request.Context(), creds)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": result.Token,
		"user": gin.H{
			"id":       result.ID,
			"username": result.Username,
			"role":     result.Role,
		},
	})
}
