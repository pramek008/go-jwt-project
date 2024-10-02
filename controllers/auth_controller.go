// controllers/auth.go
package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pramek008/go-jwt-project/database"
	"github.com/pramek008/go-jwt-project/models"
	"github.com/pramek008/go-jwt-project/utils"
)

func InitiateRegistration(c *gin.Context) {
	var userData struct {
		Nickname string `json:"nickname" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&userData); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if userData.Nickname == "" {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Nickname is required")
		return
	} else if userData.Email == "" {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Email is required")
		return
	} else if userData.Password == "" {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Password is required")
		return
	}

	var existingTempUser models.TempUser
	if err := database.DB.Db.Where("email = ? OR nickname = ?", userData.Email, userData.Nickname).First(&existingTempUser).Error; err == nil {
		if existingTempUser.Email == userData.Email {
			utils.SendErrorResponse(c, http.StatusConflict, "User with this email already exists")
			return
		}
		if existingTempUser.Nickname == userData.Nickname {
			utils.SendErrorResponse(c, http.StatusConflict, "User with this nickname already exists")
			return
		}
	}

	var existingUser models.User
	if err := database.DB.Db.Where("email = ? OR nickname = ?", userData.Email, userData.Nickname).First(&existingUser).Error; err == nil {
		if existingUser.Email == userData.Email {
			utils.SendErrorResponse(c, http.StatusConflict, "User with this email already exists")
			return
		}
		if existingUser.Nickname == userData.Nickname {
			utils.SendErrorResponse(c, http.StatusConflict, "User with this nickname already exists")
			return
		}
	}

	otpCode := utils.GenerateOTP()
	if err := utils.SaveOTP(userData.Email, otpCode); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to save OTP")
		return
	}

	//send otp to email
	subject := "Your OTP for Registration"
	body := fmt.Sprintf("<h1>Your OTP is: %s</h1><p>This OTP will expire in 15 minutes.</p>", otpCode)
	if err := utils.SendEmail(userData.Email, subject, body); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to send OTP email")
		return
	}

	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	hashedPassword, err := utils.HashPassword(userData.Password)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}
	log.Printf("Hashed password for %s: %s", userData.Email, string(hashedPassword))

	tempUser := models.TempUser{
		Nickname:  userData.Nickname,
		Email:     userData.Email,
		Password:  string(hashedPassword),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := database.DB.Db.Create(&tempUser).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create temporary user")
		return
	}

	utils.SendResponse(c, http.StatusOK, true, "OTP sent successfully to your email", gin.H{
		"email": userData.Email,
	})
}

func CompleteRegistration(c *gin.Context) {
	var verificationData struct {
		Email string `json:"email" binding:"required,email"`
		OTP   string `json:"otp" binding:"required,min=6,max=6"`
	}

	if err := c.ShouldBindJSON(&verificationData); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if !utils.ValidateOTP(verificationData.Email, verificationData.OTP) {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid OTP or expired")
		return
	}

	var tempUser models.TempUser
	if err := database.DB.Db.Where("email = ?", verificationData.Email).First(&tempUser).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	if time.Now().After(tempUser.ExpiresAt) {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Registration expired, please try again")
		return
	}

	// Create new user with hashed password from tempUser
	user := models.User{
		ID:       uuid.New(),
		Nickname: tempUser.Nickname,
		Email:    tempUser.Email,
		Password: tempUser.Password, // Use the hashed password from tempUser
	}

	log.Printf("Using hashed password for %s: %s", user.Email, user.Password)

	if err := database.DB.Db.Create(&user).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	database.DB.Db.Delete(&tempUser)

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.SendResponse(c, http.StatusOK, true, "Registration completed successfully", gin.H{
		"id":       user.ID,
		"nickname": user.Nickname,
		"email":    user.Email,
		"created":  user.CreatedAt,
		"token":    token,
	})
}

func ResendOTP(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	// Check if the user can resend OTP or needs to wait
	canResend, waitTime, err := utils.CanResendOTP(request.Email)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Server error")
		return
	}

	if !canResend {
		// Inform the user how long they need to wait
		utils.SendErrorResponse(c, http.StatusTooManyRequests, "You must wait "+waitTime.String()+" before resending OTP")
		return
	}

	// Generate and save new OTP
	otp := utils.GenerateOTP()
	if otp == "" {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to generate OTP")
		return
	}

	if err := utils.SaveOTP(request.Email, otp); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to save OTP")
		return
	}

	// Send OTP via email or other methods here (omitted for brevity)
	subject := "OTP Verification"
	body := fmt.Sprintf("<h1>Your OTP is: %s</h1><p>This OTP will expire in 15 minutes.</p>", otp)
	if err := utils.SendEmail(request.Email, subject, body); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to send OTP email")
		return
	}

	utils.SendResponse(c, http.StatusOK, true, "OTP sent successfully", gin.H{
		"email": request.Email,
	})
}

func ForgotPassword(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	var user models.User
	if err := database.DB.Db.Where("email = ?", request.Email).First(&user).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "User not found")
	}

	otp := utils.GenerateOTP()
	if otp == "" {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to generate OTP")
		return
	}

	if err := utils.SaveOTP(request.Email, otp); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to save OTP")
		return
	}

	subject := "OTP for Password Reset"
	body := fmt.Sprintf("<h1>Your OTP is: %s</h1><p>This OTP will expire in 15 minutes.</p>", otp)
	if err := utils.SendEmail(request.Email, subject, body); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to send OTP email")
		return
	}

	utils.SendResponse(c, http.StatusOK, true, "OTP for password reset sent successfully", gin.H{
		"email": request.Email,
	})
}

func ResetPassword(c *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		OTP      string `json:"otp" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	if !utils.ValidateOTP(request.Email, request.OTP) {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid or expired OTP")
	}

	var user models.User
	if err := database.DB.Db.Where("email = ?", request.Email).First(&user).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "User not found")
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user.Password = string(hashedPassword)
	if err := database.DB.Db.Save(&user).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to save user")
		return
	}

	utils.SendResponse(c, http.StatusOK, true, "Password reset successfully", gin.H{
		"email":    request.Email,
		"nickname": user.Nickname,
	})
}

// func Register(c *gin.Context) {
// 	var user models.User
// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 	hashedPassword, err := utils.HashPassword(user.Password)
// 	if err != nil {
// 		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
// 		return
// 	}
// 	user.Password = string(hashedPassword)
// 	// Generate UUID untuk pengguna baru
// 	user.ID = uuid.New()
// 	if err := database.DB.Db.Create(&user).Error; err != nil {
// 		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
// 		return
// 	}
// 	token, err := utils.GenerateToken(user.ID)
// 	if err != nil {
// 		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
// 		return
// 	}
// 	utils.SendResponse(c, http.StatusCreated, true, "User created successfully", gin.H{
// 		"id":        user.ID,
// 		"nickname":  user.Nickname,
// 		"email":     user.Email,
// 		"createdAt": user.CreatedAt,
// 		"token":     token,
// 	})
// }

func Login(c *gin.Context) {
	var user struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var foundUser models.User
	if err := database.DB.Db.Where("email = ?", user.Email).First(&foundUser).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "User not found")
		return
	} else if foundUser.Email != user.Email {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid email")
		return
	} else if foundUser.Password == "" {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid password")
		return
	} else {
		log.Printf("User found: %v", foundUser)
	}

	log.Printf("Stored hashed password for %s: %s", foundUser.Email, foundUser.Password)
	log.Printf("Attempting to compare with provided password")

	if err := utils.VerifyPassword(foundUser.Password, user.Password); err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid password")
		return
	}

	token, err := utils.GenerateToken(foundUser.ID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.SendResponse(c, http.StatusOK, true, "Login successful", gin.H{
		"id":       foundUser.ID,
		"nickname": foundUser.Nickname,
		"email":    foundUser.Email,
		"token":    token,
	})
}

func GetMe(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	userId, err := utils.ExtractUserIDFromToken(tokenString)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid token")
		return
	}

	var user models.User
	if err := database.DB.Db.Where("id = ?", userId).First(&user).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	utils.SendResponse(c, http.StatusOK, true, "User fetched successfully", gin.H{
		"id":       user.ID,
		"nickname": user.Nickname,
		"email":    user.Email,
		"created":  user.CreatedAt,
		"updated":  user.UpdatedAt,
	})
}

func Logout(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	var Token models.Token
	if err := database.DB.Db.Where("token = ?", tokenString).First(&Token).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Token not found")
		return
	}

	database.DB.Db.Delete(&Token)

	utils.SendResponse(c, http.StatusOK, true, "Logout successful", gin.H{})
}
