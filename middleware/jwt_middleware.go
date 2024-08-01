package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pramek008/go-jwt-project/utils"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			// c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Authorization header format must be Bearer {token}")
			c.Abort()
			return
		}

		token, err := utils.ValidateToken(parts[1])
		if err != nil {
			// c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			// c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid token claims")
			c.Abort()
			return
		}

		c.Set("user_id", uint32(claims["user_id"].(float64)))
		c.Next()
	}
}
