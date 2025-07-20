package route

import (
	"example/hello/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	userHandler *handler.UserHandler,
	bookHandler *handler.BookHandler,
	shortHandler *handler.ShortUrlHandler,
	webSocketHandler *handler.WebSocketHandler,
) {
	UserRoutes(r, userHandler)
	BookRoutes(r, bookHandler)
	ShortRoutes(r, shortHandler)
	WebSocketRoutes(r, webSocketHandler)
}
