package models

import (
	"time"
)

// Activity represents an art collection activity
type Activity struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	Name              string     `gorm:"not null;size:255" json:"name"`
	Deadline          *time.Time `json:"deadline"`
	Description       string     `gorm:"type:text" json:"description"`
	MaxUploadsPerUser int        `gorm:"default:5;not null" json:"max_uploads_per_user"`
	IsDeleted         bool       `gorm:"default:false;not null;index" json:"-"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	Artworks []Artwork `gorm:"foreignKey:ActivityID" json:"artworks,omitempty"`
}

// TableName specifies the table name for Activity model
func (Activity) TableName() string {
	return "activities"
}
