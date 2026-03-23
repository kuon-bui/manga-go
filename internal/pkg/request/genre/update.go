package genrerequest

type UpdateGenreRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
}
