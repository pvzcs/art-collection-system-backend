package repository

import (
	"art-collection-system/internal/models"
	"gorm.io/gorm"
)

// ArtworkRepository handles artwork data access operations
type ArtworkRepository struct {
	db *gorm.DB
}

// NewArtworkRepository creates a new artwork repository instance
func NewArtworkRepository(db *gorm.DB) *ArtworkRepository {
	return &ArtworkRepository{db: db}
}

// Create creates a new artwork in the database
func (r *ArtworkRepository) Create(artwork *models.Artwork) error {
	return r.db.Create(artwork).Error
}

// Delete deletes an artwork from the database
func (r *ArtworkRepository) Delete(id uint) error {
	return r.db.Delete(&models.Artwork{}, id).Error
}

// GetByID retrieves an artwork by ID
func (r *ArtworkRepository) GetByID(id uint) (*models.Artwork, error) {
	var artwork models.Artwork
	err := r.db.First(&artwork, id).Error
	if err != nil {
		return nil, err
	}
	return &artwork, nil
}

// GetByIDWithRelations retrieves an artwork by ID with related activity and user data
func (r *ArtworkRepository) GetByIDWithRelations(id uint) (*models.Artwork, error) {
	var artwork models.Artwork
	err := r.db.Preload("Activity").Preload("User").First(&artwork, id).Error
	if err != nil {
		return nil, err
	}
	return &artwork, nil
}

// GetByUserID retrieves all artworks by a specific user
func (r *ArtworkRepository) GetByUserID(userID uint) ([]models.Artwork, error) {
	var artworks []models.Artwork
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&artworks).Error
	if err != nil {
		return nil, err
	}
	return artworks, nil
}

// GetByUserIDWithRelations retrieves all artworks by a specific user with related data
func (r *ArtworkRepository) GetByUserIDWithRelations(userID uint) ([]models.Artwork, error) {
	var artworks []models.Artwork
	err := r.db.Preload("Activity").Where("user_id = ?", userID).Order("created_at DESC").Find(&artworks).Error
	if err != nil {
		return nil, err
	}
	return artworks, nil
}

// GetByActivityID retrieves all artworks for a specific activity
func (r *ArtworkRepository) GetByActivityID(activityID uint) ([]models.Artwork, error) {
	var artworks []models.Artwork
	err := r.db.Where("activity_id = ?", activityID).Order("created_at DESC").Find(&artworks).Error
	if err != nil {
		return nil, err
	}
	return artworks, nil
}

// GetReviewQueue retrieves artworks pending review with pagination
func (r *ArtworkRepository) GetReviewQueue(page, pageSize int) ([]models.Artwork, int64, error) {
	var artworks []models.Artwork
	var total int64

	// Count total pending artworks
	if err := r.db.Model(&models.Artwork{}).Where("review_status = ?", models.StatusPending).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Retrieve paginated pending artworks with user and activity information
	err := r.db.Preload("User").Preload("Activity").
		Where("review_status = ?", models.StatusPending).
		Order("created_at ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&artworks).Error

	if err != nil {
		return nil, 0, err
	}

	return artworks, total, nil
}

// CountByUserAndActivity counts artworks by a user in a specific activity
func (r *ArtworkRepository) CountByUserAndActivity(userID, activityID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Artwork{}).
		Where("user_id = ? AND activity_id = ?", userID, activityID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// UpdateReviewStatus updates the review status of a single artwork
func (r *ArtworkRepository) UpdateReviewStatus(id uint, status models.ReviewStatus) error {
	return r.db.Model(&models.Artwork{}).Where("id = ?", id).Update("review_status", status).Error
}

// BatchUpdateReviewStatus updates the review status of multiple artworks
func (r *ArtworkRepository) BatchUpdateReviewStatus(ids []uint, status models.ReviewStatus) error {
	return r.db.Model(&models.Artwork{}).Where("id IN ?", ids).Update("review_status", status).Error
}

// Update updates artwork information
func (r *ArtworkRepository) Update(artwork *models.Artwork) error {
	return r.db.Save(artwork).Error
}

// Exists checks if an artwork exists
func (r *ArtworkRepository) Exists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Artwork{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
