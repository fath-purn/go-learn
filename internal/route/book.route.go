package route

import (
	"example/hello/internal/handler"

	"github.com/gin-gonic/gin"
)

func BookRoutes(r *gin.Engine, bookHandler *handler.BookHandler) {
	// Create a new group for book routes
	bookGroup := r.Group("/v1")

	// Define a simple GET endpoint
	bookGroup.GET("/get-books", bookHandler.GetBooks)
	bookGroup.GET("/get-book/:id", bookHandler.GetBookById)
	bookGroup.PUT("/book/:id", bookHandler.UpdateBook)
	bookGroup.DELETE("/book/:id", bookHandler.DeleteBook)
	bookGroup.POST("/book", bookHandler.CreateBook)
}
