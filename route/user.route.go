package route

import (
	"example/hello/handler"
	"example/hello/middleware"
	"example/hello/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(r *gin.Engine, db *gorm.DB) {
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService)

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
