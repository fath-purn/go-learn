package route

import (
	"example/hello/handler"
	"example/hello/middleware"
	"example/hello/realtime"

	"github.com/gin-gonic/gin"
)

func WebSocketRoutes(r *gin.Engine, hub *realtime.Hub) {
	webSocketHandler := handler.NewWebSocketHandler(hub)
	// Rute Terlindungi (membutuhkan Bearer Token JWT)
	protected := r.Group("/v1/api")
	protected.Use(middleware.AuthMiddleware())

	protected.GET("/ws", webSocketHandler.ServeWs)
}
