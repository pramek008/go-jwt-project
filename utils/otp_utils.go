// utils/otp_utils.go
package utils

import (
	"crypto/rand"
	"time"

	"github.com/pramek008/go-jwt-project/database"
	"github.com/pramek008/go-jwt-project/models"
)

const otpExpiryDuration = 15 * time.Minute
const otpResendCooldown = 5 * time.Minute

func GenerateOTP() string {
	const otpChars = "1234567890"
	buffer := make([]byte, 6)
	_, err := rand.Read(buffer)
	if err != nil {
		return ""
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < 6; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer)
}

func SaveOTP(email, otp string) error {
	otpRecord := models.OTP{
		Email:     email,
		Code:      otp,
		ExpiresAt: time.Now().Add(otpExpiryDuration),
	}

	if err := database.DB.Db.Create(&otpRecord).Error; err != nil {
		return err
	}

	return nil
}

func CanResendOTP(email string) (bool, time.Duration, error) {
	var otpRecord models.OTP
	err := database.DB.Db.Where("email = ?", email).Order("created_at desc").First(&otpRecord).Error
	if err != nil {
		return true, 0, nil // No previous OTP, can resend
	}

	timeSinceLastOTP := time.Since(otpRecord.CreatedAt)
	if timeSinceLastOTP < otpResendCooldown {
		waitTime := otpResendCooldown - timeSinceLastOTP
		return false, waitTime, nil // Must wait to resend
	}

	return true, 0, nil // Allowed to resend
}

func ValidateOTP(email, otp string) bool {
	var otpRecord models.OTP
	if err := database.DB.Db.Where("email = ? AND code = ? AND expires_at > ?", email, otp, time.Now()).First(&otpRecord).Error; err != nil {
		return false
	}

	// Delete the OTP record after successful validation
	database.DB.Db.Delete(&otpRecord)

	return true
}
