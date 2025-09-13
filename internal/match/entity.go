package match

import "time"

type Gender string

const (
	GenderCowok Gender = "boy"
	GenderCewek Gender = "girl"
)

type Interest string

const (
	InterestBoy  Interest = "boy"
	InterestGirl Interest = "girl"
	InterestAny  Interest = "any"
)

type Match struct {
	ID         int
	UserID     int `json:"user_id"` // Foreign Key
	Age        int
	Gender     Gender
	Interested Interest
	City       string
	Name       string
	Bio        string
	ImageURL   string `json:"image_url"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
