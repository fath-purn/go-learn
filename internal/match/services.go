package match

import (
	"fmt"
	"log"
	"time"

	"github.com/patrickmn/go-cache"
)

type Service interface {
	GetAll() ([]Match, error)
	FindByID(ID int) (Match, error)
	FindByCity(city string) ([]Match, error)
	Create(matchRequest MatchRequest) (Match, error)
	Update(ID int, match MatchRequest) (Match, error)
	Delete(ID int) error
}

type service struct {
	repository Repository
	cache      *cache.Cache
}

const (
	allMatchsCacheKey       = "all_matchs"
	matchByIDCacheKeyPrefix = "match_by_url_"
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

func (s *service) Create(matchRequest MatchRequest) (Match, error) {

	data := Match{
		UserID:     matchRequest.UserID,
		Age:        matchRequest.Age,
		Gender:     matchRequest.Gender,
		Interested: matchRequest.Interested,
		City:       matchRequest.City,
		Name:       matchRequest.Name,
		Bio:        matchRequest.Bio,
		ImageURL:   matchRequest.ImageURL,
	}

	created, err := s.repository.Create(data)
	if err != nil {
		return Match{}, err
	}

	s.cache.Delete(allMatchsCacheKey)

	return created, nil
}

func (s *service) Delete(ID int) error {
	// Pastikan match ada sebelum menghapus
	_, err := s.repository.FindByID(ID)
	if err != nil {
		return fmt.Errorf("match dengan ID %d tidak ditemukan: %w", ID, err)
	}

	if err := s.repository.Delete(ID); err != nil {
		return err
	}

	// Hapus cache yang relevan
	s.cache.Delete(allMatchsCacheKey)
	return nil
}

func (s *service) FindByCity(city string) ([]Match, error) {
	matches, err := s.repository.FindByCity(city)
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func (s *service) GetAll() ([]Match, error) {
	matchs, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}
	return matchs, nil
}

func (s *service) FindByID(ID int) (Match, error) {
	cacheKey := fmt.Sprintf("%s%d", matchByIDCacheKeyPrefix, ID)
	// 1. Coba ambil dari cache
	if x, found := s.cache.Get(cacheKey); found {
		log.Printf("Service: Cache HIT for match ID %d", ID)
		matchFromCache := x.(Match)
		return matchFromCache, nil
	}

	match, err := s.repository.FindByID(ID)
	if err != nil {
		return Match{}, err
	}

	s.cache.Set(cacheKey, match, cache.DefaultExpiration)

	return match, nil
}

func (s *service) Update(ID int, match MatchRequest) (Match, error) {
	data, err := s.repository.FindByID(ID)
	if err != nil {
		return Match{}, err
	}

	data.Age = match.Age
	data.Gender = match.Gender
	data.Interested = match.Interested
	data.City = match.City
	data.Name = match.Name
	data.Bio = match.Bio

	updatedMatch, err := s.repository.Update(data)
	if err != nil {
		return Match{}, err
	}

	// Setelah operasi tulis berhasil, hapus cache yang relevan.
	// 1. Hapus cache untuk semua match.
	s.cache.Delete(allMatchsCacheKey)
	// 2. Hapus cache untuk match spesifik yang baru saja di-update.
	s.cache.Delete(fmt.Sprintf("%s%d", matchByIDCacheKeyPrefix, ID))

	return updatedMatch, nil
}
