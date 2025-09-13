package user

import (
	"example/hello/internal/auth"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(user UserRequest) (User, error)
	UserLogin(req UserLogin) (string, User, error)
	FindAll() ([]User, error)
	FindByID(ID int) (User, error)
	Update(ID int, user UserRequest) (User, error)
	Delete(ID int) error
	FindOrCreateByGoogle(input GoogleLoginInput) (User, error)
	VerifyEmail(token string) error
	ResendVerificationEmail(email string) error
	ForgotPassword(email string) error
	ResetPassword(token string, newPassword string) error
}

type service struct {
	repository Repository
	cache      *cache.Cache
}

// Cache keys
const (
	allUsersCacheKey       = "all_users"
	userByIDCacheKeyPrefix = "user_by_id_"
)

func NewService(repository Repository) *service {
	// Inisialisasi cache:
	// - 5 menit (5*time.Minute) untuk default expiration
	// - 10 menit (10*time.Minute) untuk cleanup interval (seberapa sering item kadaluarsa dihapus)
	c := cache.New(5*time.Minute, 10*time.Minute)
	return &service{
		repository: repository,
		cache:      c,
	}
}

func (s *service) FindOrCreateByGoogle(input GoogleLoginInput) (User, error) {
	// Check if user with this email already exists.
	user, err := s.repository.FindByEmail(input.Email)
	if err != nil {
		// If user not found, create a new one.
		// It's better to check for gorm.ErrRecordNotFound specifically.
		if err.Error() == "record not found" {
			newUser := User{
				Name:     input.Name,
				Email:    input.Email,
				Verivied: true, // Pengguna dari Google dianggap terverifikasi secara default
				// Password can be left empty or set to a random string for Google users.
				// This prevents them from logging in with a password.
			}
			createdUser, err := s.repository.RegisterUser(newUser)
			if err != nil {
				return User{}, fmt.Errorf("failed to create Google user: %w", err)
			}
			s.cache.Delete(allUsersCacheKey)
			return createdUser, nil
		}
		// Handle other potential errors from the repository.
		return User{}, err
	}
	// If user exists, return them.
	return user, nil
}

// FindByID implements Service.
func (s *service) FindByID(ID int) (User, error) {
	cacheKey := fmt.Sprintf("%s%d", userByIDCacheKeyPrefix, ID)

	// 1. Coba ambil dari cache
	if x, found := s.cache.Get(cacheKey); found {
		// Ingat untuk membersihkan password hash jika Anda menyimpannya dalam cache
		userFromCache := x.(User)
		userFromCache.Password = ""
		return userFromCache, nil
	}

	user, err := s.repository.FindByID(ID)
	if err != nil {
		return User{}, fmt.Errorf("error finding user by ID: %w", err)
	}

	s.cache.Set(cacheKey, user, cache.DefaultExpiration) // Pastikan *user karena FindByID repository mengembalikan pointer

	user.Password = "" // Clear password before returning
	return user, nil
}

func (s *service) RegisterUser(userRequest UserRequest) (User, error) {
	// Validate email format
	if !isValidEmail(userRequest.Email) {
		return User{}, fmt.Errorf("invalid email format")
	}

	// Check for existing user in a single query
	existingUser, err := s.repository.FindByEmail(userRequest.Email)
	if err == nil && existingUser.Email != "" {
		return User{}, fmt.Errorf("email already registered")
	}

	// Define a constant for bcrypt cost
	const bcryptCost = 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcryptCost)
	if err != nil {
		return User{}, fmt.Errorf("error hashing password: %w", err)
	}
	user := User{
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: string(hashedPassword),
		Phone:    userRequest.Phone,
		Verivied: false,
	}

	// token verifikasi
	token := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)
	user.VerificationToken = &token
	user.VerificationTokenExpiresAt = &expiresAt

	createdUser, err := s.repository.RegisterUser(user)
	if err != nil {
		return User{}, fmt.Errorf("error registering user: %w", err)
	}

	go sendVerificationEmail(createdUser)

	s.cache.Delete(allUsersCacheKey)

	return createdUser, nil
}

func (s *service) UserLogin(req UserLogin) (string, User, error) {
	user := User{
		Email:    req.Email,
		Password: req.Password,
	}

	foundUser, err := s.repository.LoginUser(user)
	if err != nil {
		return "", User{}, fmt.Errorf("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
		return "", User{}, fmt.Errorf("invalid email or password")
	}

	token, err := auth.GenerateToken(fmt.Sprintf("%d", foundUser.ID), foundUser.Verivied)
	if err != nil {
		return "", User{}, fmt.Errorf("failed to generate authentication token")
	}

	foundUser.Password = "" // Clear password before returning
	return token, foundUser, nil
}

func (s *service) FindAll() ([]User, error) {
	if x, found := s.cache.Get(allUsersCacheKey); found {
		return x.([]User), nil // Type assert ke []User
	}

	// Fetch all users from the repository
	users, err := s.repository.FindAll()
	if err != nil {
		return nil, fmt.Errorf("error finding all users: %w", err)
	}

	// Store the result in cache
	s.cache.Set(allUsersCacheKey, users, cache.DefaultExpiration)
	// Clear passwords before returning
	for i := range users {
		users[i].Password = ""
	}
	return users, nil
}

func (s *service) Update(ID int, userRequest UserRequest) (User, error) {
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return User{}, fmt.Errorf("error finding user for update: %w", err)
	}
	user.Name = userRequest.Name
	user.Email = userRequest.Email
	// Hash the password if it's being updated
	if userRequest.Password != "" {
		const bcryptCost = 10
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcryptCost)
		if err != nil {
			return User{}, fmt.Errorf("error hashing password: %w", err)
		}
		user.Password = string(hashedPassword)
	}
	user.Phone = userRequest.Phone
	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return User{}, fmt.Errorf("error updating user: %w", err)
	}

	// Setelah operasi tulis, invalidate cache yang relevan
	s.cache.Delete(allUsersCacheKey) // Cache daftar semua user tidak valid lagi
	s.cache.Delete(fmt.Sprintf("%s%d", userByIDCacheKeyPrefix, ID))

	return updatedUser, nil
}

func (s *service) Delete(ID int) error {
	if _, err := s.repository.FindByID(ID); err != nil {
		return fmt.Errorf("error finding user for deletion: %w", err)
	}

	if err := s.repository.Delete(ID); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	// Setelah operasi tulis, invalidate cache yang relevan
	s.cache.Delete(allUsersCacheKey)
	s.cache.Delete(fmt.Sprintf("%s%d", userByIDCacheKeyPrefix, ID))
	return nil
}

// Helper function to validate email format
func isValidEmail(email string) bool {
	// Improved email validation regex
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func (s *service) VerifyEmail(token string) error {
	// cari user berdasarkan token
	user, err := s.repository.FindByVerificationToken(token)
	if err != nil {
		return fmt.Errorf("token verifikasi tidak valid atau tidak bisa digunakan")
	}

	// Periksa apakah token sudah kadaluarsa
	if user.VerificationTokenExpiresAt.Before(time.Now()) {
		return fmt.Errorf("token verifikasi sudah kadaluarsa")
	}

	// Set verivied menjadi true
	user.Verivied = true
	user.VerificationToken = nil
	user.VerificationTokenExpiresAt = nil

	// Simpan perubahan
	_, err = s.repository.Update(user)
	if err != nil {
		return fmt.Errorf("gagal memperbarui status verifikasi: %w", err)
	}

	s.cache.Delete(allUsersCacheKey)
	s.cache.Delete(fmt.Sprintf("%s%d", userByIDCacheKeyPrefix, user.ID))

	return nil
}

func (s *service) ResendVerificationEmail(email string) error {
	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return fmt.Errorf("user with that email not found")
	}

	if user.Verivied {
		return fmt.Errorf("this account is already verified")
	}

	// Generate a new token and expiration
	token := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)
	user.VerificationToken = &token
	user.VerificationTokenExpiresAt = &expiresAt

	// Save the new token to the database
	_, err = s.repository.Update(user)
	if err != nil {
		return fmt.Errorf("failed to update verification token: %w", err)
	}

	go sendVerificationEmail(user)
	return nil
}

func (s *service) ForgotPassword(email string) error {
	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return fmt.Errorf("pengguna dengan email tersebut tidak ditemukan")
	}

	// Hasilkan token baru yang berisi email pengguna.
	token, err := auth.GenerateTokenPassword(user.Email)
	if err != nil {
		return fmt.Errorf("gagal memulai proses reset password")
	}

	go sendForgotPasswordEmail(user.Email, token)
	return nil
}

func (s *service) ResetPassword(token string, newPassword string) error {
	claims, err := auth.ValidateTokenPassword(token)
	if err != nil {
		return fmt.Errorf("token reset password tidak valid atau kadaluarsa: %w", err)
	}

	user, err := s.repository.FindByEmail(claims.Email)
	if err != nil {
		return fmt.Errorf("pengguna terkait dengan token ini tidak ditemukan")
	}

	const bcryptCost = 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcryptCost)
	if err != nil {
		return fmt.Errorf("gagal memproses password baru")
	}

	user.Password = string(hashedPassword)
	if _, err := s.repository.ResetPassword(user); err != nil {
		return fmt.Errorf("gagal memperbarui password: %w", err)
	}

	return nil
}

// sendVerificationEmail adalah helper untuk mengirim email menggunakan SMTP.
func sendVerificationEmail(user User) {
	if user.VerificationToken == nil {
		log.Printf("Tidak ada token verifikasi untuk pengguna %s, email tidak dikirim.", user.Email)
		return
	}

	// Ambil konfigurasi dari environment variables
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_SENDER_EMAIL")
	appURL := os.Getenv("APP_URL")

	addr := fmt.Sprintf("%s:%s", host, port)
	auth := smtp.PlainAuth("", smtpUser, pass, host)
	verificationLink := fmt.Sprintf("%s/v1/verify-email?token=%s", appURL, *user.VerificationToken)

	subject := "Subject: Verifikasi Akun Anda\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`<html><body><h2>Selamat Datang!</h2><p>Terima kasih telah mendaftar. Silakan klik link di bawah ini untuk memverifikasi alamat email Anda:</p><p><a href="%s">Verifikasi Email Saya</a></p><p>Link ini akan kedaluwarsa dalam 24 jam.</p></body></html>`, verificationLink)
	msg := []byte(subject + mime + body)

	err := smtp.SendMail(addr, auth, from, []string{user.Email}, msg)
	if err != nil {
		log.Printf("Gagal mengirim email verifikasi ke %s: %v", user.Email, err)
	} else {
		log.Printf("Email verifikasi berhasil dikirim ke %s", user.Email)
	}
}

// sendVerificationEmail adalah helper untuk mengirim email menggunakan SMTP.
func sendForgotPasswordEmail(email, token string) {
	// Ambil konfigurasi dari environment variables
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_SENDER_EMAIL")
	resetAppURL := os.Getenv("FRONTEND_RESET_URL")
	if resetAppURL == "" {
		resetAppURL = os.Getenv("APP_URL") + "/v1/reset-password"
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	auth := smtp.PlainAuth("", smtpUser, pass, host)
	resetLink := fmt.Sprintf("%s?token=%s", resetAppURL, token)

	subject := "Subject: Reset Password\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`<html><body><h2>Reset Password</h2>
	<p>Anda meminta untuk mereset password Anda. Klik link di bawah ini untuk melanjutkan:</p>
	<p><a href="%s">Reset Password</a></p>
	<p>Jika Anda tidak meminta ini, abaikan email ini. Link ini akan kedaluwarsa dalam 24 jam.</p></body></html>`, resetLink)
	msg := []byte(subject + mime + body)

	err := smtp.SendMail(addr, auth, from, []string{email}, msg)
	if err != nil {
		log.Printf("Passeord %s gagal di perbarui karena: %v", email, err)
	} else {
		log.Printf("Password %s berhasil di ganti", email)
	}
}
