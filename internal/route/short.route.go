package route

import (
	"example/hello/internal/handler"

	"github.com/gin-gonic/gin"
)

func ShortRoutes(r *gin.Engine, shortHandler *handler.ShortUrlHandler) {

	shortGroup := r.Group("/v1")

	// Define a simple GET endpoint
	shortGroup.GET("/:url", shortHandler.GetShortUrl)
	shortGroup.POST("/shorten", shortHandler.CreateShortUrl)
	shortGroup.PUT("/:id", shortHandler.UpdateShortUrl)
	shortGroup.DELETE("/:id", shortHandler.DeleteShortUrl)
	shortGroup.GET("/all", shortHandler.GetAllShortUrls)
	shortGroup.GET("/find/:id", shortHandler.GetShortUrlByID)
}
