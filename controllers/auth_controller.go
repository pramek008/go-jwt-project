package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pramek008/go-jwt-project/database"
	"github.com/pramek008/go-jwt-project/models"
	"github.com/pramek008/go-jwt-project/utils"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}
	user.Password = string(hashedPassword)

	// Generate UUID for the new user
	user.ID = uuid.New()

	if err := database.DB.Db.Create(&user).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	utils.SendResponse(c, http.StatusCreated, true, "User created successfully", user)
}

func Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var foundUser models.User
	if err := database.DB.Db.Where("email = ?", user.Email).First(&foundUser).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := utils.GenerateToken(foundUser.ID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.SendResponse(c, http.StatusOK, true, "Login successful", gin.H{"nickname": foundUser.Nickname, "email": foundUser.Email, "token": token})
}
