package handler

import (
	"example/hello/internal/match"

	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MatchHandler struct {
	matchService match.Service
}

func NewMatchHandler(matchService match.Service) *MatchHandler {
	return &MatchHandler{
		matchService: matchService,
	}
}

func (h *MatchHandler) CreateMatch(c *gin.Context) {
	// 1. Ambil file gambar dari form
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Image file is required",
			"errors":  []string{err.Error()},
		})
		return
	}

	// 2. Buat nama file unik dan simpan
	ext := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + ext
	dst := filepath.Join("assets", "images", newFileName)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to save file"})
		return
	}
	imageURL := "/assets/images/" + newFileName

	// 3. Ambil data lain dari form
	ageStr := c.PostForm("age")
	age, err := strconv.Atoi(ageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid age format"})
		return
	}

	// Ambil userID dari context yang sudah di-set oleh middleware Auth
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User not authenticated"})
		return
	}

	userID, err := strconv.Atoi(userIDVal.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid user ID in token"})
		return
	}

	// 4. Buat request object untuk service
	matchRequest := match.MatchRequest{
		UserID:     userID,
		Age:        age,
		Gender:     match.Gender(c.PostForm("gender")),
		Interested: match.Interest(c.PostForm("interested")),
		City:       c.PostForm("city"),
		Name:       c.PostForm("name"),
		Bio:        c.PostForm("bio"),
		ImageURL:   imageURL,
	}

	// 5. Panggil service untuk membuat match
	newMatch, err := h.matchService.Create(matchRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat profile",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Profile berhasil dibuat",
		"data":    newMatch,
	})
}

func (h *MatchHandler) GetAllMatchUrls(c *gin.Context) {
	matchs, err := h.matchService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve profile",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Profile retrieved successfully",
		"data":    matchs,
	})
}

func (h *MatchHandler) GetMatchByCity(c *gin.Context) {
	city := c.Param("city")

	matchs, err := h.matchService.FindByCity(city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve profile",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Profile retrieved successfully",
		"data":    matchs,
	})
}

func (h *MatchHandler) GetMatchUrlByID(c *gin.Context) {
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

	match, err := h.matchService.FindByID(intID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve profile",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Profile retrieved successfully",
		"data":    match,
	})
}

func (h *MatchHandler) UpdateMatchUrl(c *gin.Context) {
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

	var bookRequest match.MatchRequest
	if err := c.ShouldBindJSON(&bookRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"errors":  []string{err.Error()},
		})
		return
	}

	updated, err := h.matchService.Update(intID, bookRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update profile",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "MatchUrl updated successfully",
		"data":    updated,
	})
}

func (h *MatchHandler) DeleteMatchUrl(c *gin.Context) {
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

	if err := h.matchService.Delete(intID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete profile",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Profile deleted successfully",
	})
}
