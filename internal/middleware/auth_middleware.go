package middleware

import (
	"example/hello/internal/auth"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware memverifikasi Bearer Token JWT dari header Authorization.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required", "status": false})
			return
		}

		// Periksa apakah header dimulai dengan "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // Jika TrimPrefix tidak mengubah string, berarti "Bearer " tidak ada
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token format required", "status": false})
			return
		}

		// Validasi token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid or expired token: %s", err.Error()), "status": false})
			return
		}

		// Simpan informasi user dari klaim di konteks Gin
		c.Set("userID", claims.UserID)

		c.Next() // Lanjutkan ke handler berikutnya
	}
}
