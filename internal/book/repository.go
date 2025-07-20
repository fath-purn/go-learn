package book

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Book, error)
	FIndByID(ID int) (Book, error)
	Create(book Book) (Book, error)
	Update(book Book) (Book, error)
	Delete(ID int) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]Book, error) {
	var books []Book
	if err := r.db.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *repository) FIndByID(ID int) (Book, error) {
	var book Book
	if err := r.db.First(&book, ID).Error; err != nil {
		return Book{}, err
	}
	return book, nil
}

func (r *repository) Create(book Book) (Book, error) {
	if err := r.db.Create(&book).Error; err != nil {
		return Book{}, err
	}
	return book, nil
}

func (r *repository) Update(book Book) (Book, error) {
	if err := r.db.Save(&book).Error; err != nil {
		return Book{}, err
	}
	return book, nil
}

func (r *repository) Delete(ID int) error {
	if err := r.db.Delete(&Book{}, ID).Error; err != nil {
		return err
	}
	return nil
}
