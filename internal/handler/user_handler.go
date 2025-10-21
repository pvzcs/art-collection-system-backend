package handler

import (
	"art-collection-system/internal/service"
	"art-collection-system/internal/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler instance
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile retrieves the current user's profile
// GET /api/v1/user/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, 401, "未授权")
		return
	}

	// Get user information
	user, err := h.userService.GetUserByID(userID.(uint))
	if err != nil {
		utils.Error(c, 404, err.Error())
		return
	}

	// Return user info (without password)
	utils.Success(c, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"nickname":   user.Nickname,
		"role":       user.Role,
		"created_at": user.CreatedAt,
	})
}

// UpdateProfileRequest represents the request body for updating profile
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"required"`
}

// UpdateProfile updates the current user's profile
// PUT /api/v1/user/profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, 401, "未授权")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	// Validate nickname
	if err := utils.ValidateNickname(req.Nickname); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	// Update profile
	if err := h.userService.UpdateProfile(userID.(uint), req.Nickname); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "更新成功"})
}

// ChangePasswordRequest represents the request body for changing password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// ChangePassword changes the current user's password
// PUT /api/v1/user/password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, 401, "未授权")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	// Validate new password
	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	// Change password
	if err := h.userService.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword); err != nil {
		if strings.Contains(err.Error(), "旧密码") {
			utils.Error(c, 400, err.Error())
		} else {
			utils.Error(c, 500, err.Error())
		}
		return
	}

	utils.Success(c, gin.H{"message": "密码修改成功"})
}

// GetUserArtworks retrieves all artworks for a specific user
// GET /api/v1/users/:id/artworks
func (h *UserHandler) GetUserArtworks(c *gin.Context) {
	// Get target user ID from URL parameter
	targetUserIDStr := c.Param("id")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		utils.Error(c, 400, "无效的用户ID")
		return
	}

	// Get requester info from context
	requesterID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, 401, "未授权")
		return
	}

	requesterRole, exists := c.Get("user_role")
	if !exists {
		utils.Error(c, 401, "未授权")
		return
	}

	// Get user artworks with permission check
	artworks, err := h.userService.GetUserArtworks(uint(targetUserID), requesterID.(uint), requesterRole.(string))
	if err != nil {
		if strings.Contains(err.Error(), "权限不足") {
			utils.Error(c, 403, err.Error())
		} else if strings.Contains(err.Error(), "不存在") {
			utils.Error(c, 404, err.Error())
		} else {
			utils.Error(c, 500, err.Error())
		}
		return
	}

	utils.Success(c, gin.H{
		"artworks": artworks,
		"total":    len(artworks),
	})
}
