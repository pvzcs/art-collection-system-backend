package service

import (
	"art-collection-system/internal/models"
	"art-collection-system/internal/repository"
	"errors"
	"gorm.io/gorm"
)

// AdminService handles administrator-related business logic
type AdminService struct {
	userRepo *repository.UserRepository
}

// NewAdminService creates a new admin service instance
func NewAdminService(userRepo *repository.UserRepository) *AdminService {
	return &AdminService{
		userRepo: userRepo,
	}
}

// ListUsers retrieves a paginated list of users
func (s *AdminService) ListUsers(page, pageSize int) ([]models.User, int64, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20 // Default page size
	}

	users, total, err := s.userRepo.List(page, pageSize)
	if err != nil {
		return nil, 0, errors.New("获取用户列表失败")
	}

	return users, total, nil
}

// UpdateUserRole updates a user's role
func (s *AdminService) UpdateUserRole(userID uint, role string) error {
	// Validate role
	if role != "user" && role != "admin" {
		return errors.New("无效的角色，只能是 'user' 或 'admin'")
	}

	// Check if user exists
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// Update role
	user.Role = role
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("更新用户角色失败")
	}

	return nil
}

// GetUserStatistics retrieves statistics for a specific user
func (s *AdminService) GetUserStatistics(userID uint) (map[string]interface{}, error) {
	// Check if user exists
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// Count total artworks
	totalArtworks, err := s.userRepo.CountArtworks(userID)
	if err != nil {
		return nil, errors.New("获取作品统计失败")
	}

	// Count approved artworks
	approvedArtworks, err := s.userRepo.CountArtworksByStatus(userID, string(models.StatusApproved))
	if err != nil {
		return nil, errors.New("获取已审核作品统计失败")
	}

	// Count pending artworks
	pendingArtworks, err := s.userRepo.CountArtworksByStatus(userID, string(models.StatusPending))
	if err != nil {
		return nil, errors.New("获取待审核作品统计失败")
	}

	// Build statistics response
	statistics := map[string]interface{}{
		"user_id":           user.ID,
		"email":             user.Email,
		"nickname":          user.Nickname,
		"role":              user.Role,
		"total_artworks":    totalArtworks,
		"approved_artworks": approvedArtworks,
		"pending_artworks":  pendingArtworks,
		"created_at":        user.CreatedAt,
	}

	return statistics, nil
}
