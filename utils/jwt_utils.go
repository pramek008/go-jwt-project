// utils/token.go
package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pramek008/go-jwt-project/database"
	"github.com/pramek008/go-jwt-project/models"
)

var jwtKey = []byte("p@ssw0rd")

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(userID uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	// Simpan token ke database dan hapus token lama
	err = saveTokenToDB(userID, tokenString, expirationTime)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func saveTokenToDB(userID uuid.UUID, tokenString string, expirationTime time.Time) error {
	db := database.DB.Db

	// Hapus token lama
	var oldTokens []models.Token
	db.Where("user_id = ?", userID).Find(&oldTokens)
	for _, oldToken := range oldTokens {
		db.Delete(&oldToken)
	}

	// Simpan token baru
	newToken := models.Token{
		Token:     tokenString,
		UserID:    userID,
		ExpiredAt: expirationTime,
	}
	return db.Create(&newToken).Error
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
