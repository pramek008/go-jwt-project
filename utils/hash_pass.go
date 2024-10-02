package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return "", err
	}
	hashedPassword := string(hashedBytes)
	log.Printf("Original password: %s, Hashed password: %s", password, hashedPassword)
	return hashedPassword, nil
}

func VerifyPassword(hashedPassword, password string) error {
	log.Printf("Verifying - Stored hash: %s, Provided password: %s", hashedPassword, password)
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Printf("Password verification failed: %v", err)
	} else {
		log.Printf("Password verified successfully")
	}
	return err
}
