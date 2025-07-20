package realtime

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Content  string `gorm:"type:text;not null"`
	SenderID uint   `gorm:"not null"`
}
