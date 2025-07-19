package short

import (
	"fmt"
	"log"
	"time"

	"github.com/patrickmn/go-cache"
)

type Service interface {
	GetAll() ([]Short, error)
	FindByID(ID int) (Short, error)
	FindByUrl(url string) (Short, error)
	Create(shortRequest ShortRequest) (Short, error)
	Update(ID int, short ShortRequest) (Short, error)
	Delete(ID int) error
}

type service struct {
	repository Repository
	cache      *cache.Cache
}

const (
	allShortsCacheKey       = "all_shorts"
	shortByIDCacheKeyPrefix = "short_by_url_"
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

func (s *service) Create(shortRequest ShortRequest) (Short, error) {
	cacheKey := fmt.Sprintf("%s%s", shortByIDCacheKeyPrefix, shortRequest.Shortened)

	// 1. Coba ambil dari cache
	if x, found := s.cache.Get(cacheKey); found {
		log.Printf("Service: Cache HIT for %d", shortRequest.Shortened)
		// Ingat untuk membersihkan password hash jika Anda menyimpannya dalam cache
		userFromCache := x.(Short)
		return userFromCache, fmt.Errorf("shortened URL already exists")
	}

	// check if the shortened URL already exists
	existingShort, _ := s.repository.FindByUrl(shortRequest.Shortened)
	if existingShort.Shortened != "" {
		return Short{}, fmt.Errorf("shortened URL already exists")
	}

	data := Short{
		Shortened: shortRequest.Shortened,
		Original:  shortRequest.Original,
	}

	created, err := s.repository.Create(data)
	if err != nil {
		return Short{}, err
	}

	s.cache.Delete(allShortsCacheKey)

	return created, nil
}

func (s *service) FindByUrl(url string) (Short, error) {
	cacheKey := fmt.Sprintf("%s%s", shortByIDCacheKeyPrefix, url)
	// 1. Coba ambil dari cache
	if x, found := s.cache.Get(cacheKey); found {
		log.Printf("Service: Cache HIT for URL %s", url)
		shortFromCache := x.(Short)
		return shortFromCache, nil
	}

	short, err := s.repository.FindByUrl(url)
	if err != nil {
		return Short{}, err
	}

	s.cache.Set(cacheKey, short, cache.DefaultExpiration)

	return short, nil
}

func (s *service) GetAll() ([]Short, error) {
	shorts, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}
	return shorts, nil
}

func (s *service) FindByID(ID int) (Short, error) {
	cacheKey := fmt.Sprintf("%s%d", shortByIDCacheKeyPrefix, ID)
	// 1. Coba ambil dari cache
	if x, found := s.cache.Get(cacheKey); found {
		log.Printf("Service: Cache HIT for short ID %d", ID)
		shortFromCache := x.(Short)
		return shortFromCache, nil
	}

	short, err := s.repository.FindByID(ID)
	if err != nil {
		return Short{}, err
	}

	s.cache.Set(cacheKey, short, cache.DefaultExpiration)

	return short, nil
}

func (s *service) Update(ID int, short ShortRequest) (Short, error) {
	data, err := s.repository.FindByID(ID)
	if err != nil {
		return Short{}, err
	}

	data.Shortened = short.Shortened
	data.Original = short.Original

	updatedBook, err := s.repository.Update(data)
	if err != nil {
		return Short{}, err
	}
	return updatedBook, nil
}

func (s *service) Delete(ID int) error {
	if err := s.repository.Delete(ID); err != nil {
		return err
	}
	return nil
}
