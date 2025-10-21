package service

import (
	"art-collection-system/internal/models"
	"art-collection-system/internal/repository"
	"errors"
	"mime/multipart"
)

// ArtworkService handles business logic for artworks
type ArtworkService struct {
	repo            *repository.ArtworkRepository
	activityService *ActivityService
	fileService     *FileService
}

// NewArtworkService creates a new artwork service instance
func NewArtworkService(repo *repository.ArtworkRepository, activityService *ActivityService, fileService *FileService) *ArtworkService {
	return &ArtworkService{
		repo:            repo,
		activityService: activityService,
		fileService:     fileService,
	}
}

// UploadArtwork handles artwork upload with validation
// Requirements: 4.1, 4.2, 4.3, 4.4, 5.1
func (s *ArtworkService) UploadArtwork(userID, activityID uint, file multipart.File, filename string) (*models.Artwork, error) {
	// Validate activity is active
	isActive, err := s.activityService.IsActivityActive(activityID)
	if err != nil {
		return nil, err
	}
	if !isActive {
		return nil, errors.New("activity is not active or has expired")
	}

	// Check upload limit
	canUpload, err := s.CheckUploadLimit(userID, activityID)
	if err != nil {
		return nil, err
	}
	if !canUpload {
		return nil, errors.New("upload limit exceeded for this activity")
	}

	// Save file
	filePath, err := s.fileService.SaveFile(file, filename)
	if err != nil {
		return nil, err
	}

	// Create artwork record with pending status
	artwork := &models.Artwork{
		ActivityID:   activityID,
		UserID:       userID,
		FilePath:     filePath,
		FileName:     filename,
		ReviewStatus: models.StatusPending,
	}

	if err := s.repo.Create(artwork); err != nil {
		// If database creation fails, try to delete the uploaded file
		_ = s.fileService.DeleteFile(filePath)
		return nil, err
	}

	return artwork, nil
}

// GetArtwork retrieves an artwork with permission validation
// Requirements: 5.5, 6.1, 6.2, 6.3
func (s *ArtworkService) GetArtwork(artworkID, requesterID uint, requesterRole string) (interface{}, error) {
	// Retrieve artwork
	artwork, err := s.repo.GetByID(artworkID)
	if err != nil {
		return nil, errors.New("artwork not found")
	}

	// Admin can access all artworks
	if requesterRole == "admin" {
		return artwork, nil
	}

	// Pending artworks are only visible to admins
	if artwork.ReviewStatus == models.StatusPending {
		return nil, errors.New("permission denied: artwork is pending review")
	}

	// Approved artworks are only visible to the author
	if artwork.UserID != requesterID {
		return nil, errors.New("permission denied: you can only view your own artworks")
	}

	return artwork, nil
}

// ReviewArtwork updates the review status of a single artwork
// Requirements: 5.2, 5.3, 5.4
func (s *ArtworkService) ReviewArtwork(artworkID uint, approved bool) error {
	// Check if artwork exists
	exists, err := s.repo.Exists(artworkID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("artwork not found")
	}

	// Determine new status
	status := models.StatusPending
	if approved {
		status = models.StatusApproved
	}

	// Update review status
	return s.repo.UpdateReviewStatus(artworkID, status)
}

// BatchReviewArtworks updates the review status of multiple artworks
// Requirements: 5.2, 5.3, 5.4
func (s *ArtworkService) BatchReviewArtworks(artworkIDs []uint, approved bool) error {
	if len(artworkIDs) == 0 {
		return errors.New("no artwork IDs provided")
	}

	// Determine new status
	status := models.StatusPending
	if approved {
		status = models.StatusApproved
	}

	// Batch update review status
	return s.repo.BatchUpdateReviewStatus(artworkIDs, status)
}

// GetReviewQueue retrieves pending artworks sorted by upload time
// Requirements: 5.3, 5.4, 8.5
func (s *ArtworkService) GetReviewQueue(page, pageSize int) ([]models.Artwork, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	return s.repo.GetReviewQueue(page, pageSize)
}

// DeleteArtwork deletes an artwork and its associated file
// Requirements: 4.5
func (s *ArtworkService) DeleteArtwork(artworkID, userID uint) error {
	// Retrieve artwork
	artwork, err := s.repo.GetByID(artworkID)
	if err != nil {
		return errors.New("artwork not found")
	}

	// Verify user is the owner
	if artwork.UserID != userID {
		return errors.New("permission denied: you can only delete your own artworks")
	}

	// Delete the physical file
	if err := s.fileService.DeleteFile(artwork.FilePath); err != nil {
		// Log error but continue with database deletion
		// File might already be deleted or path might be invalid
	}

	// Delete artwork record from database
	return s.repo.Delete(artworkID)
}

// CheckUploadLimit checks if a user can upload more artworks to an activity
// Requirements: 4.2
func (s *ArtworkService) CheckUploadLimit(userID, activityID uint) (bool, error) {
	// Get activity to check max uploads limit
	activity, err := s.activityService.GetActivityByID(activityID)
	if err != nil {
		return false, err
	}

	// Count user's current uploads for this activity
	count, err := s.repo.CountByUserAndActivity(userID, activityID)
	if err != nil {
		return false, err
	}

	// Check if user has reached the limit
	return count < int64(activity.MaxUploadsPerUser), nil
}
