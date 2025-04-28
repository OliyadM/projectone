package controllers

import (
	"net/http"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/order"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"
	"github.com/gin-gonic/gin"
)

type AdminController struct {
	userUC  user.Usecase
	orderUC order.Usecase
}

func NewAdminController(userUC user.Usecase, orderUC order.Usecase) *AdminController {
	return &AdminController{userUC: userUC, orderUC: orderUC}
}

// GET /api/admin/users
func (a *AdminController) GetAllUsers(c *gin.Context) {
	roleParam := c.Query("role")

	// If no role is passed, return all users by calling ListByRole for each role
	if roleParam == "" {
		allRoles := []user.Role{
			user.RoleSupplier,
			user.RoleReseller,
			user.RoleConsumer,
			user.RoleAdmin,
		}

		var allUsers []*user.User
		for _, role := range allRoles {
			usersByRole, err := a.userUC.ListByRole(c.Request.Context(), role)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users by role"})
				return
			}
			allUsers = append(allUsers, usersByRole...)
		}

		c.JSON(http.StatusOK, allUsers)
		return
	}

	// Otherwise, get users by specified role
	role := user.Role(roleParam)
	users, err := a.userUC.ListByRole(c.Request.Context(), role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}
func (a *AdminController) DeleteUserIfBlacklisted(c *gin.Context) {
	userID := c.Param("userId")

	userData, err := a.userUC.GetByID(c.Request.Context(), userID)
	if err != nil || userData == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Only allow suppliers or resellers
	if userData.Role != string(user.RoleSupplier) && userData.Role != string(user.RoleReseller) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only suppliers or resellers can be blacklisted"})
		return
	}

	if userData.TrustScore >= 60 {
		c.JSON(http.StatusForbidden, gin.H{"error": "User not eligible for deletion"})
		return
	}

	// Soft delete (deactivation)
	err = a.userUC.Update(c.Request.Context(), userID, map[string]interface{}{"is_deleted": true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User successfully deactivated"})
}
func (a *AdminController) GetTrustScores(c *gin.Context) {
	roleParam := c.Query("role")

	var roles []user.Role
	if roleParam != "" {
		roles = append(roles, user.Role(roleParam))
	} else {
		roles = []user.Role{user.RoleSupplier, user.RoleReseller}
	}

	var result []gin.H

	for _, r := range roles {
		users, err := a.userUC.ListByRole(c.Request.Context(), r)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}

		for _, u := range users {
			status := "active"
			if u.TrustScore < 60 {
				status = "blacklisted"
			}
			result = append(result, gin.H{
				"userId":     u.ID,
				"name":       u.Name,
				"role":       u.Role,
				"trustScore": u.TrustScore,
				"status":     status,
			})
		}
	}

	c.JSON(http.StatusOK, result)
}

// GET /api/admin/blacklisted-users
func (a *AdminController) GetBlacklistedUsers(c *gin.Context) {
	users, err := a.userUC.GetBlacklistedUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch blacklisted users"})
		return
	}
	c.JSON(http.StatusOK, users)
}
func (a *AdminController) GetDashboardMetrics(c *gin.Context) {
	metrics, err := a.orderUC.GetAdminDashboardMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch dashboard metrics",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Dashboard metrics retrieved successfully",
		"data":    metrics,
	})
}
