package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Title     string         `gorm:"size:255;not null" json:"title"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	FileURL   string         `gorm:"size:255;" json:"fileUrl"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"userId"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

type PostResponse struct {
	ID        uuid.UUID    `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	FileURL   string       `json:"fileUrl"`
	User      UserResponse `json:"user"` // Use the custom UserResponse struct
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	DeletedAt *time.Time   `json:"deletedAt,omitempty"`
}

func (Post) TableName() string {
	return "posts"
}
