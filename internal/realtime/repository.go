package realtime

import "gorm.io/gorm"

type Repository interface {
	Save(message Message) (Message, error)
	FindByRoomID(roomID string) ([]Message, error)
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

func (r *repository) FindByRoomID(roomID string) ([]Message, error) {
	var messages []Message

	err := r.db.Where("room_id = ?", roomID).Order("created_at asc").Find(&messages).Error
	return messages, err
}
