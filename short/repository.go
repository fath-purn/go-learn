package short

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]Short, error)
	FindByID(ID int) (Short, error)
	FindByUrl(url string) (Short, error)
	Create(short Short) (Short, error)
	Update(short Short) (Short, error)
	Delete(ID int) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindByUrl(url string) (Short, error) {
	var short Short
	if err := r.db.Where("shortened = ?", url).First(&short).Error; err != nil {
		return Short{}, err
	}
	return short, nil
}

func (r *repository) Create(short Short) (Short, error) {
	if err := r.db.Create(&short).Error; err != nil {
		return Short{}, err
	}

	return short, nil
}

func (r *repository) GetAll() ([]Short, error) {
	var shorts []Short
	if err := r.db.Find(&shorts).Error; err != nil {
		return nil, err
	}
	return shorts, nil
}

func (r *repository) FindByID(ID int) (Short, error) {
	var short Short
	if err := r.db.First(&short, ID).Error; err != nil {
		return Short{}, err
	}
	return short, nil
}

func (r *repository) Update(short Short) (Short, error) {
	if err := r.db.Save(&short).Error; err != nil {
		return Short{}, err
	}
	return short, nil
}

func (r *repository) Delete(ID int) error {
	if err := r.db.Delete(&Short{}, ID).Error; err != nil {
		return err
	}
	return nil
}
