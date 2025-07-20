package short

type ShortResponse struct {
	Original  string `json:"original" validate:"required"`
	Shortened string `json:"shortened" validate:"required"`
}
