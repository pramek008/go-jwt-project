// middleware/rate_limiter.go
package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	// limiter = rate.NewLimiter(rate.Every(time.Second), 5) // 5 requests per second
	ipMap = make(map[string]*rate.Limiter)
	mu    sync.Mutex
)

func getVisitorLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	l, exists := ipMap[ip]
	if !exists {
		l = rate.NewLimiter(rate.Every(time.Second), 5) // 5 requests per second per IP
		ipMap[ip] = l
	}

	return l
}

func RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getVisitorLimiter(ip)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}
