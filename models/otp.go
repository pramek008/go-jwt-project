package models

import (
	"time"

	"github.com/google/uuid"
)

type OTP struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Email     string    `gorm:"type:varchar(255);index;not null"`
	Code      string    `gorm:"type:varchar(6);not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (o *OTP) TableName() string {
	return "otps"
}
