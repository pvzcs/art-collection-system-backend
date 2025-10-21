package handler

import (
	"art-collection-system/internal/service"
	"art-collection-system/internal/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ActivityHandler handles activity-related HTTP requests
type ActivityHandler struct {
	activityService *service.ActivityService
}

// NewActivityHandler creates a new activity handler instance
func NewActivityHandler(activityService *service.ActivityService) *ActivityHandler {
	return &ActivityHandler{
		activityService: activityService,
	}
}

// CreateActivityRequest represents the request body for creating an activity
type CreateActivityRequest struct {
	Name              string  `json:"name" binding:"required"`
	Description       string  `json:"description"`
	Deadline          *string `json:"deadline"`
	MaxUploadsPerUser int     `json:"max_uploads_per_user"`
}

// CreateActivity creates a new activity (admin only)
// POST /api/v1/admin/activities
func (h *ActivityHandler) CreateActivity(c *gin.Context) {
	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	// Parse deadline if provided
	var deadline *time.Time
	if req.Deadline != nil && *req.Deadline != "" {
		parsedTime, err := time.Parse(time.RFC3339, *req.Deadline)
		if err != nil {
			utils.Error(c, 400, "截止日期格式不正确，请使用 RFC3339 格式")
			return
		}
		deadline = &parsedTime
	}

	// Set default max uploads if not provided
	maxUploads := req.MaxUploadsPerUser
	if maxUploads <= 0 {
		maxUploads = 5
	}

	// Create activity
	activity, err := h.activityService.CreateActivity(req.Name, req.Description, deadline, maxUploads)
	if err != nil {
		utils.Error(c, 500, "创建活动失败")
		return
	}

	utils.Success(c, activity)
}

// UpdateActivityRequest represents the request body for updating an activity
type UpdateActivityRequest struct {
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Deadline          *string `json:"deadline"`
	MaxUploadsPerUser int     `json:"max_uploads_per_user"`
}

// UpdateActivity updates an existing activity (admin only)
// PUT /api/v1/admin/activities/:id
func (h *ActivityHandler) UpdateActivity(c *gin.Context) {
	// Get activity ID from URL parameter
	activityIDStr := c.Param("id")
	activityID, err := strconv.ParseUint(activityIDStr, 10, 32)
	if err != nil {
		utils.Error(c, 400, "无效的活动ID")
		return
	}

	var req UpdateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	// Parse deadline if provided
	var deadline *time.Time
	if req.Deadline != nil && *req.Deadline != "" {
		parsedTime, err := time.Parse(time.RFC3339, *req.Deadline)
		if err != nil {
			utils.Error(c, 400, "截止日期格式不正确，请使用 RFC3339 格式")
			return
		}
		deadline = &parsedTime
	}

	// Update activity
	if err := h.activityService.UpdateActivity(uint(activityID), req.Name, req.Description, deadline, req.MaxUploadsPerUser); err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.Error(c, 404, "活动不存在")
		} else {
			utils.Error(c, 500, "更新活动失败")
		}
		return
	}

	utils.Success(c, gin.H{"message": "更新成功"})
}

// DeleteActivity soft deletes an activity (admin only)
// DELETE /api/v1/admin/activities/:id
func (h *ActivityHandler) DeleteActivity(c *gin.Context) {
	// Get activity ID from URL parameter
	activityIDStr := c.Param("id")
	activityID, err := strconv.ParseUint(activityIDStr, 10, 32)
	if err != nil {
		utils.Error(c, 400, "无效的活动ID")
		return
	}

	// Delete activity
	if err := h.activityService.DeleteActivity(uint(activityID)); err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.Error(c, 404, "活动不存在")
		} else {
			utils.Error(c, 500, "删除活动失败")
		}
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}

// GetActivity retrieves a single activity by ID
// GET /api/v1/activities/:id
func (h *ActivityHandler) GetActivity(c *gin.Context) {
	// Get activity ID from URL parameter
	activityIDStr := c.Param("id")
	activityID, err := strconv.ParseUint(activityIDStr, 10, 32)
	if err != nil {
		utils.Error(c, 400, "无效的活动ID")
		return
	}

	// Get activity
	activity, err := h.activityService.GetActivityByID(uint(activityID))
	if err != nil {
		utils.Error(c, 404, "活动不存在")
		return
	}

	utils.Success(c, activity)
}

// ListActivities retrieves a paginated list of activities
// GET /api/v1/activities
func (h *ActivityHandler) ListActivities(c *gin.Context) {
	// Get pagination parameters
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	// Limit page size to prevent excessive queries
	if pageSize > 100 {
		pageSize = 100
	}

	// Get activities
	activities, total, err := h.activityService.ListActivities(page, pageSize)
	if err != nil {
		utils.Error(c, 500, "获取活动列表失败")
		return
	}

	utils.Success(c, gin.H{
		"activities": activities,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	})
}
