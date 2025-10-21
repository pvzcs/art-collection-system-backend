package handler

import (
	"art-collection-system/internal/models"
	"art-collection-system/internal/service"
	"art-collection-system/internal/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ArtworkHandler handles artwork-related HTTP requests
type ArtworkHandler struct {
	artworkService *service.ArtworkService
	fileService    *service.FileService
}

// NewArtworkHandler creates a new artwork handler instance
func NewArtworkHandler(artworkService *service.ArtworkService, fileService *service.FileService) *ArtworkHandler {
	return &ArtworkHandler{
		artworkService: artworkService,
		fileService:    fileService,
	}
}

// UploadArtwork handles artwork file upload
// POST /api/v1/artworks
func (h *ArtworkHandler) UploadArtwork(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, 401, "未授权")
		return
	}

	// Get activity ID from form
	activityIDStr := c.PostForm("activity_id")
	if activityIDStr == "" {
		utils.Error(c, 400, "活动ID不能为空")
		return
	}

	activityID, err := strconv.ParseUint(activityIDStr, 10, 32)
	if err != nil {
		utils.Error(c, 400, "无效的活动ID")
		return
	}

	// Get uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.Error(c, 400, "请上传文件")
		return
	}
	defer file.Close()

	// Validate file (size, type, and content)
	if err := utils.ValidateImageFile(header); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	// Upload artwork
	artwork, err := h.artworkService.UploadArtwork(userID.(uint), uint(activityID), file, header.Filename)
	if err != nil {
		if strings.Contains(err.Error(), "活动") {
			utils.Error(c, 400, err.Error())
		} else if strings.Contains(err.Error(), "上传数量") {
			utils.Error(c, 400, err.Error())
		} else {
			utils.Error(c, 500, "上传作品失败")
		}
		return
	}

	utils.Success(c, gin.H{
		"id":            artwork.ID,
		"activity_id":   artwork.ActivityID,
		"file_name":     artwork.FileName,
		"review_status": artwork.ReviewStatus,
		"created_at":    artwork.CreatedAt,
	})
}

// DeleteArtwork deletes an artwork
// DELETE /api/v1/artworks/:id
func (h *ArtworkHandler) DeleteArtwork(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, 401, "未授权")
		return
	}

	// Get artwork ID from URL parameter
	artworkIDStr := c.Param("id")
	artworkID, err := strconv.ParseUint(artworkIDStr, 10, 32)
	if err != nil {
		utils.Error(c, 400, "无效的作品ID")
		return
	}

	// Delete artwork
	if err := h.artworkService.DeleteArtwork(uint(artworkID), userID.(uint)); err != nil {
		if strings.Contains(err.Error(), "权限") {
			utils.Error(c, 403, err.Error())
		} else if strings.Contains(err.Error(), "不存在") {
			utils.Error(c, 404, err.Error())
		} else {
			utils.Error(c, 500, "删除作品失败")
		}
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}

// GetArtwork retrieves artwork information with permission check
// GET /api/v1/artworks/:id
func (h *ArtworkHandler) GetArtwork(c *gin.Context) {
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

	// Get artwork ID from URL parameter
	artworkIDStr := c.Param("id")
	artworkID, err := strconv.ParseUint(artworkIDStr, 10, 32)
	if err != nil {
		utils.Error(c, 400, "无效的作品ID")
		return
	}

	// Get artwork with permission check
	artwork, err := h.artworkService.GetArtwork(uint(artworkID), requesterID.(uint), requesterRole.(string))
	if err != nil {
		if strings.Contains(err.Error(), "权限") {
			utils.Error(c, 403, err.Error())
		} else if strings.Contains(err.Error(), "不存在") {
			utils.Error(c, 404, err.Error())
		} else {
			utils.Error(c, 500, "获取作品失败")
		}
		return
	}

	utils.Success(c, artwork)
}

// ServeImage serves the artwork image file with permission check
// GET /api/v1/artworks/:id/image
func (h *ArtworkHandler) ServeImage(c *gin.Context) {
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

	// Get artwork ID from URL parameter
	artworkIDStr := c.Param("id")
	artworkID, err := strconv.ParseUint(artworkIDStr, 10, 32)
	if err != nil {
		utils.Error(c, 400, "无效的作品ID")
		return
	}

	// Serve file through proxy (permission check is done inside ServeFile)
	// We need to get the artwork first to get the file path
	artworkInterface, err := h.artworkService.GetArtwork(uint(artworkID), requesterID.(uint), requesterRole.(string))
	if err != nil {
		if strings.Contains(err.Error(), "权限") || strings.Contains(err.Error(), "permission") {
			utils.Error(c, 403, err.Error())
		} else if strings.Contains(err.Error(), "不存在") || strings.Contains(err.Error(), "not found") {
			utils.Error(c, 404, err.Error())
		} else {
			utils.Error(c, 500, "获取作品失败")
		}
		return
	}

	// Type assert to get the actual artwork model
	artwork, ok := artworkInterface.(*models.Artwork)
	if !ok {
		utils.Error(c, 500, "内部错误")
		return
	}

	// Serve file through proxy
	fileData, contentType, err := h.fileService.ServeFile(artwork.FilePath, uint(artworkID), requesterID.(uint), requesterRole.(string), h.artworkService)
	if err != nil {
		if strings.Contains(err.Error(), "权限") || strings.Contains(err.Error(), "permission") {
			utils.Error(c, 403, err.Error())
		} else {
			utils.Error(c, 500, "读取文件失败")
		}
		return
	}

	// Set content type and return file
	c.Data(200, contentType, fileData)
}
