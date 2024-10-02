package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pramek008/go-jwt-project/middleware"
	"github.com/pramek008/go-jwt-project/utils"
)

func SetupRoutes(r *gin.Engine) {
	// Public routes
	r.Use(middleware.RateLimiterMiddleware())
	public := r.Group("/api")
	{
		public.GET("/", func(ctx *gin.Context) {
			utils.SendResponse[map[string]interface{}](ctx, 200, true, "Hello from the API!", nil)
		})
	}
	AuthRoute(r)
	PostRoute(r)
}
