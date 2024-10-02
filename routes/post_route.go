package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pramek008/go-jwt-project/controllers"
	"github.com/pramek008/go-jwt-project/middleware"
)

func PostRoute(r *gin.Engine) {
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
