package route

import (
	"example/hello/internal/handler"

	"github.com/gin-gonic/gin"
)

func MatchRoutes(r *gin.Engine, matchHandler *handler.MatchHandler) {

	matchGroup := r.Group("/v1/match")

	// Define a simple GET endpoint
	matchGroup.GET("/:city", matchHandler.GetMatchByCity)
	matchGroup.POST("/", matchHandler.CreateMatch)
	matchGroup.PUT("/:id", matchHandler.UpdateMatchUrl)
	matchGroup.DELETE("/:id", matchHandler.DeleteMatchUrl)
	matchGroup.GET("/all", matchHandler.GetAllMatchUrls)
	matchGroup.GET("/find/:id", matchHandler.GetMatchUrlByID)
}
