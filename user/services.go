package user

import (
	"example/hello/auth"
	"fmt"
	"log"
	"regexp"
	"time"

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
}

type service struct {
	repository Repository
	cache      *cache.Cache
}

// Cache keys
const (
	allUsersCacheKey       = "all_users"
	userByIDCacheKeyPrefix = "user_by_id_"
	// ... tambahkan kunci lain sesuai kebutuhan
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

// FindByID implements Service.
func (s *service) FindByID(ID int) (User, error) {
	cacheKey := fmt.Sprintf("%s%d", userByIDCacheKeyPrefix, ID)

	// 1. Coba ambil dari cache
	if x, found := s.cache.Get(cacheKey); found {
		log.Printf("Service: Cache HIT for user ID %d", ID)
		// Ingat untuk membersihkan password hash jika Anda menyimpannya dalam cache
		userFromCache := x.(User)
		userFromCache.Password = ""
		return userFromCache, nil
	}

	log.Printf("Service: Cache MISS for user ID %d, fetching from repository", ID)

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
	}
	createdUser, err := s.repository.RegisterUser(user)
	if err != nil {
		return User{}, fmt.Errorf("error registering user: %w", err)
	}

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
		log.Println("Error finding user:", err)
		return "", User{}, fmt.Errorf("invalid email or password")
	}
	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
		log.Println("Error verifying password:", err)
		return "", User{}, fmt.Errorf("invalid email or password")
	}
	token, err := auth.GenerateToken(fmt.Sprintf("%d", foundUser.ID))
	if err != nil {
		log.Println("Service: Error generating token:", err)
		return "", User{}, fmt.Errorf("failed to generate authentication token")
	}
	foundUser.Password = "" // Clear password before returning
	return token, foundUser, nil
}

func (s *service) FindAll() ([]User, error) {
	if x, found := s.cache.Get(allUsersCacheKey); found {
		log.Println("Service: Cache HIT for all users")
		return x.([]User), nil // Type assert ke []User
	}

	log.Println("Service: Cache MISS for all users, fetching from repository")

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
