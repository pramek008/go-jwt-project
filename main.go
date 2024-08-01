package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"github.com/pramek008/go-jwt-project/database"
	"github.com/pramek008/go-jwt-project/routes"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to database
	database.ConnectDb()

	// Set up Gin router
	r := gin.Default()

	// Set up routes
	routes.SetupRoutes(r)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	r.Run(":" + port)
}
