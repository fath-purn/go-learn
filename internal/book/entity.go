package book

import "time"

type Book struct {
	ID          int
	Title       string
	Price       int
	Synopsis    string
	Description string
	Rating      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
