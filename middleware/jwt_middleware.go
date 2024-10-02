// middleware/jwt_middleware.go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pramek008/go-jwt-project/database"
	"github.com/pramek008/go-jwt-project/models"
	"github.com/pramek008/go-jwt-project/utils"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Authorization header format must be Bearer {token}")
			c.Abort()
			return
		}

		userID, err := utils.ExtractUserIDFromToken(tokenString)
		if err != nil {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		// Periksa apakah token ada di database
		var storedToken models.Token
		db := database.DB.Db
		if err := db.Where("token = ? AND user_id = ?", tokenString, userID).First(&storedToken).Error; err != nil {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid Token not found")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
