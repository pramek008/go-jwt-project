package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func SetBaseURL() gin.HandlerFunc {
	return func(c *gin.Context) {
		var baseURL string

		// Try to get base URL from X-Base-URL header first
		baseURL = c.GetHeader("X-Base-URL")

		// If header is not set, construct base URL from request
		if baseURL == "" {
			scheme := "http"
			if c.Request.TLS != nil {
				scheme = "https"
			}
			baseURL = fmt.Sprintf("%s://%s", scheme, c.Request.Host)
		}

		// Ensure baseURL doesn't end with a slash
		baseURL = strings.TrimRight(baseURL, "/")

		// Set the base URL in the context
		c.Set("BaseURL", baseURL)
		c.Next()
	}
}
