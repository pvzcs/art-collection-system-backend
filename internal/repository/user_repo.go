package repository

import (
	"art-collection-system/internal/models"
	"gorm.io/gorm"
)

// UserRepository handles user data access operations
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user in the database
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email address
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates user information
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// UpdateFields updates specific fields of a user
func (r *UserRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(fields).Error
}

// EmailExists checks if an email already exists in the database
func (r *UserRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// List retrieves users with pagination
func (r *UserRepository) List(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Count total users
	if err := r.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Retrieve paginated users
	err := r.db.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// CountArtworks counts the total number of artworks for a user
func (r *UserRepository) CountArtworks(userID uint) (int64, error) {
	var count int64
	err := r.db.Table("artworks").Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CountArtworksByStatus counts artworks by review status for a user
func (r *UserRepository) CountArtworksByStatus(userID uint, status string) (int64, error) {
	var count int64
	err := r.db.Table("artworks").
		Where("user_id = ? AND review_status = ?", userID, status).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
