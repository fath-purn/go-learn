package realtime

import "gorm.io/gorm"

type Repository interface {
	Save(message Message) (Message, error)
	FindAll() ([]Message, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(message Message) (Message, error) {
	err := r.db.Create(&message).Error
	return message, err
}

func (r *repository) FindAll() ([]Message, error) {
	var messages []Message
	// Mengambil pesan dan mengurutkannya berdasarkan waktu pembuatan (dari yang terlama).
	err := r.db.Order("created_at asc").Find(&messages).Error
	return messages, err
}
