package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"example/hello/internal/auth"
	"example/hello/internal/user"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	oauth2api "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type AuthHandler struct {
	googleOauthConfig *oauth2.Config
	userService       user.Service
}

func NewAuthHandler(googleOauthConfig *oauth2.Config, userService user.Service) *AuthHandler {
	return &AuthHandler{
		googleOauthConfig: googleOauthConfig,
		userService:       userService,
	}
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	// Generate a random state string to prevent CSRF attacks.
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}
	state := base64.URLEncoding.EncodeToString(b)

	// Store the state in a short-lived cookie.
	c.SetCookie("oauthstate", state, 3600, "/", "localhost", false, true)

	// Redirect user to Google's consent page.
	url := h.googleOauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	// Check if the state from the cookie matches the state from the URL.
	oauthState, _ := c.Cookie("oauthstate")
	if c.Query("state") != oauthState {
		log.Println("Invalid oauth google state")
		c.Redirect(http.StatusTemporaryRedirect, "/") // Redirect to home or login page
		return
	}

	// Exchange the authorization code for a token.
	code := c.Query("code")
	token, err := h.googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Code exchange failed: %s\n", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	// Use the token to get user info from Google.
	oauth2Service, err := oauth2api.NewService(context.Background(), option.WithTokenSource(h.googleOauthConfig.TokenSource(context.Background(), token)))
	if err != nil {
		log.Printf("Failed to create oauth2 service: %s\n", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		log.Printf("Failed to get user info: %s\n", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	// Find or create a user in our database.
	userInput := user.GoogleLoginInput{
		Email: userInfo.Email,
		Name:  userInfo.Name,
	}

	loggedInUser, err := h.userService.FindOrCreateByGoogle(userInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user"})
		return
	}

	// Generate our own JWT for the user.
	jwtToken, err := auth.GenerateToken(fmt.Sprintf("%d", loggedInUser.ID), loggedInUser.Verivied)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return the JWT to the client.
	c.JSON(http.StatusOK, gin.H{"token": jwtToken})
}
