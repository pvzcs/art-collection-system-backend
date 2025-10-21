package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// ArtworkServiceInterface defines the interface for artwork service operations
// This allows FileService to work with ArtworkService without circular dependency
type ArtworkServiceInterface interface {
	GetArtwork(artworkID, requesterID uint, requesterRole string) (interface{}, error)
}

// FileService handles file storage and access operations
type FileService struct {
	uploadPath string
}

// NewFileService creates a new file service instance
func NewFileService(uploadPath string) *FileService {
	return &FileService{
		uploadPath: uploadPath,
	}
}

// SaveFile saves an uploaded file to the server with a unique filename
// File path structure: uploads/{year}/{month}/{uuid}_{original_filename}
// Requirements: 11.1, 11.2
func (s *FileService) SaveFile(file multipart.File, filename string) (string, error) {
	// Generate unique filename using UUID + original filename
	uniqueFilename := fmt.Sprintf("%s_%s", uuid.New().String(), filename)

	// Create directory structure based on current year and month
	now := time.Now()
	year := fmt.Sprintf("%d", now.Year())
	month := fmt.Sprintf("%02d", now.Month())
	
	dirPath := filepath.Join(s.uploadPath, year, month)
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Full file path
	filePath := filepath.Join(dirPath, uniqueFilename)

	// Create the file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Return relative path from upload root
	relativePath := filepath.Join(year, month, uniqueFilename)
	return relativePath, nil
}

// ServeFile reads and returns file content after validating permissions
// Requirements: 11.3, 11.4, 6.4
func (s *FileService) ServeFile(filePath string, artworkID, requesterID uint, requesterRole string, artworkService ArtworkServiceInterface) ([]byte, string, error) {
	// Validate permissions by calling ArtworkService.GetArtwork
	_, err := artworkService.GetArtwork(artworkID, requesterID, requesterRole)
	if err != nil {
		return nil, "", fmt.Errorf("permission denied: %w", err)
	}

	// Construct full file path
	fullPath := filePath
	if !filepath.IsAbs(fullPath) {
		fullPath = filepath.Join(s.uploadPath, filePath)
	}

	// Read file content
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read file: %w", err)
	}

	// Determine content type based on file extension
	contentType := getContentType(fullPath)

	return data, contentType, nil
}

// DeleteFile deletes a physical file from the server
// Requirements: 4.5
func (s *FileService) DeleteFile(filePath string) error {
	// Construct full file path
	fullPath := filePath
	if !filepath.IsAbs(fullPath) {
		fullPath = filepath.Join(s.uploadPath, filePath)
	}

	// Delete the file
	if err := os.Remove(fullPath); err != nil {
		// If file doesn't exist, consider it already deleted (not an error)
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// getContentType determines the MIME type based on file extension
func getContentType(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".bmp":
		return "image/bmp"
	case ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}
