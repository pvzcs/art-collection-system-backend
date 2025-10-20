package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Password  string    `gorm:"not null;size:255" json:"-"`
	Nickname  string    `gorm:"not null;size:100" json:"nickname"`
	Role      string    `gorm:"type:enum('user','admin');default:'user';not null" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Artworks []Artwork `gorm:"foreignKey:UserID" json:"artworks,omitempty"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}
