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

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Authorization header format must be Bearer {token}")
			c.Abort()
			return
		}

		token, err := utils.ValidateToken(parts[1])
		if err != nil {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*utils.Claims); ok && token.Valid {
			// Periksa apakah token ada di database
			var storedToken models.Token
			db := database.DB.Db
			if err := db.Where("token = ?", parts[1]).First(&storedToken).Error; err != nil {
				utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid Token not found")
				c.Abort()
				return
			}

			c.Set("user_id", claims.UserID)
		} else {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid token claims")
			c.Abort()
			return
		}

		c.Next()
	}
}
