package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"
)

type UpdateProfileRequest struct {
	Name     string `json:"name,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

type UserController struct {
	userUsecase user.Usecase
}

func NewUserController(userUsecase user.Usecase) *UserController {
	return &UserController{
		userUsecase: userUsecase,
	}
}

// GetUserByID handles GET /api/users/:id
func (c *UserController) GetUserByID(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	user, err := c.userUsecase.GetByID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Only return public user information
	ctx.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})
}

func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	var req UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Create update map with only non-empty fields
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.ImageURL != "" {
		updates["image_url"] = req.ImageURL
	}

	if len(updates) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	// Get current user to check username and email
	currentUser, err := c.userUsecase.GetByID(ctx, userIDStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user information"})
		return
	}

	// Check if username is already taken if it's being updated
	if req.Username != "" && req.Username != currentUser.Username {
		// Get all users to check for duplicate username
		users, err := c.userUsecase.ListByRole(ctx, user.Role(currentUser.Role))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check username availability"})
			return
		}

		for _, u := range users {
			if u.Username == req.Username && u.ID != userIDStr {
				ctx.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
				return
			}
		}
	}

	// Check if email is already taken if it's being updated
	if req.Email != "" && req.Email != currentUser.Email {
		// Get all users to check for duplicate email
		users, err := c.userUsecase.ListByRole(ctx, user.Role(currentUser.Role))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check email availability"})
			return
		}

		for _, u := range users {
			if u.Email == req.Email && u.ID != userIDStr {
				ctx.JSON(http.StatusConflict, gin.H{"error": "email already taken"})
				return
			}
		}
	}

	err = c.userUsecase.Update(ctx, userIDStr, updates)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
} 