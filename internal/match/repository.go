package match

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]Match, error)
	FindByID(ID int) (Match, error)
	FindByCity(city string) ([]Match, error)
	Create(match Match) (Match, error)
	Update(match Match) (Match, error)
	Delete(ID int) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(match Match) (Match, error) {
	if err := r.db.Create(&match).Error; err != nil {
		return Match{}, err
	}

	return match, nil
}

func (r *repository) GetAll() ([]Match, error) {
	var matches []Match
	if err := r.db.Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

func (r *repository) FindByID(ID int) (Match, error) {
	var match Match
	if err := r.db.First(&match, ID).Error; err != nil {
		return Match{}, err
	}
	return match, nil
}

func (r *repository) FindByCity(city string) ([]Match, error) {
	var matches []Match
	if err := r.db.Where("city = ?", city).Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

func (r *repository) Update(match Match) (Match, error) {
	if err := r.db.Save(&match).Error; err != nil {
		return Match{}, err
	}
	return match, nil
}

func (r *repository) Delete(ID int) error {
	if err := r.db.Delete(&Match{}, ID).Error; err != nil {
		return err
	}
	return nil
}
