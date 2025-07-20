package handler

import (
	"example/hello/internal/short"
	"strconv"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ShortUrlHandler struct {
	shortService short.Service
}

func NewShortUrlHandler(shortService short.Service) *ShortUrlHandler {
	return &ShortUrlHandler{shortService: shortService}
}

func (h *ShortUrlHandler) GetShortUrl(c *gin.Context) {
	url := c.Param("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID is required",
		})
		return
	}

	short, err := h.shortService.FindByUrl(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve short",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.Redirect(http.StatusFound, short.Original)
}

func (h *ShortUrlHandler) CreateShortUrl(c *gin.Context) {
	var shortRequest short.ShortRequest

	if err := c.ShouldBindJSON(&shortRequest); err != nil {
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

	short, err := h.shortService.Create(shortRequest)
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
			"original":  short.Original,
			"shortened": short.Shortened,
		},
	})
}

func (h *ShortUrlHandler) GetAllShortUrls(c *gin.Context) {
	shorts, err := h.shortService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve shorts",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "ShortUrls retrieved successfully",
		"data":    shorts,
	})
}

func (h *ShortUrlHandler) GetShortUrlByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID is required",
		})
		return
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid ID format",
			"errors":  []string{err.Error()},
		})
		return
	}

	short, err := h.shortService.FindByID(intID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve short",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "ShortUrl retrieved successfully",
		"data":    short,
	})
}

func (h *ShortUrlHandler) UpdateShortUrl(c *gin.Context) {
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

	var bookRequest short.ShortRequest
	if err := c.ShouldBindJSON(&bookRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"errors":  []string{err.Error()},
		})
		return
	}

	updated, err := h.shortService.Update(intID, bookRequest)
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
		"message": "ShortUrl updated successfully",
		"data":    updated,
	})
}

func (h *ShortUrlHandler) DeleteShortUrl(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID is required",
		})
		return
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid ID format",
			"errors":  []string{err.Error()},
		})
		return
	}

	if err := h.shortService.Delete(intID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete short",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "ShortUrl deleted successfully",
	})
}
