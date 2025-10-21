package service

import (
	"art-collection-system/internal/models"
	"art-collection-system/internal/repository"
	"art-collection-system/internal/utils"
	"errors"
	"gorm.io/gorm"
)

// UserService handles user-related business logic
type UserService struct {
	userRepo    *repository.UserRepository
	artworkRepo *repository.ArtworkRepository
}

// NewUserService creates a new user service instance
func NewUserService(userRepo *repository.UserRepository, artworkRepo *repository.ArtworkRepository) *UserService {
	return &UserService{
		userRepo:    userRepo,
		artworkRepo: artworkRepo,
	}
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return user, nil
}

// UpdateProfile updates user profile information (nickname)
func (s *UserService) UpdateProfile(userID uint, nickname string) error {
	// Validate nickname
	if nickname == "" {
		return errors.New("昵称不能为空")
	}

	// Check if user exists
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// Update nickname
	user.Nickname = nickname
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("更新用户信息失败")
	}

	return nil
}

// ChangePassword changes user password
func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	// Validate passwords
	if oldPassword == "" || newPassword == "" {
		return errors.New("密码不能为空")
	}

	if len(newPassword) < 8 {
		return errors.New("新密码长度至少为8位")
	}

	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// Verify old password
	if err := utils.ComparePassword(user.Password, oldPassword); err != nil {
		return errors.New("旧密码不正确")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// Update password
	user.Password = hashedPassword
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("更新密码失败")
	}

	return nil
}

// GetUserArtworks retrieves all artworks for a specific user with permission check
// Only the user themselves or an administrator can access
func (s *UserService) GetUserArtworks(userID uint, requesterID uint, requesterRole string) ([]models.Artwork, error) {
	// Permission check: only the user themselves or admin can access
	if requesterRole != "admin" && requesterID != userID {
		return nil, errors.New("权限不足，仅本人或管理员可访问")
	}

	// Check if user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// Get user's artworks with related activity information
	artworks, err := s.artworkRepo.GetByUserIDWithRelations(userID)
	if err != nil {
		return nil, errors.New("获取用户作品失败")
	}

	return artworks, nil
}
