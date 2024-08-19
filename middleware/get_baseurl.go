package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func SetBaseURL() gin.HandlerFunc {
	return func(c *gin.Context) {
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		baseURL := fmt.Sprintf("%s://%s", scheme, c.Request.Host)
		c.Set("BaseURL", baseURL)
		c.Next()
	}
}
