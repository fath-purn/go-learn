package auth

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// jwtSecret adalah kunci rahasia untuk menandatangani token.
// Variabel ini diinisialisasi di dalam fungsi init().
var jwtSecret []byte

// init dieksekusi secara otomatis saat paket 'auth' diimpor.
// Fungsi ini bertanggung jawab untuk memuat konfigurasi yang diperlukan.
func init() {
	// Memuat variabel dari file .env, berguna untuk pengembangan lokal.
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: file .env tidak ditemukan, akan membaca environment variables dari sistem.")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("FATAL: Environment variable JWT_SECRET tidak di-set.")
	}
	jwtSecret = []byte(secret)
}

// MyClaims mendefinisikan struktur klaim kustom untuk JWT kita.
type MyClaims struct {
	UserID   string `json:"user_id"`
	Verified bool   `json:"verified"`
	jwt.RegisteredClaims
}

// GenerateToken membuat JWT baru.
// Ini adalah "sign" token seperti di jsonwebtoken.sign()
func GenerateToken(userID string, verified bool) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token berlaku 24 jam

	claims := &MyClaims{
		UserID:   userID,
		Verified: verified,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
			Issuer:    "MyApplication",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // Menggunakan HS256

	tokenString, err := token.SignedString(jwtSecret) // Menandatangani token
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken memverifikasi tanda tangan token dan mengurai klaim.
func ValidateToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Pastikan algoritma penandatanganan adalah yang kita harapkan (HMAC).
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	// Kita bisa langsung mengembalikannya untuk memberikan konteks yang lebih baik.
	if err != nil {
		return nil, fmt.Errorf("validasi token gagal: %w", err)
	}

	claims, ok := token.Claims.(*MyClaims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token tidak valid")
}
