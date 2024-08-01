package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `gorm:"primaryKey"`
	Title     string         `gorm:"size:255;not null"`
	Content   string         `gorm:"type:text;not null"`
	UserID    uint           `gorm:"not null"`
	User      User           `gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Post) TableName() string {
	return "posts"
}
