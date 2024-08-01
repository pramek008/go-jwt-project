package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pramek008/go-jwt-project/controllers"
	"github.com/pramek008/go-jwt-project/middleware"
	"github.com/pramek008/go-jwt-project/utils"
)

func SetupRoutes(r *gin.Engine) {

	r.GET("/", func(ctx *gin.Context) {
		// ctx.JSON(200, gin.H{
		// 	"message": "Hello from the API!",
		// })
		utils.SendResponse(ctx, 200, true, "Success", "Hello from the API!")
	})

	// Public routes
	public := r.Group("/api")
	{
		public.POST("/register", controllers.Register)
		public.POST("/login", controllers.Login)
	}

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.JWTMiddleware())
	{
		protected.POST("/posts", controllers.CreatePost)
		protected.GET("/posts/:id", controllers.GetPost)
		protected.PUT("/posts/:id", controllers.UpdatePost)
		protected.DELETE("/posts/:id", controllers.DeletePost)
		protected.GET("/posts", controllers.ListPosts)
	}
}
