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

	// Rute Terlindungi (membutuhkan Bearer Token JWT)
	protected := r.Group("/v1/api")
	protected.Use(middleware.AuthMiddleware())

	protected.GET("/user", userHandler.GetUsers)
	protected.GET("/user/:id", userHandler.GetUserById)
	protected.PUT("/user/:id", userHandler.UpdateUser)
	protected.DELETE("/user/:id", userHandler.DeleteUser)

}
