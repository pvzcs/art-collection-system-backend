package handler

import (
	"art-collection-system/internal/service"
	"art-collection-system/internal/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdminHandler handles administrator-related HTTP requests
type AdminHandler struct {
	artworkService *service.ArtworkService
	adminService   *service.AdminService
}

// NewAdminHandler creates a new admin handler instance
func NewAdminHandler(artworkService *service.ArtworkService, adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{
		artworkService: artworkService,
		adminService:   adminService,
	}
}

// GetReviewQueue retrieves the list of artworks pending review
// GET /api/v1/admin/review-queue
func (h *AdminHandler) GetReviewQueue(c *gin.Context) {
	// Get pagination parameters
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 20
	}

	// Limit page size
	if pageSize > 100 {
		pageSize = 100
	}

	// Get review queue
	artworks, total, err := h.artworkService.GetReviewQueue(page, pageSize)
	if err != nil {
		utils.Error(c, 500, "获取审核队列失败")
		return
	}

	utils.Success(c, gin.H{
		"artworks":  artworks,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ReviewArtworkRequest represents the request body for reviewing an artwork
type ReviewArtworkRequest struct {
	Approved bool `json:"approved"`
}

// ReviewArtwork reviews a single artwork
// PUT /api/v1/admin/artworks/:id/review
func (h *AdminHandler) ReviewArtwork(c *gin.Context) {
	// Get artwork ID from URL parameter
	artworkIDStr := c.Param("id")
	artworkID, err := strconv.ParseUint(artworkIDStr, 10, 32)
	if err != nil {
		utils.Error(c, 400, "无效的作品ID")
		return
	}

	var req ReviewArtworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	// Review artwork
	if err := h.artworkService.ReviewArtwork(uint(artworkID), req.Approved); err != nil {
		if strings.Contains(err.Error(), "不存在") {
			utils.Error(c, 404, err.Error())
		} else {
			utils.Error(c, 500, "审核作品失败")
		}
		return
	}

	status := "未审核"
	if req.Approved {
		status = "已审核"
	}

	utils.Success(c, gin.H{
		"message": "审核成功",
		"status":  status,
	})
}

// BatchReviewArtworksRequest represents the request body for batch reviewing artworks
type BatchReviewArtworksRequest struct {
	ArtworkIDs []uint `json:"artwork_ids" binding:"required"`
	Approved   bool   `json:"approved"`
}

// BatchReviewArtworks reviews multiple artworks at once
// PUT /api/v1/admin/artworks/batch-review
func (h *AdminHandler) BatchReviewArtworks(c *gin.Context) {
	var req BatchReviewArtworksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	// Validate artwork IDs
	if len(req.ArtworkIDs) == 0 {
		utils.Error(c, 400, "作品ID列表不能为空")
		return
	}

	// Limit batch size
	if len(req.ArtworkIDs) > 100 {
		utils.Error(c, 400, "批量审核数量不能超过100个")
		return
	}

	// Batch review artworks
	if err := h.artworkService.BatchReviewArtworks(req.ArtworkIDs, req.Approved); err != nil {
		utils.Error(c, 500, "批量审核失败")
		return
	}

	status := "未审核"
	if req.Approved {
		status = "已审核"
	}

	utils.Success(c, gin.H{
		"message": "批量审核成功",
		"count":   len(req.ArtworkIDs),
		"status":  status,
	})
}

// ListUsers retrieves a paginated list of users
// GET /api/v1/admin/users
func (h *AdminHandler) ListUsers(c *gin.Context) {
	// Get pagination parameters
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 20
	}

	// Limit page size
	if pageSize > 100 {
		pageSize = 100
	}

	// Get users
	users, total, err := h.adminService.ListUsers(page, pageSize)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"users":     users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// UpdateUserRoleRequest represents the request body for updating user role
type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

// UpdateUserRole updates a user's role
// PUT /api/v1/admin/users/:id/role
func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	// Get user ID from URL parameter
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.Error(c, 400, "无效的用户ID")
		return
	}

	var req UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	// Update user role
	if err := h.adminService.UpdateUserRole(uint(userID), req.Role); err != nil {
		if strings.Contains(err.Error(), "无效的角色") {
			utils.Error(c, 400, err.Error())
		} else if strings.Contains(err.Error(), "不存在") {
			utils.Error(c, 404, err.Error())
		} else {
			utils.Error(c, 500, err.Error())
		}
		return
	}

	utils.Success(c, gin.H{"message": "更新角色成功"})
}

// GetUserStatistics retrieves statistics for a specific user
// GET /api/v1/admin/users/:id/statistics
func (h *AdminHandler) GetUserStatistics(c *gin.Context) {
	// Get user ID from URL parameter
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.Error(c, 400, "无效的用户ID")
		return
	}

	// Get user statistics
	statistics, err := h.adminService.GetUserStatistics(uint(userID))
	if err != nil {
		if strings.Contains(err.Error(), "不存在") {
			utils.Error(c, 404, err.Error())
		} else {
			utils.Error(c, 500, err.Error())
		}
		return
	}

	utils.Success(c, statistics)
}
