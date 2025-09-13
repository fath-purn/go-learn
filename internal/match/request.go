package match

type MatchRequest struct {
	UserID     int      // Akan diisi oleh handler, bukan dari request body
	Age        int      `json:"age" binding:"required"`
	Gender     Gender   `json:"gender" binding:"required,oneof=boy girl"`
	Interested Interest `json:"interested" binding:"required,oneof=boy girl any"`
	City       string   `json:"city" binding:"required"`
	Name       string   `json:"name" binding:"required"`
	Bio        string   `json:"bio" binding:"required"`
	ImageURL   string   `json:"image_url"`
}
