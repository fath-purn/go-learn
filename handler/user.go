package handler

import (
	"example/hello/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService: userService}
}
func (h *userHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to retrieve users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "data": convertToUserResponses(users)})
}
func (h *userHandler) GetUserById(c *gin.Context) {
	intID, err := h.validateID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}
	user, err := h.userService.FindByID(intID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "data": convertToUserResponse(user)})
}
func (h *userHandler) RegisterUser(c *gin.Context) {
	var userRequest user.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": h.getValidationErrors(err)})
		return
	}
	newUser, err := h.userService.RegisterUser(userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Registration failed"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": true, "data": convertToUserResponse(newUser)})
}
func (h *userHandler) LoginUser(c *gin.Context) {
	var userRequest user.UserLogin
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid input data"})
		return
	}
	token, _, err := h.userService.UserLogin(userRequest)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Login failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "token": token})
}
func (h *userHandler) UpdateUser(c *gin.Context) {
	intID, err := h.validateID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}
	var userRequest user.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid input data"})
		return
	}
	updatedUser, err := h.userService.Update(intID, userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to update user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "data": convertToUserResponse(updatedUser)})
}
func (h *userHandler) DeleteUser(c *gin.Context) {
	intID, err := h.validateID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}
	if err := h.userService.Delete(intID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "User deleted successfully"})
}
func (h *userHandler) validateID(c *gin.Context) (int, error) {
	ID := c.Param("id")
	if ID == "" {
		return 0, fmt.Errorf("ID is required")
	}
	intID, err := strconv.Atoi(ID)
	if err != nil {
		return 0, fmt.Errorf("ID must be a valid integer")
	}
	return intID, nil
}
func (h *userHandler) getValidationErrors(err error) []string {
	var errorMessages []string
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			errorMessages = append(errorMessages, fmt.Sprintf("Error on field '%s', condition: '%s'", fieldErr.Field(), fieldErr.ActualTag()))
		}
	} else {
		errorMessages = append(errorMessages, "Invalid JSON format")
	}
	return errorMessages
}
func convertToUserResponse(b user.User) user.UserResponse {
	return user.UserResponse{
		ID:    b.ID,
		Name:  b.Name,
		Email: b.Email,
		Phone: b.Phone,
	}
}
func convertToUserResponses(users []user.User) []user.UserResponse {
	var userResponses []user.UserResponse
	for _, b := range users {
		userResponses = append(userResponses, convertToUserResponse(b))
	}
	return userResponses
}
