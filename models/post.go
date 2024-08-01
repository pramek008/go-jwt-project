package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;" json:"title"`
	Content   string    `gorm:"type:text;not null;" json:"content"`
	UserID    uint32    `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignkey:UserID" json:"user"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
