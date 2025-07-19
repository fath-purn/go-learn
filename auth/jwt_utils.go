package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Pastikan secret key diambil dari environment variable di produksi!
var jwtSecret = []byte("your_super_secret_jwt_key_here_for_prod") // Ganti dengan secret key yang kuat!

// MyClaims mendefinisikan struktur klaim kustom untuk JWT kita.
type MyClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken membuat JWT baru.
// Ini adalah "sign" token seperti di jsonwebtoken.sign()
func GenerateToken(userID string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token berlaku 24 jam

	claims := &MyClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
			Issuer:    "MyApplication", // Aplikasi yang mengeluarkan token
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // Menggunakan HS256

	tokenString, err := token.SignedString(jwtSecret) // Menandatangani token
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token is expired")
		}
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, fmt.Errorf("invalid token signature")
		}
		return nil, fmt.Errorf("failed to parse or validate token: %w", err)
	}

	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token claims")
	}
}
