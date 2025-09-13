package route

import (
	"example/hello/internal/handler"
	"example/hello/internal/middleware"

	"github.com/gin-gonic/gin"
)

func WebSocketRoutes(r *gin.Engine, webSocketHandler *handler.WebSocketHandler) {
	// Rute Terlindungi (membutuhkan Bearer Token JWT)
	protected := r.Group("/v1")
	protected.Use(middleware.AuthMiddleware())

	protected.GET("/ws", webSocketHandler.ServeWs)
}
