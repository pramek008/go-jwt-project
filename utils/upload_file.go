package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadFile(c *gin.Context, file *multipart.FileHeader) (string, error) {
	// Generate a unique filename
	filename := uuid.New().String() + filepath.Ext(file.Filename)

	// Ensure the upload directory exists
	uploadDir := "/usr/src/app/uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Save the file
	dst := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Get the base URL from environment variable
	baseURL, _ := c.Get("BaseURL")

	// Return the full URL
	return fmt.Sprintf("%s/uploads/%s", baseURL, filename), nil
}
