package book

type BookRequest struct {
	Title       string `json:"title" binding:"required"`
	Price       int    `json:"price" binding:"required,number"`
	Synopsis    string `json:"synopsis"`
	Description string `json:"description"`
	Rating      int    `json:"rating" binding:"required,number"`
}
