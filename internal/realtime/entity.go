package realtime

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	RoomID   string `gorm:"type:varchar(100);not null"`
	Content  string `gorm:"type:text;not null"`
	SenderID uint   `gorm:"not null"`
}
