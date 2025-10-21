package service

import (
	"art-collection-system/internal/models"
	"art-collection-system/internal/repository"
	"errors"
	"time"
)

// ActivityService handles business logic for activities
type ActivityService struct {
	repo *repository.ActivityRepository
}

// NewActivityService creates a new activity service instance
func NewActivityService(repo *repository.ActivityRepository) *ActivityService {
	return &ActivityService{repo: repo}
}

// CreateActivity creates a new activity
// Requirements: 3.1
func (s *ActivityService) CreateActivity(name, description string, deadline *time.Time, maxUploads int) (*models.Activity, error) {
	if name == "" {
		return nil, errors.New("activity name is required")
	}

	if maxUploads <= 0 {
		maxUploads = 5 // Default value
	}

	activity := &models.Activity{
		Name:              name,
		Description:       description,
		Deadline:          deadline,
		MaxUploadsPerUser: maxUploads,
		IsDeleted:         false,
	}

	if err := s.repo.Create(activity); err != nil {
		return nil, err
	}

	return activity, nil
}

// UpdateActivity updates an existing activity
// Requirements: 3.2
func (s *ActivityService) UpdateActivity(id uint, name, description string, deadline *time.Time, maxUploads int) error {
	// Check if activity exists
	activity, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("activity not found")
	}

	// Update fields
	if name != "" {
		activity.Name = name
	}
	activity.Description = description
	activity.Deadline = deadline
	if maxUploads > 0 {
		activity.MaxUploadsPerUser = maxUploads
	}

	return s.repo.Update(activity)
}

// DeleteActivity soft deletes an activity
// Requirements: 3.3
func (s *ActivityService) DeleteActivity(id uint) error {
	// Check if activity exists
	exists, err := s.repo.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("activity not found")
	}

	return s.repo.SoftDelete(id)
}

// GetActivityByID retrieves an activity by ID
// Requirements: 3.4
func (s *ActivityService) GetActivityByID(id uint) (*models.Activity, error) {
	activity, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("activity not found")
	}
	return activity, nil
}

// ListActivities retrieves a paginated list of activities
// Requirements: 3.4
func (s *ActivityService) ListActivities(page, pageSize int) ([]models.Activity, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	return s.repo.List(page, pageSize)
}

// IsActivityActive checks if an activity is active (exists, not deleted, not expired)
// Requirements: 3.5
func (s *ActivityService) IsActivityActive(id uint) (bool, error) {
	// Check if activity exists and is not deleted
	activity, err := s.repo.GetByID(id)
	if err != nil {
		return false, nil // Activity doesn't exist or is deleted
	}

	// Check if activity has expired
	if activity.Deadline != nil && time.Now().After(*activity.Deadline) {
		return false, nil // Activity has expired
	}

	return true, nil
}
