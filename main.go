package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pramek008/go-jwt-project/database"
	"github.com/pramek008/go-jwt-project/middleware"
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

	r.Use(middleware.SetBaseURL())

	r.Static("/uploads", "./uploads")

	// Set up routes
	routes.SetupRoutes(r)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "1500"
	}
	r.Run("0.0.0.0:" + port)
}
