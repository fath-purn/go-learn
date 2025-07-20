package short

import "time"

type Short struct {
	ID        int
	Original  string
	Shortened string
	CreatedAt time.Time
	UpdatedAt time.Time
}
