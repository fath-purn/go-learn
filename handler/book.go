package handler

import (
	"example/hello/book"

	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type bookHandler struct {
	bookService book.Service
}

func NewBookHandler(bookService book.Service) *bookHandler {
	return &bookHandler{bookService: bookService}
}

func (h *bookHandler) GetBooks(c *gin.Context) {
	books, err := h.bookService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve books",
			"errors":  []string{err.Error()},
		})
		return
	}

	var bookResponse []book.BookResponse
	for _, b := range books {
		bookResponse = append(bookResponse, convertToBookResponse(b))
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Books retrieved successfully",
		"data":    bookResponse,
	})
}

func (h *bookHandler) GetBookById(c *gin.Context) {
	ID := c.Param("id")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID is required",
		})
		return
	}

	// Convert ID from string to int
	intID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID must be a valid integer",
		})
		return
	}

	book, err := h.bookService.FIndByID(intID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve book",
			"errors":  []string{err.Error()},
		})
		return
	}

	var bookResponse = convertToBookResponse(book)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Book retrieved successfully",
		"data":    bookResponse,
	})
}

func (h *bookHandler) CreateBook(c *gin.Context) {
	var bookRequest book.BookRequest

	if err := c.ShouldBindJSON(&bookRequest); err != nil {
		// Inisialisasi slice untuk menampung pesan error
		errorMessages := []string{}

		// Periksa apakah error adalah ValidationError dari validator
		if _, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range err.(validator.ValidationErrors) {
				errorMessage := fmt.Sprintf("Error pada kolom '%s', kondisi: '%s'", fieldErr.Field(), fieldErr.ActualTag())
				errorMessages = append(errorMessages, errorMessage)
			}
		} else {
			errorMessages = append(errorMessages, fmt.Sprintf("Format JSON tidak valid: %s", err.Error()))
		}

		// Kirim semua pesan error dalam satu respons JSON
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Input tidak valid",
			"errors":  errorMessages, // Mengirim slice pesan error
		})
		return
	}

	book, err := h.bookService.Create(bookRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat buku",
			"errors":  []string{err.Error()},
		})
		return
	}

	// Jika tidak ada error, lanjutkan proses dan kirim respons sukses
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Buku berhasil dibuat",
		"data": gin.H{
			"title":       book.Title,
			"synopsis":    book.Synopsis,
			"price":       book.Price,
			"description": book.Description,
			"rating":      book.Rating,
		},
	})
}

func (h *bookHandler) UpdateBook(c *gin.Context) {
	ID := c.Param("id")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID is required",
		})
		return
	}

	intID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID must be a valid integer",
		})
		return
	}

	var bookRequest book.BookRequest
	if err := c.ShouldBindJSON(&bookRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"errors":  []string{err.Error()},
		})
		return
	}

	book, err := h.bookService.Update(intID, bookRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update book",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Book updated successfully",
		"data":    convertToBookResponse(book),
	})
}

func (h *bookHandler) DeleteBook(c *gin.Context) {
	ID := c.Param("id")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID is required",
		})
		return
	}

	intID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID must be a valid integer",
		})
		return
	}

	if err := h.bookService.Delete(intID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete book",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Book deleted successfully",
	})
}

func convertToBookResponse(b book.Book) book.BookResponse {
	return book.BookResponse{
		ID:          b.ID,
		Title:       b.Title,
		Price:       b.Price,
		Synopsis:    b.Synopsis,
		Description: b.Description,
		Rating:      b.Rating,
	}
}
