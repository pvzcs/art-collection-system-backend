package repository

import (
	"art-collection-system/internal/models"
	"gorm.io/gorm"
)

// ActivityRepository handles activity data access operations
type ActivityRepository struct {
	db *gorm.DB
}

// NewActivityRepository creates a new activity repository instance
func NewActivityRepository(db *gorm.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

// Create creates a new activity in the database
func (r *ActivityRepository) Create(activity *models.Activity) error {
	return r.db.Create(activity).Error
}

// Update updates activity information
func (r *ActivityRepository) Update(activity *models.Activity) error {
	return r.db.Save(activity).Error
}

// UpdateFields updates specific fields of an activity
func (r *ActivityRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&models.Activity{}).Where("id = ? AND is_deleted = ?", id, false).Updates(fields).Error
}

// SoftDelete marks an activity as deleted (soft delete)
func (r *ActivityRepository) SoftDelete(id uint) error {
	return r.db.Model(&models.Activity{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// GetByID retrieves an activity by ID (excluding deleted activities)
func (r *ActivityRepository) GetByID(id uint) (*models.Activity, error) {
	var activity models.Activity
	err := r.db.Where("is_deleted = ?", false).First(&activity, id).Error
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

// GetByIDIncludeDeleted retrieves an activity by ID (including deleted activities)
func (r *ActivityRepository) GetByIDIncludeDeleted(id uint) (*models.Activity, error) {
	var activity models.Activity
	err := r.db.First(&activity, id).Error
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

// List retrieves a paginated list of activities (excluding deleted activities)
func (r *ActivityRepository) List(page, pageSize int) ([]models.Activity, int64, error) {
	var activities []models.Activity
	var total int64

	// Count total non-deleted activities
	if err := r.db.Model(&models.Activity{}).Where("is_deleted = ?", false).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Retrieve paginated activities
	err := r.db.Where("is_deleted = ?", false).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&activities).Error

	if err != nil {
		return nil, 0, err
	}

	return activities, total, nil
}

// ListAll retrieves all activities (excluding deleted activities)
func (r *ActivityRepository) ListAll() ([]models.Activity, error) {
	var activities []models.Activity
	err := r.db.Where("is_deleted = ?", false).Order("created_at DESC").Find(&activities).Error
	if err != nil {
		return nil, err
	}
	return activities, nil
}

// Exists checks if an activity exists and is not deleted
func (r *ActivityRepository) Exists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Activity{}).Where("id = ? AND is_deleted = ?", id, false).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
