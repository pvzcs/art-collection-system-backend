package models

import (
	"time"
)

// ReviewStatus represents the review status of an artwork
type ReviewStatus string

const (
	StatusPending  ReviewStatus = "pending"
	StatusApproved ReviewStatus = "approved"
)

// Artwork represents an artwork submitted by a user
type Artwork struct {
	ID           uint         `gorm:"primaryKey" json:"id"`
	ActivityID   uint         `gorm:"not null;index:idx_activity_id" json:"activity_id"`
	UserID       uint         `gorm:"not null;index:idx_user_id,priority:1;index:idx_user_activity,priority:1" json:"user_id"`
	FilePath     string       `gorm:"not null;size:500" json:"-"`
	FileName     string       `gorm:"not null;size:255" json:"file_name"`
	ReviewStatus ReviewStatus `gorm:"type:enum('pending','approved');default:'pending';not null;index:idx_review_status" json:"review_status"`
	CreatedAt    time.Time    `gorm:"index" json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`

	Activity Activity `gorm:"foreignKey:ActivityID" json:"activity,omitempty"`
	User     User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for Artwork model
func (Artwork) TableName() string {
	return "artworks"
}
