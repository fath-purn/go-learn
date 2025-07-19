package route

import (
	"example/hello/handler"
	"example/hello/short"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ShortRoutes(r *gin.Engine, db *gorm.DB) {
	shortRepository := short.NewRepository(db)
	shortService := short.NewService(shortRepository)
	shortHandler := handler.NewShortUrlHandler(shortService)

	shortGroup := r.Group("/v1")

	// Define a simple GET endpoint
	shortGroup.GET("/:url", shortHandler.GetShortUrl)
	shortGroup.POST("/shorten", shortHandler.CreateShortUrl)
	shortGroup.PUT("/:id", shortHandler.UpdateShortUrl)
	shortGroup.DELETE("/:id", shortHandler.DeleteShortUrl)
	shortGroup.GET("/all", shortHandler.GetAllShortUrls)
	shortGroup.GET("/find/:id", shortHandler.GetShortUrlByID)
}
