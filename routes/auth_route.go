package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pramek008/go-jwt-project/controllers"
	"github.com/pramek008/go-jwt-project/middleware"
)

func AuthRoute(r *gin.Engine) {
	auth := r.Group("/api/auth")
	{
		// auth.POST("/register", controllers.Register)
		auth.POST("/register-initiate", controllers.InitiateRegistration)
		auth.POST("/register-complete", controllers.CompleteRegistration)
		auth.POST("/resend-otp", controllers.ResendOTP)
		auth.POST("/forgot-password", controllers.ForgotPassword)
		auth.POST("/reset-password", controllers.ResetPassword)
		auth.POST("/login", controllers.Login)
	}
	protected := r.Group("/api/auth")
	protected.Use(middleware.JWTMiddleware())
	{
		protected.POST("/logout", controllers.Logout)
		protected.GET("/me", controllers.GetMe)
	}

}
