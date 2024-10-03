package database

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pramek008/go-jwt-project/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func ConnectDb() {
	// First, connect to the default 'postgres' database to create our app's database
	defaultDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
	)

	defaultDB, err := gorm.Open(postgres.Open(defaultDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to default database: ", err)
		os.Exit(2)
	}

	// Create the application database if it doesn't exist
	dbName := os.Getenv("DB_NAME")
	err = defaultDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error
	if err != nil {
		// If the database already exists, this is not a fatal error
		if !strings.Contains(err.Error(), "already exists") {
			log.Fatal("Failed to create database: ", err)
			os.Exit(2)
		}
	}

	// Close the connection to the default database
	sqlDB, err := defaultDB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance: ", err)
		os.Exit(2)
	}
	sqlDB.Close()

	// Now connect to the application database
	appDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		dbName,
		os.Getenv("DB_PORT"),
	)

	var db *gorm.DB
	retries := 5
	for i := 0; i < retries; i++ {
		db, err = gorm.Open(postgres.Open(appDSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to application database. Retrying in 5 seconds... (%d/%d)\n", i+1, retries)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to application database after multiple retries: ", err)
		os.Exit(2)
	}

	log.Println("Connected to application database")
	db.Logger = logger.Default.LogMode(logger.Info)

	// Ensure UUID extension is created
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Fatal("Failed to create UUID extension: ", err)
		os.Exit(2)
	}

	log.Println("Running Migrations")
	err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.Token{}, &models.TempUser{}, &models.OTP{})
	if err != nil {
		log.Fatal("Failed to auto migrate: ", err)
		os.Exit(2)
	}
	log.Println("Migrations completed")

	// Run seeder
	if err := seedData(db); err != nil {
		log.Fatalf("Failed to seed data. Error: %v\n", err)
		os.Exit(2)
	}

	DB = Dbinstance{
		Db: db,
	}
}

func seedData(db *gorm.DB) error {
	// Check if data already exists
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		log.Println("Data already seeded")
		return nil
	}

	// Create sample users with UUIDs
	users := []models.User{
		{ID: uuid.New(), Nickname: "user1", Email: "user1@example.com", Password: HashPassword("password1")},
		{ID: uuid.New(), Nickname: "user2", Email: "user2@example.com", Password: HashPassword("password2")},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			return err
		}
		log.Printf("Created user: %v", user)
	}

	// Create sample posts
	posts := []models.Post{
		{Title: "First Post", Content: "This is the content of the first post", UserID: users[0].ID},
		{Title: "Second Post", Content: "This is the content of the second post", UserID: users[1].ID},
		{Title: "Another Post", Content: "This is another post by user 1", UserID: users[0].ID},
	}

	for _, post := range posts {
		if err := db.Create(&post).Error; err != nil {
			log.Printf("Failed to create post: %v, error: %v", post, err)
			return err
		}
		log.Printf("Created post: %v", post)
	}

	log.Println("Data seeded successfully")
	return nil
}

func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password. Error: %v\n", err)
	}
	return string(hashedPassword)
}
