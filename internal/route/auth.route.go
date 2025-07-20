package route

import (
	"example/hello/internal/handler"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine, authHandler *handler.AuthHandler) {
	authGroup := r.Group("/v1/auth")

	authGroup.GET("/google/login", authHandler.GoogleLogin)
	authGroup.GET("/google/callback", authHandler.GoogleCallback)
}
