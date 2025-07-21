package route

import (
	"example/hello/internal/handler"
	"example/hello/internal/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, userHandler *handler.UserHandler) {
	// Create a new group for user routes
	userGroup := r.Group("/v1")

	userGroup.POST("/register", userHandler.RegisterUser)
	userGroup.POST("/login", userHandler.LoginUser)
	userGroup.GET("/verify-email", userHandler.VerifyEmail)
	userGroup.POST("/resend-verification", userHandler.ResendVerificationEmail)
	userGroup.POST("/forgot-password", userHandler.ForgotPassword)
	userGroup.POST("/reset-password", userHandler.ResetPassword)

	// Rute Terlindungi (membutuhkan Bearer Token JWT)
	protected := r.Group("/v1/user")
	protected.Use(middleware.AuthMiddleware())

	protected.GET("/", userHandler.GetUsers)
	protected.GET("/me", userHandler.MyAccount)
	protected.GET("/:id", userHandler.GetUserById)
	protected.PUT("/:id", userHandler.UpdateUser)
	protected.DELETE("/:id", userHandler.DeleteUser)
}
