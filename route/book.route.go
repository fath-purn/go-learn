package route

import (
	"example/hello/book"
	"example/hello/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func BookRoutes(r *gin.Engine, db *gorm.DB) {
	bookRepository := book.NewRepository(db)
	bookService := book.NewService(bookRepository)
	bookHandler := handler.NewBookHandler(bookService)

	// Create a new group for book routes
	bookGroup := r.Group("/v1")

	// Define a simple GET endpoint
	bookGroup.GET("/get-books", bookHandler.GetBooks)
	bookGroup.GET("/get-book/:id", bookHandler.GetBookById)
	bookGroup.PUT("/book/:id", bookHandler.UpdateBook)
	bookGroup.DELETE("/book/:id", bookHandler.DeleteBook)
	bookGroup.POST("/book", bookHandler.CreateBook)

}
