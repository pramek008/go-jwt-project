package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Title     string         `gorm:"size:255;not null"`
	Content   string         `gorm:"type:text;not null"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null"`
	User      User           `gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Post) TableName() string {
	return "posts"
}
