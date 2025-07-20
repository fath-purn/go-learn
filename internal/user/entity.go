package user

import "time"

type User struct {
	ID                         int
	Name                       string
	Email                      string
	Password                   string
	Phone                      string
	Verivied                   bool       `gorm:"default:false"`
	VerificationToken          *string    `gorm:"uniqueIndex"` // Pointer agar bisa NULL
	VerificationTokenExpiresAt *time.Time // Pointer agar bisa NULL
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}
