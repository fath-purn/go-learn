package short

type ShortRequest struct {
	Original  string `json:"original" binding:"required"`
	Shortened string `json:"shortened" binding:"required"`
}
