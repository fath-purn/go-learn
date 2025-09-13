package handler

import (
	"example/hello/internal/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *UserHandler {
	return &UserHandler{userService: userService}
}
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to retrieve users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "data": convertToUserResponses(users)})
}
func (h *UserHandler) GetUserById(c *gin.Context) {
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
func (h *UserHandler) MyAccount(c *gin.Context) {
	// ambil userID
	userIDVal, exists := c.Get("userID")
	if !exists {
		return
	}

	// Lakukan type assertion dari interface{} ke string.
	userIDStr, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Invalid user ID format in context"})
		return
	}

	// Konversi string ID ke integer.
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Invalid user ID in token"})
		return
	}

	user, err := h.userService.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "data": convertToUserResponse(user)})
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var userRequest user.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": h.getValidationErrors(err)})
		return
	}

	newUser, err := h.userService.RegisterUser(userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Registration failed", "err": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": true, "data": convertToUserResponse(newUser)})
}

func (h *UserHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Token is required"})
		return
	}

	err := h.userService.VerifyEmail(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Email verified successfully"})
}

func (h *UserHandler) ResendVerificationEmail(c *gin.Context) {
	var input user.ResendVerificationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Email is required"})
		return
	}

	err := h.userService.ResendVerificationEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "message": "A new verification email has been sent. Please check your inbox."})
}

func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var input user.ResendVerificationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": h.getValidationErrors(err)})
		return
	}

	err := h.userService.ForgotPassword(input.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "message": "A password reset has been sent. Please check your inbox."})
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var newPassword user.ResetPassword
	if err := c.ShouldBindJSON(&newPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": h.getValidationErrors(err)})
		return
	}

	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Token is required"})
		return
	}

	err := h.userService.ResetPassword(token, newPassword.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Password reset successfully"})
}

func (h *UserHandler) LoginUser(c *gin.Context) {
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

func (h *UserHandler) UpdateUser(c *gin.Context) {
	intID, err := h.validateID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	var userRequest user.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid input data", "err": err.Error()})
		return
	}

	updatedUser, err := h.userService.Update(intID, userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to update user", "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "data": convertToUserResponse(updatedUser)})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	intID, err := h.validateID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	if err := h.userService.Delete(intID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to update user", "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "User deleted successfully"})
}

func (h *UserHandler) validateID(c *gin.Context) (int, error) {
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

func (h *UserHandler) getValidationErrors(err error) []string {
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
